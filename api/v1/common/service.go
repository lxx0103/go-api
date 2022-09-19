package common

import (
	"errors"
	"fmt"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type commonService struct {
}

func NewCommonService() *commonService {
	return &commonService{}
}

func (s *commonService) NewHistory(info HistoryNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	var history History
	history.HistoryID = "his-" + xid.New().String()
	history.OrganizationID = info.OrganizationID
	history.HistoryBy = info.HistoryBy
	history.HistoryTime = info.HistoryTime
	history.HistoryType = info.HistoryType
	history.ReferenceID = info.ReferenceID
	history.Description = info.Description
	history.Status = 1
	history.Created = time.Now()
	history.CreatedBy = info.User
	history.Updated = time.Now()
	history.UpdatedBy = info.User

	err = repo.CreateHistory(history)
	if err != nil {
		msg := "create history error: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &history.HistoryID, err
}

func (s *commonService) GetHistoryList(filter HistoryFilter) (int, *[]HistoryResponse, error) {
	db := database.RDB()
	query := NewCommonQuery(db)
	count, err := query.GetHistoryCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetHistoryList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *commonService) GetNextNumber(filter NumberFilter) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	lastNumber, err := repo.GetLastNumber(filter)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			lastNumber = 0
		} else {
			return nil, err
		}
	}
	currentNumber := lastNumber + 1
	prefix := ""
	switch filter.NumberType {
	case "purchaseorder":
		prefix = "PO"
	case "salesorder":
		prefix = "SO"
	case "purchasereceive":
		prefix = "PR"
	case "pickingorder":
		prefix = "PIC"
	case "package":
		prefix = "PAC"
	case "shippingorder":
		prefix = "SHIP"
	case "invoice":
		prefix = "INV"
	case "paymentreceived":
		prefix = "PAYR"
	case "bill":
		prefix = "BIL"
	case "paymentmade":
		prefix = "PAYM"
	default:
		msg := "number type error"
		return nil, errors.New(msg)
	}
	res := prefix + fmt.Sprintf("%05d", currentNumber)
	if lastNumber == 0 {
		err = repo.CreateNumber(filter)
		if err != nil {
			msg := "create number error"
			return nil, errors.New(msg)
		}
	} else {
		err = repo.UpdateNumber(filter, currentNumber)
		if err != nil {
			msg := "update number error"
			return nil, errors.New(msg)
		}
	}
	tx.Commit()
	return &res, nil
}
