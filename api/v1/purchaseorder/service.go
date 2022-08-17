package purchaseorder

import (
	"encoding/json"
	"errors"
	"go-api/api/v1/history"
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
	purchaseorder.CreatedBy = info.Email
	purchaseorder.Updated = time.Now()
	purchaseorder.UpdatedBy = info.Email

	err = repo.CreatePurchaseorder(purchaseorder)
	if err != nil {
		msg := "create purchaseorder error: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent history.NewHistoryCreated
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
	var newEvent history.NewHistoryCreated
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
	_, err = repo.GetPurchaseorderByID(organizationID, purchaseorderID)
	if err != nil {
		msg := "Purchaseorder not exist"
		return errors.New(msg)
	}
	err = repo.DeletePurchaseorder(purchaseorderID, email)
	if err != nil {
		return err
	}
	var newEvent history.NewHistoryCreated
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
	var newEvent history.NewHistoryCreated
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
					quantityToReceive = quantityToReceive - itemRow.Quantity
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
		err = itemRepo.UpdateItemStock(itemRow.ItemID, itemInfo.StockOnHand+itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: " + err.Error()
			return nil, errors.New(msg)
		}

		var newBatchEvent item.NewBatchCreated
		newBatchEvent.Type = "NewReceive"
		newBatchEvent.Quantity = itemRow.Quantity
		newBatchEvent.Balance = itemRow.Quantity
		newBatchEvent.ReferenceID = receiveItemID
		newBatchEvent.ItemID = itemRow.ItemID
		newBatchEvent.OrganizationID = info.OrganizationID
		newBatchEvent.Email = info.Email
		msg, _ := json.Marshal(newBatchEvent)
		msgs = append(msgs, msg)
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
	var newEvent history.NewHistoryCreated
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
