package organization

import (
	"database/sql"
)

type organizationRepository struct {
	tx *sql.Tx
}

func NewOrganizationRepository(transaction *sql.Tx) *organizationRepository {
	return &organizationRepository{
		tx: transaction,
	}
}

func (r *organizationRepository) CheckConfict(owner string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_organizations WHERE owner = ? AND status > 0 ", owner)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *organizationRepository) CreateOrganization(info Organization) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO s_organizations
		(
			organization_id,
			name,
			owner,
			owner_email,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Owner, info.OwnerEmail, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
