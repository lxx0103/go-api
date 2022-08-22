package item

import (
	"encoding/json"
	"fmt"
	"go-api/core/database"
	"go-api/core/queue"
	"time"

	"github.com/rs/xid"
	"github.com/streadway/amqp"
)

type NewBatchCreated struct {
	Type           string `json:"type" binding:"required,min=6,max=64"`
	Quantity       int    `json:"quantity" binding:"required"`
	Balance        int    `json:"balance" binding:"required"`
	ItemID         string `json:"item_id" binding:"required"`
	ReferenceID    string `json:"reference_id" binding:"required"`
	LocationID     string `json:"location_id" binding:"required"`
	OrganizationID string `json:"organiztion_id" binding:"required"`
	User           string `json:"user"  binding:"required,max=64"`
	Email          string `json:"email" binding:"required,max=255"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("CreateBatch", "NewBatchCreated", CreateBatch)
}

func CreateBatch(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var info NewBatchCreated
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
	repo := NewItemRepository(tx)
	var batch ItemBatch
	batch.OrganizationID = info.OrganizationID
	batch.BatchID = "bat-" + xid.New().String()
	batch.ItemID = info.ItemID
	batch.Type = info.Type
	batch.ReferenceID = info.ReferenceID
	batch.Quantity = info.Quantity
	batch.Balance = info.Balance
	batch.Status = 1
	batch.Created = time.Now()
	batch.CreatedBy = info.Email
	batch.Updated = time.Now()
	batch.UpdatedBy = info.Email

	err = repo.CreateItemBatch(batch)
	if err != nil {
		msg := "create batch error: " + err.Error()
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
