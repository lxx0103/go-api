package crm

import (
	"go-api/core/request"
)

type LeadNew struct {
	Source         string `json:"source" binding:"required"`
	Company        string `json:"company" binding:"required"`
	Salutation     string `json:"salutation" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	LeadEmail      string `json:"lead_email" binding:"omitempty,email"`
	Phone          string `json:"phone" binding:"required"`
	Mobile         string `json:"mobile" binding:"omitempty"`
	Fax            string `json:"fax" binding:"omitempty"`
	Country        string `json:"country" binding:"omitempty"`
	State          string `json:"state" binding:"omitempty"`
	City           string `json:"city" binding:"omitempty"`
	Address1       string `json:"address1" binding:"omitempty"`
	Address2       string `json:"address2" binding:"omitempty"`
	Zip            string `json:"zip" binding:"omitempty"`
	Status         int    `json:"status" binding:"omitempty,min=1"`
	Notes          string `json:"notes" binding:"omitempty"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
	Email          string `json:"email" swaggerignore:"true"`
}

type LeadFilter struct {
	Company        string `form:"company" binding:"omitempty,max=255,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type LeadResponse struct {
	OrganizationID string `db:"organization_id" json:"organization_id"`
	LeadID         string `db:"lead_id" json:"lead_id"`
	Source         string `db:"source" json:"source"`
	Company        string `db:"company" json:"company"`
	Salutation     string `db:"salutation" json:"salutation"`
	FirstName      string `db:"first_name" json:"first_name"`
	LastName       string `db:"last_name" json:"last_name"`
	Email          string `db:"email" json:"email"`
	Phone          string `db:"phone" json:"phone"`
	Mobile         string `db:"mobile" json:"mobile"`
	Fax            string `db:"fax" json:"fax"`
	Country        string `db:"country" json:"country"`
	State          string `db:"state" json:"state"`
	City           string `db:"city" json:"city"`
	Address1       string `db:"address1" json:"address1"`
	Address2       string `db:"address2" json:"address2"`
	Zip            string `db:"zip" json:"zip"`
	Notes          string `db:"notes" json:"notes"`
	Status         int    `db:"status" json:"status"`
}

type LeadID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type LeadConvertNew struct {
	Type           string `json:"type" binding:"required,oneof=customer vendor"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
	Email          string `json:"email" swaggerignore:"true"`
}
