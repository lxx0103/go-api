package history

import "time"

type History struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	HistoryID      string    `db:"history_id" json:"history_id"`
	HistoryType    string    `db:"history_type" json:"history_type"`
	HistoryTime    string    `db:"history_date" json:"history_date"`
	HistoryBy      string    `db:"history_by" json:"history_by"`
	Description    string    `db:"description" json:"description"`
	ReferenceID    string    `db:"reference_id" json:"reference_id"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
