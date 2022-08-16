package purchaseorder

import (
	"database/sql"
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
	row := r.tx.QueryRow("SELECT count(1) FROM p_purchaseorders WHERE organization_id = ? AND purchaseorder_id != ? AND purchaseorder_number = ? AND status > 0 ", organizationID, purchaseorder_id, PONumber)
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
			item_id,
			quantity,
			rate,
			amount,
			quantity_received,
			quantity_billed,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderItemID, info.PurchaseorderID, info.ItemID, info.Quantity, info.Rate, info.Amount, info.QuantityReceived, info.QuantityBilled, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r purchaseorderRepository) UpdatePurchaseorderItem(id string, info PurchaseorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE p_purchaseorder_items set
		quantity = ?,
		rate = ?,
		amount = ?,
		status = ?,
		updated = ?,
		updated_by =?
		WHERE purchaseorder_item_id = ?
	`, info.Quantity, info.Rate, info.Amount, info.Status, info.Updated, info.UpdatedBy, id)
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
			receive_status,
			billing_status,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderID, info.PurchaseorderNumber, info.PurchaseorderDate, info.ExpectedDeliveryDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.ShippingFee, info.Total, info.Notes, info.ReceiveStatus, info.BillingStatus, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetPurchaseorderByID(organizationID, purchaseorderID string) (*PurchaseorderResponse, error) {
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
		receive_status,
		billing_status,
		status
		FROM p_purchaseorders WHERE organization_id = ? AND purchaseorder_id = ? AND status > 0 LIMIT 1
	`, organizationID, purchaseorderID)
	err := row.Scan(&res.PurchaseorderID, &res.OrganizationID, &res.PurchaseorderNumber, &res.PurchaseorderDate, &res.ExpectedDeliveryDate, &res.VendorID, &res.ItemCount, &res.Subtotal, &res.DiscountType, &res.DiscountValue, &res.ShippingFee, &res.Total, &res.Notes, &res.ReceiveStatus, &res.BillingStatus, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) GetPurchaseorderItemByID(organizationID, purchaseorderID, itemID string) (*PurchaseorderItemResponse, error) {
	var res PurchaseorderItemResponse
	row := r.tx.QueryRow(`
		SELECT
		p.organization_id,
		p.purchaseorder_item_id,
		p.purchaseorder_id,
		p.item_id,
		i.name as item_name,
		i.sku as sku,
		p.quantity,
		p.rate,
		p.amount,
		p.quantity_received,
		p.status
		FROM p_purchaseorder_items p
		LEFT JOIN i_items i
		ON p.item_id = i.item_id
		WHERE p.organization_id = ? AND p.purchaseorder_id = ? AND p.item_id = ? AND p.status > 0 LIMIT 1
	`, organizationID, purchaseorderID, itemID)
	err := row.Scan(&res.OrganizationID, &res.PurchaseorderItemID, &res.PurchaseorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.Amount, &res.QuantityReceived, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) UpdatePurchaseorder(id string, info Purchaseorder) error {
	_, err := r.tx.Exec(`
		Update p_purchaseorders SET
		purchaseorder_number = ?,
		purchaseorder_date = ?,
		expected_delivery_date = ?,
		vendor_id = ?,
		item_count = ?,
		sub_total = ?,
		discount_type = ?,
		discount_value = ?,
		shipping_fee = ?,
		total = ?,
		notes = ?,
		receive_status = ?,
		billing_status = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, info.PurchaseorderNumber, info.PurchaseorderDate, info.ExpectedDeliveryDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.ShippingFee, info.Total, info.Notes, info.ReceiveStatus, info.BillingStatus, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *purchaseorderRepository) DeletePurchaseorder(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_purchaseorders SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		UPDATE p_purchaseorder_items SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *purchaseorderRepository) GetPurchaseorderItemByIDAll(organizationID, purchaseorderID, purchaseorderItemID string) (*PurchaseorderItemResponse, error) {
	var res PurchaseorderItemResponse
	row := r.tx.QueryRow(`
		SELECT
		p.organization_id,
		p.purchaseorder_item_id,
		p.purchaseorder_id,
		p.item_id,
		i.name as item_name,
		i.sku as sku,
		p.quantity,
		p.rate,
		p.amount,
		p.quantity_received,
		p.quantity_billed,
		p.status
		FROM p_purchaseorder_items p
		LEFT JOIN i_items i
		ON p.item_id = i.item_id
		WHERE p.purchaseorder_item_id = ? AND p.organization_id = ? AND p.purchaseorder_id = ? LIMIT 1
	`, purchaseorderItemID, organizationID, purchaseorderID)
	err := row.Scan(&res.OrganizationID, &res.PurchaseorderItemID, &res.PurchaseorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.Amount, &res.QuantityReceived, &res.QuantityBilled, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) UpdatePurchaseorderStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_purchaseorders SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}

func (r *purchaseorderRepository) UpdatePurchaseorderReceiveStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_purchaseorders SET
		receive_status = ?,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}
func (r *purchaseorderRepository) CheckPOItem(purchaseorder_id, organizationID string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM p_purchaseorder_items WHERE organization_id = ? AND purchaseorder_id = ? AND status = -1  AND (quantity_received > 0 OR quantity_billed > 0)", organizationID, purchaseorder_id)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

//receive

func (r *purchaseorderRepository) CheckReceiveNumberConfict(purchasereceive_id, organizationID, receiveNumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM p_purchasereceives WHERE organization_id = ? AND purchasereceive_id != ? AND purchasereceive_number = ? AND status > 0 ", organizationID, purchasereceive_id, receiveNumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r purchaseorderRepository) ReceivePurchaseorderItem(info PurchaseorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE p_purchaseorder_items set
		quantity_received = ?,
		updated = ?,
		updated_by =?
		WHERE purchaseorder_item_id = ?
	`, info.QuantityReceived, info.Updated, info.UpdatedBy, info.PurchaseorderItemID)
	return err
}

func (r purchaseorderRepository) CreatePurchasereceiveItem(info PurchasereceiveItem) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_purchasereceive_items 
		(
			organization_id,
			purchasereceive_id,
			purchaseorder_item_id,
			purchasereceive_item_id,
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
	`, info.OrganizationID, info.PurchasereceiveID, info.PurchaseorderItemID, info.PurchasereceiveItemID, info.ItemID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r purchaseorderRepository) CreatePurchasereceiveDetail(info PurchasereceiveDetail) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_purchasereceive_details 
		(
			organization_id,
			purchasereceive_detail_id,
			purchasereceive_item_id,
			purchaseorder_item_id,
			purchasereceive_id,
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
	`, info.OrganizationID, info.PurchasereceiveDetailID, info.PurchasereceiveItemID, info.PurchaseorderItemID, info.PurchasereceiveID, info.LocationID, info.ItemID, info.Quantity, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetPurchaseorderReceivedCount(organizationID, purchaseorderID string) (float64, error) {
	var sum float64
	row := r.tx.QueryRow("SELECT SUM(quantity_received) FROM p_purchaseorder_items WHERE organization_id = ? AND purchaseorder_id = ? AND status > 0", organizationID, purchaseorderID)
	err := row.Scan(&sum)
	return sum, err
}
func (r purchaseorderRepository) CreatePurchasereceive(info Purchasereceive) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_purchasereceives 
		(
			organization_id,
			purchaseorder_id,
			purchasereceive_id,
			purchasereceive_number,
			purchasereceive_date,
			notes,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderID, info.PurchasereceiveID, info.PurchasereceiveNumber, info.PurchasereceiveDate, info.Notes, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}
