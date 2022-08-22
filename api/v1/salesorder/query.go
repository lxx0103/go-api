package salesorder

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type salesorderQuery struct {
	conn *sqlx.DB
}

func NewSalesorderQuery(connection *sqlx.DB) *salesorderQuery {
	return &salesorderQuery{
		conn: connection,
	}
}

func (r *salesorderQuery) GetSalesorderByID(organizationID, id string) (*SalesorderResponse, error) {
	var salesorder SalesorderResponse
	err := r.conn.Get(&salesorder, `
	SELECT 
	s.salesorder_id,
	s.organization_id,
	s.salesorder_number, 
	s.salesorder_date, 
	s.expected_shipment_date, 
	s.customer_id,
	c.name as customer_name, 
	s.item_count,
	s.sub_total,
	s.tax_total,
	s.discount_type,
	s.discount_value,
	s.shipping_fee,
	s.total,
	s.notes,
	s.invoice_status,
	s.picking_status,
	s.packing_status,
	s.shipping_status,
	s.status
	FROM s_salesorders s
	LEFT JOIN s_customers c
	ON s.customer_id = c.customer_id
	WHERE s.organization_id = ? AND s.salesorder_id = ? AND s.status > 0
	`, organizationID, id)
	return &salesorder, err
}

func (r *salesorderQuery) GetSalesorderCount(filter SalesorderFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.SalesorderNumber; v != "" {
		where, args = append(where, "salesorder_number like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_salesorders
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *salesorderQuery) GetSalesorderList(filter SalesorderFilter) (*[]SalesorderResponse, error) {
	where, args := []string{"s.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "s.organization_id = ?"), append(args, v)
	}
	if v := filter.SalesorderNumber; v != "" {
		where, args = append(where, "s.salesorder_number like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var salesorders []SalesorderResponse
	err := r.conn.Select(&salesorders, `
		SELECT 
		s.salesorder_id,
		s.organization_id,
		s.salesorder_number, 
		s.salesorder_date, 
		s.expected_shipment_date, 
		s.customer_id,
		c.name as customer_name, 
		s.item_count,
		s.sub_total,
		s.tax_total,
		s.discount_type,
		s.discount_value,
		s.shipping_fee,
		s.total,
		s.notes,
		s.invoice_status,
		s.picking_status,
		s.packing_status,
		s.shipping_status,
		s.status
		FROM s_salesorders s
		LEFT JOIN s_customers c
		ON s.customer_id = c.customer_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &salesorders, err
}

func (r *salesorderQuery) GetSalesorderItemList(salesorderID string) (*[]SalesorderItemResponse, error) {
	var salesorders []SalesorderItemResponse
	err := r.conn.Select(&salesorders, `
		SELECT
		s.organization_id,
		s.salesorder_item_id,
		s.salesorder_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.rate,
		s.tax_id,
		s.tax_value,
		s.tax_amount,
		s.amount,
		s.quantity_invoiced,
		s.quantity_picked,
		s.quantity_packed,
		s.quantity_shipped,
		s.status
		FROM s_salesorder_items s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.salesorder_id = ? AND s.status > 0 
	`, salesorderID)
	return &salesorders, err
}
