package item

import (
	"database/sql"
	"fmt"
	"time"
)

type itemRepository struct {
	tx *sql.Tx
}

func NewItemRepository(tx *sql.Tx) *itemRepository {
	return &itemRepository{tx: tx}
}

func (r *itemRepository) CheckSKUConfict(item_id, organizationID, SKU string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM i_items WHERE organization_id = ? AND item_id != ? AND sku = ? AND status > 0 ", organizationID, item_id, SKU)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r itemRepository) CreateItem(info Item) error {
	_, err := r.tx.Exec(`
		INSERT INTO i_items 
		(
			organization_id,
			item_id,
			sku,
			name,
			unit_id,
			manufacturer_id,
			brand_id,
			weight_unit,
			weight,
			dimension_unit,
			length,
			width,
			height,
			selling_price,
			cost_price,
			openning_stock,
			openning_stock_rate,
			reorder_stock,
			default_vendor_id,
			description,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.ItemID, info.SKU, info.Name, info.UnitID, info.ManufacturerID, info.BrandID, info.WeightUnit, info.Weight, info.DimensionUnit, info.Length, info.Width, info.Height, info.SellingPrice, info.CostPrice, info.OpenningStock, info.OpenningStockRate, info.ReorderStock, info.DefaultVendorID, info.Description, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *itemRepository) GetItemByID(itemID string) (*ItemResponse, error) {
	var res ItemResponse
	row := r.tx.QueryRow(`
		SELECT
		item_id,
		organization_id,
		sku,
		name,
		unit_id,
		manufacturer_id,
		brand_id,
		weight_unit,
		weight,
		dimension_unit,
		length,
		width,
		height,
		selling_price,
		cost_price,
		openning_stock,
		openning_stock_rate,
		reorder_stock,
		default_vendor_id,
		description,
		status
		FROM i_items WHERE item_id = ? AND status > 0 LIMIT 1
	`, itemID)
	err := row.Scan(&res.ItemID, &res.OrganizationID, &res.SKU, &res.Name, &res.UnitID, &res.ManufacturerID, &res.BrandID, &res.WeightUnit, &res.Weight, &res.DimensionUnit, &res.Length, &res.Width, &res.Height, &res.SellingPrice, &res.CostPrice, &res.OpenningStock, &res.OpenningStockRate, &res.ReorderStock, &res.DefaultVendorID, &res.Description, &res.Status)
	return &res, err
}

func (r *itemRepository) UpdateItem(id string, info Item) error {
	_, err := r.tx.Exec(`
		Update i_items SET
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
		WHERE item_id = ?
	`, info.SKU, info.Name, info.UnitID, info.ManufacturerID, info.BrandID, info.WeightUnit, info.Weight, info.DimensionUnit, info.Length, info.Width, info.Height, info.SellingPrice, info.CostPrice, info.OpenningStock, info.OpenningStockRate, info.ReorderStock, info.DefaultVendorID, info.Description, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *itemRepository) DeleteItem(id, byUser string) error {
	fmt.Println(id)
	_, err := r.tx.Exec(`
		Update i_items SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE item_id = ?
	`, time.Now(), byUser, id)
	return err
}

//Barcode

func (r *itemRepository) CheckBarcodeConfict(barcodeID, organizationID, code string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM i_barcodes WHERE organization_id = ? AND barcode_id != ? AND code = ? AND status > 0 ", organizationID, barcodeID, code)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r itemRepository) CreateBarcode(info Barcode) error {
	_, err := r.tx.Exec(`
		INSERT INTO i_barcodes 
		(
			organization_id,
			barcode_id,
			code,
			item_id,
			quantity,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BarcodeID, info.Code, info.ItemID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *itemRepository) GetBarcodeByID(barcodeID string) (*BarcodeResponse, error) {
	var res BarcodeResponse
	row := r.tx.QueryRow(`
		SELECT 
		b.barcode_id, 
		b.organization_id,
		b.item_id,
		i.name as item_name, 
		b.code, 
		i.sku, 
		u.name as unit, 
		b.quantity,
		b.status
		FROM i_barcodes b
		LEFT JOIN i_items i
		ON b.item_id = i.item_id
		LEFT JOIN s_units u
		ON i.unit_id = u.unit_id
		WHERE b.barcode_id = ? LIMIT 1
	`, barcodeID)
	err := row.Scan(&res.BarcodeID, &res.OrganizationID, &res.ItemID, &res.ItemName, &res.Code, &res.SKU, &res.Unit, &res.Quantity, &res.Status)
	return &res, err
}

func (r *itemRepository) UpdateBarcode(id string, info Barcode) error {
	_, err := r.tx.Exec(`
		Update i_barcodes SET
		code = ?,
		item_id = ?,
		quantity = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE barcode_id = ?
	`, info.Code, info.ItemID, info.Quantity, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *itemRepository) DeleteBarcode(id, byUser string) error {
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
