package crm

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type crmQuery struct {
	conn *sqlx.DB
}

func NewCrmQuery(connection *sqlx.DB) *crmQuery {
	return &crmQuery{
		conn: connection,
	}
}

func (r *crmQuery) GetLeadByID(organizationID, id string) (*LeadResponse, error) {
	var crm LeadResponse
	err := r.conn.Get(&crm, `
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
	status,
	converted_to,
	notes
	FROM c_leads
	WHERE organization_id = ? AND lead_id = ? AND status > 0
	`, organizationID, id)
	return &crm, err
}

func (r *crmQuery) GetLeadCount(filter LeadFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Company; v != "" {
		where, args = append(where, "company like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM c_leads
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *crmQuery) GetLeadList(filter LeadFilter) (*[]LeadResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Company; v != "" {
		where, args = append(where, "company like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var crms []LeadResponse
	err := r.conn.Select(&crms, `
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
		status,
		converted_to,
		notes
		FROM c_leads
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &crms, err
}
