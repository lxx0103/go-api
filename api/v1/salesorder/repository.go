package salesorder

import (
	"database/sql"
	"time"
)

type salesorderRepository struct {
	tx *sql.Tx
}

func NewSalesorderRepository(tx *sql.Tx) *salesorderRepository {
	return &salesorderRepository{tx: tx}
}

func (r *salesorderRepository) CheckSONumberConfict(salesorder_id, organizationID, SONumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_salesorders WHERE organization_id = ? AND salesorder_id != ? AND salesorder_number = ? AND status > 0 ", organizationID, salesorder_id, SONumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r salesorderRepository) CreateSalesorderItem(info SalesorderItem) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_salesorder_items 
		(
			organization_id,
			salesorder_item_id,
			salesorder_id,
			item_id,
			quantity,
			rate,
			amount,
			tax_id,
			tax_value,
			tax_amount,
			quantity_invoiced,
			quantity_picked,
			quantity_packed,
			quantity_shipped,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.SalesorderItemID, info.SalesorderID, info.ItemID, info.Quantity, info.Rate, info.Amount, info.TaxID, info.TaxValue, info.TaxAmount, info.QuantityInvoiced, info.QuantityPicked, info.QuantityPacked, info.QuantityShipped, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r salesorderRepository) UpdateSalesorderItem(id string, info SalesorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE s_salesorder_items set
		quantity = ?,
		rate = ?,
		tax_id = ?,
		tax_value = ?,
		tax_amount = ?,
		amount = ?,
		status = ?,
		updated = ?,
		updated_by =?
		WHERE salesorder_item_id = ?
	`, info.Quantity, info.Rate, info.TaxID, info.TaxValue, info.TaxAmount, info.Amount, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r salesorderRepository) CreateSalesorder(info Salesorder) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_salesorders 
		(
			organization_id,
			salesorder_id,
			salesorder_number,
			salesorder_date,
			expected_shipment_date,
			customer_id,
			item_count,
			sub_total,
			discount_type,
			discount_value,
			tax_total,
			shipping_fee,
			total,
			notes,
			invoice_status,
			picking_status,
			packing_status,
			shipping_status,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.SalesorderID, info.SalesorderNumber, info.SalesorderDate, info.ExpectedShipmentDate, info.CustomerID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.TaxTotal, info.ShippingFee, info.Total, info.Notes, info.InvoiceStatus, info.PickingStatus, info.PackingStatus, info.ShippingStatus, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *salesorderRepository) GetSalesorderByID(organizationID, salesorderID string) (*SalesorderResponse, error) {
	var res SalesorderResponse
	row := r.tx.QueryRow(`
		SELECT
		salesorder_id,
		organization_id,
		salesorder_number,
		salesorder_date,
		expected_shipment_date,
		customer_id,
		item_count,
		tax_total,
		sub_total,
		discount_type,
		discount_value,
		shipping_fee,
		total,
		notes,
		invoice_status,
		picking_status,
		packing_status,
		shipping_status,
		status
		FROM s_salesorders WHERE organization_id = ? AND salesorder_id = ? AND status > 0 LIMIT 1
	`, organizationID, salesorderID)
	err := row.Scan(&res.SalesorderID, &res.OrganizationID, &res.SalesorderNumber, &res.SalesorderDate, &res.ExpectedShipmentDate, &res.CustomerID, &res.ItemCount, &res.TaxTotal, &res.Subtotal, &res.DiscountType, &res.DiscountValue, &res.ShippingFee, &res.Total, &res.Notes, &res.InvoiceStatus, &res.PickingStatus, &res.PackingStatus, &res.ShippingStatus, &res.Status)
	return &res, err
}

func (r *salesorderRepository) GetSalesorderItemByID(organizationID, salesorderID, itemID string) (*SalesorderItemResponse, error) {
	var res SalesorderItemResponse
	row := r.tx.QueryRow(`
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
		WHERE s.organization_id = ? AND s.salesorder_id = ? AND s.item_id = ? AND s.status > 0 LIMIT 1
	`, organizationID, salesorderID, itemID)
	err := row.Scan(&res.OrganizationID, &res.SalesorderItemID, &res.SalesorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.QuantityInvoiced, &res.QuantityPicked, &res.QuantityPacked, &res.QuantityShipped, &res.Status)
	return &res, err
}

func (r *salesorderRepository) GetSalesorderItemList(organizationID, salesorderID string) (*[]SalesorderItemResponse, error) {
	var salesorders []SalesorderItemResponse
	rows, err := r.tx.Query(`
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
		WHERE s.organization_id = ? AND s.salesorder_id = ? AND s.status > 0 
	`, organizationID, salesorderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res SalesorderItemResponse
		err = rows.Scan(&res.OrganizationID, &res.SalesorderItemID, &res.SalesorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.QuantityInvoiced, &res.QuantityPicked, &res.QuantityPacked, &res.QuantityShipped, &res.Status)
		salesorders = append(salesorders, res)
		if err != nil {
			return nil, err
		}
	}
	return &salesorders, err
}

func (r *salesorderRepository) UpdateSalesorder(id string, info Salesorder) error {
	_, err := r.tx.Exec(`
		Update s_salesorders SET
		salesorder_number = ?,
		salesorder_date = ?,
		expected_shipment_date = ?,
		customer_id = ?,
		item_count = ?,
		sub_total = ?,
		tax_total = ?,
		discount_type = ?,
		discount_value = ?,
		shipping_fee = ?,
		total = ?,
		notes = ?,
		invoice_status = ?,
		picking_status = ?,
		packing_status = ?,
		shipping_status = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE salesorder_id = ?
	`, info.SalesorderNumber, info.SalesorderDate, info.ExpectedShipmentDate, info.CustomerID, info.ItemCount, info.Subtotal, info.TaxTotal, info.DiscountType, info.DiscountValue, info.ShippingFee, info.Total, info.Notes, info.InvoiceStatus, info.PickingStatus, info.PackingStatus, info.ShippingStatus, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *salesorderRepository) DeleteSalesorder(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_salesorders SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE salesorder_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		UPDATE s_salesorder_items SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE salesorder_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *salesorderRepository) GetSalesorderItemByIDAll(organizationID, salesorderID, salesorderItemID string) (*SalesorderItemResponse, error) {
	var res SalesorderItemResponse
	row := r.tx.QueryRow(`
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
		WHERE s.salesorder_item_id = ? AND s.organization_id = ? AND s.salesorder_id = ? LIMIT 1
	`, salesorderItemID, organizationID, salesorderID)
	err := row.Scan(&res.OrganizationID, &res.SalesorderItemID, &res.SalesorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.QuantityInvoiced, &res.QuantityPicked, &res.QuantityPacked, &res.QuantityShipped, &res.Status)
	return &res, err
}

func (r *salesorderRepository) UpdateSalesorderStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_salesorders SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE salesorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}

func (r *salesorderRepository) UpdateSalesorderPickingStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_salesorders SET
		picking_status = ?,
		updated = ?,
		updated_by = ?
		WHERE salesorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}
func (r *salesorderRepository) CheckSOItem(salesorder_id, organizationID string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_salesorder_items WHERE organization_id = ? AND salesorder_id = ? AND status = -1  AND (quantity_invoiced > 0 OR quantity_picked > 0 OR quantity_packed > 0 OR quantity_shipped > 0)", organizationID, salesorder_id)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

//receive

func (r *salesorderRepository) CheckPickingorderNumberConfict(pickingorderID, organizationID, pickingorderNumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_pickingorders WHERE organization_id = ? AND pickingorder_id != ? AND pickingorder_number = ? AND status > 0 ", organizationID, pickingorderID, pickingorderNumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r salesorderRepository) PickSalesorderItem(info SalesorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE s_salesorder_items set
		quantity_picked = ?,
		updated = ?,
		updated_by =?
		WHERE salesorder_item_id = ?
	`, info.QuantityPicked, info.Updated, info.UpdatedBy, info.SalesorderItemID)
	return err
}

func (r salesorderRepository) CreatePickingorderItem(info PickingorderItem) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_pickingorder_items 
		(
			organization_id,
			pickingorder_id,
			salesorder_item_id,
			pickingorder_item_id,
			item_id,
			quantity,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PickingorderID, info.SalesorderItemID, info.PickingorderItemID, info.ItemID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r salesorderRepository) CreatePickingorderLog(info PickingorderLog) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_pickingorder_logs 
		(
			organization_id,
			pickingorder_log_id,
			pickingorder_item_id,
			salesorder_item_id,
			pickingorder_id,
			location_id,
			item_id,
			quantity,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PickingorderLogID, info.PickingorderItemID, info.SalesorderItemID, info.PickingorderID, info.LocationID, info.ItemID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}
func (r salesorderRepository) CreatePickingorderDetail(info PickingorderDetail) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_pickingorder_details 
		(
			organization_id,
			pickingorder_detail_id,
			pickingorder_id,
			location_id,
			item_id,
			quantity,
			quantity_picked,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PickingorderDetailID, info.PickingorderID, info.LocationID, info.ItemID, info.Quantity, info.QuantityPicked, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *salesorderRepository) GetSalesorderPickedCount(organizationID, salesorderID string) (float64, error) {
	var sum float64
	row := r.tx.QueryRow("SELECT SUM(quantity_picked) FROM s_salesorder_items WHERE organization_id = ? AND salesorder_id = ? AND status > 0", organizationID, salesorderID)
	err := row.Scan(&sum)
	return sum, err
}
func (r salesorderRepository) CreatePickingorder(info Pickingorder) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_pickingorders 
		(
			organization_id,
			salesorder_id,
			pickingorder_id,
			pickingorder_number,
			pickingorder_date,
			notes,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.SalesorderID, info.PickingorderID, info.PickingorderNumber, info.PickingorderDate, info.Notes, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r salesorderRepository) GetPickingorderLogSum(pickingorderID string) (*[]PickingorderLogResponse, error) {

	var pickingorderLogs []PickingorderLogResponse
	rows, err := r.tx.Query(`
	SELECT
	s.organization_id,
	s.pickingorder_id,
	s.location_id,
	IFNULL(l.code, "") as location_code,
	s.item_id,
	i.name as item_name,
	i.sku as sku,
	sum(s.quantity) as quantity
	FROM s_pickingorder_logs s
	LEFT JOIN i_items i
	ON s.item_id = i.item_id
	LEFT JOIN w_locations l
	ON s.location_id = l.location_id
	WHERE s.pickingorder_id = ? AND s.status > 0 
	GROUP BY s.organization_id, s.pickingorder_id,s.location_id, s.item_id 
	`, pickingorderID)
	for rows.Next() {
		var res PickingorderLogResponse
		rows.Scan(&res.OrganizationID, &res.PickingorderID, &res.LocationID, &res.LocationCode, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity)
		pickingorderLogs = append(pickingorderLogs, res)
	}
	return &pickingorderLogs, err
}

func (r *salesorderRepository) GetPickingorderDetailByLocationID(organizationID, pickingorderID, locationID string) (*PickingorderDetailResponse, error) {
	var res PickingorderDetailResponse
	row := r.tx.QueryRow(`
		SELECT
		s.organization_id,
		s.pickingorder_detail_id,
		s.pickingorder_id,
		s.location_id,
		l.code as location_code,
		s.item_id,
		i.name as item_name,
		i.sku,
		s.quantity,
		s.quantity_picked,
		s.status
		FROM s_pickingorder_details s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		LEFT JOIN w_locations l
		ON s.location_id = l.location_id
		WHERE s.organization_id = ? AND s.pickingorder_id = ? AND s.location_id = ? AND s.status > 0 LIMIT 1
	`, organizationID, pickingorderID, locationID)
	err := row.Scan(&res.OrganizationID, &res.PickingorderDetailID, &res.PickingorderID, &res.LocationID, &res.LocationCode, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.QuantityPicked, &res.Status)
	return &res, err
}

func (r *salesorderRepository) UpdatePickingorderPicked(id string, quantity int, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_pickingorder_details SET
		quantity_picked = quantity_picked + ?,
		updated = ?,
		updated_by = ?
		WHERE pickingorder_detail_id = ?
	`, quantity, time.Now(), byUser, id)
	return err
}

func (r *salesorderRepository) CheckPickingorderPicked(organizationID, pickingorderID string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_pickingorder_details WHERE organization_id = ? AND pickingorder_id = ? AND quantity > quantity_picked AND status > 0 ", organizationID, pickingorderID)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *salesorderRepository) UpdatePickingorderStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_pickingorders SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE pickingorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}
