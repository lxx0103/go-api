package purchaseorder

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
		var poItem PurchaseorderItem
		poItem.OrganizationID = info.OrganizationID
		poItem.PurchaseorderID = poID
		poItem.PurchaseorderItemID = "poi-" + xid.New().String()
		poItem.ItemID = item.ItemID
		poItem.Quantity = item.Quantity
		poItem.Rate = item.Rate
		poItem.TaxID = item.TaxID
		poItem.TaxValue = taxValue
		poItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
		poItem.Amount = float64(item.Quantity) * item.Rate
		poItem.QuantityReceived = 0
		poItem.QuantityBilled = 0
		poItem.Status = 1
		poItem.Created = time.Now()
		poItem.CreatedBy = info.Email
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

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
	purchaseorder.TaxTotal = taxTotal
	purchaseorder.DiscountType = info.DiscountType
	purchaseorder.DiscountValue = info.DiscountValue
	purchaseorder.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		purchaseorder.Total = itemTotal + info.ShippingFee + taxTotal
	}
	purchaseorder.Notes = info.Notes
	purchaseorder.Status = 1        //Draft
	purchaseorder.ReceiveStatus = 1 //no receive
	purchaseorder.BillingStatus = 1 //unbilled
	purchaseorder.Created = time.Now()
	purchaseorder.CreatedBy = info.Email
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.Email

	err = repo.CreatePurchaseorder(purchaseorder)
	if err != nil {
		msg := "create purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = poID
	newEvent.Description = "Purchase Order Created"
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
			poItem.TaxID = item.TaxID
			poItem.TaxValue = taxValue
			poItem.Amount = float64(item.Quantity) * item.Rate
			poItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
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
			poItem.TaxID = item.TaxID
			poItem.TaxValue = taxValue
			poItem.Amount = float64(item.Quantity) * item.Rate
			poItem.TaxAmount = float64(item.Quantity) * item.Rate * taxValue / 100
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
		taxTotal += item.Rate * float64(item.Quantity) * taxValue / 100
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
	purchaseorder.TaxTotal = taxTotal
	purchaseorder.DiscountType = info.DiscountType
	purchaseorder.DiscountValue = info.DiscountValue
	purchaseorder.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		purchaseorder.Total = itemTotal + taxTotal - info.DiscountValue + info.ShippingFee
	} else {
		purchaseorder.Total = itemTotal + taxTotal + info.ShippingFee
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
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = purchaseorderID
	newEvent.Description = "Purchase Order Updated"
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

func (s *purchaseorderService) DeletePurchaseorder(purchaseorderID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	po, err := repo.GetPurchaseorderByID(organizationID, purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return errors.New(msg)
	}
	if po.BillingStatus != 1 {
		msg := "Purchaseorder billed can not be deleted"
		return errors.New(msg)
	}
	if po.ReceiveStatus != 1 {
		msg := "Purchaseorder received can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeletePurchaseorder(purchaseorderID, email)
	if err != nil {
		return err
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = purchaseorderID
	newEvent.Description = "Purchase Order Deleted"
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

func (s *purchaseorderService) IssuePurchaseorder(purchaseorderID, organizationID, user, email string) error {
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
	err = repo.UpdatePurchaseorderStatus(purchaseorderID, 2, email) //ISSUED
	if err != nil {
		msg := "update purchaseorder error: " + err.Error()
		return errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = purchaseorderID
	newEvent.Description = "Purchase Order Issued"
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

// receive

func (s *purchaseorderService) NewPurchasereceive(purchaseorderID string, info PurchasereceiveNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckReceiveNumberConfict("", info.OrganizationID, info.PurchasereceiveNumber)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "purchase receive number exists"
		return nil, errors.New(msg)
	}
	var msgs [][]byte
	receiveID := "rec-" + xid.New().String()
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldPoItem, err := repo.GetPurchaseorderItemByID(info.OrganizationID, purchaseorderID, itemRow.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return nil, errors.New(msg)
		}
		itemInfo, err := itemRepo.GetItemByID(itemRow.ItemID, info.OrganizationID)
		if err != nil {
			return nil, err
		}
		receiveItemID := "rei-" + xid.New().String()
		if itemInfo.TrackLocation == 1 {
			warehouseRepo := warehouse.NewWarehouseRepository(tx)
			canReceived, err := warehouseRepo.GetItemAvailable(itemRow.ItemID, info.OrganizationID)
			if err != nil {
				msg := "get available location error"
				return nil, errors.New(msg)
			}
			if canReceived < itemRow.Quantity {
				msg := "no enough space to receive item"
				return nil, errors.New(msg)
			}
			quantityToReceive := itemRow.Quantity
			for quantityToReceive > 0 {
				nextLocation, err := warehouseRepo.GetNextLocation(itemRow.ItemID, info.OrganizationID)
				if err != nil {
					msg := "get next location error" + err.Error()
					return nil, errors.New(msg)
				}
				if nextLocation.Available >= quantityToReceive {
					err = warehouseRepo.ReceiveItem(nextLocation.LocationID, quantityToReceive, info.Email)
					if err != nil {
						msg := "receive item to location error"
						return nil, errors.New(msg)
					}
					var purchasereceiveDetail PurchasereceiveDetail
					purchasereceiveDetail.PurchasereceiveDetailID = "prd-" + xid.New().String()
					purchasereceiveDetail.OrganizationID = info.OrganizationID
					purchasereceiveDetail.PurchasereceiveID = receiveID
					purchasereceiveDetail.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
					purchasereceiveDetail.PurchasereceiveItemID = receiveItemID
					purchasereceiveDetail.LocationID = nextLocation.LocationID
					purchasereceiveDetail.ItemID = itemRow.ItemID
					purchasereceiveDetail.Quantity = quantityToReceive
					purchasereceiveDetail.Status = 1
					purchasereceiveDetail.Created = time.Now()
					purchasereceiveDetail.CreatedBy = info.Email
					purchasereceiveDetail.Updated = time.Now()
					purchasereceiveDetail.UpdatedBy = info.Email
					err = repo.CreatePurchasereceiveDetail(purchasereceiveDetail)
					if err != nil {
						msg := "create purchase receive detail error" + err.Error()
						return nil, errors.New(msg)
					}

					var newBatchEvent item.NewBatchCreated
					newBatchEvent.Type = "NewReceive"
					newBatchEvent.Quantity = quantityToReceive
					newBatchEvent.Balance = quantityToReceive
					newBatchEvent.ReferenceID = receiveItemID
					newBatchEvent.ItemID = itemRow.ItemID
					newBatchEvent.LocationID = nextLocation.LocationID
					newBatchEvent.OrganizationID = info.OrganizationID
					newBatchEvent.Rate = oldPoItem.Rate
					newBatchEvent.Email = info.Email
					msg, _ := json.Marshal(newBatchEvent)
					msgs = append(msgs, msg)

					quantityToReceive = 0
				} else {
					err = warehouseRepo.ReceiveItem(nextLocation.LocationID, nextLocation.Available, info.Email)
					if err != nil {
						msg := "receive item to location error"
						return nil, errors.New(msg)
					}
					var purchasereceiveDetail PurchasereceiveDetail
					purchasereceiveDetail.PurchasereceiveDetailID = "prd-" + xid.New().String()
					purchasereceiveDetail.OrganizationID = info.OrganizationID
					purchasereceiveDetail.PurchasereceiveID = receiveID
					purchasereceiveDetail.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
					purchasereceiveDetail.PurchasereceiveItemID = receiveItemID
					purchasereceiveDetail.LocationID = nextLocation.LocationID
					purchasereceiveDetail.ItemID = itemRow.ItemID
					purchasereceiveDetail.Quantity = nextLocation.Available
					purchasereceiveDetail.Status = 1
					purchasereceiveDetail.Created = time.Now()
					purchasereceiveDetail.CreatedBy = info.Email
					purchasereceiveDetail.Updated = time.Now()
					purchasereceiveDetail.UpdatedBy = info.Email
					err = repo.CreatePurchasereceiveDetail(purchasereceiveDetail)
					if err != nil {
						msg := "create purchase receive detail error"
						return nil, errors.New(msg)
					}

					var newBatchEvent item.NewBatchCreated
					newBatchEvent.Type = "NewReceive"
					newBatchEvent.Quantity = nextLocation.Available
					newBatchEvent.Balance = nextLocation.Available
					newBatchEvent.ReferenceID = receiveItemID
					newBatchEvent.ItemID = itemRow.ItemID
					newBatchEvent.LocationID = nextLocation.LocationID
					newBatchEvent.OrganizationID = info.OrganizationID
					newBatchEvent.Email = info.Email
					msg, _ := json.Marshal(newBatchEvent)
					msgs = append(msgs, msg)

					quantityToReceive = quantityToReceive - nextLocation.Available
				}
			}
		}
		if oldPoItem.Quantity < oldPoItem.QuantityReceived+itemRow.Quantity {
			msg := "receive quantity greater than Unreceived"
			return nil, errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		poItem.QuantityReceived = oldPoItem.QuantityReceived + itemRow.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

		err = repo.ReceivePurchaseorderItem(poItem)
		if err != nil {
			msg := "receive purchaseorder item error: " + err.Error()
			return nil, errors.New(msg)
		}
		var prItem PurchasereceiveItem
		prItem.OrganizationID = info.OrganizationID
		prItem.PurchasereceiveID = receiveID
		prItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		prItem.PurchasereceiveItemID = receiveItemID
		prItem.ItemID = oldPoItem.ItemID
		prItem.Quantity = itemRow.Quantity
		prItem.Status = 1
		prItem.CreatedBy = info.Email
		prItem.Created = time.Now()
		prItem.Updated = time.Now()
		prItem.UpdatedBy = info.Email

		err = repo.CreatePurchasereceiveItem(prItem)
		if err != nil {
			msg := "create purchase receive item error: " + err.Error()
			return nil, errors.New(msg)
		}
		err = itemRepo.UpdateItemStock(itemRow.ItemID, itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	var purchasereceive Purchasereceive
	purchasereceive.PurchaseorderID = purchaseorderID
	purchasereceive.PurchasereceiveID = receiveID
	purchasereceive.PurchasereceiveNumber = info.PurchasereceiveNumber
	purchasereceive.PurchasereceiveDate = info.PurchasereceiveDate
	purchasereceive.OrganizationID = info.OrganizationID
	purchasereceive.Notes = info.Notes
	purchasereceive.Status = 1
	purchasereceive.Created = time.Now()
	purchasereceive.CreatedBy = info.Email
	purchasereceive.Updated = time.Now()
	purchasereceive.UpdatedBy = info.Email
	err = repo.CreatePurchasereceive(purchasereceive)
	if err != nil {
		msg := "create purchase receive error: " + err.Error()
		return nil, errors.New(msg)
	}
	po, err := repo.GetPurchaseorderByID(info.OrganizationID, purchaseorderID)
	if err != nil {
		msg := "get purchase order error: " + err.Error()
		return nil, errors.New(msg)
	}
	receivedCount, err := repo.GetPurchaseorderReceivedCount(info.OrganizationID, purchaseorderID)
	if err != nil {
		msg := "get purchase order received count error: " + err.Error()
		return nil, errors.New(msg)
	}
	receivedStatus := 1
	if po.ItemCount == receivedCount {
		receivedStatus = 3
	} else {
		receivedStatus = 2
	}
	err = repo.UpdatePurchaseorderReceiveStatus(purchaseorderID, receivedStatus, info.Email)
	if err != nil {
		msg := "update purchase order receive status error: " + err.Error()
		return nil, errors.New(msg)
	}
	if po.BillingStatus == 3 && receivedStatus == 3 {
		err = repo.UpdatePurchaseorderStatus(purchaseorderID, 3, info.Email)
		if err != nil {
			msg := "update purchase order status error: " + err.Error()
			return nil, errors.New(msg)
		}
	} else {
		err = repo.UpdatePurchaseorderStatus(purchaseorderID, 2, info.Email)
		if err != nil {
			msg := "update purchase order status error: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = purchaseorderID
	newEvent.Description = "Purchase Receive Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewBatchCreated", msgRow)
		if err != nil {
			msg := "create event NewBatchCreated error"
			return nil, errors.New(msg)
		}
	}
	return &receiveID, err
}

func (s *purchaseorderService) GetPurchasereceiveList(filter PurchasereceiveFilter) (int, *[]PurchasereceiveResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	count, err := query.GetPurchasereceiveCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPurchasereceiveList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *purchaseorderService) GetPurchaseReceiveItemList(purchasereceiveID, organizationID string) (*[]PurchasereceiveItemResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	_, err := query.GetPurchasereceiveByID(organizationID, purchasereceiveID)
	if err != nil {
		msg := "get purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	list, err := query.GetPurchasereceiveItemList(purchasereceiveID)
	return list, err
}

func (s *purchaseorderService) GetPurchaseReceiveDetailList(purchasereceiveID, organizationID string) (*[]PurchasereceiveDetailResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	_, err := query.GetPurchasereceiveByID(organizationID, purchasereceiveID)
	if err != nil {
		msg := "get purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	list, err := query.GetPurchasereceiveDetailList(purchasereceiveID)
	return list, err
}

func (s *purchaseorderService) DeletePurchasereceive(purchasereceiveID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	itemRepo := item.NewItemRepository(tx)
	var msgs [][]byte
	var itemHistorys []string
	oldPurchasereceive, err := repo.GetPurchasereceiveByID(organizationID, purchasereceiveID)
	if err != nil {
		msg := "get purchase receive error"
		return errors.New(msg)
	}
	details, err := repo.GetPurchasereceiveDetailList(organizationID, purchasereceiveID)
	if err != nil {
		msg := "get purchase receive  detail error"
		return errors.New(msg)
	}
	for _, detail := range *details {
		oldPoItem, err := repo.GetPurchaseorderItemByID(organizationID, oldPurchasereceive.PurchaseorderID, detail.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return errors.New(msg)
		}
		itemInfo, err := itemRepo.GetItemByID(detail.ItemID, organizationID)
		if err != nil {
			msg := " item not exist"
			return errors.New(msg)
		}
		if itemInfo.StockAvailable < detail.Quantity {
			msg := "item stock available not enough"
			return errors.New(msg)
		}
		warehouseRepo := warehouse.NewWarehouseRepository(tx)
		batch, err := itemRepo.GetItemBatchByReferenceID(detail.ItemID, detail.PurchasereceiveItemID, organizationID)
		if err != nil {
			msg := " batch not exist"
			return errors.New(msg)
		}
		if batch.Balance != batch.Quantity {
			msg := "this batch of item has been used"
			return errors.New(msg)
		}
		err = itemRepo.DeleteItemBatch(batch.BatchID, email)
		if err != nil {
			msg := " delete batch error"
			return errors.New(msg)
		}
		location, err := warehouseRepo.GetLocationByID(detail.LocationID, organizationID)
		if err != nil {
			msg := " location not exist"
			return errors.New(msg)
		}
		if location.CanPick < detail.Quantity {
			msg := "location stock not enough"
			return errors.New(msg)
		}
		err = warehouseRepo.ReceiveItem(detail.LocationID, -detail.Quantity, email)
		if err != nil {
			msg := "get item from location error"
			return errors.New(msg)
		}

		if oldPoItem.QuantityReceived < detail.Quantity {
			msg := "received quantity error"
			return errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		poItem.QuantityReceived = oldPoItem.QuantityReceived - detail.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = email

		err = repo.ReceivePurchaseorderItem(poItem)
		if err != nil {
			msg := "cancel purchaseorder item error: " + err.Error()
			return errors.New(msg)
		}
		err = itemRepo.UpdateItemStock(detail.ItemID, -detail.Quantity, email)
		if err != nil {
			msg := "update item stock error: " + err.Error()
			return errors.New(msg)
		}
		itemUpdated := false
		for _, itemID := range itemHistorys {
			if itemID == detail.ItemID {
				itemUpdated = true
				continue
			}
		}
		if !itemUpdated {
			itemHistorys = append(itemHistorys, detail.ItemID)
			var newEvent common.NewHistoryCreated
			newEvent.HistoryType = "item"
			newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
			newEvent.HistoryBy = user
			newEvent.ReferenceID = detail.ItemID
			newEvent.Description = "Purchase Receive Deleted"
			newEvent.OrganizationID = organizationID
			newEvent.Email = email
			msg, _ := json.Marshal(newEvent)
			msgs = append(msgs, msg)
		}
	}
	err = repo.DeletePurchasereceive(purchasereceiveID, email)
	if err != nil {
		msg := "delete purchase receive error: " + err.Error()
		return errors.New(msg)
	}
	_, err = repo.GetPurchaseorderByID(organizationID, oldPurchasereceive.PurchaseorderID)
	if err != nil {
		msg := "get purchase order error: " + err.Error()
		return errors.New(msg)
	}
	receivedCount, err := repo.GetPurchaseorderReceivedCount(organizationID, oldPurchasereceive.PurchaseorderID)
	if err != nil {
		msg := "get purchase order received count error: " + err.Error()
		return errors.New(msg)
	}
	receivedStatus := 1
	if receivedCount == 0 {
		receivedStatus = 1
	} else {
		receivedStatus = 2
	}
	err = repo.UpdatePurchaseorderReceiveStatus(oldPurchasereceive.PurchaseorderID, receivedStatus, email)
	if err != nil {
		msg := "update purchase order receive status error: " + err.Error()
		return errors.New(msg)
	}

	err = repo.UpdatePurchaseorderStatus(oldPurchasereceive.PurchaseorderID, 2, email)
	if err != nil {
		msg := "update purchase order status error: " + err.Error()
		return errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldPurchasereceive.PurchaseorderID
	newEvent.Description = "Purchase Receive Deleted"
	newEvent.OrganizationID = organizationID
	newEvent.Email = email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	var newEvent2 common.NewHistoryCreated
	newEvent2.HistoryType = "purchasereceive"
	newEvent2.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent2.HistoryBy = user
	newEvent2.ReferenceID = purchasereceiveID
	newEvent2.Description = "Purchase Receive Deleted"
	newEvent2.OrganizationID = organizationID
	newEvent2.Email = email
	msg2, _ := json.Marshal(newEvent2)
	err = rabbit.Publish("NewHistoryCreated", msg2)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return errors.New(msg)
	}
	for _, msgRow := range msgs {
		err = rabbit.Publish("NewHistoryCreated", msgRow)
		if err != nil {
			msg := "create event NewBatchCreated error"
			return errors.New(msg)
		}
	}
	return err
}

func (s *purchaseorderService) NewBill(purchaseorderID string, info BillNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckBillNumberConfict("", info.OrganizationID, info.BillNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "bill number exists"
		return nil, errors.New(msg)
	}
	billID := "bil-" + xid.New().String()
	settingRepo := setting.NewSettingRepository(tx)
	_, err = settingRepo.GetVendorByID(info.VendorID, info.OrganizationID)
	if err != nil {
		msg := "vendor not exists"
		return nil, errors.New(msg)
	}
	itemCount := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldPoItem, err := repo.GetPurchaseorderItemByID(info.OrganizationID, purchaseorderID, itemRow.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return nil, errors.New(msg)
		}
		if oldPoItem.PurchaseorderItemID != itemRow.PurchaseorderItemID {
			msg := "purchase order item id error"
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

		billItemID := "bili-" + xid.New().String()
		if oldPoItem.Quantity < oldPoItem.QuantityBilled+itemRow.Quantity {
			msg := "invoicing quantity greater than not billd"
			return nil, errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		poItem.QuantityBilled = oldPoItem.QuantityBilled + itemRow.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

		err = repo.BillPurchaseorderItem(poItem)
		if err != nil {
			msg := "bill purchaseorder item error: "
			return nil, errors.New(msg)
		}
		var billItem BillItem
		billItem.OrganizationID = info.OrganizationID
		billItem.BillID = billID
		billItem.BillItemID = billItemID
		billItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		billItem.ItemID = oldPoItem.ItemID
		billItem.Quantity = itemRow.Quantity
		billItem.Rate = itemRow.Rate
		billItem.TaxID = itemRow.TaxID
		billItem.TaxValue = taxValue
		billItem.TaxAmount = float64(itemRow.Quantity) * itemRow.Rate * taxValue / 100
		billItem.Amount = float64(itemRow.Quantity) * itemRow.Rate
		billItem.Status = 1
		billItem.CreatedBy = info.Email
		billItem.Created = time.Now()
		billItem.Updated = time.Now()
		billItem.UpdatedBy = info.Email

		err = repo.CreateBillItem(billItem)
		if err != nil {
			msg := "create bill item error: "
			return nil, errors.New(msg)
		}
	}
	var bill Bill
	bill.OrganizationID = info.OrganizationID
	bill.BillID = billID
	bill.PurchaseorderID = purchaseorderID
	bill.BillNumber = info.BillNumber
	bill.BillDate = info.BillDate
	bill.DueDate = info.DueDate
	bill.VendorID = info.VendorID
	bill.ItemCount = itemCount
	bill.Subtotal = itemTotal
	bill.DiscountType = info.DiscountType
	bill.DiscountValue = info.DiscountValue
	bill.TaxTotal = taxTotal
	bill.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		bill.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		bill.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		bill.Total = itemTotal + info.ShippingFee + taxTotal
	}
	bill.Notes = info.Notes
	bill.Status = 1
	bill.Created = time.Now()
	bill.CreatedBy = info.Email
	bill.Updated = time.Now()
	bill.UpdatedBy = info.Email
	err = repo.CreateBill(bill)
	if err != nil {
		msg := "create bill error: "
		return nil, errors.New(msg)
	}
	so, err := repo.GetPurchaseorderByID(info.OrganizationID, purchaseorderID)
	if err != nil {
		msg := "get purchase order error: "
		return nil, errors.New(msg)
	}
	billdCount, err := repo.GetPurchaseorderBilledCount(info.OrganizationID, purchaseorderID)
	if err != nil {
		msg := "get purchase order billd count error: "
		return nil, errors.New(msg)
	}
	billStatus := 1
	if so.ItemCount == billdCount {
		billStatus = 3
	} else {
		billStatus = 2
	}
	err = repo.UpdatePurchaseorderBillStatus(purchaseorderID, billStatus, info.Email)
	if err != nil {
		msg := "update purchase order bill status error: "
		return nil, errors.New(msg)
	}
	soStatus := 1
	if so.ReceiveStatus == 3 && billStatus == 3 {
		soStatus = 3
	} else {
		soStatus = 2
	}
	err = repo.UpdatePurchaseorderStatus(purchaseorderID, soStatus, info.Email)
	if err != nil {
		msg := "update purchase order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = purchaseorderID
	newEvent.Description = "Bill Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &billID, err
}

func (s *purchaseorderService) GetBillList(filter BillFilter) (int, *[]BillResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	count, err := query.GetBillCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBillList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *purchaseorderService) GetBillItemList(billID, organizationID string) (*[]BillItemResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	_, err := query.GetBillByID(organizationID, billID)
	if err != nil {
		msg := "get pickingorder error: "
		return nil, errors.New(msg)
	}
	list, err := query.GetBillItemList(billID)
	return list, err
}

func (s *purchaseorderService) UpdateBill(billID string, info BillNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckBillNumberConfict(billID, info.OrganizationID, info.BillNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "bill number exists"
		return nil, errors.New(msg)
	}
	settingRepo := setting.NewSettingRepository(tx)
	_, err = settingRepo.GetVendorByID(info.VendorID, info.OrganizationID)
	if err != nil {
		msg := "vendor not exists"
		return nil, errors.New(msg)
	}
	oldBill, err := repo.GetBillByID(info.OrganizationID, billID)
	if err != nil {
		msg := "get bill error"
		return nil, errors.New(msg)
	}
	oldBillItems, err := repo.GetBillItemList(info.OrganizationID, billID)
	if err != nil {
		msg := "get bill item error"
		return nil, errors.New(msg)
	}
	for _, oldBillItem := range *oldBillItems {
		oldPoItem, err := repo.GetPurchaseorderItemByID(info.OrganizationID, oldBill.PurchaseorderID, oldBillItem.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return nil, errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldBillItem.PurchaseorderItemID
		poItem.QuantityBilled = oldPoItem.QuantityBilled - oldBillItem.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

		err = repo.BillPurchaseorderItem(poItem)
		if err != nil {
			msg := "cancel bill purchaseorder item error: "
			return nil, errors.New(msg)
		}

	}
	err = repo.DeleteBill(billID, info.User)
	if err != nil {
		msg := "Bill Update error"
		return nil, errors.New(msg)
	}
	itemCount := 0
	itemTotal := 0.0
	taxTotal := 0.0
	itemRepo := item.NewItemRepository(tx)
	for _, itemRow := range info.Items {
		oldPoItem, err := repo.GetPurchaseorderItemByID(info.OrganizationID, oldBill.PurchaseorderID, itemRow.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return nil, errors.New(msg)
		}
		if oldPoItem.PurchaseorderItemID != itemRow.PurchaseorderItemID {
			msg := "purchase order item id error"
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

		billItemID := "bili-" + xid.New().String()
		if oldPoItem.Quantity < oldPoItem.QuantityBilled+itemRow.Quantity {
			fmt.Println(oldPoItem.Quantity, oldPoItem.QuantityBilled, itemRow.Quantity)
			msg := "invoicing quantity greater than not billd"
			return nil, errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		poItem.QuantityBilled = oldPoItem.QuantityBilled + itemRow.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = info.Email

		err = repo.BillPurchaseorderItem(poItem)
		if err != nil {
			msg := "bill purchaseorder item error: "
			return nil, errors.New(msg)
		}
		var billItem BillItem
		billItem.OrganizationID = info.OrganizationID
		billItem.BillID = billID
		billItem.BillItemID = billItemID
		billItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		billItem.ItemID = oldPoItem.ItemID
		billItem.Quantity = itemRow.Quantity
		billItem.Rate = itemRow.Rate
		billItem.TaxID = itemRow.TaxID
		billItem.TaxValue = taxValue
		billItem.TaxAmount = float64(itemRow.Quantity) * itemRow.Rate * taxValue / 100
		billItem.Amount = float64(itemRow.Quantity) * itemRow.Rate
		billItem.Status = 1
		billItem.CreatedBy = info.Email
		billItem.Created = time.Now()
		billItem.Updated = time.Now()
		billItem.UpdatedBy = info.Email

		err = repo.CreateBillItem(billItem)
		if err != nil {
			msg := "create bill item error: "
			return nil, errors.New(msg)
		}
	}
	var bill Bill
	bill.BillNumber = info.BillNumber
	bill.BillDate = info.BillDate
	bill.DueDate = info.DueDate
	bill.VendorID = info.VendorID
	bill.ItemCount = itemCount
	bill.Subtotal = itemTotal
	bill.DiscountType = info.DiscountType
	bill.DiscountValue = info.DiscountValue
	bill.TaxTotal = taxTotal
	bill.ShippingFee = info.ShippingFee
	if info.DiscountType == 1 {
		if info.DiscountValue < 0 || info.DiscountValue > 100 {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		bill.Total = (itemTotal+taxTotal)*(1-info.DiscountValue/100) + info.ShippingFee
	} else if info.DiscountType == 2 {
		if info.DiscountValue > (itemTotal + taxTotal + info.ShippingFee) {
			msg := "discount value error"
			return nil, errors.New(msg)
		}
		bill.Total = itemTotal - info.DiscountValue + info.ShippingFee + taxTotal
	} else {
		bill.Total = itemTotal + info.ShippingFee + taxTotal
	}
	bill.Notes = info.Notes
	bill.Status = 1
	bill.Updated = time.Now()
	bill.UpdatedBy = info.Email
	err = repo.UpdateBill(billID, bill)
	if err != nil {
		msg := "update bill error: "
		return nil, errors.New(msg)
	}
	po, err := repo.GetPurchaseorderByID(info.OrganizationID, oldBill.PurchaseorderID)
	if err != nil {
		msg := "get purchase order error: "
		return nil, errors.New(msg)
	}
	billdCount, err := repo.GetPurchaseorderBilledCount(info.OrganizationID, oldBill.PurchaseorderID)
	if err != nil {
		msg := "get purchase order billd count error: "
		return nil, errors.New(msg)
	}
	billStatus := 1
	if po.ItemCount == billdCount {
		billStatus = 3
	} else {
		billStatus = 2
	}
	err = repo.UpdatePurchaseorderBillStatus(oldBill.PurchaseorderID, billStatus, info.Email)
	if err != nil {
		msg := "update purchase order bill status error: "
		return nil, errors.New(msg)
	}
	poStatus := 1
	if po.ReceiveStatus == 3 && billStatus == 3 {
		poStatus = 3
	} else {
		poStatus = 2
	}
	err = repo.UpdatePurchaseorderStatus(oldBill.PurchaseorderID, poStatus, info.Email)
	if err != nil {
		msg := "update purchase order status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = oldBill.PurchaseorderID
	newEvent.Description = "Bill Updated"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &billID, err
}

func (s *purchaseorderService) DeleteBill(billID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	oldBill, err := repo.GetBillByID(organizationID, billID)
	if err != nil {
		msg := "get picking order error"
		return errors.New(msg)
	}
	if oldBill.Status != 1 {
		msg := " bill status error"
		return errors.New(msg)
	}
	billItems, err := repo.GetBillItemList(organizationID, billID)
	if err != nil {
		msg := "get picking order log error"
		return errors.New(msg)
	}
	for _, billItem := range *billItems {
		oldPoItem, err := repo.GetPurchaseorderItemByID(organizationID, oldBill.PurchaseorderID, billItem.ItemID)
		if err != nil {
			msg := "purchase order item not exist"
			return errors.New(msg)
		}
		if oldPoItem.QuantityBilled < billItem.Quantity {
			msg := "purchase order itme billd quantity error"
			return errors.New(msg)
		}
		var poItem PurchaseorderItem
		poItem.PurchaseorderItemID = oldPoItem.PurchaseorderItemID
		poItem.QuantityBilled = oldPoItem.QuantityBilled - billItem.Quantity
		poItem.Updated = time.Now()
		poItem.UpdatedBy = email

		err = repo.BillPurchaseorderItem(poItem)
		if err != nil {
			msg := "cancel bill purchaseorder item error: "
			return errors.New(msg)
		}
	}
	err = repo.DeleteBill(billID, email)
	if err != nil {
		msg := "delete picking order error: "
		return errors.New(msg)
	}
	billedCount, err := repo.GetPurchaseorderBilledCount(organizationID, oldBill.PurchaseorderID)
	if err != nil {
		msg := "get purchase order billed count error: "
		return errors.New(msg)
	}
	billStatus := 1
	if billedCount == 0 {
		billStatus = 1
	} else {
		billStatus = 2
	}
	err = repo.UpdatePurchaseorderBillStatus(oldBill.PurchaseorderID, billStatus, email)
	if err != nil {
		msg := "update purchase order receive status error: "
		return errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "purchaseorder"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldBill.PurchaseorderID
	newEvent.Description = "Bill Deleted"
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

func (s *purchaseorderService) NewPaymentMade(billID string, info PaymentMadeNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckPaymentMadeNumberConfict("", info.OrganizationID, info.PaymentMadeNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "payment number exists"
		return nil, errors.New(msg)
	}
	paymentMadeID := "paym-" + xid.New().String()
	settingRepo := setting.NewSettingRepository(tx)

	bill, err := repo.GetBillByID(info.OrganizationID, billID)
	if err != nil {
		msg := "get bill error: "
		return nil, errors.New(msg)
	}
	billdPaid, err := repo.GetBillPaidCount(info.OrganizationID, billID)
	if err != nil {
		msg := "get bill paid count error: "
		return nil, errors.New(msg)
	}
	if bill.Total < billdPaid+info.Amount {
		msg := "pay too much error: "
		return nil, errors.New(msg)
	}
	_, err = settingRepo.GetPaymentMethodByID(info.OrganizationID, info.PaymentMethodID)
	if err != nil {
		msg := "payment method not exists"
		return nil, errors.New(msg)
	}
	var paymentMade PaymentMade
	paymentMade.OrganizationID = info.OrganizationID
	paymentMade.BillID = billID
	paymentMade.VendorID = bill.VendorID
	paymentMade.PaymentMadeID = paymentMadeID
	paymentMade.PaymentMadeNumber = info.PaymentMadeNumber
	paymentMade.PaymentMadeDate = info.PaymentMadeDate
	paymentMade.PaymentMethodID = info.PaymentMethodID
	paymentMade.Amount = info.Amount
	paymentMade.Notes = info.Notes
	paymentMade.Status = 1
	paymentMade.Created = time.Now()
	paymentMade.CreatedBy = info.Email
	paymentMade.Updated = time.Now()
	paymentMade.UpdatedBy = info.Email
	err = repo.CreatePaymentMade(paymentMade)
	if err != nil {
		msg := "create payment error: "
		return nil, errors.New(msg)
	}
	billdPaid, err = repo.GetBillPaidCount(info.OrganizationID, billID)
	if err != nil {
		msg := "get bill paid count error: "
		return nil, errors.New(msg)
	}
	billStatus := 1
	if bill.Total == billdPaid {
		billStatus = 3
	} else {
		billStatus = 2
	}
	err = repo.UpdateBillStatus(billID, billStatus, info.Email)
	if err != nil {
		msg := "update bill status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "bill"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = billID
	newEvent.Description = "Payment Made Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &paymentMadeID, err
}

func (s *purchaseorderService) GeBillPaymentMade(organizationID, billID string) (float64, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	res, err := query.GeBillPaymentMade(organizationID, billID)
	return res, err
}

func (s *purchaseorderService) GetPaymentMadeList(filter PaymentMadeFilter) (int, *[]PaymentMadeResponse, error) {
	db := database.RDB()
	query := NewPurchaseorderQuery(db)
	count, err := query.GetPaymentMadeCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPaymentMadeList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *purchaseorderService) UpdatePaymentMade(paymentMadeID string, info PaymentMadeNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	isConflict, err := repo.CheckPaymentMadeNumberConfict(paymentMadeID, info.OrganizationID, info.PaymentMadeNumber)
	if err != nil {
		msg := "check conflict error: "
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "payment number exists"
		return nil, errors.New(msg)
	}
	settingRepo := setting.NewSettingRepository(tx)
	oldPayment, err := repo.GetPaymentMadeByID(info.OrganizationID, paymentMadeID)
	if err != nil {
		msg := "payment not exist"
		return nil, errors.New(msg)
	}
	bill, err := repo.GetBillByID(info.OrganizationID, oldPayment.BillID)
	if err != nil {
		msg := "get bill error: "
		return nil, errors.New(msg)
	}
	_, err = settingRepo.GetPaymentMethodByID(info.OrganizationID, info.PaymentMethodID)
	if err != nil {
		msg := "payment method not exists"
		return nil, errors.New(msg)
	}
	var paymentMade PaymentMade
	paymentMade.PaymentMadeNumber = info.PaymentMadeNumber
	paymentMade.PaymentMadeDate = info.PaymentMadeDate
	paymentMade.PaymentMethodID = info.PaymentMethodID
	paymentMade.Amount = info.Amount
	paymentMade.Notes = info.Notes
	paymentMade.Status = 1
	paymentMade.Updated = time.Now()
	paymentMade.UpdatedBy = info.Email
	err = repo.UpdatePaymentMade(paymentMadeID, paymentMade)
	if err != nil {
		msg := "create payment error: "
		return nil, errors.New(msg)
	}
	billdPaid, err := repo.GetBillPaidCount(info.OrganizationID, oldPayment.BillID)
	if err != nil {
		msg := "get bill paid count error: "
		return nil, errors.New(msg)
	}
	billStatus := 1
	if bill.Total < billdPaid {
		msg := "pay too much"
		return nil, errors.New(msg)
	} else if bill.Total == billdPaid {
		billStatus = 3
	} else {
		billStatus = 2
	}
	err = repo.UpdateBillStatus(oldPayment.BillID, billStatus, info.Email)
	if err != nil {
		msg := "update bill status error: "
		return nil, errors.New(msg)
	}
	tx.Commit()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "bill"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = oldPayment.BillID
	newEvent.Description = "Payment Made Created"
	newEvent.OrganizationID = info.OrganizationID
	newEvent.Email = info.Email
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewHistoryCreated", msg)
	if err != nil {
		msg := "create event NewHistoryCreated error"
		return nil, errors.New(msg)
	}
	return &paymentMadeID, err
}

func (s *purchaseorderService) DeletePaymentMade(paymentMadeID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewPurchaseorderRepository(tx)
	oldPaymentMade, err := repo.GetPaymentMadeByID(organizationID, paymentMadeID)
	if err != nil {
		msg := "get payment error"
		return errors.New(msg)
	}
	err = repo.DeletePaymentMade(paymentMadeID, email)
	if err != nil {
		msg := "delete picking order error: "
		return errors.New(msg)
	}
	billdPaid, err := repo.GetBillPaidCount(organizationID, oldPaymentMade.BillID)
	if err != nil {
		msg := "get bill error: "
		return errors.New(msg)
	}
	billStatus := 1
	if billdPaid == 0 {
		billStatus = 1
	} else {
		billStatus = 2
	}
	err = repo.UpdateBillStatus(oldPaymentMade.BillID, billStatus, email)
	if err != nil {
		msg := "update bill status error: "
		return errors.New(msg)
	}
	tx.Commit()
	rabbit, _ := queue.GetConn()
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "bill"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = oldPaymentMade.BillID
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
