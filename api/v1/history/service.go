package history

import (
	"errors"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type historyService struct {
}

func NewHistoryService() *historyService {
	return &historyService{}
}

func (s *historyService) NewHistory(info HistoryNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewHistoryRepository(tx)
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

func (s *historyService) GetHistoryList(filter HistoryFilter) (int, *[]HistoryResponse, error) {
	db := database.RDB()
	query := NewHistoryQuery(db)
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
