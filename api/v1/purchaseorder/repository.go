package purchaseorder

import (
	"database/sql"
	"fmt"
	"time"
)

type purchaseorderRepository struct {
	tx *sql.Tx
}

func NewPurchaseorderRepository(tx *sql.Tx) *purchaseorderRepository {
	return &purchaseorderRepository{tx: tx}
}

func (r *purchaseorderRepository) CheckPONumberConfict(purchaseorder_id, organizationID, PONumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM p_purchaseorders WHERE organization_id = ? AND purchaseorder_id != ? AND purchase_order_number = ? AND status > 0 ", organizationID, purchaseorder_id, PONumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r purchaseorderRepository) CreatePurchaseorderItem(info PurchaseorderItem) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_purchaseorder_items 
		(
			organization_id,
			purchaseorder_item_id,
			purchaseorder_id,
			quantity,
			rate,
			amount,
			quantity_received,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderItemID, info.PurchaseorderID, info.Quantity, info.Rate, info.Amount, info.QuantityReceived, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r purchaseorderRepository) CreatePurchaseorder(info Purchaseorder) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_purchaseorders 
		(
			organization_id,
			purchaseorder_id,
			purchaseorder_number,
			purchaseorder_date,
			expected_delivery_date,
			vendor_id,
			item_count,
			sub_total,
			discount_type,
			discount_value,
			shipping_fee,
			total,
			notes,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderID, info.PurchaseorderNumber, info.PurchaseorderDate, info.ExpectedDeliveryDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.ShippingFee, info.Total, info.Notes, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetPurchaseorderByID(purchaseorderID string) (*PurchaseorderResponse, error) {
	var res PurchaseorderResponse
	row := r.tx.QueryRow(`
		SELECT
		purchaseorder_id,
		organization_id,
		purchaseorder_number,
		purchaseorder_date,
		expected_delivery_date,
		vendor_id,
		item_count,
		sub_total,
		discount_type,
		discount_value,
		shipping_fee,
		total,
		notes,
		status
		FROM p_purchaseorders WHERE purchaseorder_id = ? AND status > 0 LIMIT 1
	`, purchaseorderID)
	err := row.Scan(&res.PurchaseorderID, &res.OrganizationID, &res.PurchaseorderNumber, &res.PurchaseorderDate, &res.ExpectedDeliveryDate, &res.VendorID, &res.ItemCount, &res.Subtotal, &res.DiscountType, &res.DiscountValue, &res.ShippingFee, &res.Total, &res.Notes, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) UpdatePurchaseorder(id string, info Purchaseorder) error {
	_, err := r.tx.Exec(`
		Update i_purchaseorders SET
		sku = ?,
		name = ?,
		unit_id = ?,
		manufacturer_id = ?,
		brand_id = ?,
		weight_unit = ?,
		weight = ?,
		dimension_unit = ?,
		length = ?,
		width = ?,
		height = ?,
		selling_price = ?,
		cost_price = ?,
		openning_stock = ?,
		openning_stock_rate = ?,
		reorder_stock = ?,
		default_vendor_id = ?,
		description = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, info.SKU, info.Name, info.UnitID, info.ManufacturerID, info.BrandID, info.WeightUnit, info.Weight, info.DimensionUnit, info.Length, info.Width, info.Height, info.SellingPrice, info.CostPrice, info.OpenningStock, info.OpenningStockRate, info.ReorderStock, info.DefaultVendorID, info.Description, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *purchaseorderRepository) DeletePurchaseorder(id, byUser string) error {
	fmt.Println(id)
	_, err := r.tx.Exec(`
		Update i_purchaseorders SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, time.Now(), byUser, id)
	return err
}

//Barcode

func (r *purchaseorderRepository) CheckBarcodeConfict(barcodeID, organizationID, code string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM i_barcodes WHERE organization_id = ? AND barcode_id != ? AND code = ? AND status > 0 ", organizationID, barcodeID, code)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r purchaseorderRepository) CreateBarcode(info Barcode) error {
	_, err := r.tx.Exec(`
		INSERT INTO i_barcodes 
		(
			organization_id,
			barcode_id,
			code,
			purchaseorder_id,
			quantity,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BarcodeID, info.Code, info.PurchaseorderID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetBarcodeByID(barcodeID string) (*BarcodeResponse, error) {
	var res BarcodeResponse
	row := r.tx.QueryRow(`
		SELECT 
		b.barcode_id, 
		b.organization_id,
		b.purchaseorder_id,
		i.name as purchaseorder_name, 
		b.code, 
		i.sku, 
		u.name as unit, 
		b.quantity,
		b.status
		FROM i_barcodes b
		LEFT JOIN i_purchaseorders i
		ON b.purchaseorder_id = i.purchaseorder_id
		LEFT JOIN s_units u
		ON i.unit_id = u.unit_id
		WHERE b.barcode_id = ? LIMIT 1
	`, barcodeID)
	err := row.Scan(&res.BarcodeID, &res.OrganizationID, &res.PurchaseorderID, &res.PurchaseorderName, &res.Code, &res.SKU, &res.Unit, &res.Quantity, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) UpdateBarcode(id string, info Barcode) error {
	_, err := r.tx.Exec(`
		Update i_barcodes SET
		code = ?,
		purchaseorder_id = ?,
		quantity = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE barcode_id = ?
	`, info.Code, info.PurchaseorderID, info.Quantity, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *purchaseorderRepository) DeleteBarcode(id, byUser string) error {
	fmt.Println(id)
	_, err := r.tx.Exec(`
		Update i_barcodes SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE barcode_id = ?
	`, time.Now(), byUser, id)
	return err
}
