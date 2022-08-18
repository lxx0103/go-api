package setting

import "time"

type Unit struct {
	ID             int64     `db:"id" json:"id"`
	UnitID         string    `db:"unit_id" json:"unit_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type Manufacturer struct {
	ID             int64     `db:"id" json:"id"`
	ManufacturerID string    `db:"manufacturer_id" json:"manufacturer_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type Brand struct {
	ID             int64     `db:"id" json:"id"`
	BrandID        string    `db:"brand_id" json:"brand_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type Vendor struct {
	ID                int64     `db:"id" json:"id"`
	VendorID          string    `db:"vendor_id" json:"vendor_id"`
	OrganizationID    string    `db:"organization_id" json:"organization_id"`
	Name              string    `db:"name" json:"name"`
	ContactSalutation string    `db:"contact_salutation" json:"contact_salutation"`
	ContactFirstName  string    `db:"contact_first_name" json:"contact_first_name"`
	ContactLastName   string    `db:"contact_last_name" json:"contact_last_name"`
	ContactEmail      string    `db:"contact_email" json:"contact_email"`
	ContactPhone      string    `db:"contact_phone" json:"contact_phone"`
	Country           string    `db:"country" json:"country"`
	State             string    `db:"state" json:"state"`
	City              string    `db:"city" json:"city"`
	Address1          string    `db:"address1" json:"address1"`
	Address2          string    `db:"address2" json:"address2"`
	Zip               string    `db:"zip" json:"zip"`
	Phone             string    `db:"phone" json:"phone"`
	Fax               string    `db:"fax" json:"fax"`
	Status            int       `db:"status" json:"status"`
	Created           time.Time `db:"created" json:"created"`
	CreatedBy         string    `db:"created_by" json:"created_by"`
	Updated           time.Time `db:"updated" json:"updated"`
	UpdatedBy         string    `db:"updated_by" json:"updated_by"`
}

type Tax struct {
	ID             int64     `db:"id" json:"id"`
	TaxID          string    `db:"tax_id" json:"tax_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	TaxValue       float64   `db:"tax_value" json:"tax_value"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
