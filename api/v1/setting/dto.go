package setting

import "go-api/core/request"

type UnitFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	UnitType       string `form:"unit_type" binding:"required,oneof=weight length custom"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type UnitResponse struct {
	UnitID         string `db:"unit_id" json:"unit_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Name           string `db:"name" json:"name"`
	Status         int    `db:"status" json:"status"`
}

type UnitNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}

type UnitID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type ManufacturerFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type ManufacturerResponse struct {
	ManufacturerID string `db:"manufacturer_id" json:"manufacturer_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Name           string `db:"name" json:"name"`
	Status         int    `db:"status" json:"status"`
}

type ManufacturerNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}

type ManufacturerID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

//brand

type BrandFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type BrandResponse struct {
	BrandID        string `db:"brand_id" json:"brand_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Name           string `db:"name" json:"name"`
	Status         int    `db:"status" json:"status"`
}

type BrandNew struct {
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}

type BrandID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

//vendor

type VendorFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type VendorResponse struct {
	VendorID          string `db:"vendor_id" json:"vendor_id"`
	OrganizationID    string `db:"organization_id" json:"organization_id"`
	Name              string `db:"name" json:"name"`
	ContactSalutation string `db:"contact_salutation" json:"contact_salutation"`
	ContactFirstName  string `db:"contact_first_name" json:"contact_first_name"`
	ContactLastName   string `db:"contact_last_name" json:"contact_last_name"`
	ContactEmail      string `db:"contact_email" json:"contact_email"`
	ContactPhone      string `db:"contact_phone" json:"contact_phone"`
	Country           string `db:"country" json:"country"`
	State             string `db:"state" json:"state"`
	City              string `db:"city" json:"city"`
	Address1          string `db:"address1" json:"address1"`
	Address2          string `db:"address2" json:"address2"`
	Zip               string `db:"zip" json:"zip"`
	Phone             string `db:"phone" json:"phone"`
	Fax               string `db:"fax" json:"fax"`
	Status            int    `db:"status" json:"status"`
}

type VendorNew struct {
	Name              string `json:"name" binding:"required,min=1,max=64"`
	ContactSalutation string `json:"contact_salutation" binding:"omitempty,max=64"`
	ContactFirstName  string `json:"contact_first_name" binding:"omitempty,max=64"`
	ContactLastName   string `json:"contact_last_name" binding:"omitempty,max=64"`
	ContactEmail      string `json:"contact_email" binding:"omitempty,email,max=64"`
	ContactPhone      string `json:"contact_phone" binding:"omitempty,max=64"`
	Country           string `json:"country" binding:"omitempty,max=64"`
	State             string `json:"state" binding:"omitempty,max=64"`
	City              string `json:"city" binding:"omitempty,max=64"`
	Address1          string `json:"address1" binding:"omitempty,max=255"`
	Address2          string `json:"address2" binding:"omitempty,max=255"`
	Zip               string `json:"zip" binding:"omitempty,max=64"`
	Phone             string `json:"phone" binding:"omitempty,max=64"`
	Fax               string `json:"fax" binding:"omitempty,max=64"`
	Status            int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID    string `json:"organiztion_id" swaggerignore:"true"`
	User              string `json:"user" swaggerignore:"true"`
}

type VendorID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type TaxNew struct {
	Name           string  `json:"name" binding:"required,min=1,max=64"`
	TaxValue       float64 `json:"tax_value" binding:"required,max=100"`
	Status         int     `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string  `json:"organiztion_id" swaggerignore:"true"`
	User           string  `json:"user" swaggerignore:"true"`
}

type TaxID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type TaxFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type TaxResponse struct {
	TaxID          string  `db:"tax_id" json:"tax_id"`
	OrganizationID string  `db:"organization_id" json:"organization_id"`
	Name           string  `db:"name" json:"name"`
	TaxValue       float64 `db:"tax_value" json:"tax_value"`
	Status         int     `db:"status" json:"status"`
}
