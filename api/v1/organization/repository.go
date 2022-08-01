package organization

import (
	"database/sql"
	"time"
)

type organizationRepository struct {
	tx *sql.Tx
}

func NewOrganizationRepository(transaction *sql.Tx) OrganizationRepository {
	return &organizationRepository{
		tx: transaction,
	}
}

type OrganizationRepository interface {
	//Organization Management
	CheckConfict(owner string) (bool, error)
	CreateOrganization(info OrganizationNew) (int64, error)
	// UpdateOrganization(id int64, info OrganizationNew) (int64, error)
	// GetOrganizationByID(id int64) (*Organization, error)
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

func (r *organizationRepository) CreateOrganization(info OrganizationNew) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO s_organizations
		(
			name,
			owner,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, 2, ?, "SIGNUP", ?, "SIGNUP")
	`, info.Name, info.Email, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (r *organizationRepository) UpdateOrganization(id int64, info OrganizationNew) (int64, error) {
// 	result, err := r.tx.Exec(`
// 		Update organizations SET
// 		name = ?,
// 		status = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Status, time.Now(), info.User, id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return affected, nil
// }

// func (r *organizationRepository) GetOrganizationByID(id int64) (*Organization, error) {
// 	var res Organization
// 	row := r.tx.QueryRow(`SELECT id, name, status, created, created_by, updated, updated_by FROM organizations WHERE id = ? LIMIT 1`, id)
// 	err := row.Scan(&res.ID, &res.Name, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }
