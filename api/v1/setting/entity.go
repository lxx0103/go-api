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
