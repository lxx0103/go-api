package purchaseorder

import (
	"errors"
	"go-api/api/v1/item"
	"go-api/api/v1/setting"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type purchaseorderService struct {
}

func NewPurchaseorderService() *purchaseorderService {
	return &purchaseorderService{}
}

func (s *purchaseorderService) NewPurchaseorder(info PurchaseorderNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckPONumberConfict("", info.OrganizationID, info.PurchaseorderNumber)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "purchaseorder number exists"
		return nil, errors.New(msg)
	}
	poID := "po-" + xid.New().String()
	settingService := setting.NewSettingService()
	_, err = settingService.GetVendorByID(info.OrganizationID, info.VendorID)
	if err != nil {
		return nil, err
	}
	itemCount := 0
	itemTotal := 0.0
	itemService := item.NewItemService()
	for _, item := range info.Items {
		_, err = itemService.GetItemByID(info.OrganizationID, item.ItemID)
		if err != nil {
			return nil, err
		}
		itemCount += item.Quantity
		itemTotal += item.Rate * float64(item.Quantity)
		var poItem PurchaseorderItem
		poItem.OrganizationID = info.OrganizationID
		poItem.PurchaseorderID = poID
		poItem.PurchaseorderItemID = "poi-" + xid.New().String()
		poItem.ItemID = item.ItemID
		poItem.Quantity = item.Quantity
		poItem.Rate = item.Rate
		poItem.Amount = float64(item.Quantity) * item.Rate
		poItem.QuantityReceived = 0
		poItem.Status = 1
		poItem.Created = time.Now()
		poItem.CreatedBy = info.User
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.User

		err = repo.CreatePurchaseorderItem(poItem)
		if err != nil {
			msg := "create purchaseorder item error: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	var purchaseorder Purchaseorder
	purchaseorder.PurchaseorderID = poID
	purchaseorder.OrganizationID = info.OrganizationID
	purchaseorder.PurchaseorderNumber = info.PurchaseorderNumber
	purchaseorder.PurchaseorderDate = info.PurchaseorderDate
	purchaseorder.ExpectedDeliveryDate = info.ExpectedDeliveryDate
	purchaseorder.VendorID = info.VendorID
	purchaseorder.Subtotal = itemTotal
	purchaseorder.DiscountType = info.DiscountType
	purchaseorder.DiscountValue = info.DiscountValue
	purchaseorder.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = itemTotal*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = itemTotal - info.DiscountValue + info.ShippingFee
	} else {
		purchaseorder.Total = itemTotal + info.ShippingFee
	}
	purchaseorder.Notes = info.Notes
	purchaseorder.Status = info.Status
	purchaseorder.Created = time.Now()
	purchaseorder.CreatedBy = info.User
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.User

	err = repo.CreatePurchaseorder(purchaseorder)
	if err != nil {
		msg := "create purchaseordererror: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &purchaseorder.PurchaseorderID, err
}

func (s *purchaseorderService) GetPurchaseorderList(filter PurchaseorderFilter) (int, *[]PurchaseorderResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	count, err := query.GetPurchaseorderCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPurchaseorderList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *purchaseorderService) UpdatePurchaseorder(purchaseorderID string, info PurchaseorderNew) (*PurchaseorderResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckPONumberConfict("", info.OrganizationID, info.PurchaseorderNumber)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "purchaseorder SKU conflict"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetVendorByID(info.OrganizationID, info.VendorID)
	if err != nil {
		return nil, err
	}
	oldPurchaseorder, err := repo.GetPurchaseorderByID(purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	if oldPurchaseorder.OrganizationID != info.OrganizationID {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	var purchaseorder Purchaseorder
	purchaseorder.SKU = info.SKU
	purchaseorder.Name = info.Name
	purchaseorder.UnitID = info.UnitID
	purchaseorder.ManufacturerID = info.ManufacturerID
	purchaseorder.BrandID = info.BrandID
	purchaseorder.WeightUnit = info.WeightUnit
	purchaseorder.Weight = info.Weight
	purchaseorder.DimensionUnit = info.DimensionUnit
	purchaseorder.Length = info.Length
	purchaseorder.Width = info.Width
	purchaseorder.Height = info.Height
	purchaseorder.SellingPrice = info.SellingPrice
	purchaseorder.CostPrice = info.CostPrice
	purchaseorder.OpenningStock = info.OpenningStock
	purchaseorder.OpenningStockRate = info.OpenningStockRate
	purchaseorder.ReorderStock = info.ReorderStock
	purchaseorder.DefaultVendorID = info.DefaultVendorID
	purchaseorder.Description = info.Description
	purchaseorder.Status = info.Status
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.User
	err = repo.UpdatePurchaseorder(purchaseorderID, purchaseorder)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetPurchaseorderByID(purchaseorderID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, err
}

func (s *purchaseorderService) GetPurchaseorderByID(organizationID, id string) (*PurchaseorderResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	unit, err := query.GetPurchaseorderByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *purchaseorderService) DeletePurchaseorder(purchaseorderID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	oldPurchaseorder, err := repo.GetPurchaseorderByID(purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return errors.New(msg)
	}
	if oldPurchaseorder.OrganizationID != organizationID {
		msg := "Purchaseorder not exist"
		return errors.New(msg)
	}
	err = repo.DeletePurchaseorder(purchaseorderID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//Barcode

func (s *purchaseorderService) GetBarcodeList(filter BarcodeFilter) (int, *[]BarcodeResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
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

func (s *purchaseorderService) NewBarcode(info BarcodeNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckBarcodeConfict("", info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "barcode conflict"
		return nil, errors.New(msg)
	}
	purchaseorder, err := repo.GetPurchaseorderByID(info.PurchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	if purchaseorder.OrganizationID != info.OrganizationID {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	var barcode Barcode
	barcode.BarcodeID = "bar-" + xid.New().String()
	barcode.OrganizationID = info.OrganizationID
	barcode.Code = info.Code
	barcode.PurchaseorderID = info.PurchaseorderID
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

func (s *purchaseorderService) UpdateBarcode(barcodeID string, info BarcodeNew) (*BarcodeResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckBarcodeConfict(barcodeID, info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "barcode conflict"
		return nil, errors.New(msg)
	}
	purchaseorder, err := repo.GetPurchaseorderByID(info.PurchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	if purchaseorder.OrganizationID != info.OrganizationID {
		msg := "Purchaseorder not exist"
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
	barcode.PurchaseorderID = info.PurchaseorderID
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

func (s *purchaseorderService) GetBarcodeByID(organizationID, id string) (*BarcodeResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	unit, err := query.GetBarcodeByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *purchaseorderService) DeleteBarcode(barcodeID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
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
