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
		msg := "check conflict error: " + err.Error()
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
			msg := "create salesorder item error: " + err.Error()
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
		msg := "create salesorder error: " + err.Error()
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
		msg := "check conflict error: " + err.Error()
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
				msg := "update salesorder item error: " + err.Error()
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
				msg := "create salesorder item error: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		itemCount += item.Quantity
		itemTotal += item.Rate * float64(item.Quantity)
		taxTotal += item.Rate * float64(item.Quantity) * taxValue / 100
	}
	itemDeletedError, err := repo.CheckSOItem(salesorderID, info.OrganizationID)
	if err != nil {
		msg := "check salesorder item error: " + err.Error()
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
		msg := "update salesorder error: " + err.Error()
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
		msg := "get salesorder error: " + err.Error()
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
	_, err = repo.GetSalesorderByID(organizationID, salesorderID)
	if err != nil {
		msg := "Salesorder not exist"
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
		msg := "get salesorder error: " + err.Error()
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
		msg := "update salesorder error: " + err.Error()
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
		msg := "check conflict error: " + err.Error()
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
					msg := "get next batch error" + err.Error()
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
					pickingorderLog.SalesorderItemID = oldSoItem.SalesorderItemID
					pickingorderLog.PickingorderItemID = pickingorderItemID
					pickingorderLog.LocationID = nextBatch.LocationID
					pickingorderLog.ItemID = itemRow.ItemID
					pickingorderLog.Quantity = quantityToPick
					pickingorderLog.Status = 1
					pickingorderLog.Created = time.Now()
					pickingorderLog.CreatedBy = info.Email
					pickingorderLog.Updated = time.Now()
					pickingorderLog.UpdatedBy = info.Email
					err = repo.CreatePickingorderLog(pickingorderLog)
					if err != nil {
						msg := "create picking order log error1" + err.Error()
						return nil, errors.New(msg)
					}
					err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, quantityToPick, info.Email)
					if err != nil {
						msg := "update location canpick error: " + err.Error()
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
					pickingorderLog.SalesorderItemID = oldSoItem.SalesorderItemID
					pickingorderLog.PickingorderItemID = pickingorderItemID
					pickingorderLog.LocationID = nextBatch.LocationID
					pickingorderLog.ItemID = itemRow.ItemID
					pickingorderLog.Quantity = nextBatch.Balance
					pickingorderLog.Status = 1
					pickingorderLog.Created = time.Now()
					pickingorderLog.CreatedBy = info.Email
					pickingorderLog.Updated = time.Now()
					pickingorderLog.UpdatedBy = info.Email
					err = repo.CreatePickingorderLog(pickingorderLog)
					if err != nil {
						msg := "create picking order log error" + err.Error()
						return nil, errors.New(msg)
					}
					err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, nextBatch.Balance, info.Email)
					if err != nil {
						msg := "update location canpick error: " + err.Error()
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
			msg := "pick salesorder item error: " + err.Error()
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
			msg := "create picking order item error: " + err.Error()
			return nil, errors.New(msg)
		}
		err = itemRepo.UpdateItemPickingStock(itemRow.ItemID, itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	logs, err := repo.GetPickingorderLogSum(pickingorderID)
	if err != nil {
		msg := "get picking order logs  error: " + err.Error()
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
			msg := "create picking order detail error: " + err.Error()
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
		msg := "create picking order error: " + err.Error()
		return nil, errors.New(msg)
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order error: " + err.Error()
		return nil, errors.New(msg)
	}
	receivedCount, err := repo.GetSalesorderPickedCount(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order received count error: " + err.Error()
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
		msg := "update sales order receive status error: " + err.Error()
		return nil, errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(salesorderID, 2, info.Email)
	if err != nil {
		msg := "update sales order status error: " + err.Error()
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
		msg := "get pickingorder error: " + err.Error()
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
		msg := "get pickingorder error: " + err.Error()
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
		msg := "check conflict error: " + err.Error()
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
			msg := "get salesorder error: " + err.Error()
			return nil, errors.New(msg)
		}
		items, err := repo.GetSalesorderItemList(info.OrganizationID, soID)
		if err != nil {
			msg := "get salesorder items error: " + err.Error()
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
						msg := "get next batch error" + err.Error()
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
						pickingorderLog.SalesorderItemID = itemRow.SalesorderItemID
						pickingorderLog.PickingorderItemID = pickingorderItemID
						pickingorderLog.LocationID = nextBatch.LocationID
						pickingorderLog.ItemID = itemRow.ItemID
						pickingorderLog.Quantity = quantityToPick
						pickingorderLog.Status = 1
						pickingorderLog.Created = time.Now()
						pickingorderLog.CreatedBy = info.Email
						pickingorderLog.Updated = time.Now()
						pickingorderLog.UpdatedBy = info.Email
						err = repo.CreatePickingorderLog(pickingorderLog)
						if err != nil {
							msg := "create picking order detail error1" + err.Error()
							return nil, errors.New(msg)
						}
						err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, quantityToPick, info.Email)
						if err != nil {
							msg := "update location canpick error: " + err.Error()
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
						pickingorderLog.SalesorderItemID = itemRow.SalesorderItemID
						pickingorderLog.PickingorderItemID = pickingorderItemID
						pickingorderLog.LocationID = nextBatch.LocationID
						pickingorderLog.ItemID = itemRow.ItemID
						pickingorderLog.Quantity = nextBatch.Balance
						pickingorderLog.Status = 1
						pickingorderLog.Created = time.Now()
						pickingorderLog.CreatedBy = info.Email
						pickingorderLog.Updated = time.Now()
						pickingorderLog.UpdatedBy = info.Email
						err = repo.CreatePickingorderLog(pickingorderLog)
						if err != nil {
							msg := "create picking order detail error" + err.Error()
							return nil, errors.New(msg)
						}
						err = warehouseRepo.UpdateLocationCanPick(nextBatch.LocationID, nextBatch.Balance, info.Email)
						if err != nil {
							msg := "update location canpick error: " + err.Error()
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
				msg := "pick salesorder item error: " + err.Error()
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
				msg := "create picking order item error: " + err.Error()
				return nil, errors.New(msg)
			}
			err = itemRepo.UpdateItemPickingStock(itemRow.ItemID, toPick, info.Email)
			if err != nil {
				msg := "update item stock error: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		picking_status := 3
		err = repo.UpdateSalesorderPickingStatus(soID, picking_status, info.Email)
		if err != nil {
			msg := "update sales order picking status error: " + err.Error()
			return nil, errors.New(msg)
		}
		err = repo.UpdateSalesorderStatus(soID, 2, info.Email)
		if err != nil {
			msg := "update sales order status error: " + err.Error()
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
		msg := "get picking order logs  error: " + err.Error()
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
			msg := "create picking order detail error: " + err.Error()
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
		msg := "create picking order error: " + err.Error()
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
		msg := "picking order detail not exist" + err.Error()
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
		msg := "update pickingorder error: " + err.Error()
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
		msg := "check conflict error: " + err.Error()
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
		soItem.QuantityPicked = oldSoItem.QuantityPicked - itemRow.Quantity
		soItem.QuantityPacked = oldSoItem.QuantityPacked + itemRow.Quantity
		soItem.Updated = time.Now()
		soItem.UpdatedBy = info.Email

		err = repo.PackSalesorderItem(soItem)
		if err != nil {
			msg := "pack salesorder item error: " + err.Error()
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
			msg := "create package item error: " + err.Error()
			return nil, errors.New(msg)
		}
		err = itemRepo.UpdateItemPackedStock(itemRow.ItemID, itemRow.Quantity, info.Email)
		if err != nil {
			msg := "update item stock error: " + err.Error()
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
		msg := "create package error: " + err.Error()
		return nil, errors.New(msg)
	}
	so, err := repo.GetSalesorderByID(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order error: " + err.Error()
		return nil, errors.New(msg)
	}
	packedCount, err := repo.GetSalesorderPackedCount(info.OrganizationID, salesorderID)
	if err != nil {
		msg := "get sales order packed count error: " + err.Error()
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
		msg := "update sales order packing status error: " + err.Error()
		return nil, errors.New(msg)
	}
	err = repo.UpdateSalesorderStatus(salesorderID, 2, info.Email)
	if err != nil {
		msg := "update sales order status error: " + err.Error()
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
		msg := "get package error: " + err.Error()
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
		msg := "check conflict error: " + err.Error()
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
			msg := "get package error: " + err.Error()
			return nil, errors.New(msg)
		}
		if packageInfo.Status != 1 {
			msg := "package status error for package: " + packageInfo.PackageNumber
			return nil, errors.New(msg)
		}
		items, err := repo.GetPackageItemList(info.OrganizationID, packageID)
		if err != nil {
			msg := "get package items error: " + err.Error()
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
			soItem.QuantityPacked = salesorderItem.QuantityPacked - itemRow.Quantity
			soItem.QuantityShipped = salesorderItem.QuantityShipped + itemRow.Quantity
			soItem.Updated = time.Now()
			soItem.UpdatedBy = info.Email

			err = repo.ShipSalesorderItem(soItem)
			if err != nil {
				msg := "ship salesorder item error: " + err.Error()
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
				msg := "create shipping order item error: " + err.Error()
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
				newEvent.ReferenceID = packageInfo.SalesorderID
				newEvent.Description = "Shipping Order Created"
				newEvent.OrganizationID = info.OrganizationID
				newEvent.Email = info.Email
				msg, _ := json.Marshal(newEvent)
				msgs = append(msgs, msg)
			}
			shippedCount, err := repo.GetSalesorderShippedCount(info.OrganizationID, packageInfo.SalesorderID)
			if err != nil {
				msg := "get sales order packed count error: " + err.Error()
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
				msg := "update salesorder shipping status error: " + err.Error()
				return nil, errors.New(msg)
			}
			if salesorderInfo.InvoiceStatus == 3 && shippingStatus == 3 {
				err = repo.UpdateSalesorderStatus(packageInfo.SalesorderID, 3, info.Email)
				if err != nil {
					msg := "update sales order status error: " + err.Error()
					return nil, errors.New(msg)
				}
			}
		}
		err = repo.UpdatePackageStatus(packageID, 2, info.Email)
		if err != nil {
			msg := "update package status error: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	details, err := repo.GetShippingorderDetailSum(shippingorderID)
	if err != nil {
		msg := "get picking order detail  error: " + err.Error()
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
			msg := "create picking order item error: " + err.Error()
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
		msg := "create shipping order error: " + err.Error()
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
		msg := "get package error: " + err.Error()
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
		msg := "get package error: " + err.Error()
		return nil, errors.New(msg)
	}
	list, err := query.GetShippingorderDetailList(salesorderID)
	return list, err
}
