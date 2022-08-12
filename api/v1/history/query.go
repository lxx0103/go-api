package history

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type historyQuery struct {
	conn *sqlx.DB
}

func NewHistoryQuery(connection *sqlx.DB) *historyQuery {
	return &historyQuery{
		conn: connection,
	}
}

func (r *historyQuery) GetHistoryByID(organizationID, id string) (*HistoryResponse, error) {
	var history HistoryResponse
	err := r.conn.Get(&history, `
	SELECT 
	p.history_id,
	p.organization_id,
	p.history_number, 
	p.history_date, 
	p.expected_delivery_date, 
	v.name as vendor_name, 
	p.item_count,
	p.sub_total,
	p.discount_type,
	p.discount_value,
	p.shipping_fee,
	p.total,
	p.notes,
	p.receive_status,
	p.billing_status,
	p.status
	FROM p_historys p
	LEFT JOIN s_vendors v
	ON p.vendor_id = v.vendor_id
	WHERE p.organization_id = ? AND p.history_id = ? AND p.status > 0
	`, organizationID, id)
	return &history, err
}

func (r *historyQuery) GetHistoryCount(filter HistoryFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.HistoryType; v != "" {
		where, args = append(where, "history_type = ?"), append(args, v)
	}
	if v := filter.ReferenceID; v != "" {
		where, args = append(where, "reference_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_historys
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *historyQuery) GetHistoryList(filter HistoryFilter) (*[]HistoryResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.HistoryType; v != "" {
		where, args = append(where, "history_type = ?"), append(args, v)
	}
	if v := filter.ReferenceID; v != "" {
		where, args = append(where, "reference_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var historys []HistoryResponse
	err := r.conn.Select(&historys, `
		SELECT 
		history_id,
		organization_id,
		history_type, 
		history_time, 
		history_by, 
		reference_id, 
		description
		FROM s_historys
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &historys, err
}
