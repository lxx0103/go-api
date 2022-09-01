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

func (r *salesorderQuery) GetPickingorderCount(filter PickingorderFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PickingorderNumber; v != "" {
		where, args = append(where, "pickingorder_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.SalesorderID; v != "" {
		where, args = append(where, "salesorder_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_pickingorders
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *salesorderQuery) GetPickingorderList(filter PickingorderFilter) (*[]PickingorderResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.PickingorderNumber; v != "" {
		where, args = append(where, "p.pickingorder_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.SalesorderID; v != "" {
		where, args = append(where, "p.salesorder_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var salesorders []PickingorderResponse
	err := r.conn.Select(&salesorders, `
		SELECT 
		p.organization_id,
		p.salesorder_id,
		IFNULL(s.salesorder_number, "") as salesorder_number, 
		p.pickingorder_id, 
		p.pickingorder_number, 
		p.pickingorder_date,
		p.notes,
		p.status
		FROM s_pickingorders p
		LEFT JOIN s_salesorders s
		ON s.salesorder_id = p.salesorder_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &salesorders, err
}

func (r *salesorderQuery) GetPickingorderByID(organizationID, id string) (*PickingorderResponse, error) {
	var pickingorder PickingorderResponse
	err := r.conn.Get(&pickingorder, `
	SELECT 
	p.organization_id,
	p.salesorder_id,
	IFNULL(s.salesorder_number, "") as salesorder_number, 
	p.pickingorder_id, 
	p.pickingorder_number, 
	p.pickingorder_date,
	p.notes,
	p.status
	FROM s_pickingorders p
	LEFT JOIN s_salesorders s
	ON s.salesorder_id = p.salesorder_id
	WHERE p.organization_id = ? AND p.pickingorder_id = ? AND p.status > 0
	`, organizationID, id)
	return &pickingorder, err
}

func (r *salesorderQuery) GetPickingorderItemList(salesorderID string) (*[]PickingorderItemResponse, error) {
	var pickingorderItems []PickingorderItemResponse
	err := r.conn.Select(&pickingorderItems, `
		SELECT
		s.organization_id,
		s.pickingorder_id,
		s.salesorder_item_id,
		s.pickingorder_item_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.status
		FROM s_pickingorder_items s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.pickingorder_id = ? AND s.status > 0 
	`, salesorderID)
	return &pickingorderItems, err
}

func (r *salesorderQuery) GetPickingorderDetailList(salesorderID string) (*[]PickingorderDetailResponse, error) {
	var pickingorderDetails []PickingorderDetailResponse
	err := r.conn.Select(&pickingorderDetails, `
	SELECT
	s.organization_id,
	s.pickingorder_id,
	s.pickingorder_detail_id,
	s.location_id,
	IFNULL(l.code, "") as location_code,
	s.item_id,
	i.name as item_name,
	i.sku as sku,
	s.quantity,
	s.quantity_picked,
	s.status
	FROM s_pickingorder_details s
	LEFT JOIN i_items i
	ON s.item_id = i.item_id
	LEFT JOIN w_locations l
	ON s.location_id = l.location_id
	WHERE s.pickingorder_id = ? AND s.status > 0 
	`, salesorderID)
	return &pickingorderDetails, err
}

func (r *salesorderQuery) GetPackageCount(filter PackageFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.PackageNumber; v != "" {
		where, args = append(where, "package_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.SalesorderID; v != "" {
		where, args = append(where, "salesorder_id = ?"), append(args, v)
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_packages
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *salesorderQuery) GetPackageList(filter PackageFilter) (*[]PackageResponse, error) {
	where, args := []string{"p.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "p.organization_id = ?"), append(args, v)
	}
	if v := filter.PackageNumber; v != "" {
		where, args = append(where, "p.package_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.SalesorderID; v != "" {
		where, args = append(where, "p.salesorder_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var salesorders []PackageResponse
	err := r.conn.Select(&salesorders, `
		SELECT 
		p.organization_id,
		p.salesorder_id,
		IFNULL(s.salesorder_number, "") as salesorder_number, 
		p.package_id, 
		p.package_number, 
		p.package_date,
		p.notes,
		p.status
		FROM s_packages p
		LEFT JOIN s_salesorders s
		ON s.salesorder_id = p.salesorder_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &salesorders, err
}

func (r *salesorderQuery) GetPackageByID(organizationID, id string) (*PackageResponse, error) {
	var packageRes PackageResponse
	err := r.conn.Get(&packageRes, `
	SELECT 
	p.organization_id,
	p.salesorder_id,
	IFNULL(s.salesorder_number, "") as salesorder_number, 
	p.package_id, 
	p.package_number, 
	p.package_date,
	p.notes,
	p.status
	FROM s_packages p
	LEFT JOIN s_salesorders s
	ON s.salesorder_id = p.salesorder_id
	WHERE p.organization_id = ? AND p.package_id = ? AND p.status > 0
	`, organizationID, id)
	return &packageRes, err
}

func (r *salesorderQuery) GetPackageItemList(salesorderID string) (*[]PackageItemResponse, error) {
	var packageItems []PackageItemResponse
	err := r.conn.Select(&packageItems, `
		SELECT
		s.organization_id,
		s.package_id,
		s.salesorder_item_id,
		s.package_item_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.status
		FROM s_package_items s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.package_id = ? AND s.status > 0 
	`, salesorderID)
	return &packageItems, err
}

func (r *salesorderQuery) GetShippingorderCount(filter ShippingorderFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.ShippingorderNumber; v != "" {
		where, args = append(where, "shippingorder_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.PackageID; v != "" {
		where, args = append(where, "package_id like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_shippingorders
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}
func (r *salesorderQuery) GetShippingorderList(filter ShippingorderFilter) (*[]ShippingorderResponse, error) {
	where, args := []string{"s.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "s.organization_id = ?"), append(args, v)
	}
	if v := filter.ShippingorderNumber; v != "" {
		where, args = append(where, "s.shippingorder_number like ?"), append(args, "%"+v+"%")
	}
	if v := filter.PackageID; v != "" {
		where, args = append(where, "s.package_id like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var shippingorders []ShippingorderResponse
	err := r.conn.Select(&shippingorders, `
		SELECT 
		s.organization_id,
		s.shippingorder_id,
		s.package_id,
		IFNULL(p.package_number, "") as package_number,
		s.shippingorder_number, 
		s.shippingorder_date, 
		s.carrier_id,
		IFNULL(c.name, "") as carrier_name,
		s.notes,
		s.status
		FROM s_shippingorders s
		LEFT JOIN s_packages p
		ON s.package_id = p.package_id
		LEFT JOIN s_carriers c
		ON s.carrier_id = c.carrier_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &shippingorders, err
}

func (r *salesorderQuery) GetShippingorderByID(organizationID, id string) (*ShippingorderResponse, error) {
	var shippingorder ShippingorderResponse
	err := r.conn.Get(&shippingorder, `
	SELECT 
	s.organization_id,
	s.shippingorder_id,
	s.package_id,
	IFNULL(p.package_number, "") as package_number,
	s.shippingorder_number, 
	s.shippingorder_date, 
	s.carrier_id,
	IFNULL(c.name, "") as carrier_name,
	s.notes,
	s.status
	FROM s_shippingorders s
	LEFT JOIN s_packages p
	ON s.package_id = p.package_id
	LEFT JOIN s_carriers c
	ON s.carrier_id = c.carrier_id
	WHERE s.organization_id = ? AND s.shippingorder_id = ? AND s.status > 0
	`, organizationID, id)
	return &shippingorder, err
}

func (r *salesorderQuery) GetShippingorderItemList(shippingorderID string) (*[]ShippingorderItemResponse, error) {
	var shippingorderItems []ShippingorderItemResponse
	err := r.conn.Select(&shippingorderItems, `
		SELECT
		s.organization_id,
		s.shippingorder_id,
		s.shippingorder_item_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.status
		FROM s_shippingorder_items s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.shippingorder_id = ? AND s.status > 0 
	`, shippingorderID)
	return &shippingorderItems, err
}

func (r *salesorderQuery) GetShippingorderDetailList(shippingorderID string) (*[]ShippingorderDetailResponse, error) {
	var shippingorderDetails []ShippingorderDetailResponse
	err := r.conn.Select(&shippingorderDetails, `
		SELECT
		s.organization_id,
		s.shippingorder_id,
		s.shippingorder_detail_id,
		s.package_id,
		s.package_item_id,
		s.item_id,
		i.name as item_name,
		i.sku as sku,
		s.quantity,
		s.status
		FROM s_shippingorder_details s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		WHERE s.shippingorder_id = ? AND s.status > 0 
	`, shippingorderID)
	return &shippingorderDetails, err
}
