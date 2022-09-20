package report

import (
	"database/sql"
)

type reportRepository struct {
	tx *sql.Tx
}

func NewReportRepository(tx *sql.Tx) *reportRepository {
	return &reportRepository{tx: tx}
}
