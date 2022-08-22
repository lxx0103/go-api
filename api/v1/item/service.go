package item

import (
	"encoding/json"
	"errors"
	"go-api/api/v1/history"
	"go-api/api/v1/setting"
	"go-api/core/database"
	"go-api/core/queue"
	"strconv"
	"time"

	"github.com/rs/xid"
)

type itemService struct {
}

func NewItemService() *itemService {
	return &itemService{}
}

func (s *itemService) NewItem(info ItemNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckSKUConfict("", info.OrganizationID, info.SKU)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "item SKU exists"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetUnitByID(info.OrganizationID, info.UnitID)
	if err != nil {
		return nil, err
	}
	if info.BrandID != "" {
		_, err = settingService.GetBrandByID(info.OrganizationID, info.BrandID)
		if err != nil {
			return nil, err
		}
	}
	if info.ManufacturerID != "" {
		_, err = settingService.GetManufacturerByID(info.OrganizationID, info.ManufacturerID)
		if err != nil {
			return nil, err
		}
	}
	if info.DefaultVendorID != "" {
		_, err = settingService.GetVendorByID(info.OrganizationID, info.DefaultVendorID)
		if err != nil {
			return nil, err
		}
	}
	var item Item
	item.ItemID = "item-" + xid.New().String()
	item.OrganizationID = info.OrganizationID
	item.SKU = info.SKU
	item.Name = info.Name
	item.UnitID = info.UnitID
	item.ManufacturerID = info.ManufacturerID
	item.BrandID = info.BrandID
	item.WeightUnit = info.WeightUnit
	item.Weight = info.Weight
	item.DimensionUnit = info.DimensionUnit
	item.Length = info.Length
	item.Width = info.Width
	item.Height = info.Height
	item.SellingPrice = info.SellingPrice
	item.CostPrice = info.CostPrice
	item.OpenningStock = info.OpenningStock
	item.OpenningStockRate = info.OpenningStockRate
	item.ReorderStock = info.ReorderStock
	item.StockOnHand = info.OpenningStock
	item.StockAvailable = info.OpenningStock
	item.StockPicking = 0
	item.StockPacking = 0
	item.DefaultVendorID = info.DefaultVendorID
	item.Description = info.Description
	item.TrackLocation = info.TrackLocation
	item.Status = info.Status
	item.Created = time.Now()
	item.CreatedBy = info.Email
	item.Updated = time.Now()
	item.UpdatedBy = info.Email

	err = repo.CreateItem(item)
	if err != nil {
		msg := "create itemerror: " + err.Error()
		return nil, errors.New(msg)
	}
	rabbit, _ := queue.GetConn()
	var newEvent history.NewHistoryCreated
	newEvent.HistoryType = "item"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = item.ItemID
	newEvent.Description = "Item Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	var newBatchEvent NewBatchCreated
	newBatchEvent.Type = "NewItem"
	newBatchEvent.Quantity = item.OpenningStock
	newBatchEvent.Balance = item.OpenningStock
	newBatchEvent.ReferenceID = item.ItemID
	newBatchEvent.ItemID = item.ItemID
	newBatchEvent.OrganizationID = info.OrganizationID
	newBatchEvent.Email = info.Email
	msg2, _ := json.Marshal(newBatchEvent)
	err = rabbit.Publish("NewBatchCreated", msg2)
	if err != nil {
		msg := "create event NewBatchCreated error"
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &item.ItemID, err
}

func (s *itemService) GetItemList(filter ItemFilter) (int, *[]ItemResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	count, err := query.GetItemCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetItemList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *itemService) UpdateItem(itemID string, info ItemNew) (*ItemResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckSKUConfict(itemID, info.OrganizationID, info.SKU)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "item SKU conflict"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetUnitByID(info.OrganizationID, info.UnitID)
	if err != nil {
		return nil, err
	}
	if info.BrandID != "" {
		_, err = settingService.GetBrandByID(info.OrganizationID, info.BrandID)
		if err != nil {
			return nil, err
		}
	}
	if info.ManufacturerID != "" {
		_, err = settingService.GetManufacturerByID(info.OrganizationID, info.ManufacturerID)
		if err != nil {
			return nil, err
		}
	}
	if info.DefaultVendorID != "" {
		_, err = settingService.GetVendorByID(info.OrganizationID, info.DefaultVendorID)
		if err != nil {
			return nil, err
		}
	}
	oldItem, err := repo.GetItemByID(itemID, info.OrganizationID)
	if err != nil {
		msg := "Item not exist"
		return nil, errors.New(msg)
	}
	if oldItem.OpenningStock != 0 {
		openningBatch, err := repo.GetItemOpenningBatch(itemID, info.OrganizationID)
		if err != nil {
			msg := "get Item Openning Batch error" + err.Error()
			return nil, errors.New(msg)
		}
		newBalance := openningBatch.Balance - oldItem.OpenningStock + info.OpenningStock
		if newBalance < 0 {
			msg := "openning stock already sold, minimum openning stock should be " + strconv.Itoa(openningBatch.Quantity-openningBatch.Balance)
			return nil, errors.New(msg)
		}
		var itemBatch ItemBatch
		itemBatch.Quantity = info.OpenningStock
		itemBatch.Balance = newBalance
		itemBatch.UpdatedBy = info.Email
		itemBatch.Updated = time.Now()
		if newBalance == 0 {
			itemBatch.Status = 2
		} else {
			itemBatch.Status = 1
		}
		err = repo.UpdateItemBatch(openningBatch.BatchID, itemBatch)
		if err != nil {
			msg := "update item batch error"
			return nil, errors.New(msg)
		}
	}
	var item Item
	item.SKU = info.SKU
	item.Name = info.Name
	item.UnitID = info.UnitID
	item.ManufacturerID = info.ManufacturerID
	item.BrandID = info.BrandID
	item.WeightUnit = info.WeightUnit
	item.Weight = info.Weight
	item.DimensionUnit = info.DimensionUnit
	item.Length = info.Length
	item.Width = info.Width
	item.Height = info.Height
	item.SellingPrice = info.SellingPrice
	item.CostPrice = info.CostPrice
	item.OpenningStock = info.OpenningStock
	item.OpenningStockRate = info.OpenningStockRate
	item.ReorderStock = info.ReorderStock
	item.StockOnHand = oldItem.StockOnHand - oldItem.OpenningStock + info.OpenningStock
	item.StockAvailable = oldItem.StockAvailable - oldItem.OpenningStock + info.OpenningStock
	item.DefaultVendorID = info.DefaultVendorID
	item.Description = info.Description
	item.TrackLocation = info.TrackLocation
	item.Status = info.Status
	item.Updated = time.Now()
	item.UpdatedBy = info.User
	err = repo.UpdateItem(itemID, item)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetItemByID(itemID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	var newEvent history.NewHistoryCreated
	newEvent.HistoryType = "item"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = itemID
	newEvent.Description = "Item Updated"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	tx.Commit()
	return res, err
}

func (s *itemService) GetItemByID(organizationID, id string) (*ItemResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	item, err := query.GetItemByID(organizationID, id)
	if err != nil {
		msg := "get item error: " + err.Error()
		return nil, errors.New(msg)
	}
	return item, nil
}

func (s *itemService) DeleteItem(itemID, organizationID, email, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	oldItem, err := repo.GetItemByID(itemID, organizationID)
	if err != nil {
		msg := "Item not exist"
		return errors.New(msg)
	}
	if oldItem.OpenningStock != 0 {
		openningBatch, err := repo.GetItemOpenningBatch(itemID, organizationID)
		if err != nil {
			msg := "get Item Openning Batch error" + err.Error()
			return errors.New(msg)
		}
		if openningBatch.Quantity != openningBatch.Balance {
			msg := "openning stock already sold, can not delete"
			return errors.New(msg)
		}
		var itemBatch ItemBatch
		itemBatch.Quantity = openningBatch.Quantity
		itemBatch.Balance = openningBatch.Balance
		itemBatch.Status = -1
		itemBatch.UpdatedBy = email
		itemBatch.Updated = time.Now()
		err = repo.UpdateItemBatch(openningBatch.BatchID, itemBatch)
		if err != nil {
			msg := "update item batch error"
			return errors.New(msg)
		}
	}
	err = repo.DeleteItem(itemID, email)
	if err != nil {
		return err
	}
	var newEvent history.NewHistoryCreated
	newEvent.HistoryType = "item"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = itemID
	newEvent.Description = "Item Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

//Barcode

func (s *itemService) GetBarcodeList(filter BarcodeFilter) (int, *[]BarcodeResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	count, err := query.GetBarcodeCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBarcodeList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *itemService) NewBarcode(info BarcodeNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckBarcodeConfict("", info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "barcode conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetItemByID(info.ItemID, info.OrganizationID)
	if err != nil {
		msg := "Item not exist"
		return nil, errors.New(msg)
	}
	var barcode Barcode
	barcode.BarcodeID = "bar-" + xid.New().String()
	barcode.OrganizationID = info.OrganizationID
	barcode.Code = info.Code
	barcode.ItemID = info.ItemID
	barcode.Quantity = info.Quantity
	barcode.Status = info.Status
	barcode.Created = time.Now()
	barcode.CreatedBy = info.User
	barcode.Updated = time.Now()
	barcode.UpdatedBy = info.User

	err = repo.CreateBarcode(barcode)
	if err != nil {
		msg := "create barcode error: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &barcode.BarcodeID, err
}

func (s *itemService) UpdateBarcode(barcodeID string, info BarcodeNew) (*BarcodeResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckBarcodeConfict(barcodeID, info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "barcode conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetItemByID(info.ItemID, info.OrganizationID)
	if err != nil {
		msg := "Item not exist"
		return nil, errors.New(msg)
	}
	oldBarcode, err := repo.GetBarcodeByID(barcodeID)
	if err != nil {
		msg := "Barcode not exist"
		return nil, errors.New(msg)
	}
	if oldBarcode.OrganizationID != info.OrganizationID {
		msg := "Barcode not exist"
		return nil, errors.New(msg)
	}
	var barcode Barcode
	barcode.Code = info.Code
	barcode.ItemID = info.ItemID
	barcode.Quantity = info.Quantity
	barcode.Status = info.Status
	barcode.Updated = time.Now()
	barcode.UpdatedBy = info.User
	err = repo.UpdateBarcode(barcodeID, barcode)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetBarcodeByID(barcodeID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, err
}

func (s *itemService) GetBarcodeByID(organizationID, id string) (*BarcodeResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	unit, err := query.GetBarcodeByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *itemService) DeleteBarcode(barcodeID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	oldBarcode, err := repo.GetBarcodeByID(barcodeID)
	if err != nil {
		msg := "Barcode not exist"
		return errors.New(msg)
	}
	if oldBarcode.OrganizationID != organizationID {
		msg := "Barcode not exist"
		return errors.New(msg)
	}
	err = repo.DeleteBarcode(barcodeID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
