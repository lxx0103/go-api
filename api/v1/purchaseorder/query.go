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
