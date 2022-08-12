package history

import (
	"database/sql"
)

type historyRepository struct {
	tx *sql.Tx
}

func NewHistoryRepository(tx *sql.Tx) *historyRepository {
	return &historyRepository{tx: tx}
}

func (r historyRepository) CreateHistory(info History) error {
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
