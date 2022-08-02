package auth

import (
	"encoding/json"
	"fmt"
	"go-api/core/database"
	"go-api/core/log"
	"go-api/core/queue"
	"time"

	"github.com/rs/xid"
	"github.com/streadway/amqp"
)

type NewOrganizationCreated struct {
	OrganizationID string `json:"organization_id"`
	Owner          string `json:"owner"`
	Password       string `json:"password"`
}

type NewAuthCreated struct {
	AuthID     int64  `json:"auth_id"`
	AuthType   int    `json:"auth_type"`
	Identifier string `json:"identifier"`
	Credential string `json:"credential"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type NewProfileCreated struct {
	AuthID int64 `json:"auth_id"`
	UserID int64 `json:"user_id"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("CreateOrganizationOwner", "NewOrganizationCreated", CreateOrganizationOwner)
}

func CreateOrganizationOwner(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var event NewOrganizationCreated
	err := json.Unmarshal(d.Body, &event)
	if err != nil {
		fmt.Println(err)
		return false
	}
	hashed, err := hashPassword(event.Password)
	if err != nil {
		// msg := "hash password error"
		return false
	}
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		// msg := "begin transaction error"
		return false
	}
	defer tx.Rollback()
	repo := NewAuthRepository(tx)
	var roleInfo Role
	roleInfo.RoleID = "role-" + xid.New().String()
	roleInfo.OrganizationID = event.OrganizationID
	roleInfo.Name = "owner"
	roleInfo.Priority = 99
	roleInfo.IsAdmin = 1
	roleInfo.IsDefault = 1
	roleInfo.Created = time.Now()
	roleInfo.CreatedBy = "SIGNUP"
	roleInfo.Updated = time.Now()
	roleInfo.UpdatedBy = "SIGNUP"
	roleInfo.Status = 2
	err = repo.CreateRole(roleInfo)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var info User
	info.UserID = "user-" + xid.New().String()
	info.Email = event.Owner
	info.Password = hashed
	info.OrganizationID = event.OrganizationID
	info.Status = 2
	info.RoleID = roleInfo.RoleID
	info.Created = time.Now()
	info.CreatedBy = "SIGNUP"
	info.Updated = time.Now()
	info.UpdatedBy = "SIGNUP"
	userID, err := repo.CreateUser(info)
	if err != nil {
		fmt.Println(err)
		return false
	}
	log.Debug(fmt.Sprintf("%d", userID))
	tx.Commit()
	return true
}
