package crm

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-api/api/v1/common"
	"go-api/api/v1/setting"
	"go-api/core/database"
	"go-api/core/queue"
	"time"

	"github.com/rs/xid"
)

type crmService struct {
}

func NewCrmService() *crmService {
	return &crmService{}
}

func (s *crmService) NewLead(info LeadNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewCrmRepository(tx)
	leadID := "lead-" + xid.New().String()
	var lead Lead
	lead.LeadID = leadID
	lead.OrganizationID = info.OrganizationID
	lead.Source = info.Source
	lead.Company = info.Company
	lead.Salutation = info.Salutation
	lead.FirstName = info.FirstName
	lead.LastName = info.LastName
	lead.Email = info.LeadEmail
	lead.Phone = info.Phone
	lead.Mobile = info.Mobile
	lead.Fax = info.Fax
	lead.Country = info.Country
	lead.State = info.State
	lead.City = info.City
	lead.Address1 = info.Address1
	lead.Address2 = info.Address2
	lead.Zip = info.Zip
	lead.Notes = info.Notes
	if info.Status == 0 {
		lead.Status = 1 //New
	} else {
		lead.Status = info.Status
	}
	lead.Created = time.Now()
	lead.CreatedBy = info.Email
	lead.Updated = time.Now()
	lead.UpdatedBy = info.Email
	fmt.Println(lead.Notes)
	err = repo.CreateLead(lead)
	if err != nil {
		msg := "create lead error: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "crm"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = leadID
	newEvent.Description = "Lead Created"
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
	return &leadID, err
}

func (s *crmService) GetLeadList(filter LeadFilter) (int, *[]LeadResponse, error) {
	db := database.RDB()
	query := NewCrmQuery(db)
	count, err := query.GetLeadCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetLeadList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *crmService) UpdateLead(LeadID string, info LeadNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewCrmRepository(tx)
	oldLead, err := repo.GetLeadByID(info.OrganizationID, LeadID)
	if err != nil {
		msg := "Lead not exist"
		return nil, errors.New(msg)
	}
	var lead Lead
	lead.Source = info.Source
	lead.Company = info.Company
	lead.Salutation = info.Salutation
	lead.FirstName = info.FirstName
	lead.LastName = info.LastName
	lead.Email = info.LeadEmail
	lead.Phone = info.Phone
	lead.Mobile = info.Mobile
	lead.Fax = info.Fax
	lead.Country = info.Country
	lead.State = info.State
	lead.City = info.City
	lead.Address1 = info.Address1
	lead.Address2 = info.Address2
	lead.Zip = info.Zip
	lead.Status = oldLead.Status
	lead.Notes = info.Notes
	lead.Updated = time.Now()
	lead.UpdatedBy = info.Email
	err = repo.UpdateLead(LeadID, lead)
	if err != nil {
		msg := "update lead error: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "crm"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = LeadID
	newEvent.Description = "Lead Updated"
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
	return &LeadID, err
}

func (s *crmService) GetLeadByID(organizationID, id string) (*LeadResponse, error) {
	db := database.RDB()
	query := NewCrmQuery(db)
	crm, err := query.GetLeadByID(organizationID, id)
	if err != nil {
		msg := "get crm error: " + err.Error()
		return nil, errors.New(msg)
	}
	return crm, nil
}

func (s *crmService) DeleteLead(crmID, organizationID, user, email string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCrmRepository(tx)
	_, err = repo.GetLeadByID(organizationID, crmID)
	if err != nil {
		msg := "Lead not exist"
		return errors.New(msg)
	}
	err = repo.DeleteLead(crmID, email)
	if err != nil {
		return err
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "crm"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = user
	newEvent.ReferenceID = crmID
	newEvent.Description = "Lead Deleted"
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

func (s *crmService) ConvertLead(LeadID string, info LeadConvertNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewCrmRepository(tx)
	settingRepo := setting.NewSettingRepository(tx)
	oldLead, err := repo.GetLeadByID(info.OrganizationID, LeadID)
	if err != nil {
		msg := "Lead not exist"
		return nil, errors.New(msg)
	}
	if info.Type == "customer" {
		var customer setting.Customer
		customer.CustomerID = "cus-" + xid.New().String()
		customer.OrganizationID = info.OrganizationID
		customer.Name = oldLead.Company
		customer.ContactSalutation = oldLead.Salutation
		customer.ContactFirstName = oldLead.FirstName
		customer.ContactLastName = oldLead.LastName
		customer.ContactEmail = oldLead.Email
		customer.ContactPhone = oldLead.Mobile
		customer.Country = oldLead.Country
		customer.State = oldLead.State
		customer.City = oldLead.City
		customer.Address1 = oldLead.Address1
		customer.Address2 = oldLead.Address2
		customer.Zip = oldLead.Zip
		customer.Phone = oldLead.Phone
		customer.Fax = oldLead.Fax
		customer.Status = 1
		customer.Created = time.Now()
		customer.CreatedBy = info.Email
		customer.Updated = time.Now()
		customer.UpdatedBy = info.Email
		err = settingRepo.CreateCustomer(customer)
		if err != nil {
			msg := "convert to customer error"
			return nil, errors.New(msg)
		}
	} else {
		var vendor setting.Vendor
		vendor.VendorID = "ven-" + xid.New().String()
		vendor.OrganizationID = info.OrganizationID
		vendor.Name = oldLead.Company
		vendor.ContactSalutation = oldLead.Salutation
		vendor.ContactFirstName = oldLead.FirstName
		vendor.ContactLastName = oldLead.LastName
		vendor.ContactEmail = oldLead.Email
		vendor.ContactPhone = oldLead.Phone
		vendor.Country = oldLead.Country
		vendor.State = oldLead.State
		vendor.City = oldLead.City
		vendor.Address1 = oldLead.Address1
		vendor.Address2 = oldLead.Address2
		vendor.Zip = oldLead.Zip
		vendor.Phone = oldLead.Mobile
		vendor.Fax = oldLead.Fax
		vendor.Status = 1
		vendor.Created = time.Now()
		vendor.CreatedBy = info.Email
		vendor.Updated = time.Now()
		vendor.UpdatedBy = info.Email
		err = settingRepo.CreateVendor(vendor)
		if err != nil {
			msg := "convert to vendor error"
			return nil, errors.New(msg)
		}
	}
	err = repo.UpdateLeadStatus(LeadID, 2, info.Email)
	if err != nil {
		msg := "update lead error: " + err.Error()
		return nil, errors.New(msg)
	}
	err = repo.UpdateLeadConverted(LeadID, info.Type, info.Email)
	if err != nil {
		msg := "update lead error: " + err.Error()
		return nil, errors.New(msg)
	}
	var newEvent common.NewHistoryCreated
	newEvent.HistoryType = "crm"
	newEvent.HistoryTime = time.Now().Format("2006-01-02 15:04:05")
	newEvent.HistoryBy = info.User
	newEvent.ReferenceID = LeadID
	newEvent.Description = "Lead Convert to " + info.Type
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
	return &LeadID, err
}
