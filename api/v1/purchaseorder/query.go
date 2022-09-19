package purchaseorder

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type purchaseorderQuery struct {
	conn *sqlx.DB
}

func NewPurchaseorderQuery(connection *sqlx.DB) *purchaseorderQuery {
	return &purchaseorderQuery{
		conn: connection,
	}
}

func (r *purchaseorderQuery) GetPurchaseorderByID(organizationID, id string) (*PurchaseorderResponse, error) {
	var purchaseorder PurchaseorderResponse
	err := r.conn.Get(&purchaseorder, `
	SELECT 
	p.purchaseorder_id,
	p.organization_id,
	p.purchaseorder_number, 
	p.purchaseorder_date, 
	p.expected_delivery_date, 
	p.vendor_id,
	v.name as vendor_name, 
	p.item_count,
	p.sub_total,
	p.tax_total,
	p.discount_type,
	p.discount_value,
	p.shipping_fee,
	p.total,
	p.notes,
	p.receive_status,
	p.billing_status,
	p.status
	FROM p_purchaseorders p
	LEFT JOIN s_vendors v
	ON p.vendor_id = v.vendor_id
	WHERE p.organization_id = ? AND p.purchaseorder_id = ? AND p.status > 0
	`, organizationID, id)
	return &purchaseorder, err
}

func (r *purchaseorderQuery) GetPurchaseorderCount(filter PurchaseorderFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PurchaseorderNumber; v != "" {
		where, args = append(where, "purchaseorder_number like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM p_purchaseorders
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *purchaseorderQuery) GetPurchaseorderList(filter PurchaseorderFilter) (*[]PurchaseorderResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.PurchaseorderNumber; v != "" {
		where, args = append(where, "p.purchaseorder_number like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var purchaseorders []PurchaseorderResponse
	err := r.conn.Select(&purchaseorders, `
		SELECT 
		p.purchaseorder_id,
		p.organization_id,
		p.purchaseorder_number, 
		p.purchaseorder_date, 
		p.expected_delivery_date, 
		p.vendor_id,
		v.name as vendor_name, 
		p.item_count,
		p.sub_total,
		p.tax_total,
		p.discount_type,
		p.discount_value,
		p.shipping_fee,
		p.total,
		p.notes,
		p.receive_status,
		p.billing_status,
		p.status
		FROM p_purchaseorders p
		LEFT JOIN s_vendors v
		ON p.vendor_id = v.vendor_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &purchaseorders, err
}

func (r *purchaseorderQuery) GetPurchaseorderItemList(purchaseorderID string) (*[]PurchaseorderItemResponse, error) {
	var purchaseorders []PurchaseorderItemResponse
	err := r.conn.Select(&purchaseorders, `
		SELECT
		p.organization_id,
		p.purchaseorder_item_id,
		p.purchaseorder_id,
		p.item_id,
		i.name as item_name,
		i.sku as sku,
		p.quantity,
		p.rate,
		p.tax_id,
		p.tax_value,
		p.tax_amount,
		p.amount,
		p.quantity_received,
		p.quantity_billed,
		p.status
		FROM p_purchaseorder_items p
		LEFT JOIN i_items i
		ON p.item_id = i.item_id
		WHERE p.purchaseorder_id = ? AND p.status > 0 
	`, purchaseorderID)
	return &purchaseorders, err
}

func (r *purchaseorderQuery) GetPurchasereceiveCount(filter PurchasereceiveFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PurchasereceiveNumber; v != "" {
		where, args = append(where, "purchasereceive_number like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM p_purchasereceives
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *purchaseorderQuery) GetPurchasereceiveList(filter PurchasereceiveFilter) (*[]PurchasereceiveResponse, error) {
	where, args := []string{"r.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "r.organization_id = ?"), append(args, v)
	}
	if v := filter.PurchasereceiveNumber; v != "" {
		where, args = append(where, "r.purchasereceive_number like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var purchasereceives []PurchasereceiveResponse
	err := r.conn.Select(&purchasereceives, `
		SELECT 
		r.purchasereceive_id,
		r.purchaseorder_id,
		p.purchaseorder_number,
		r.organization_id,
		r.purchasereceive_number, 
		r.purchasereceive_date, 
		r.notes,
		r.status
		FROM p_purchasereceives r
		LEFT JOIN p_purchaseorders p
		ON p.purchaseorder_id = r.purchaseorder_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &purchasereceives, err
}

func (r *purchaseorderQuery) GetPurchasereceiveByID(organizationID, id string) (*PurchasereceiveResponse, error) {
	var purchasereceive PurchasereceiveResponse
	err := r.conn.Get(&purchasereceive, `
	SELECT 
	r.purchasereceive_id,
	r.purchaseorder_id,
	p.purchaseorder_number,
	r.organization_id,
	r.purchasereceive_number, 
	r.purchasereceive_date, 
	r.notes,
	r.status
	FROM p_purchasereceives r
	LEFT JOIN p_purchaseorders p
	ON p.purchaseorder_id = r.purchaseorder_id
	WHERE r.organization_id = ? AND r.purchasereceive_id = ? AND r.status > 0
	`, organizationID, id)
	return &purchasereceive, err
}

func (r *purchaseorderQuery) GetPurchasereceiveItemList(purchasereceiveID string) (*[]PurchasereceiveItemResponse, error) {
	var purchasereceives []PurchasereceiveItemResponse
	err := r.conn.Select(&purchasereceives, `
		SELECT
		p.organization_id,
		p.purchasereceive_item_id,
		p.purchaseorder_item_id,
		p.purchasereceive_id,
		p.item_id,
		i.name as item_name,
		i.sku as sku,
		p.quantity,
		p.status
		FROM p_purchasereceive_items p
		LEFT JOIN i_items i
		ON p.item_id = i.item_id
		WHERE p.purchasereceive_id = ? AND p.status > 0 
	`, purchasereceiveID)
	return &purchasereceives, err
}

func (r *purchaseorderQuery) GetPurchasereceiveDetailList(purchasereceiveID string) (*[]PurchasereceiveDetailResponse, error) {
	var purchasereceives []PurchasereceiveDetailResponse
	err := r.conn.Select(&purchasereceives, `
		SELECT
		p.organization_id,
		p.purchasereceive_item_id,
		p.purchasereceive_detail_id,
		p.purchaseorder_item_id,
		p.purchasereceive_id,
		p.location_id,
		l.code as location_code,
		p.item_id,
		i.name as item_name,
		i.sku as sku,
		p.quantity,
		p.status
		FROM p_purchasereceive_details p
		LEFT JOIN i_items i
		ON p.item_id = i.item_id
		LEFT JOIN w_locations l
		ON p.location_id = l.location_id
		WHERE p.purchasereceive_id = ? AND p.status > 0 
	`, purchasereceiveID)
	return &purchasereceives, err
}

func (r *purchaseorderQuery) GetBillCount(filter BillFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.BillNumber; v != "" {
		where, args = append(where, "bill_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.PurchaseorderID; v != "" {
		where, args = append(where, "purchaseorder_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM p_bills
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *purchaseorderQuery) GetBillList(filter BillFilter) (*[]BillResponse, error) {
	where, args := []string{"i.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "i.organization_id = ?"), append(args, v)
	}
	if v := filter.BillNumber; v != "" {
		where, args = append(where, "i.bill_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.PurchaseorderID; v != "" {
		where, args = append(where, "i.purchaseorder_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var purchaseorders []BillResponse
	err := r.conn.Select(&purchaseorders, `
		SELECT 
		i.organization_id,
		i.purchaseorder_id,
		IFNULL(s.purchaseorder_number, "") as purchaseorder_number, 
		i.bill_id, 
		i.bill_number, 
		i.bill_date,
		i.due_date,
		i.vendor_id,
		IFNULL(c.name, "") as vendor_name,
		i.item_count,
		i.sub_total,
		i.discount_type,
		i.discount_value,
		i.tax_total,
		i.shipping_fee,
		i.total,
		i.notes,
		i.status
		FROM p_bills i
		LEFT JOIN p_purchaseorders s
		ON s.purchaseorder_id = i.purchaseorder_id
		LEFT JOIN s_vendors c
		ON i.vendor_id = c.vendor_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &purchaseorders, err
}

func (r *purchaseorderQuery) GetBillByID(organizationID, id string) (*BillResponse, error) {
	var purchaseorder BillResponse
	err := r.conn.Get(&purchaseorder, `
	SELECT 
	i.organization_id,
	i.purchaseorder_id,
	IFNULL(s.purchaseorder_number, "") as purchaseorder_number, 
	i.bill_id, 
	i.bill_number, 
	i.bill_date,
	i.due_date,
	i.vendor_id,
	IFNULL(c.name, "") as vendor_name,
	i.item_count,
	i.sub_total,
	i.discount_type,
	i.discount_value,
	i.tax_total,
	i.shipping_fee,
	i.total,
	i.notes,
	i.status
	FROM p_bills i
	LEFT JOIN p_purchaseorders s
	ON s.purchaseorder_id = i.purchaseorder_id
	LEFT JOIN s_vendors c
	ON i.vendor_id = c.vendor_id
	WHERE i.organization_id = ? AND i.bill_id = ? AND s.status > 0
	`, organizationID, id)
	return &purchaseorder, err
}

func (r *purchaseorderQuery) GetBillItemList(billID string) (*[]BillItemResponse, error) {
	var billItems []BillItemResponse
	err := r.conn.Select(&billItems, `
		SELECT
		s.organization_id,
		s.bill_id,
		s.purchaseorder_item_id,
		s.bill_item_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.rate,
		s.tax_id,
    	s.tax_value,
    	s.tax_amount,
    	s.amount,
		s.status
		FROM p_bill_items s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.bill_id = ? AND i.status > 0 
	`, billID)
	return &billItems, err
}

func (r *purchaseorderQuery) GeBillPaymentMade(organizationID, billID string) (float64, error) {
	var count float64
	err := r.conn.Get(&count, `
		SELECT IFNULL(SUM(amount),0) FROM p_payment_mades 
		WHERE organization_id = ? AND bill_id = ? AND status > 0 
		`, organizationID, billID)
	return count, err
}

func (r *purchaseorderQuery) GetPaymentMadeCount(filter PaymentMadeFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentMadeNumber; v != "" {
		where, args = append(where, "payment_made_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.BillID; v != "" {
		where, args = append(where, "bill_id = ?"), append(args, v)
	}
	if v := filter.PaymentMethodID; v != "" {
		where, args = append(where, "payment_method_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM p_payment_mades
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *purchaseorderQuery) GetPaymentMadeList(filter PaymentMadeFilter) (*[]PaymentMadeResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.PaymentMadeNumber; v != "" {
		where, args = append(where, "p.payment_made_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.BillID; v != "" {
		where, args = append(where, "p.bill_id = ?"), append(args, v)
	}
	if v := filter.PaymentMethodID; v != "" {
		where, args = append(where, "p.payment_method_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var paymentReceiveds []PaymentMadeResponse
	err := r.conn.Select(&paymentReceiveds, `
		SELECT 
		p.organization_id,
		p.bill_id,
		IFNULL(i.bill_number, "") as bill_number, 
		p.vendor_id,
		IFNULL(c.name, "") as vendor_name,
		p.payment_made_id, 
		p.payment_made_number, 
		p.payment_made_date,
		p.payment_method_id,
		IFNULL(pm.name, "") as payment_method_name,
		p.amount,
		p.notes,
		p.status
		FROM p_payment_mades p
		LEFT JOIN p_bills i
		ON p.bill_id = i.bill_id
		LEFT JOIN s_vendors c
		ON p.vendor_id = c.vendor_id
		LEFT JOIN s_payment_methods pm
		ON p.payment_method_id = pm.payment_method_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &paymentReceiveds, err
}
