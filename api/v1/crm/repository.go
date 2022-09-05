package crm

import (
	"database/sql"
	"time"
)

type crmRepository struct {
	tx *sql.Tx
}

func NewCrmRepository(tx *sql.Tx) *crmRepository {
	return &crmRepository{tx: tx}
}

func (r crmRepository) CreateLead(info Lead) error {
	_, err := r.tx.Exec(`
		INSERT INTO c_leads 
		(
			organization_id,
			lead_id,
			source,
			company,
			salutation,
			first_name,
			last_name,
			email,
			phone,
			mobile,
			fax,
			country,
			state,
			city,
			address1,
			address2,
			zip,
			status,
			notes,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.LeadID, info.Source, info.Company, info.Salutation, info.FirstName, info.LastName, info.Email, info.Phone, info.Mobile, info.Fax, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Status, info.Notes, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *crmRepository) GetLeadByID(organizationID, leadID string) (*LeadResponse, error) {
	var res LeadResponse
	row := r.tx.QueryRow(`
		SELECT
		organization_id, 
		lead_id, 
		source, 
		company, 
		salutation, 
		first_name, 
		last_name, 
		email, 
		phone, 
		mobile, 
		fax, 
		country, 
		state, 
		city, 
		address1, 
		address2, 
		zip, 
		notes,
		status
		FROM c_leads WHERE organization_id = ? AND lead_id = ? AND status > 0 LIMIT 1
	`, organizationID, leadID)
	err := row.Scan(&res.OrganizationID, &res.LeadID, &res.Source, &res.Company, &res.Salutation, &res.FirstName, &res.LastName, &res.Email, &res.Phone, &res.Mobile, &res.Fax, &res.Country, &res.State, &res.City, &res.Address1, &res.Address2, &res.Zip, &res.Notes, &res.Status)
	return &res, err
}

func (r *crmRepository) UpdateLead(id string, info Lead) error {
	_, err := r.tx.Exec(`
		Update c_leads SET
		source = ?, 
		company = ?, 
		salutation = ?, 
		first_name = ?, 
		last_name = ?, 
		email = ?, 
		phone = ?, 
		mobile = ?, 
		fax = ?, 
		country = ?, 
		state = ?, 
		city = ?, 
		address1 = ?, 
		address2 = ?, 
		zip = ?, 
		status = ?,
		notes = ?,
		updated = ?, 
		updated_by = ?
		WHERE lead_id = ?
	`, info.Source, info.Company, info.Salutation, info.FirstName, info.LastName, info.Email, info.Phone, info.Mobile, info.Fax, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Status, info.Notes, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *crmRepository) DeleteLead(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update c_leads SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE lead_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *crmRepository) UpdateLeadStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update c_leads SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE lead_id = ?
	`, status, time.Now(), byUser, id)
	return err
}
