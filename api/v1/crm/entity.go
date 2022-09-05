package crm

import "time"

type Lead struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	LeadID         string    `db:"lead_id" json:"lead_id"`
	Source         string    `db:"source" json:"source"`
	Company        string    `db:"company" json:"company"`
	Salutation     string    `db:"salutation" json:"salutation"`
	FirstName      string    `db:"first_name" json:"first_name"`
	LastName       string    `db:"last_name" json:"last_name"`
	Email          string    `db:"email" json:"email"`
	Phone          string    `db:"phone" json:"phone"`
	Mobile         string    `db:"mobile" json:"mobile"`
	Fax            string    `db:"fax" json:"fax"`
	Country        string    `db:"country" json:"country"`
	State          string    `db:"state" json:"state"`
	City           string    `db:"city" json:"city"`
	Address1       string    `db:"address1" json:"address1"`
	Address2       string    `db:"address2" json:"address2"`
	Zip            string    `db:"zip" json:"zip"`
	Notes          string    `db:"notes" json:"notes"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
