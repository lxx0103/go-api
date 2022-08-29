package common

import (
	"encoding/json"
	"fmt"
	"go-api/core/database"
	"go-api/core/queue"
	"time"

	"github.com/rs/xid"
	"github.com/streadway/amqp"
)

type NewHistoryCreated struct {
	HistoryType    string `json:"history_type" binding:"required,min=6,max=64"`
	HistoryTime    string `json:"history_time" binding:"required,datetime=2006-01-02 15:04:05"`
	HistoryBy      string `json:"history_by" binding:"required"`
	ReferenceID    string `json:"reference_id" binding:"required"`
	Description    string `json:"description" binding:"required,max=255"`
	OrganizationID string `json:"organiztion_id" binding:"required,max=64"`
	User           string `json:"user"  binding:"required,max=64"`
	Email          string `json:"email" binding:"required,max=255"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("CreateNewHistory", "NewHistoryCreated", CreateNewHistory)
}

func CreateNewHistory(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var info NewHistoryCreated
	err := json.Unmarshal(d.Body, &info)
	if err != nil {
		fmt.Println(err)
		return false
	}
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		// msg := "begin transaction error"
		return false
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
	history.CreatedBy = info.Email
	history.Updated = time.Now()
	history.UpdatedBy = info.Email

	err = repo.CreateHistory(history)
	if err != nil {
		msg := "create history error: " + err.Error()
		fmt.Println(msg)
		return false
	}
	if err != nil {
		fmt.Println(err)
		return false
	}
	tx.Commit()
	return true
}
