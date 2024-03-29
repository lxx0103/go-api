package organization

import (
	"encoding/json"
	"errors"
	"go-api/core/database"
	"go-api/core/queue"
	"time"

	"github.com/rs/xid"
)

type organizationService struct {
}

func NewOrganizationService() *organizationService {
	return &organizationService{}
}

func (s *organizationService) NewOrganization(info OrganizationNew) (*int64, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	isConflict, err := repo.CheckConfict(info.Email)
	if err != nil {
		msg := "check conflict error"
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "organization owner exists"
		return nil, errors.New(msg)
	}
	var organization Organization
	organization.OrganizationID = "org-" + xid.New().String()
	organization.Name = info.Name
	organization.Owner = info.UserName
	organization.OwnerEmail = info.Email
	organization.Status = 2
	organization.Created = time.Now()
	organization.CreatedBy = "SIGNUP"
	organization.Updated = time.Now()
	organization.UpdatedBy = "SIGNUP"
	organizationID, err := repo.CreateOrganization(organization)
	if err != nil {
		msg := "create organizationerror: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent NewOrganizationCreated
	newEvent.OrganizationID = organization.OrganizationID
	newEvent.Owner = info.UserName
	newEvent.OwnerEmail = info.Email
	newEvent.Password = info.Password
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewOrganizationCreated", msg)
	if err != nil {
		msg := "create event NewOrganizationCreated error"
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &organizationID, err
}

// func (s *organizationService) GetOrganizationList(filter OrganizationFilter) (int, *[]Organization, error) {
// 	db := database.RDB()
// 	query := NewOrganizationQuery(db)
// 	count, err := query.GetOrganizationCount(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	list, err := query.GetOrganizationList(filter)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	return count, list, err
// }

// func (s *organizationService) UpdateOrganization(organizationID int64, info OrganizationNew) (*Organization, error) {
// 	db := database.WDB()
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()
// 	repo := NewOrganizationRepository(tx)
// 	_, err = repo.UpdateOrganization(organizationID, info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	organization, err := repo.GetOrganizationByID(organizationID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tx.Commit()
// 	return organization, err
// }
