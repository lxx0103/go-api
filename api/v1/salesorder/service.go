package salesorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-api/api/v1/common"
	"go-api/api/v1/item"
	"go-api/api/v1/setting"
	"go-api/api/v1/warehouse"
	"go-api/core/database"
	"go-api/core/queue"
	"strings"
	"time"

	"github.com/rs/xid"
)

type salesorderService struct {
}

func NewSalesorderService() *salesorderService {
	return &salesorderService{}
}

func (s *salesorderService) NewSalesorder(info SalesorderNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckSONumberConfict("", info.OrganizationID, info.SalesorderNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "salesorder number exists"
		return nil, errors.New(msg)
	}
	soID := "so-" + xid.New().String()
	settingService := setting.NewSettingService()
	_, err = settingService.GetCustomerByID(info.OrganizationID, info.CustomerID)
	if err != nil {
		return nil, err
	}
	itemCount := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemService := item.NewItemService()
	for _, item := range info.Items {
		_, err = itemService.GetItemByID(info.OrganizationID, item.ItemID)
		if err != nil {
			return nil, err
		}
		taxValue := 0.0
		if item.TaxID != "" {
			tax, err := settingService.GetTaxByID(info.OrganizationID, item.TaxID)
			if err != nil {
				return nil, err
			}
			taxValue = tax.TaxValue
		}
		itemCount += item.Quantity
		itemTotal += item.Rate * float64(item.Quantity)
		taxTotal += item.Rate * float64(item.Quantity) * taxValue / 100
		var soItem SalesorderItem
		soItem.OrganizationID = info.OrganizationID
		soItem.SalesorderID = soID
		soItem.SalesorderItemID = "soi-" + xid.New().String()
		soItem.ItemID = item.ItemID
		soItem.Quantity = item.Quantity
		soItem.Rate = item.Rate
		soItem.TaxID = item.TaxID
		soItem.TaxValue = taxValue
		soItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
		soItem.Amount = float64(item.Quantity) * item.Rate
		soItem.QuantityInvoiced = 0
		soItem.QuantityPicked = 0
		soItem.QuantityPacked = 0
		soItem.QuantityShipped = 0
		soItem.Status = 1
		soItem.Created = time.Now()
		soItem.CreatedBy = info.Email
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.CreateSalesorderItem(soItem)
		if err != nil {
			msg := "create salesorder item error: "
			return nil, errors.New(msg)
		}
	}
	var salesorder Salesorder
	salesorder.SalesorderID = soID
	salesorder.OrganizationID = info.OrganizationID
	salesorder.SalesorderNumber = info.SalesorderNumber
	salesorder.SalesorderDate = info.SalesorderDate
	salesorder.ExpectedShipmentDate = info.ExpectedShipmentDate
	salesorder.CustomerID = info.CustomerID
	salesorder.ItemCount = itemCount
	salesorder.Subtotal = itemTotal
	salesorder.TaxTotal = taxTotal
	salesorder.DiscountType = info.DiscountType
	salesorder.DiscountValue = info.DiscountValue
	salesorder.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		salesorder.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		salesorder.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		salesorder.Total = itemTotal + info.ShippingFee + taxTotal
	}
	salesorder.Notes = info.Notes
	salesorder.Status = 1         //Draft
	salesorder.InvoiceStatus = 1  //not invoiced
	salesorder.PickingStatus = 1  //not picked
	salesorder.PackingStatus = 1  //not packed
	salesorder.ShippingStatus = 1 //not shipped
	salesorder.Created = time.Now()
	salesorder.CreatedBy = info.Email
	salesorder.Updated = time.Now()
	salesorder.UpdatedBy = info.Email

	err = repo.CreateSalesorder(salesorder)
	if err != nil {
		msg := "create salesorder error: "
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = soID
	newEvent.Description = "Sales Order Created"
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
	return &salesorder.SalesorderID, err
}

func (s *salesorderService) GetSalesorderList(filter SalesorderFilter) (int, *[]SalesorderResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetSalesorderCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetSalesorderList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) UpdateSalesorder(salesorderID string, info SalesorderNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckSONumberConfict(salesorderID, info.OrganizationID, info.SalesorderNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "salesorder number conflict"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetCustomerByID(info.OrganizationID, info.CustomerID)
	if err != nil {
		return nil, err
	}
	oldSalesorder, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "Salesorder not exist"
		return nil, errors.New(msg)
	}
	err = repo.DeleteSalesorder(salesorderID, info.User)
	if err != nil {
		msg := "Salesorder Update error"
		return nil, errors.New(msg)
	}
	itemCount := 0
	quantityInvoiced := 0
	quantityPicked := 0
	quantityPacked := 0
	quantityShipped := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemService := item.NewItemService()
	for _, item := range info.Items {
		_, err = itemService.GetItemByID(info.OrganizationID, item.ItemID)
		if err != nil {
			return nil, err
		}
		taxValue := 0.0
		if item.TaxID != "" {
			tax, err := settingService.GetTaxByID(info.OrganizationID, item.TaxID)
			if err != nil {
				return nil, err
			}
			taxValue = tax.TaxValue
		}
		if item.SalesorderItemID != "" {
			oldItem, err := repo.GetSalesorderItemByIDAll(info.OrganizationID, salesorderID, item.SalesorderItemID)
			if err != nil {
				msg := "Salesorder Item not exist"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityInvoiced > item.Quantity {
				msg := "can not set quantity lower than quantity invoiced"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityPicked > item.Quantity {
				msg := "can not set quantity lower than quantity picked"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityPacked > item.Quantity {
				msg := "can not set quantity lower than quantity packed"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityShipped > item.Quantity {
				msg := "can not set quantity lower than quantity shipped"
				return nil, errors.New(msg)
			}
			quantityInvoiced += oldItem.QuantityInvoiced
			quantityPicked += oldItem.QuantityPicked
			quantityPacked += oldItem.QuantityPacked
			quantityShipped += oldItem.QuantityShipped
			var soItem SalesorderItem
			soItem.Quantity = item.Quantity
			soItem.Rate = item.Rate
			soItem.TaxID = item.TaxID
			soItem.TaxValue = taxValue
			soItem.Amount = float64(item.Quantity) * item.Rate
			soItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
			soItem.Status = 1
			soItem.Updated = time.Now()
			soItem.UpdatedBy = info.Email
			err = repo.UpdateSalesorderItem(item.SalesorderItemID, soItem)
			if err != nil {
				msg := "update salesorder item error: "
				return nil, errors.New(msg)
			}
		} else {
			var soItem SalesorderItem
			soItem.OrganizationID = info.OrganizationID
			soItem.SalesorderID = salesorderID
			soItem.SalesorderItemID = "soi-" + xid.New().String()
			soItem.ItemID = item.ItemID
			soItem.Quantity = item.Quantity
			soItem.Rate = item.Rate
			soItem.TaxID = item.TaxID
			soItem.TaxValue = taxValue
			soItem.Amount = float64(item.Quantity) * item.Rate
			soItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
			soItem.QuantityInvoiced = 0
			soItem.QuantityPicked = 0
			soItem.QuantityPacked = 0
			soItem.QuantityShipped = 0
			soItem.Status = 1
			soItem.Created = time.Now()
			soItem.CreatedBy = info.Email
			soItem.Updated = time.Now()
			soItem.UpdatedBy = info.Email
			err = repo.CreateSalesorderItem(soItem)
			if err != nil {
				msg := "create salesorder item error: "
				return nil, errors.New(msg)
			}
		}
		itemCount += item.Quantity
		itemTotal += item.Rate * float64(item.Quantity)
		taxTotal += item.Rate * float64(item.Quantity) * taxValue / 100
	}
	itemDeletedError, err := repo.CheckSOItem(salesorderID, info.OrganizationID)
	if err != nil {
		msg := "check salesorder item error: "
		return nil, errors.New(msg)
	}
	if itemDeletedError {
		msg := "item invoiced or picked or packed or shipped can not be delete"
		return nil, errors.New(msg)
	}
	var salesorder Salesorder
	salesorder.SalesorderNumber = info.SalesorderNumber
	salesorder.SalesorderDate = info.SalesorderDate
	salesorder.ExpectedShipmentDate = info.ExpectedShipmentDate
	salesorder.CustomerID = info.CustomerID
	salesorder.ItemCount = itemCount
	salesorder.Subtotal = itemTotal
	salesorder.TaxTotal = taxTotal
	salesorder.DiscountType = info.DiscountType
	salesorder.DiscountValue = info.DiscountValue
	salesorder.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		salesorder.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		salesorder.Total = itemTotal + taxTotal - info.DiscountValue + info.ShippingFee
	} else {
		salesorder.Total = itemTotal + taxTotal + info.ShippingFee
	}
	salesorder.Notes = info.Notes
	if quantityInvoiced > 0 {
		if quantityInvoiced == itemCount {
			salesorder.InvoiceStatus = 3
		} else {
			salesorder.InvoiceStatus = 2
		}
	} else {
		salesorder.InvoiceStatus = 1
	}
	if quantityPicked > 0 {
		if quantityPicked == itemCount {
			salesorder.PickingStatus = 3
		} else {
			salesorder.PickingStatus = 2
		}
	} else {
		salesorder.PickingStatus = 1
	}
	if quantityPacked > 0 {
		if quantityPacked == itemCount {
			salesorder.PackingStatus = 3
		} else {
			salesorder.PackingStatus = 2
		}
	} else {
		salesorder.PackingStatus = 1
	}
	if quantityShipped > 0 {
		if quantityPicked == itemCount {
			salesorder.ShippingStatus = 3
		} else {
			salesorder.ShippingStatus = 2
		}
	} else {
		salesorder.ShippingStatus = 1
	}
	if salesorder.InvoiceStatus == 3 && salesorder.PickingStatus == 3 && salesorder.PackingStatus == 3 && salesorder.ShippingStatus == 3 {
		salesorder.Status = 3 //CLOSE
	} else {
		if oldSalesorder.Status == 3 {
			salesorder.Status = 2
		} else {
			salesorder.Status = oldSalesorder.Status
		}
	}
	salesorder.Updated = time.Now()
	salesorder.UpdatedBy = info.Email

	err = repo.UpdateSalesorder(salesorderID, salesorder)
	if err != nil {
		msg := "update salesorder error: "
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Sales Order Updated"
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
	return &salesorderID, err
}

func (s *salesorderService) GetSalesorderByID(organizationID, id string) (*SalesorderResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	salesorder, err := query.GetSalesorderByID(organizationID, id)
	if err != nil {
		msg := "get salesorder error: "
		return nil, errors.New(msg)
	}
	return salesorder, nil
}

func (s *salesorderService) DeleteSalesorder(salesorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	so, err := repo.GetSalesorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "Salesorder not exist"
		return errors.New(msg)
	}
	if so.InvoiceStatus != 1 {
		msg := "Salesorder invoiced can not be deleted"
		return errors.New(msg)
	}
	if so.PickingStatus != 1 {
		msg := "Salesorder picked can not be deleted"
		return errors.New(msg)
	}
	if so.PackingStatus != 1 {
		msg := "Salesorder packed can not be deleted"
		return errors.New(msg)
	}
	if so.ShippingStatus != 1 {
		msg := "Salesorder shipped can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteSalesorder(salesorderID, email)
	if err != nil {
		return err
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Sales Order Deleted"
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

func (s *salesorderService) GetSalesorderItemList(salesorderID, organizationID string) (*[]SalesorderItemResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetSalesorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "get salesorder error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetSalesorderItemList(salesorderID)
	return list, err
}

func (s *salesorderService) ConfirmSalesorder(salesorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	oldSalesorder, err := repo.GetSalesorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "Salesorder not exist"
		return errors.New(msg)
	}
	if oldSalesorder.Status != 1 {
		msg := "Salesorder status error"
		return errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(salesorderID, 2, email) //CONFIRMED
	if err != nil {
		msg := "update salesorder error: "
		return errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Sales Order Confirmed"
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
	return err
}

// picking

func (s *salesorderService) NewPickingorder(salesorderID string, info PickingorderNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckPickingorderNumberConfict("", info.OrganizationID, info.PickingorderNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "picking order number exists"
		return nil, errors.New(msg)
	}
	pickingorderID := "pic-" + xid.New().String()
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	for _, itemRow := range info.Items {
		oldSoItem, err := repo.GetSalesorderItemByID(info.OrganizationID, salesorderID, itemRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return nil, errors.New(msg)
		}
		itemInfo, err := itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
		if err != nil {
			msg := "item not exist"
			return nil, errors.New(msg)
		}
		pickingorderItemID := "pii-" + xid.New().String()
		if itemInfo.TrackLocation == 1 {
			if itemInfo.StockAvailable < itemRow.Quantity {
				msg := "no enough stock to pick"
				return nil, errors.New(msg)
			}
			quantityToPick := itemRow.Quantity
			for quantityToPick > 0 {
				nextBatch, err := itemRepo.GetItemNextBatch(itemRow.ItemID, info.OrganizationID)
				if err != nil {
					msg := "get next batch error"
					return nil, errors.New(msg)
				}
				if nextBatch.Balance >= quantityToPick {
					err = itemRepo.PickItem(nextBatch.BatchID, quantityToPick, info.Email)
					if err != nil {
						msg := "pick item from batch error"
						return nil, errors.New(msg)
					}
					var pickingorderLog PickingorderLog
					pickingorderLog.PickingorderLogID = "pil-" + xid.New().String()
					pickingorderLog.OrganizationID = info.OrganizationID
					pickingorderLog.PickingorderID = pickingorderID
					pickingorderLog.SalesorderID = salesorderID
					pickingorderLog.SalesorderItemID = oldSoItem.SalesorderItemID
					pickingorderLog.PickingorderItemID = pickingorderItemID
					pickingorderLog.LocationID = nextBatch.LocationID
					pickingorderLog.BatchID = nextBatch.BatchID
					pickingorderLog.ItemID = itemRow.ItemID
					pickingorderLog.Quantity = quantityToPick
					pickingorderLog.Status = 1
					pickingorderLog.Created = time.Now()
					pickingorderLog.CreatedBy = info.Email
					pickingorderLog.Updated = time.Now()
					pickingorderLog.UpdatedBy = info.Email
					err = repo.CreatePickingorderLog(pickingorderLog)
					if err != nil {
						msg := "create picking order log error1"
						return nil, errors.New(msg)
					}
					err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, quantityToPick, info.Email)
					if err != nil {
						msg := "update location canpick error: "
						return nil, errors.New(msg)
					}
					quantityToPick = 0
				} else {
					err = itemRepo.PickItem(nextBatch.BatchID, nextBatch.Balance, info.Email)
					if err != nil {
						msg := "pick item from batch error"
						return nil, errors.New(msg)
					}
					var pickingorderLog PickingorderLog
					pickingorderLog.PickingorderLogID = "pil-" + xid.New().String()
					pickingorderLog.OrganizationID = info.OrganizationID
					pickingorderLog.PickingorderID = pickingorderID
					pickingorderLog.SalesorderID = salesorderID
					pickingorderLog.SalesorderItemID = oldSoItem.SalesorderItemID
					pickingorderLog.PickingorderItemID = pickingorderItemID
					pickingorderLog.LocationID = nextBatch.LocationID
					pickingorderLog.BatchID = nextBatch.BatchID
					pickingorderLog.ItemID = itemRow.ItemID
					pickingorderLog.Quantity = nextBatch.Balance
					pickingorderLog.Status = 1
					pickingorderLog.Created = time.Now()
					pickingorderLog.CreatedBy = info.Email
					pickingorderLog.Updated = time.Now()
					pickingorderLog.UpdatedBy = info.Email
					err = repo.CreatePickingorderLog(pickingorderLog)
					if err != nil {
						msg := "create picking order log error"
						return nil, errors.New(msg)
					}
					err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, nextBatch.Balance, info.Email)
					if err != nil {
						msg := "update location canpick error: "
						return nil, errors.New(msg)
					}
					quantityToPick = quantityToPick - nextBatch.Balance
				}
			}
		}
		if oldSoItem.Quantity < oldSoItem.QuantityPicked+itemRow.Quantity {
			msg := "pick quantity greater than not picked"
			return nil, errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityPicked = oldSoItem.QuantityPicked + itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.PickSalesorderItem(soItem)
		if err != nil {
			msg := "pick salesorder item error: "
			return nil, errors.New(msg)
		}
		var poItem PickingorderItem
		poItem.OrganizationID = info.OrganizationID
		poItem.PickingorderID = pickingorderID
		poItem.SalesorderItemID = oldSoItem.SalesorderItemID
		poItem.PickingorderItemID = pickingorderItemID
		poItem.ItemID = oldSoItem.ItemID
		poItem.Quantity = itemRow.Quantity
		poItem.Status = 1
		poItem.CreatedBy = info.Email
		poItem.Created = time.Now()
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

		err = repo.CreatePickingorderItem(poItem)
		if err != nil {
			msg := "create picking order item error: "
			return nil, errors.New(msg)
		}
		err = itemRepo.UpdateItemPickingStock(itemRow.ItemID, itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: "
			return nil, errors.New(msg)
		}
	}
	logs, err := repo.GetPickingorderLogSum(pickingorderID)
	if err != nil {
		msg := "get picking order logs  error: "
		return nil, errors.New(msg)
	}
	fmt.Println(logs)
	for _, logRow := range *logs {
		var pickingorderDetail PickingorderDetail
		pickingorderDetail.PickingorderDetailID = "pid-" + xid.New().String()
		pickingorderDetail.OrganizationID = logRow.OrganizationID
		pickingorderDetail.PickingorderID = logRow.PickingorderID
		pickingorderDetail.ItemID = logRow.ItemID
		pickingorderDetail.LocationID = logRow.LocationID
		pickingorderDetail.Quantity = logRow.Quantity
		pickingorderDetail.QuantityPicked = 0
		pickingorderDetail.Status = 1
		pickingorderDetail.CreatedBy = info.Email
		pickingorderDetail.Created = time.Now()
		pickingorderDetail.Updated = time.Now()
		pickingorderDetail.UpdatedBy = info.Email
		err = repo.CreatePickingorderDetail(pickingorderDetail)
		if err != nil {
			msg := "create picking order detail error: "
			return nil, errors.New(msg)
		}
	}
	var pickingorder Pickingorder
	pickingorder.SalesorderID = salesorderID
	pickingorder.PickingorderID = pickingorderID
	pickingorder.PickingorderNumber = info.PickingorderNumber
	pickingorder.PickingorderDate = info.PickingorderDate
	pickingorder.OrganizationID = info.OrganizationID
	pickingorder.Notes = info.Notes
	pickingorder.Status = 1
	pickingorder.Created = time.Now()
	pickingorder.CreatedBy = info.Email
	pickingorder.Updated = time.Now()
	pickingorder.UpdatedBy = info.Email
	err = repo.CreatePickingorder(pickingorder)
	if err != nil {
		msg := "create picking order error: "
		return nil, errors.New(msg)
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order error: "
		return nil, errors.New(msg)
	}
	receivedCount, err := repo.GetSalesorderPickedCount(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order received count error: "
		return nil, errors.New(msg)
	}
	picking_status := 1
	if so.ItemCount == receivedCount {
		picking_status = 3
	} else {
		picking_status = 2
	}
	err = repo.UpdateSalesorderPickingStatus(salesorderID, picking_status, info.Email)
	if err != nil {
		msg := "update sales order receive status error: "
		return nil, errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(salesorderID, 2, info.Email)
	if err != nil {
		msg := "update sales order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Picking Order Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &pickingorderID, err
}

func (s *salesorderService) GetPickingorderList(filter PickingorderFilter) (int, *[]PickingorderResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetPickingorderCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPickingorderList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) GetPickingorderItemList(salesorderID, organizationID string) (*[]PickingorderItemResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetPickingorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "get pickingorder error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetPickingorderItemList(salesorderID)
	return list, err
}

func (s *salesorderService) GetPickingorderDetailList(pickingorderID, organizationID string) (*[]PickingorderDetailResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetPickingorderByID(organizationID, pickingorderID)
	if err != nil {
		msg := "get pickingorder error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetPickingorderDetailList(pickingorderID)
	return list, err
}

func (s *salesorderService) BatchPickingorder(info PickingorderBatch) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckPickingorderNumberConfict("", info.OrganizationID, info.PickingorderNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "picking order number exists"
		return nil, errors.New(msg)
	}
	pickingorderID := "pic-" + xid.New().String()
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	var msgs [][]byte
	for _, soID := range info.SOID {
		salesorder, err := repo.GetSalesorderByID(info.OrganizationID, soID)
		if err != nil {
			msg := "get salesorder error: "
			return nil, errors.New(msg)
		}
		items, err := repo.GetSalesorderItemList(info.OrganizationID, soID)
		if err != nil {
			msg := "get salesorder items error: "
			return nil, errors.New(msg)
		}
		for _, itemRow := range *items {
			toPick := itemRow.Quantity - itemRow.QuantityPicked
			itemInfo, err := itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
			if err != nil {
				msg := "item not exist"
				return nil, errors.New(msg)
			}
			pickingorderItemID := "pii-" + xid.New().String()
			if itemInfo.TrackLocation == 1 {
				if itemInfo.StockAvailable < toPick {
					msg := "no enough stock for item: " + itemInfo.Name + " in salesorder :" + salesorder.SalesorderNumber
					return nil, errors.New(msg)
				}
				quantityToPick := toPick
				for quantityToPick > 0 {
					nextBatch, err := itemRepo.GetItemNextBatch(itemRow.ItemID, info.OrganizationID)
					if err != nil {
						msg := "get next batch error"
						return nil, errors.New(msg)
					}
					if nextBatch.Balance >= quantityToPick {
						err = itemRepo.PickItem(nextBatch.BatchID, quantityToPick, info.Email)
						if err != nil {
							msg := "pick item from batch error"
							return nil, errors.New(msg)
						}
						var pickingorderLog PickingorderLog
						pickingorderLog.PickingorderLogID = "pil-" + xid.New().String()
						pickingorderLog.OrganizationID = info.OrganizationID
						pickingorderLog.PickingorderID = pickingorderID
						pickingorderLog.SalesorderID = soID
						pickingorderLog.SalesorderItemID = itemRow.SalesorderItemID
						pickingorderLog.PickingorderItemID = pickingorderItemID
						pickingorderLog.LocationID = nextBatch.LocationID
						pickingorderLog.BatchID = nextBatch.BatchID
						pickingorderLog.ItemID = itemRow.ItemID
						pickingorderLog.Quantity = quantityToPick
						pickingorderLog.Status = 1
						pickingorderLog.Created = time.Now()
						pickingorderLog.CreatedBy = info.Email
						pickingorderLog.Updated = time.Now()
						pickingorderLog.UpdatedBy = info.Email
						err = repo.CreatePickingorderLog(pickingorderLog)
						if err != nil {
							msg := "create picking order detail error1"
							return nil, errors.New(msg)
						}
						err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, quantityToPick, info.Email)
						if err != nil {
							msg := "update location canpick error: "
							return nil, errors.New(msg)
						}
						quantityToPick = 0
					} else {
						err = itemRepo.PickItem(nextBatch.BatchID, nextBatch.Balance, info.Email)
						if err != nil {
							msg := "pick item from batch error"
							return nil, errors.New(msg)
						}
						var pickingorderLog PickingorderLog
						pickingorderLog.PickingorderLogID = "pil-" + xid.New().String()
						pickingorderLog.OrganizationID = info.OrganizationID
						pickingorderLog.PickingorderID = pickingorderID
						pickingorderLog.SalesorderID = soID
						pickingorderLog.SalesorderItemID = itemRow.SalesorderItemID
						pickingorderLog.PickingorderItemID = pickingorderItemID
						pickingorderLog.LocationID = nextBatch.LocationID
						pickingorderLog.BatchID = nextBatch.BatchID
						pickingorderLog.ItemID = itemRow.ItemID
						pickingorderLog.Quantity = nextBatch.Balance
						pickingorderLog.Status = 1
						pickingorderLog.Created = time.Now()
						pickingorderLog.CreatedBy = info.Email
						pickingorderLog.Updated = time.Now()
						pickingorderLog.UpdatedBy = info.Email
						err = repo.CreatePickingorderLog(pickingorderLog)
						if err != nil {
							msg := "create picking order detail error"
							return nil, errors.New(msg)
						}
						err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, nextBatch.Balance, info.Email)
						if err != nil {
							msg := "update location canpick error: "
							return nil, errors.New(msg)
						}
						quantityToPick = quantityToPick - nextBatch.Balance
					}
				}
			}
			var soItem SalesorderItem
			soItem.SalesorderItemID = itemRow.SalesorderItemID
			soItem.QuantityPicked = itemRow.Quantity
			soItem.Updated = time.Now()
			soItem.UpdatedBy = info.Email

			err = repo.PickSalesorderItem(soItem)
			if err != nil {
				msg := "pick salesorder item error: "
				return nil, errors.New(msg)
			}
			var poItem PickingorderItem
			poItem.OrganizationID = info.OrganizationID
			poItem.PickingorderID = pickingorderID
			poItem.SalesorderItemID = itemRow.SalesorderItemID
			poItem.PickingorderItemID = pickingorderItemID
			poItem.ItemID = itemRow.ItemID
			poItem.Quantity = toPick
			poItem.Status = 1
			poItem.CreatedBy = info.Email
			poItem.Created = time.Now()
			poItem.Updated = time.Now()
			poItem.UpdatedBy = info.Email

			err = repo.CreatePickingorderItem(poItem)
			if err != nil {
				msg := "create picking order item error: "
				return nil, errors.New(msg)
			}
			err = itemRepo.UpdateItemPickingStock(itemRow.ItemID, toPick, info.Email)
			if err != nil {
				msg := "update item stock error: "
				return nil, errors.New(msg)
			}
		}
		picking_status := 3
		err = repo.UpdateSalesorderPickingStatus(soID, picking_status, info.Email)
		if err != nil {
			msg := "update sales order picking status error: "
			return nil, errors.New(msg)
		}
		err = repo.UpdateSalesorderStatus(soID, 2, info.Email)
		if err != nil {
			msg := "update sales order status error: "
			return nil, errors.New(msg)
		}
		var newEvent common.NewHistoryCreated
		newEvent.HistoryType = "salesorder"
		newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
		newEvent.HistoryBy = info.User
		newEvent.ReferenceID = soID
		newEvent.Description = "Picking Order Created"
		newEvent.OrganizationID = info.OrganizationID
		newEvent.Email = info.Email
		msg, _ := json.Marshal(newEvent)
		msgs = append(msgs, msg)
	}
	fmt.Println(pickingorderID)
	logs, err := repo.GetPickingorderLogSum(pickingorderID)
	if err != nil {
		msg := "get picking order logs  error: "
		return nil, errors.New(msg)
	}
	fmt.Println(logs)
	for _, logRow := range *logs {
		var pickingorderDetail PickingorderDetail
		pickingorderDetail.PickingorderDetailID = "pid-" + xid.New().String()
		pickingorderDetail.OrganizationID = logRow.OrganizationID
		pickingorderDetail.PickingorderID = logRow.PickingorderID
		pickingorderDetail.ItemID = logRow.ItemID
		pickingorderDetail.LocationID = logRow.LocationID
		pickingorderDetail.Quantity = logRow.Quantity
		pickingorderDetail.QuantityPicked = 0
		pickingorderDetail.Status = 1
		pickingorderDetail.CreatedBy = info.Email
		pickingorderDetail.Created = time.Now()
		pickingorderDetail.Updated = time.Now()
		pickingorderDetail.UpdatedBy = info.Email
		err = repo.CreatePickingorderDetail(pickingorderDetail)
		if err != nil {
			msg := "create picking order detail error: "
			return nil, errors.New(msg)
		}
	}
	var pickingorder Pickingorder
	pickingorder.SalesorderID = strings.Join(info.SOID[:], ",")
	pickingorder.PickingorderID = pickingorderID
	pickingorder.PickingorderNumber = info.PickingorderNumber
	pickingorder.PickingorderDate = info.PickingorderDate
	pickingorder.OrganizationID = info.OrganizationID
	pickingorder.Notes = info.Notes
	pickingorder.Status = 1
	pickingorder.Created = time.Now()
	pickingorder.CreatedBy = info.Email
	pickingorder.Updated = time.Now()
	pickingorder.UpdatedBy = info.Email
	err = repo.CreatePickingorder(pickingorder)
	if err != nil {
		msg := "create picking order error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewHistoryCreated error"
			return nil, errors.New(msg)
		}
	}
	return &pickingorderID, err
}

func (s *salesorderService) NewPickingFromLocation(pickingorderID string, info PickingFromLocationNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	pickingorderDetail, err := repo.GetPickingorderDetailByLocationID(info.OrganizationID, pickingorderID, info.LocationID)
	if err != nil {
		msg := "picking order detail not exist"
		return nil, errors.New(msg)
	}
	itemInfo, err := itemRepo.GetItemByID(pickingorderDetail.ItemID, info.OrganizationID)
	if err != nil {
		msg := "item not exist"
		return nil, errors.New(msg)
	}
	if itemInfo.StockPicking < info.Quantity {
		msg := "item pick too many"
		return nil, errors.New(msg)
	}
	location, err := warehouseRepo.GetLocationByID(pickingorderDetail.LocationID, info.OrganizationID)
	if err != nil {
		msg := "location not exist"
		return nil, errors.New(msg)
	}
	if location.Quantity-location.CanPick < info.Quantity {
		msg := "location pick too many"
		return nil, errors.New(msg)
	}
	err = repo.UpdatePickingorderPicked(pickingorderDetail.PickingorderDetailID, info.Quantity, info.Email)
	if err != nil {
		msg := "update picking order picked error"
		return nil, errors.New(msg)
	}
	err = itemRepo.UpdateItemPackingStock(pickingorderDetail.ItemID, info.Quantity, info.Email)
	if err != nil {
		msg := "update item stock error"
		return nil, errors.New(msg)
	}
	err = warehouseRepo.UpdateLocationPicked(info.LocationID, info.Quantity, info.Email)
	if err != nil {
		msg := "update location stock error"
		return nil, errors.New(msg)
	}
	pickingorderStatus := 2
	picked, err := repo.CheckPickingorderPicked(info.OrganizationID, pickingorderID)
	if err != nil {
		msg := "check picking order status error"
		return nil, errors.New(msg)
	}
	if !picked {
		pickingorderStatus = 3
	}
	err = repo.UpdatePickingorderStatus(pickingorderID, pickingorderStatus, info.Email)
	if err != nil {
		msg := "update picking order status error"
		return nil, errors.New(msg)
	}

	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "location"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = info.LocationID
	newEvent.Description = "Item Picked From Location"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &pickingorderID, err
}

func (s *salesorderService) UpdatePickingorderPicked(pickingorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	oldPickingorder, err := repo.GetPickingorderByID(organizationID, pickingorderID)
	if err != nil {
		msg := "Pickingorder not exist"
		return errors.New(msg)
	}
	if oldPickingorder.Status != 1 && oldPickingorder.Status != 2 {
		msg := "Pickingorder status error"
		return errors.New(msg)
	}
	pickingorderDetails, err := repo.GetPickingorderDetailList(organizationID, pickingorderID)
	if err != nil {
		msg := "Pickingorder detail not exist"
		return errors.New(msg)
	}
	for _, pickingorderDetail := range *pickingorderDetails {
		topick := pickingorderDetail.Quantity - pickingorderDetail.QuantityPicked
		itemInfo, err := itemRepo.GetItemByID(pickingorderDetail.ItemID, organizationID)
		if err != nil {
			msg := "item not exist"
			return errors.New(msg)
		}
		if itemInfo.StockPicking < topick {
			msg := "item pick too many"
			return errors.New(msg)
		}
		location, err := warehouseRepo.GetLocationByID(pickingorderDetail.LocationID, organizationID)
		if err != nil {
			msg := "location not exist"
			return errors.New(msg)
		}
		if location.Quantity-location.CanPick < topick {
			msg := "location pick too many"
			return errors.New(msg)
		}
		err = repo.UpdatePickingorderPicked(pickingorderDetail.PickingorderDetailID, topick, email)
		if err != nil {
			msg := "update picking order picked error"
			return errors.New(msg)
		}
		err = itemRepo.UpdateItemPackingStock(pickingorderDetail.ItemID, topick, email)
		if err != nil {
			msg := "update item stock error"
			return errors.New(msg)
		}
		err = warehouseRepo.UpdateLocationPicked(pickingorderDetail.LocationID, topick, email)
		if err != nil {
			msg := "update location stock error"
			return errors.New(msg)
		}
	}
	err = repo.UpdatePickingorderStatus(pickingorderID, 3, email) //CONFIRMED
	if err != nil {
		msg := "update pickingorder error: "
		return errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "pickingorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = pickingorderID
	newEvent.Description = "Picking Order Fully Picked"
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
	return err
}

func (s *salesorderService) UpdatePickingorderUnPicked(pickingorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	oldPickingorder, err := repo.GetPickingorderByID(organizationID, pickingorderID)
	if err != nil {
		msg := "Pickingorder not exist"
		return errors.New(msg)
	}
	if oldPickingorder.Status != 3 && oldPickingorder.Status != 2 {
		msg := "Pickingorder status error"
		return errors.New(msg)
	}
	pickingorderDetails, err := repo.GetPickingorderDetailList(organizationID, pickingorderID)
	if err != nil {
		msg := "Pickingorder detail not exist"
		return errors.New(msg)
	}
	for _, pickingorderDetail := range *pickingorderDetails {
		itemInfo, err := itemRepo.GetItemByID(pickingorderDetail.ItemID, organizationID)
		if err != nil {
			msg := "item not exist"
			return errors.New(msg)
		}
		if itemInfo.StockPacking < pickingorderDetail.QuantityPicked {
			msg := "item packing quantity error"
			return errors.New(msg)
		}
		location, err := warehouseRepo.GetLocationByID(pickingorderDetail.LocationID, organizationID)
		if err != nil {
			msg := "location not exist"
			return errors.New(msg)
		}
		if location.Available < pickingorderDetail.QuantityPicked {
			msg := "location space not enough"
			return errors.New(msg)
		}
		err = repo.UpdatePickingorderPicked(pickingorderDetail.PickingorderDetailID, -pickingorderDetail.QuantityPicked, email)
		if err != nil {
			msg := "update picking order picked error"
			return errors.New(msg)
		}
		err = itemRepo.UpdateItemPackingStock(pickingorderDetail.ItemID, -pickingorderDetail.QuantityPicked, email)
		if err != nil {
			msg := "update item stock error"
			return errors.New(msg)
		}
		err = warehouseRepo.UpdateLocationPicked(pickingorderDetail.LocationID, -pickingorderDetail.QuantityPicked, email)
		if err != nil {
			msg := "update location stock error"
			return errors.New(msg)
		}
	}
	err = repo.UpdatePickingorderStatus(pickingorderID, 1, email) //CONFIRMED
	if err != nil {
		msg := "update pickingorder error: "
		return errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "pickingorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = pickingorderID
	newEvent.Description = "Picking Order Marked As  UnPicked"
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
	return err
}

//package

func (s *salesorderService) NewPackage(salesorderID string, info PackageNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckPackageNumberConfict("", info.OrganizationID, info.PackageNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "picking order number exists"
		return nil, errors.New(msg)
	}
	packageID := "pac-" + xid.New().String()
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldSoItem, err := repo.GetSalesorderItemByID(info.OrganizationID, salesorderID, itemRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return nil, errors.New(msg)
		}
		itemInfo, err := itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
		if err != nil {
			msg := "item not exist"
			return nil, errors.New(msg)
		}
		packageItemID := "pai-" + xid.New().String()
		if itemInfo.StockPacking < itemRow.Quantity {
			msg := "no enough stock to pack"
			return nil, errors.New(msg)
		}
		if oldSoItem.QuantityPicked < itemRow.Quantity {
			msg := "packing quantity greater than not packed"
			return nil, errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityPacked = oldSoItem.QuantityPacked + itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.PackSalesorderItem(soItem)
		if err != nil {
			msg := "pack salesorder item error: "
			return nil, errors.New(msg)
		}
		var packageItem PackageItem
		packageItem.OrganizationID = info.OrganizationID
		packageItem.PackageID = packageID
		packageItem.SalesorderItemID = oldSoItem.SalesorderItemID
		packageItem.PackageItemID = packageItemID
		packageItem.ItemID = oldSoItem.ItemID
		packageItem.Quantity = itemRow.Quantity
		packageItem.Status = 1
		packageItem.CreatedBy = info.Email
		packageItem.Created = time.Now()
		packageItem.Updated = time.Now()
		packageItem.UpdatedBy = info.Email

		err = repo.CreatePackageItem(packageItem)
		if err != nil {
			msg := "create package item error: "
			return nil, errors.New(msg)
		}
		err = itemRepo.UpdateItemPackedStock(itemRow.ItemID, itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: "
			return nil, errors.New(msg)
		}
	}
	var newPackage Package
	newPackage.SalesorderID = salesorderID
	newPackage.PackageID = packageID
	newPackage.PackageNumber = info.PackageNumber
	newPackage.PackageDate = info.PackageDate
	newPackage.OrganizationID = info.OrganizationID
	newPackage.Notes = info.Notes
	newPackage.Status = 1
	newPackage.Created = time.Now()
	newPackage.CreatedBy = info.Email
	newPackage.Updated = time.Now()
	newPackage.UpdatedBy = info.Email
	err = repo.CreatePackage(newPackage)
	if err != nil {
		msg := "create package error: "
		return nil, errors.New(msg)
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order error: "
		return nil, errors.New(msg)
	}
	packedCount, err := repo.GetSalesorderPackedCount(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order packed count error: "
		return nil, errors.New(msg)
	}
	packing_status := 1
	if so.ItemCount == packedCount {
		packing_status = 3
	} else {
		packing_status = 2
	}
	err = repo.UpdateSalesorderPackingStatus(salesorderID, packing_status, info.Email)
	if err != nil {
		msg := "update sales order packing status error: "
		return nil, errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(salesorderID, 2, info.Email)
	if err != nil {
		msg := "update sales order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Package Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &packageID, err
}

func (s *salesorderService) GetPackageList(filter PackageFilter) (int, *[]PackageResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetPackageCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPackageList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) GetPackageItemList(salesorderID, organizationID string) (*[]PackageItemResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetPackageByID(organizationID, salesorderID)
	if err != nil {
		msg := "get package error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetPackageItemList(salesorderID)
	return list, err
}

func (s *salesorderService) BatchShippingorder(info ShippingorderBatch) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckShippingorderNumberConfict("", info.OrganizationID, info.ShippingorderNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "picking order number exists"
		return nil, errors.New(msg)
	}
	shippingorderID := "shi-" + xid.New().String()
	itemRepo := item.NewItemRepository(tx)
	settingRepo := setting.NewSettingRepository(tx)
	if info.CarrierID != "" {
		_, err = settingRepo.GetCarrierByID(info.OrganizationID, info.CarrierID)
		if err != nil {
			msg := "carrier not exist"
			return nil, errors.New(msg)
		}
	}
	var msgs [][]byte
	var salesorders []string
	var itemHistorys []string
	for _, packageID := range info.PackageID {
		packageInfo, err := repo.GetPackageByID(info.OrganizationID, packageID)
		if err != nil {
			msg := "get package error: "
			return nil, errors.New(msg)
		}
		if packageInfo.Status != 1 {
			msg := "package status error for package: " + packageInfo.PackageNumber
			return nil, errors.New(msg)
		}
		items, err := repo.GetPackageItemList(info.OrganizationID, packageID)
		if err != nil {
			msg := "get package items error: "
			return nil, errors.New(msg)
		}
		for _, itemRow := range *items {
			itemInfo, err := itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
			if err != nil {
				msg := "item not exist"
				return nil, errors.New(msg)
			}
			shippingorderDetailID := "shid-" + xid.New().String()
			salesorderInfo, err := repo.GetSalesorderByID(info.OrganizationID, packageInfo.SalesorderID)
			if err != nil {
				msg := "salesorder not exist"
				return nil, errors.New(msg)
			}
			salesorderItem, err := repo.GetSalesorderItemByID(info.OrganizationID, salesorderInfo.SalesorderID, itemInfo.ItemID)
			fmt.Println(salesorderItem.QuantityPacked, "------")
			if err != nil {
				msg := "salesorder item not exist"
				return nil, errors.New(msg)
			}
			if salesorderItem.QuantityPacked < itemRow.Quantity {
				msg := "no enough stock to ship for item: " + itemInfo.Name + " in salesorder :" + salesorderInfo.SalesorderNumber
				return nil, errors.New(msg)
			}
			var soItem SalesorderItem
			soItem.SalesorderItemID = salesorderItem.SalesorderItemID
			soItem.QuantityShipped = salesorderItem.QuantityShipped + itemRow.Quantity
			soItem.Updated = time.Now()
			soItem.UpdatedBy = info.Email

			err = repo.ShipSalesorderItem(soItem)
			if err != nil {
				msg := "ship salesorder item error: "
				return nil, errors.New(msg)
			}
			var shippingorderDetail ShippingorderDetail
			shippingorderDetail.OrganizationID = info.OrganizationID
			shippingorderDetail.ShippingorderID = shippingorderID
			shippingorderDetail.ShippingorderDetailID = shippingorderDetailID
			shippingorderDetail.PackageID = packageID
			shippingorderDetail.PackageItemID = itemRow.PackageItemID
			shippingorderDetail.ItemID = itemRow.ItemID
			shippingorderDetail.Quantity = itemRow.Quantity
			shippingorderDetail.Status = 1
			shippingorderDetail.CreatedBy = info.Email
			shippingorderDetail.Created = time.Now()
			shippingorderDetail.Updated = time.Now()
			shippingorderDetail.UpdatedBy = info.Email

			err = repo.CreateShippingorderDetail(shippingorderDetail)
			if err != nil {
				msg := "create shipping order item error: "
				return nil, errors.New(msg)
			}
			salesorderUpdated := false
			for _, salesorderID := range salesorders {
				if salesorderID == packageInfo.SalesorderID {
					salesorderUpdated = true
					continue
				}
			}
			itemUpdated := false
			for _, itemID := range itemHistorys {
				if itemID == itemRow.ItemID {
					itemUpdated = true
					continue
				}
			}
			if !salesorderUpdated {
				salesorders = append(salesorders, packageInfo.SalesorderID)
				var newEvent common.NewHistoryCreated
				newEvent.HistoryType = "salesorder"
				newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
				newEvent.HistoryBy = info.User
				newEvent.ReferenceID = packageInfo.SalesorderID
				newEvent.Description = "Shipping Order Created"
				newEvent.OrganizationID = info.OrganizationID
				newEvent.Email = info.Email
				msg, _ := json.Marshal(newEvent)
				msgs = append(msgs, msg)
			}
			if !itemUpdated {
				itemHistorys = append(itemHistorys, itemRow.ItemID)
				var newEvent common.NewHistoryCreated
				newEvent.HistoryType = "item"
				newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
				newEvent.HistoryBy = info.User
				newEvent.ReferenceID = itemRow.ItemID
				newEvent.Description = "Shipping Order Created"
				newEvent.OrganizationID = info.OrganizationID
				newEvent.Email = info.Email
				msg, _ := json.Marshal(newEvent)
				msgs = append(msgs, msg)
			}
			shippedCount, err := repo.GetSalesorderShippedCount(info.OrganizationID, packageInfo.SalesorderID)
			if err != nil {
				msg := "get sales order packed count error: "
				return nil, errors.New(msg)
			}
			shippingStatus := 1
			if salesorderInfo.ItemCount == shippedCount {
				shippingStatus = 3
			} else {
				shippingStatus = 2
			}
			err = repo.UpdateSalesorderShippingStatus(packageInfo.SalesorderID, shippingStatus, info.Email)
			if err != nil {
				msg := "update salesorder shipping status error: "
				return nil, errors.New(msg)
			}
			if salesorderInfo.InvoiceStatus == 3 && shippingStatus == 3 {
				err = repo.UpdateSalesorderStatus(packageInfo.SalesorderID, 3, info.Email)
				if err != nil {
					msg := "update sales order status error: "
					return nil, errors.New(msg)
				}
			}
		}
		err = repo.UpdatePackageStatus(packageID, 2, info.Email)
		if err != nil {
			msg := "update package status error: "
			return nil, errors.New(msg)
		}
	}
	details, err := repo.GetShippingorderDetailSum(shippingorderID)
	if err != nil {
		msg := "get picking order detail  error: "
		return nil, errors.New(msg)
	}
	for _, detailRow := range *details {
		var shippingorderItem ShippingorderItem
		shippingorderItem.ShippingorderItemID = "shii-" + xid.New().String()
		shippingorderItem.OrganizationID = detailRow.OrganizationID
		shippingorderItem.ShippingorderID = detailRow.ShippingorderID
		shippingorderItem.ItemID = detailRow.ItemID
		shippingorderItem.Quantity = detailRow.Quantity
		shippingorderItem.Status = 1
		shippingorderItem.CreatedBy = info.Email
		shippingorderItem.Created = time.Now()
		shippingorderItem.Updated = time.Now()
		shippingorderItem.UpdatedBy = info.Email
		err = repo.CreateShippingorderItem(shippingorderItem)
		if err != nil {
			msg := "create picking order item error: "
			return nil, errors.New(msg)
		}
	}
	var shippingorder Shippingorder
	shippingorder.OrganizationID = info.OrganizationID
	shippingorder.ShippingorderID = shippingorderID
	shippingorder.PackageID = strings.Join(info.PackageID[:], ",")
	shippingorder.ShippingorderNumber = info.ShippingorderNumber
	shippingorder.ShippingorderDate = info.ShippingorderDate
	shippingorder.CarrierID = info.CarrierID
	shippingorder.TrackingNumber = info.TrackingNumber
	shippingorder.Notes = info.Notes
	shippingorder.Status = 1
	shippingorder.Created = time.Now()
	shippingorder.CreatedBy = info.Email
	shippingorder.Updated = time.Now()
	shippingorder.UpdatedBy = info.Email
	err = repo.CreateShippingorder(shippingorder)
	if err != nil {
		msg := "create shipping order error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewHistoryCreated error"
			return nil, errors.New(msg)
		}
	}
	return &shippingorderID, err
}

func (s *salesorderService) GetShippingorderList(filter ShippingorderFilter) (int, *[]ShippingorderResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetShippingorderCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetShippingorderList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) GetShippingorderItemList(salesorderID, organizationID string) (*[]ShippingorderItemResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetShippingorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "get package error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetShippingorderItemList(salesorderID)
	return list, err
}

func (s *salesorderService) GetShippingorderDetailList(salesorderID, organizationID string) (*[]ShippingorderDetailResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetShippingorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "get package error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetShippingorderDetailList(salesorderID)
	return list, err
}

func (s *salesorderService) GetRequisitionList(filter RequsitionFilter) (*[]RequsitionResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	if filter.StartDate == "" {
		filter.StartDate = time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	}
	if filter.EndDate == "" {
		filter.EndDate = time.Now().Format("2006-01-02")
	}
	if filter.TargetDay == 0 {
		filter.TargetDay = 30
	}
	fmt.Println(filter.EndDate)
	endTime, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		msg := "end date error"
		return nil, errors.New(msg)
	}
	fmt.Println(filter.StartDate)
	startTime, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		msg := "start date error"
		return nil, errors.New(msg)
	}
	filter.Period = int(endTime.Sub(startTime).Hours() / 24)
	fmt.Println(filter.Period)
	list, err := query.GetRequisitionList(filter)
	if err != nil {
		return nil, err
	}
	var res []RequsitionResponse
	for _, item := range *list {
		if item.Quantity > 0 {
			res = append(res, item)
		}
	}
	return &res, err
}

func (s *salesorderService) DeleteShippingorder(shippingorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	_, err = repo.GetShippingorderByID(shippingorderID, organizationID)
	if err != nil {
		msg := "shipping order not exist"
		return errors.New(msg)
	}
	shippingorderDetails, err := repo.GetShippingorderDetailList(shippingorderID)
	if err != nil {
		msg := "get shipping order details error"
		return errors.New(msg)
	}
	var msgs [][]byte
	var salesorders []string
	var itemHistorys []string
	for _, detail := range *shippingorderDetails {
		packageInfo, err := repo.GetPackageByID(organizationID, detail.PackageID)
		if err != nil {
			msg := "get package error: "
			return errors.New(msg)
		}
		// fmt.Println(packageInfo.Status)
		// if packageInfo.Status != 2 {
		// 	msg := "package status error for package: " + packageInfo.PackageNumber
		// 	return errors.New(msg)
		// }
		itemInfo, err := itemRepo.GetItemByID(detail.ItemID, organizationID)
		if err != nil {
			msg := "item not exist"
			return errors.New(msg)
		}
		salesorderInfo, err := repo.GetSalesorderByID(organizationID, packageInfo.SalesorderID)
		if err != nil {
			msg := "salesorder not exist"
			return errors.New(msg)
		}
		salesorderItem, err := repo.GetSalesorderItemByID(organizationID, salesorderInfo.SalesorderID, itemInfo.ItemID)
		if err != nil {
			msg := "salesorder item not exist"
			return errors.New(msg)
		}
		if salesorderItem.QuantityShipped < detail.Quantity {
			msg := "shipped item error for " + itemInfo.Name + " in salesorder :" + salesorderInfo.SalesorderNumber
			return errors.New(msg)
		}
		fmt.Println(salesorderItem.QuantityShipped, detail.Quantity)
		var soItem SalesorderItem
		soItem.SalesorderItemID = salesorderItem.SalesorderItemID
		soItem.QuantityShipped = salesorderItem.QuantityShipped - detail.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = email

		err = repo.ShipSalesorderItem(soItem)
		if err != nil {
			msg := "ship salesorder item error: "
			return errors.New(msg)
		}

		salesorderUpdated := false
		for _, salesorderID := range salesorders {
			if salesorderID == packageInfo.SalesorderID {
				salesorderUpdated = true
				continue
			}
		}
		itemUpdated := false
		for _, itemID := range itemHistorys {
			if itemID == detail.ItemID {
				itemUpdated = true
				continue
			}
		}
		if !salesorderUpdated {
			salesorders = append(salesorders, packageInfo.SalesorderID)
			var newEvent common.NewHistoryCreated
			newEvent.HistoryType = "salesorder"
			newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
			newEvent.HistoryBy = user
			newEvent.ReferenceID = packageInfo.SalesorderID
			newEvent.Description = "Shipping Order Deleted"
			newEvent.OrganizationID = organizationID
			newEvent.Email = email
			msg, _ := json.Marshal(newEvent)
			msgs = append(msgs, msg)
		}
		if !itemUpdated {
			itemHistorys = append(itemHistorys, detail.ItemID)
			var newEvent common.NewHistoryCreated
			newEvent.HistoryType = "item"
			newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
			newEvent.HistoryBy = user
			newEvent.ReferenceID = detail.ItemID
			newEvent.Description = "Shipping Order Deleted"
			newEvent.OrganizationID = organizationID
			newEvent.Email = email
			msg, _ := json.Marshal(newEvent)
			msgs = append(msgs, msg)
		}
		shippedCount, err := repo.GetSalesorderShippedCount(organizationID, packageInfo.SalesorderID)
		if err != nil {
			msg := "get sales order shipped count error: "
			return errors.New(msg)
		}
		shippingStatus := 2
		if shippedCount == 0 {
			shippingStatus = 1
		} else {
			shippingStatus = 2
		}
		err = repo.UpdateSalesorderShippingStatus(packageInfo.SalesorderID, shippingStatus, email)
		if err != nil {
			msg := "update salesorder shipping status error: "
			return errors.New(msg)
		}
		err = repo.UpdateSalesorderStatus(packageInfo.SalesorderID, 2, email)
		if err != nil {
			msg := "update sales order status error: "
			return errors.New(msg)
		}
		err = repo.UpdatePackageStatus(detail.PackageID, 1, email)
		if err != nil {
			msg := "update package status error: "
			return errors.New(msg)
		}
	}
	err = repo.DeleteShippingorder(shippingorderID, email)
	if err != nil {
		return err
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "shippingorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = shippingorderID
	newEvent.Description = "Shipping order Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewHistoryCreated error"
			return errors.New(msg)
		}
	}
	return nil
}

func (s *salesorderService) DeletePackage(packageID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	oldPackage, err := repo.GetPackageByID(organizationID, packageID)
	if err != nil {
		msg := "package not exist"
		return errors.New(msg)
	}
	packageItems, err := repo.GetPackageItemList(organizationID, packageID)
	if err != nil {
		msg := "get package items error"
		return errors.New(msg)
	}
	var msgs [][]byte
	for _, itemRow := range *packageItems {
		oldSoItem, err := repo.GetSalesorderItemByID(organizationID, oldPackage.SalesorderID, itemRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return errors.New(msg)
		}
		_, err = itemRepo.GetItemByID(itemRow.ItemID, organizationID)
		if err != nil {
			msg := "item not exist"
			return errors.New(msg)
		}
		if oldSoItem.QuantityPacked < itemRow.Quantity {
			msg := "sales order packed quantity error"
			return errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityPacked = oldSoItem.QuantityPacked - itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = email

		err = repo.PackSalesorderItem(soItem)
		if err != nil {
			msg := "update salesorder item packed error: "
			return errors.New(msg)
		}
		err = itemRepo.UpdateItemPackedStock(itemRow.ItemID, -itemRow.Quantity, email)
		if err != nil {
			msg := "update item stock error: "
			return errors.New(msg)
		}
		var newEvent common.NewHistoryCreated
		newEvent.HistoryType = "item"
		newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
		newEvent.HistoryBy = user
		newEvent.ReferenceID = oldPackage.SalesorderID
		newEvent.Description = "Package Deleted"
		newEvent.OrganizationID = organizationID
		newEvent.Email = email
		msg, _ := json.Marshal(newEvent)
		msgs = append(msgs, msg)
	}
	packedCount, err := repo.GetSalesorderPackedCount(organizationID, oldPackage.SalesorderID)
	if err != nil {
		msg := "get sales order packed count error: "
		return errors.New(msg)
	}
	packing_status := 1
	if packedCount == 0 {
		packing_status = 1
	} else {
		packing_status = 2
	}
	err = repo.UpdateSalesorderPackingStatus(oldPackage.SalesorderID, packing_status, email)
	if err != nil {
		msg := "update sales order packing status error: "
		return errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(oldPackage.SalesorderID, 2, email)
	if err != nil {
		msg := "update sales order status error: "
		return errors.New(msg)
	}
	err = repo.DeletePackage(packageID, email)
	if err != nil {
		msg := "delete package error: "
		return errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldPackage.SalesorderID
	newEvent.Description = "Package Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewHistoryCreated error"
			return errors.New(msg)
		}
	}
	return nil
}

func (s *salesorderService) DeletePickingorder(pickingorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	warehouseRepo := warehouse.NewWarehouseRepository(tx)
	oldPickingorder, err := repo.GetPickingorderByID(organizationID, pickingorderID)
	if err != nil {
		msg := "get picking order error"
		return errors.New(msg)
	}
	if oldPickingorder.Status != 1 {
		msg := " picking order status error"
		return errors.New(msg)
	}
	pickingorderLogs, err := repo.GetPickingorderLogList(organizationID, pickingorderID)
	if err != nil {
		msg := "get picking order log error"
		return errors.New(msg)
	}
	var msgs [][]byte
	var salesorders []string
	var itemHistorys []string
	for _, logRow := range *pickingorderLogs {
		oldSoItem, err := repo.GetSalesorderItemByID(organizationID, logRow.SalesorderID, logRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return errors.New(msg)
		}
		itemInfo, err := itemRepo.GetItemByID(logRow.ItemID, organizationID)
		if err != nil {
			msg := "item not exist"
			return errors.New(msg)
		}
		if itemInfo.StockPicking < logRow.Quantity {
			msg := "item picking quantity error"
			return errors.New(msg)
		}
		err = itemRepo.PickItem(logRow.BatchID, -logRow.Quantity, email)
		if err != nil {
			msg := "return item back to batch error"
			return errors.New(msg)
		}
		err = warehouseRepo.UpdateLocationCanPick(logRow.LocationID, -logRow.Quantity, email)
		if err != nil {
			msg := "update location canpick error: "
			return errors.New(msg)
		}
		if oldSoItem.QuantityPicked < logRow.Quantity {
			msg := "sales order itme picked quantity error"
			return errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityPicked = oldSoItem.QuantityPicked - logRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = email

		err = repo.PickSalesorderItem(soItem)
		if err != nil {
			msg := "unpick salesorder item error: "
			return errors.New(msg)
		}
		err = itemRepo.UpdateItemPickingStock(logRow.ItemID, -logRow.Quantity, email)
		if err != nil {
			msg := "update item stock error: "
			return errors.New(msg)
		}

		salesorderUpdated := false
		for _, salesorderID := range salesorders {
			if salesorderID == logRow.SalesorderID {
				salesorderUpdated = true
				continue
			}
		}
		itemUpdated := false
		for _, itemID := range itemHistorys {
			if itemID == logRow.ItemID {
				itemUpdated = true
				continue
			}
		}
		if !salesorderUpdated {
			salesorders = append(salesorders, logRow.SalesorderID)
			var newEvent common.NewHistoryCreated
			newEvent.HistoryType = "salesorder"
			newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
			newEvent.HistoryBy = user
			newEvent.ReferenceID = logRow.SalesorderID
			newEvent.Description = "Picking Order Deleted"
			newEvent.OrganizationID = organizationID
			newEvent.Email = email
			msg, _ := json.Marshal(newEvent)
			msgs = append(msgs, msg)
		}
		if !itemUpdated {
			itemHistorys = append(itemHistorys, logRow.ItemID)
			var newEvent common.NewHistoryCreated
			newEvent.HistoryType = "item"
			newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
			newEvent.HistoryBy = user
			newEvent.ReferenceID = logRow.ItemID
			newEvent.Description = "Picking Order Deleted"
			newEvent.OrganizationID = organizationID
			newEvent.Email = email
			msg, _ := json.Marshal(newEvent)
			msgs = append(msgs, msg)
		}
	}
	err = repo.DeletePickingorder(pickingorderID, email)
	if err != nil {
		msg := "delete picking order error: "
		return errors.New(msg)
	}
	for _, so := range salesorders {
		pickedCount, err := repo.GetSalesorderPickedCount(organizationID, so)
		if err != nil {
			msg := "get sales order received count error: "
			return errors.New(msg)
		}
		pickingStatus := 1
		if pickedCount == 0 {
			pickingStatus = 1
		} else {
			pickingStatus = 2
		}
		err = repo.UpdateSalesorderPickingStatus(so, pickingStatus, email)
		if err != nil {
			msg := "update sales order receive status error: "
			return errors.New(msg)
		}

	}
	tx.Commit()

	rabbit, _ := queue.GetConn()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "pickingorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = pickingorderID
	newEvent.Description = "Picking Order Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewHistoryCreated error"
			return errors.New(msg)
		}
	}
	return err
}

func (s *salesorderService) NewInvoice(salesorderID string, info InvoiceNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckInvoiceNumberConfict("", info.OrganizationID, info.InvoiceNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "invoice number exists"
		return nil, errors.New(msg)
	}
	invoiceID := "inv-" + xid.New().String()
	settingRepo := setting.NewSettingRepository(tx)
	itemCount := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldSoItem, err := repo.GetSalesorderItemByID(info.OrganizationID, salesorderID, itemRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return nil, errors.New(msg)
		}
		if oldSoItem.SalesorderItemID != itemRow.SalesorderItemID {
			msg := "sales order item id error"
			return nil, errors.New(msg)
		}
		_, err = itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
		if err != nil {
			msg := "item not exist"
			return nil, errors.New(msg)
		}

		taxValue := 0.0
		if itemRow.TaxID != "" {
			tax, err := settingRepo.GetTaxByID(itemRow.TaxID, info.OrganizationID)
			if err != nil {
				msg := "tax not exist"
				return nil, errors.New(msg)
			}
			taxValue = tax.TaxValue
		}
		itemCount += itemRow.Quantity
		itemTotal += itemRow.Rate * float64(itemRow.Quantity)
		taxTotal += itemRow.Rate * float64(itemRow.Quantity) * taxValue / 100

		invoiceItemID := "invi-" + xid.New().String()
		if oldSoItem.Quantity < oldSoItem.QuantityInvoiced+itemRow.Quantity {
			msg := "invoicing quantity greater than not invoiced"
			return nil, errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityInvoiced = oldSoItem.QuantityInvoiced + itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.InvoiceSalesorderItem(soItem)
		if err != nil {
			msg := "invoice salesorder item error: "
			return nil, errors.New(msg)
		}
		var invoiceItem InvoiceItem
		invoiceItem.OrganizationID = info.OrganizationID
		invoiceItem.InvoiceID = invoiceID
		invoiceItem.InvoiceItemID = invoiceItemID
		invoiceItem.SalesorderItemID = oldSoItem.SalesorderItemID
		invoiceItem.ItemID = oldSoItem.ItemID
		invoiceItem.Quantity = itemRow.Quantity
		invoiceItem.Rate = itemRow.Rate
		invoiceItem.TaxID = itemRow.TaxID
		invoiceItem.TaxValue = taxValue
		invoiceItem.TaxAmount = float64(itemRow.Quantity) * itemRow.Rate * taxValue / 100
		invoiceItem.Amount = float64(itemRow.Quantity) * itemRow.Rate
		invoiceItem.Status = 1
		invoiceItem.CreatedBy = info.Email
		invoiceItem.Created = time.Now()
		invoiceItem.Updated = time.Now()
		invoiceItem.UpdatedBy = info.Email

		err = repo.CreateInvoiceItem(invoiceItem)
		if err != nil {
			msg := "create invoice item error: "
			return nil, errors.New(msg)
		}
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order error: "
		return nil, errors.New(msg)
	}
	var invoice Invoice
	invoice.OrganizationID = info.OrganizationID
	invoice.InvoiceID = invoiceID
	invoice.SalesorderID = salesorderID
	invoice.InvoiceNumber = info.InvoiceNumber
	invoice.InvoiceDate = info.InvoiceDate
	invoice.DueDate = info.DueDate
	invoice.CustomerID = so.CustomerID
	invoice.ItemCount = itemCount
	invoice.Subtotal = itemTotal
	invoice.DiscountType = info.DiscountType
	invoice.DiscountValue = info.DiscountValue
	invoice.TaxTotal = taxTotal
	invoice.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		invoice.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		invoice.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		invoice.Total = itemTotal + info.ShippingFee + taxTotal
	}
	invoice.Notes = info.Notes
	invoice.Status = 1
	invoice.Created = time.Now()
	invoice.CreatedBy = info.Email
	invoice.Updated = time.Now()
	invoice.UpdatedBy = info.Email
	err = repo.CreateInvoice(invoice)
	if err != nil {
		msg := "create invoice error: "
		return nil, errors.New(msg)
	}
	invoicedCount, err := repo.GetSalesorderInvoicedCount(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order invoiced count error: "
		return nil, errors.New(msg)
	}
	invoiceStatus := 1
	if so.ItemCount == invoicedCount {
		invoiceStatus = 3
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateSalesorderInvoiceStatus(salesorderID, invoiceStatus, info.Email)
	if err != nil {
		msg := "update sales order invoice status error: "
		return nil, errors.New(msg)
	}
	soStatus := 1
	if so.ShippingStatus == 3 && invoiceStatus == 3 {
		soStatus = 3
	} else {
		soStatus = 2
	}
	err = repo.UpdateSalesorderStatus(salesorderID, soStatus, info.Email)
	if err != nil {
		msg := "update sales order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = salesorderID
	newEvent.Description = "Invoice Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &invoiceID, err
}

func (s *salesorderService) GetInvoiceList(filter InvoiceFilter) (int, *[]InvoiceResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetInvoiceCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetInvoiceList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) GetInvoiceItemList(invoiceID, organizationID string) (*[]InvoiceItemResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	_, err := query.GetInvoiceByID(organizationID, invoiceID)
	if err != nil {
		msg := "get pickingorder error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetInvoiceItemList(invoiceID)
	return list, err
}

func (s *salesorderService) UpdateInvoice(invoiceID string, info InvoiceNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckInvoiceNumberConfict(invoiceID, info.OrganizationID, info.InvoiceNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "invoice number exists"
		return nil, errors.New(msg)
	}
	settingRepo := setting.NewSettingRepository(tx)
	oldInvoice, err := repo.GetInvoiceByID(info.OrganizationID, invoiceID)
	if err != nil {
		msg := "get invoice error"
		return nil, errors.New(msg)
	}
	oldInvoiceItems, err := repo.GetInvoiceItemList(info.OrganizationID, invoiceID)
	if err != nil {
		msg := "get invoice item error"
		return nil, errors.New(msg)
	}
	for _, oldInvoiceItem := range *oldInvoiceItems {
		oldSoItem, err := repo.GetSalesorderItemByID(info.OrganizationID, oldInvoice.SalesorderID, oldInvoiceItem.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return nil, errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldInvoiceItem.SalesorderItemID
		soItem.QuantityInvoiced = oldSoItem.QuantityInvoiced - oldInvoiceItem.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.InvoiceSalesorderItem(soItem)
		if err != nil {
			msg := "cancel invoice salesorder item error: "
			return nil, errors.New(msg)
		}

	}
	err = repo.DeleteInvoice(invoiceID, info.User)
	if err != nil {
		msg := "Invoice Update error"
		return nil, errors.New(msg)
	}
	itemCount := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldSoItem, err := repo.GetSalesorderItemByID(info.OrganizationID, oldInvoice.SalesorderID, itemRow.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return nil, errors.New(msg)
		}
		if oldSoItem.SalesorderItemID != itemRow.SalesorderItemID {
			msg := "sales order item id error"
			return nil, errors.New(msg)
		}
		_, err = itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
		if err != nil {
			msg := "item not exist"
			return nil, errors.New(msg)
		}

		taxValue := 0.0
		if itemRow.TaxID != "" {
			tax, err := settingRepo.GetTaxByID(itemRow.TaxID, info.OrganizationID)
			if err != nil {
				msg := "tax not exist"
				return nil, errors.New(msg)
			}
			taxValue = tax.TaxValue
		}
		itemCount += itemRow.Quantity
		itemTotal += itemRow.Rate * float64(itemRow.Quantity)
		taxTotal += itemRow.Rate * float64(itemRow.Quantity) * taxValue / 100

		invoiceItemID := "invi-" + xid.New().String()
		if oldSoItem.Quantity < oldSoItem.QuantityInvoiced+itemRow.Quantity {
			fmt.Println(oldSoItem.Quantity, oldSoItem.QuantityInvoiced, itemRow.Quantity)
			msg := "invoicing quantity greater than not invoiced"
			return nil, errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityInvoiced = oldSoItem.QuantityInvoiced + itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.InvoiceSalesorderItem(soItem)
		if err != nil {
			msg := "invoice salesorder item error: "
			return nil, errors.New(msg)
		}
		var invoiceItem InvoiceItem
		invoiceItem.OrganizationID = info.OrganizationID
		invoiceItem.InvoiceID = invoiceID
		invoiceItem.InvoiceItemID = invoiceItemID
		invoiceItem.SalesorderItemID = oldSoItem.SalesorderItemID
		invoiceItem.ItemID = oldSoItem.ItemID
		invoiceItem.Quantity = itemRow.Quantity
		invoiceItem.Rate = itemRow.Rate
		invoiceItem.TaxID = itemRow.TaxID
		invoiceItem.TaxValue = taxValue
		invoiceItem.TaxAmount = float64(itemRow.Quantity) * itemRow.Rate * taxValue / 100
		invoiceItem.Amount = float64(itemRow.Quantity) * itemRow.Rate
		invoiceItem.Status = 1
		invoiceItem.CreatedBy = info.Email
		invoiceItem.Created = time.Now()
		invoiceItem.Updated = time.Now()
		invoiceItem.UpdatedBy = info.Email

		err = repo.CreateInvoiceItem(invoiceItem)
		if err != nil {
			msg := "create invoice item error: "
			return nil, errors.New(msg)
		}
	}
	var invoice Invoice
	invoice.InvoiceNumber = info.InvoiceNumber
	invoice.InvoiceDate = info.InvoiceDate
	invoice.DueDate = info.DueDate
	invoice.CustomerID = oldInvoice.CustomerID
	invoice.ItemCount = itemCount
	invoice.Subtotal = itemTotal
	invoice.DiscountType = info.DiscountType
	invoice.DiscountValue = info.DiscountValue
	invoice.TaxTotal = taxTotal
	invoice.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		invoice.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		invoice.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		invoice.Total = itemTotal + info.ShippingFee + taxTotal
	}
	invoice.Notes = info.Notes
	invoice.Status = 1
	invoice.Updated = time.Now()
	invoice.UpdatedBy = info.Email
	err = repo.UpdateInvoice(invoiceID, invoice)
	if err != nil {
		msg := "update invoice error: "
		return nil, errors.New(msg)
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, oldInvoice.SalesorderID)
	if err != nil {
		msg := "get sales order error: "
		return nil, errors.New(msg)
	}
	invoicedCount, err := repo.GetSalesorderInvoicedCount(info.OrganizationID, oldInvoice.SalesorderID)
	if err != nil {
		msg := "get sales order invoiced count error: "
		return nil, errors.New(msg)
	}
	invoiceStatus := 1
	if so.ItemCount == invoicedCount {
		invoiceStatus = 3
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateSalesorderInvoiceStatus(oldInvoice.SalesorderID, invoiceStatus, info.Email)
	if err != nil {
		msg := "update sales order invoice status error: "
		return nil, errors.New(msg)
	}
	soStatus := 1
	if so.ShippingStatus == 3 && invoiceStatus == 3 {
		soStatus = 3
	} else {
		soStatus = 2
	}
	err = repo.UpdateSalesorderStatus(oldInvoice.SalesorderID, soStatus, info.Email)
	if err != nil {
		msg := "update sales order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = oldInvoice.SalesorderID
	newEvent.Description = "Invoice Updated"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &invoiceID, err
}

func (s *salesorderService) DeleteInvoice(invoiceID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	oldInvoice, err := repo.GetInvoiceByID(organizationID, invoiceID)
	if err != nil {
		msg := "get picking order error"
		return errors.New(msg)
	}
	if oldInvoice.Status != 1 {
		msg := " invoice status error"
		return errors.New(msg)
	}
	invoiceItems, err := repo.GetInvoiceItemList(organizationID, invoiceID)
	if err != nil {
		msg := "get picking order log error"
		return errors.New(msg)
	}
	for _, invoiceItem := range *invoiceItems {
		oldSoItem, err := repo.GetSalesorderItemByID(organizationID, oldInvoice.SalesorderID, invoiceItem.ItemID)
		if err != nil {
			msg := "sales order item not exist"
			return errors.New(msg)
		}
		if oldSoItem.QuantityInvoiced < invoiceItem.Quantity {
			msg := "sales order itme invoiced quantity error"
			return errors.New(msg)
		}
		var soItem SalesorderItem
		soItem.SalesorderItemID = oldSoItem.SalesorderItemID
		soItem.QuantityInvoiced = oldSoItem.QuantityInvoiced - invoiceItem.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = email

		err = repo.InvoiceSalesorderItem(soItem)
		if err != nil {
			msg := "cancel invoice salesorder item error: "
			return errors.New(msg)
		}
	}
	err = repo.DeleteInvoice(invoiceID, email)
	if err != nil {
		msg := "delete picking order error: "
		return errors.New(msg)
	}
	invoicedCount, err := repo.GetSalesorderInvoicedCount(organizationID, oldInvoice.SalesorderID)
	if err != nil {
		msg := "get sales order received count error: "
		return errors.New(msg)
	}
	invoiceStatus := 1
	if invoicedCount == 0 {
		invoiceStatus = 1
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateSalesorderInvoiceStatus(oldInvoice.SalesorderID, invoiceStatus, email)
	if err != nil {
		msg := "update sales order receive status error: "
		return errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "salesorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldInvoice.SalesorderID
	newEvent.Description = "Invoice Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	return err
}

func (s *salesorderService) NewPaymentReceived(invoiceID string, info PaymentReceivedNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckPaymentReceivedNumberConfict("", info.OrganizationID, info.PaymentReceivedNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "payment number exists"
		return nil, errors.New(msg)
	}
	paymentReceivedID := "payr-" + xid.New().String()
	settingRepo := setting.NewSettingRepository(tx)

	invoice, err := repo.GetInvoiceByID(info.OrganizationID, invoiceID)
	if err != nil {
		msg := "get invoice error: "
		return nil, errors.New(msg)
	}
	invoicedPaid, err := repo.GetInvoicePaidCount(info.OrganizationID, invoiceID)
	if err != nil {
		msg := "get invoice paid count error: "
		return nil, errors.New(msg)
	}
	if invoice.Total < invoicedPaid+info.Amount {
		msg := "pay too much error: "
		return nil, errors.New(msg)
	}
	_, err = settingRepo.GetPaymentMethodByID(info.OrganizationID, info.PaymentMethodID)
	if err != nil {
		msg := "payment method not exists"
		return nil, errors.New(msg)
	}
	var paymentReceived PaymentReceived
	paymentReceived.OrganizationID = info.OrganizationID
	paymentReceived.InvoiceID = invoiceID
	paymentReceived.CustomerID = invoice.CustomerID
	paymentReceived.PaymentReceivedID = paymentReceivedID
	paymentReceived.PaymentReceivedNumber = info.PaymentReceivedNumber
	paymentReceived.PaymentReceivedDate = info.PaymentReceivedDate
	paymentReceived.PaymentMethodID = info.PaymentMethodID
	paymentReceived.Amount = info.Amount
	paymentReceived.Notes = info.Notes
	paymentReceived.Status = 1
	paymentReceived.Created = time.Now()
	paymentReceived.CreatedBy = info.Email
	paymentReceived.Updated = time.Now()
	paymentReceived.UpdatedBy = info.Email
	err = repo.CreatePaymentReceived(paymentReceived)
	if err != nil {
		msg := "create payment error: "
		return nil, errors.New(msg)
	}
	invoicedPaid, err = repo.GetInvoicePaidCount(info.OrganizationID, invoiceID)
	if err != nil {
		msg := "get invoice paid count error: "
		return nil, errors.New(msg)
	}
	invoiceStatus := 1
	if invoice.Total == invoicedPaid {
		invoiceStatus = 3
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateInvoiceStatus(invoiceID, invoiceStatus, info.Email)
	if err != nil {
		msg := "update invoice status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "invoice"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = invoiceID
	newEvent.Description = "Payment Received Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &paymentReceivedID, err
}

func (s *salesorderService) GeInvoicePaymentReceived(organizationID, invoiceID string) (float64, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	res, err := query.GeInvoicePaymentReceived(organizationID, invoiceID)
	return res, err
}

func (s *salesorderService) GetPaymentReceivedList(filter PaymentReceivedFilter) (int, *[]PaymentReceivedResponse, error) {
	db := database.RDB()
	query := NewSalesorderQuery(db)
	count, err := query.GetPaymentReceivedCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPaymentReceivedList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *salesorderService) UpdatePaymentReceived(paymentReceivedID string, info PaymentReceivedNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	isConflict, err := repo.CheckPaymentReceivedNumberConfict(paymentReceivedID, info.OrganizationID, info.PaymentReceivedNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "payment number exists"
		return nil, errors.New(msg)
	}
	settingRepo := setting.NewSettingRepository(tx)
	oldPayment, err := repo.GetPaymentReceivedByID(info.OrganizationID, paymentReceivedID)
	if err != nil {
		msg := "payment not exist"
		return nil, errors.New(msg)
	}
	invoice, err := repo.GetInvoiceByID(info.OrganizationID, oldPayment.InvoiceID)
	if err != nil {
		msg := "get invoice error: "
		return nil, errors.New(msg)
	}
	_, err = settingRepo.GetPaymentMethodByID(info.OrganizationID, info.PaymentMethodID)
	if err != nil {
		msg := "payment method not exists"
		return nil, errors.New(msg)
	}
	var paymentReceived PaymentReceived
	paymentReceived.PaymentReceivedNumber = info.PaymentReceivedNumber
	paymentReceived.PaymentReceivedDate = info.PaymentReceivedDate
	paymentReceived.PaymentMethodID = info.PaymentMethodID
	paymentReceived.Amount = info.Amount
	paymentReceived.Notes = info.Notes
	paymentReceived.Status = 1
	paymentReceived.Updated = time.Now()
	paymentReceived.UpdatedBy = info.Email
	err = repo.UpdatePaymentReceived(paymentReceivedID, paymentReceived)
	if err != nil {
		msg := "create payment error: "
		return nil, errors.New(msg)
	}
	invoicedPaid, err := repo.GetInvoicePaidCount(info.OrganizationID, oldPayment.InvoiceID)
	if err != nil {
		msg := "get invoice paid count error: "
		return nil, errors.New(msg)
	}
	invoiceStatus := 1
	if invoice.Total < invoicedPaid {
		msg := "pay too much"
		return nil, errors.New(msg)
	} else if invoice.Total == invoicedPaid {
		invoiceStatus = 3
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateInvoiceStatus(oldPayment.InvoiceID, invoiceStatus, info.Email)
	if err != nil {
		msg := "update invoice status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "invoice"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = oldPayment.InvoiceID
	newEvent.Description = "Payment Received Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &paymentReceivedID, err
}

func (s *salesorderService) DeletePaymentReceived(paymentReceivedID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewSalesorderRepository(tx)
	oldPaymentReceived, err := repo.GetPaymentReceivedByID(organizationID, paymentReceivedID)
	if err != nil {
		msg := "get payment error"
		return errors.New(msg)
	}
	err = repo.DeletePaymentReceived(paymentReceivedID, email)
	if err != nil {
		msg := "delete picking order error: "
		return errors.New(msg)
	}
	invoicedPaid, err := repo.GetInvoicePaidCount(organizationID, oldPaymentReceived.InvoiceID)
	if err != nil {
		msg := "get invoice error: "
		return errors.New(msg)
	}
	invoiceStatus := 1
	if invoicedPaid == 0 {
		invoiceStatus = 1
	} else {
		invoiceStatus = 2
	}
	err = repo.UpdateInvoiceStatus(oldPaymentReceived.InvoiceID, invoiceStatus, email)
	if err != nil {
		msg := "update invoice status error: "
		return errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "invoice"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldPaymentReceived.InvoiceID
	newEvent.Description = "Payment Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	return err
}
