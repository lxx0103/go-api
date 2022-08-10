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
		poItem.QuantityBilled = 0
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
	purchaseorder.ItemCount = itemCount
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
	purchaseorder.Status = 1        //Draft
	purchaseorder.ReceiveStatus = 1 //no receive
	purchaseorder.BillingStatus = 1 //unbilled
	purchaseorder.Created = time.Now()
	purchaseorder.CreatedBy = info.User
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.User

	err = repo.CreatePurchaseorder(purchaseorder)
	if err != nil {
		msg := "create purchaseorder error: " + err.Error()
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

func (s *purchaseorderService) UpdatePurchaseorder(purchaseorderID string, info PurchaseorderNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckPONumberConfict(purchaseorderID, info.OrganizationID, info.PurchaseorderNumber)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "purchaseorder number conflict"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetVendorByID(info.OrganizationID, info.VendorID)
	if err != nil {
		return nil, err
	}
	oldPurchaseorder, err := repo.GetPurchaseorderByID(info.OrganizationID, purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return nil, errors.New(msg)
	}
	err = repo.DeletePurchaseorder(purchaseorderID, info.User)
	if err != nil {
		msg := "Purchaseorder Update error"
		return nil, errors.New(msg)
	}
	itemCount := 0
	quantityBilled := 0
	quantityReceived := 0
	itemTotal := 0.0
	itemService := item.NewItemService()
	for _, item := range info.Items {
		_, err = itemService.GetItemByID(info.OrganizationID, item.ItemID)
		if err != nil {
			return nil, err
		}
		if item.PurchaseorderItemID != "" {
			oldItem, err := repo.GetPurchaseorderItemByIDAll(info.OrganizationID, purchaseorderID, item.PurchaseorderItemID)
			if err != nil {
				msg := "Purchaseorder Item not exist"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityBilled > item.Quantity {
				msg := "can not set quantity lower than quantity billed"
				return nil, errors.New(msg)
			}
			if oldItem.QuantityReceived > item.Quantity {
				msg := "can not set quantity lower than quantity received"
				return nil, errors.New(msg)
			}
			quantityBilled += oldItem.QuantityBilled
			quantityReceived += oldItem.QuantityReceived
			var poItem PurchaseorderItem
			poItem.Quantity = item.Quantity
			poItem.Rate = item.Rate
			poItem.Amount = float64(item.Quantity) * item.Rate
			poItem.Status = 1
			poItem.Updated = time.Now()
			poItem.UpdatedBy = info.User
			err = repo.UpdatePurchaseorderItem(item.PurchaseorderItemID, poItem)
			if err != nil {
				msg := "update purchaseorder item error: " + err.Error()
				return nil, errors.New(msg)
			}
		} else {
			var poItem PurchaseorderItem
			poItem.OrganizationID = info.OrganizationID
			poItem.PurchaseorderID = purchaseorderID
			poItem.PurchaseorderItemID = "poi-" + xid.New().String()
			poItem.ItemID = item.ItemID
			poItem.Quantity = item.Quantity
			poItem.Rate = item.Rate
			poItem.Amount = float64(item.Quantity) * item.Rate
			poItem.QuantityReceived = 0
			poItem.QuantityBilled = 0
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
		itemCount += item.Quantity
		itemTotal += item.Rate * float64(item.Quantity)
	}
	itemDeletedError, err := repo.CheckPOItem(purchaseorderID, info.OrganizationID)
	if err != nil {
		msg := "check purchaseorder item error: " + err.Error()
		return nil, errors.New(msg)
	}
	if itemDeletedError {
		msg := "item received or billed can not be delete"
		return nil, errors.New(msg)
	}
	var purchaseorder Purchaseorder
	purchaseorder.PurchaseorderNumber = info.PurchaseorderNumber
	purchaseorder.PurchaseorderDate = info.PurchaseorderDate
	purchaseorder.ExpectedDeliveryDate = info.ExpectedDeliveryDate
	purchaseorder.VendorID = info.VendorID
	purchaseorder.ItemCount = itemCount
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
	if quantityBilled > 0 {
		if quantityBilled == itemCount {
			purchaseorder.BillingStatus = 3
		} else {
			purchaseorder.BillingStatus = 2
		}
	} else {
		purchaseorder.BillingStatus = 1
	}
	if quantityReceived > 0 {
		if quantityReceived == itemCount {
			purchaseorder.ReceiveStatus = 3
		} else {
			purchaseorder.ReceiveStatus = 2
		}
	} else {
		purchaseorder.ReceiveStatus = 1
	}
	if purchaseorder.BillingStatus == 3 && purchaseorder.ReceiveStatus == 3 {
		purchaseorder.Status = 3 //Draft
	} else {
		if oldPurchaseorder.Status == 3 {
			purchaseorder.Status = 2
		} else {
			purchaseorder.Status = oldPurchaseorder.Status
		}
	}
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.User

	err = repo.UpdatePurchaseorder(purchaseorderID, purchaseorder)
	if err != nil {
		msg := "update purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &purchaseorderID, err
}

func (s *purchaseorderService) GetPurchaseorderByID(organizationID, id string) (*PurchaseorderResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	purchaseorder, err := query.GetPurchaseorderByID(organizationID, id)
	if err != nil {
		msg := "get purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	return purchaseorder, nil
}

func (s *purchaseorderService) DeletePurchaseorder(purchaseorderID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	_, err = repo.GetPurchaseorderByID(organizationID, purchaseorderID)
	if err != nil {
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

func (s *purchaseorderService) GetPurchaseorderItemList(purchaseorderID, organizationID string) (*[]PurchaseorderItemResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	_, err := query.GetPurchaseorderByID(organizationID, purchaseorderID)
	if err != nil {
		msg := "get purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	list, err := query.GetPurchaseorderItemList(purchaseorderID)
	return list, err
}

func (s *purchaseorderService) IssuePurchaseorder(purchaseorderID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	oldPurchaseorder, err := repo.GetPurchaseorderByID(organizationID, purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return errors.New(msg)
	}
	if oldPurchaseorder.Status != 1 {
		msg := "Purchaseorder status error"
		return errors.New(msg)
	}
	err = repo.UpdatePurchaseorderStatus(purchaseorderID, 2, user) //ISSUED
	if err != nil {
		msg := "update purchaseorder error: " + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	return err
}
