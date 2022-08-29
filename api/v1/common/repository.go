package common

import (
	"database/sql"
	"time"
)

type commonRepository struct {
	tx *sql.Tx
}

func NewCommonRepository(tx *sql.Tx) *commonRepository {
	return &commonRepository{tx: tx}
}

func (r commonRepository) CreateHistory(info History) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_historys 
		(
			history_id,
			organization_id,
			history_type,
			history_time,
			history_by,
			description,
			reference_id,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.HistoryID, info.OrganizationID, info.HistoryType, info.HistoryTime, info.HistoryBy, info.Description, info.ReferenceID, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *commonRepository) GetLastNumber(filter NumberFilter) (int, error) {
	var value int
	row := r.tx.QueryRow(`
		SELECT number_value
		FROM s_order_numbers
		WHERE organization_id = ? 
		AND number_type = ?
	`, filter.OrganizationID, filter.NumberType)
	err := row.Scan(&value)
	return value, err
}

func (r commonRepository) CreateNumber(info NumberFilter) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_order_numbers 
		(
			organization_id,
			number_type,
			number_value,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.NumberType, 1, 1, time.Now(), "SYSTEM", time.Now(), "SYSTEM")
	return err
}

func (r commonRepository) UpdateNumber(info NumberFilter, value int) error {
	_, err := r.tx.Exec(`
		UPDATE s_order_numbers set 
		number_value = ?,
		updated = ?,
		updated_by = ?
		where number_type = ?
	`, value, time.Now(), "SYSTEM", info.NumberType)
	return err
}
