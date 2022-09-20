package report

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type reportQuery struct {
	conn *sqlx.DB
}

func NewReportQuery(connection *sqlx.DB) *reportQuery {
	return &reportQuery{
		conn: connection,
	}
}

//sales
func (r *reportQuery) GetSalesReport(filter SalesReportFilter) (*[]InvoiceReportResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.CustomerID; v != "" {
		where, args = append(where, "customer_id = ?"), append(args, v)
	}
	if v := filter.DateFrom; v != "" {
		where, args = append(where, "invoice_date >= ?"), append(args, v)
	}
	if v := filter.DateTo; v != "" {
		where, args = append(where, "invoice_date <= ?"), append(args, v)
	}
	var res []InvoiceReportResponse
	err := r.conn.Select(&res, `
		SELECT invoice_number, invoice_date, customer_id, status, (total - tax_total) as total, tax_total
		FROM s_invoices
		WHERE `+strings.Join(where, " AND "), args...)
	return &res, err
}

//purchase
func (r *reportQuery) GetPurchaseReport(filter PurchaseReportFilter) (*[]BillReportResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.VendorID; v != "" {
		where, args = append(where, "vendor_id = ?"), append(args, v)
	}
	if v := filter.DateFrom; v != "" {
		where, args = append(where, "bill_date >= ?"), append(args, v)
	}
	if v := filter.DateTo; v != "" {
		where, args = append(where, "bill_date <= ?"), append(args, v)
	}
	var res []BillReportResponse
	err := r.conn.Select(&res, `
		SELECT bill_number, bill_date, vendor_id, status, (total - tax_total) as total, tax_total
		FROM p_bills
		WHERE `+strings.Join(where, " AND "), args...)
	return &res, err
}

//adjustment
func (r *reportQuery) GetAdjustmentReport(filter AdjustmentReportFilter) (*[]AdjustmentReportResponse, error) {
	where, args := []string{"a.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "a.organization_id = ?"), append(args, v)
	}
	if v := filter.ItemID; v != "" {
		where, args = append(where, "a.item_id = ?"), append(args, v)
	}
	if v := filter.DateFrom; v != "" {
		where, args = append(where, "a.adjustment_date >= ?"), append(args, v)
	}
	if v := filter.DateTo; v != "" {
		where, args = append(where, "a.adjustment_date <= ?"), append(args, v)
	}
	var res []AdjustmentReportResponse
	err := r.conn.Select(&res, `
		SELECT 
		l.code as location_code,
		i.name as item_name,
		i.sku,
		a.quantity, 
		a.original_quantity,
		a.new_quantity,
		a.adjustment_date,
		IFNULL(r.name, "") as adjustment_reason_name
		FROM i_adjustments a
		LEFT JOIN w_locations l
		ON l.location_id = a.location_id
		LEFT join i_items i
		ON a.item_id = i.item_id
		LEFT JOIN s_adjustment_reasons r
		ON a.adjustment_reason_id = r.adjustment_reason_id
		WHERE `+strings.Join(where, " AND "), args...)
	return &res, err
}

//item
func (r *reportQuery) GetItemReport(filter ItemReportFilter) (*[]ItemReportResponse, error) {
	where, args := []string{"i.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "i.organization_id = ?"), append(args, v)
	}
	if v := filter.ItemID; v != "" {
		where, args = append(where, "i.item_id = ?"), append(args, v)
	}
	var res []ItemReportResponse
	err := r.conn.Select(&res, `
		SELECT 
		i.name as item_name,
		i.sku,
		u.name as unit,
		i.cost_price,
		i.selling_price,
		i.stock_on_hand
		FROM i_items i
		LEFT JOIN s_units u
		ON i.unit_id = u.unit_id
		WHERE `+strings.Join(where, " AND "), args...)
	return &res, err
}
