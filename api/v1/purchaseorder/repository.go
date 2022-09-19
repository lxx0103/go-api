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
			tax_id,
			tax_value,
			tax_amount,
			quantity_received,
			quantity_billed,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderItemID, info.PurchaseorderID, info.ItemID, info.Quantity, info.Rate, info.Amount, info.TaxID, info.TaxValue, info.TaxAmount, info.QuantityReceived, info.QuantityBilled, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r purchaseorderRepository) UpdatePurchaseorderItem(id string, info PurchaseorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE p_purchaseorder_items set
		quantity = ?,
		rate = ?,
		tax_id = ?,
		tax_value = ?,
		tax_amount = ?,
		amount = ?,
		status = ?,
		updated = ?,
		updated_by =?
		WHERE purchaseorder_item_id = ?
	`, info.Quantity, info.Rate, info.TaxID, info.TaxValue, info.TaxAmount, info.Amount, info.Status, info.Updated, info.UpdatedBy, id)
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
			tax_total,
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
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.PurchaseorderID, info.PurchaseorderNumber, info.PurchaseorderDate, info.ExpectedDeliveryDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.TaxTotal, info.ShippingFee, info.Total, info.Notes, info.ReceiveStatus, info.BillingStatus, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
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
		tax_total,
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
	err := row.Scan(&res.PurchaseorderID, &res.OrganizationID, &res.PurchaseorderNumber, &res.PurchaseorderDate, &res.ExpectedDeliveryDate, &res.VendorID, &res.ItemCount, &res.TaxTotal, &res.Subtotal, &res.DiscountType, &res.DiscountValue, &res.ShippingFee, &res.Total, &res.Notes, &res.ReceiveStatus, &res.BillingStatus, &res.Status)
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
		WHERE p.organization_id = ? AND p.purchaseorder_id = ? AND p.item_id = ? AND p.status > 0 LIMIT 1
	`, organizationID, purchaseorderID, itemID)
	err := row.Scan(&res.OrganizationID, &res.PurchaseorderItemID, &res.PurchaseorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.QuantityReceived, &res.QuantityBilled, &res.Status)
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
		tax_total = ?,
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
	`, info.PurchaseorderNumber, info.PurchaseorderDate, info.ExpectedDeliveryDate, info.VendorID, info.ItemCount, info.Subtotal, info.TaxTotal, info.DiscountType, info.DiscountValue, info.ShippingFee, info.Total, info.Notes, info.ReceiveStatus, info.BillingStatus, info.Status, info.Updated, info.UpdatedBy, id)
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
		WHERE p.purchaseorder_item_id = ? AND p.organization_id = ? AND p.purchaseorder_id = ? LIMIT 1
	`, purchaseorderItemID, organizationID, purchaseorderID)
	err := row.Scan(&res.OrganizationID, &res.PurchaseorderItemID, &res.PurchaseorderID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.QuantityReceived, &res.QuantityBilled, &res.Status)
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

func (r *purchaseorderRepository) GetPurchasereceiveDetailList(organizationID, purchasereceiveID string) (*[]PurchasereceiveDetailResponse, error) {
	var purchasereceiveDetails []PurchasereceiveDetailResponse
	rows, err := r.tx.Query(`
		SELECT
		s.organization_id,
		s.purchasereceive_detail_id,
		s.purchasereceive_id,
		s.purchaseorder_item_id,
		s.purchasereceive_item_id,
		s.location_id,
		l.code as location_code,
		s.item_id,
		i.name as item_name,
		i.sku,
		s.quantity,
		s.status
		FROM p_purchasereceive_details s
		LEFT JOIN i_items i
		ON s.item_id = i.item_id
		LEFT JOIN w_locations l
		ON s.location_id = l.location_id
		WHERE s.organization_id = ? AND s.purchasereceive_id = ? AND s.status > 0
	`, organizationID, purchasereceiveID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res PurchasereceiveDetailResponse
		err = rows.Scan(&res.OrganizationID, &res.PurchasereceiveDetailID, &res.PurchasereceiveID, &res.PurchaseorderItemID, &res.PurchasereceiveItemID, &res.LocationID, &res.LocationCode, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Status)
		purchasereceiveDetails = append(purchasereceiveDetails, res)
		if err != nil {
			return nil, err
		}
	}
	return &purchasereceiveDetails, err
}

func (r *purchaseorderRepository) GetPurchasereceiveByID(organizationID, purchasereceiveID string) (*PurchasereceiveResponse, error) {
	var res PurchasereceiveResponse
	row := r.tx.QueryRow(`
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
		WHERE r.organization_id = ? AND r.purchasereceive_id = ? AND r.status > 0 LIMIT 1
	`, organizationID, purchasereceiveID)
	err := row.Scan(&res.PurchasereceiveID, &res.PurchaseorderID, &res.PurchaseorderNumber, &res.OrganizationID, &res.PurchasereceiveNumber, &res.PurchasereceiveDate, &res.Notes, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) DeletePurchasereceive(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_purchasereceives SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchasereceive_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		UPDATE p_purchasereceive_items SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchasereceive_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		UPDATE p_purchasereceive_details SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE purchasereceive_id = ?
	`, time.Now(), byUser, id)
	return err
}

// bill

func (r *purchaseorderRepository) CheckBillNumberConfict(billID, organizationID, billNumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM p_bills WHERE organization_id = ? AND bill_id != ? AND bill_number = ? AND status > 0 ", organizationID, billID, billNumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r purchaseorderRepository) BillPurchaseorderItem(info PurchaseorderItem) error {
	_, err := r.tx.Exec(`
		UPDATE p_purchaseorder_items set
		quantity_billed = ?,
		updated = ?,
		updated_by =?
		WHERE purchaseorder_item_id = ?
	`, info.QuantityBilled, info.Updated, info.UpdatedBy, info.PurchaseorderItemID)
	return err
}

func (r purchaseorderRepository) CreateBillItem(info BillItem) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_bill_items 
		(
			organization_id,
			bill_id,
			bill_item_id,
			purchaseorder_item_id,
			item_id,
			quantity,
			rate,
			tax_id,
			tax_value,
			tax_amount,
			amount,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BillID, info.BillItemID, info.PurchaseorderItemID, info.ItemID, info.Quantity, info.Rate, info.TaxID, info.TaxValue, info.TaxAmount, info.Amount, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r purchaseorderRepository) CreateBill(info Bill) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_bills 
		(
			organization_id,
			bill_id,
			purchaseorder_id,
			bill_number,
			bill_date,
			due_date,
			vendor_id,
			item_count,
			sub_total,
			discount_type,
			discount_value,
			tax_total,
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
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BillID, info.PurchaseorderID, info.BillNumber, info.BillDate, info.DueDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.TaxTotal, info.ShippingFee, info.Total, info.Notes, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetPurchaseorderBilledCount(organizationID, purchaseorderID string) (float64, error) {
	var sum float64
	row := r.tx.QueryRow("SELECT SUM(quantity_billed) FROM p_purchaseorder_items WHERE organization_id = ? AND purchaseorder_id = ? AND status > 0", organizationID, purchaseorderID)
	err := row.Scan(&sum)
	return sum, err
}

func (r *purchaseorderRepository) UpdatePurchaseorderBillStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_purchaseorders SET
		billing_status = ?,
		updated = ?,
		updated_by = ?
		WHERE purchaseorder_id = ?
	`, status, time.Now(), byUser, id)
	return err
}

func (r *purchaseorderRepository) GetBillByID(organizationID, id string) (*BillResponse, error) {
	var res BillResponse
	row := r.tx.QueryRow(`
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
		WHERE i.organization_id = ? AND i.bill_id = ? AND s.status > 0  LIMIT 1
	`, organizationID, id)
	err := row.Scan(&res.OrganizationID, &res.PurchaseorderID, &res.PurchaseorderNumber, &res.BillID, &res.BillNumber, &res.BillDate, &res.DueDate, &res.VendorID, &res.VendorName, &res.ItemCount, &res.Subtotal, &res.DiscountType, &res.DiscountValue, &res.TaxTotal, &res.ShippingFee, &res.Total, &res.Notes, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) DeleteBill(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_bills SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE bill_id = ?
	`, time.Now(), byUser, id)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(`
		UPDATE p_bill_items SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE bill_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *purchaseorderRepository) GetBillItemByIDAll(organizationID, billID, billItemID string) (*BillItemResponse, error) {
	var res BillItemResponse
	row := r.tx.QueryRow(`
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
		WHERE s.bill_item_id = ? AND s.organization_id = ? AND s.bill_id = ? LIMIT 1
	`, billItemID, organizationID, billID)
	err := row.Scan(&res.OrganizationID, &res.BillID, &res.PurchaseorderItemID, &res.BillItemID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.Status)
	return &res, err
}

func (r *purchaseorderRepository) GetBillItemList(organizationID, billID string) (*[]BillItemResponse, error) {
	var billItems []BillItemResponse
	rows, err := r.tx.Query(`
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
		WHERE s.organization_id = ? AND s.bill_id = ? AND s.status > 0 
	`, organizationID, billID)
	for rows.Next() {
		var res BillItemResponse
		err = rows.Scan(&res.OrganizationID, &res.BillID, &res.PurchaseorderItemID, &res.BillItemID, &res.ItemID, &res.ItemName, &res.SKU, &res.Quantity, &res.Rate, &res.TaxID, &res.TaxValue, &res.TaxAmount, &res.Amount, &res.Status)
		billItems = append(billItems, res)
		if err != nil {
			return nil, err
		}
	}
	return &billItems, err
}

func (r purchaseorderRepository) UpdateBill(id string, info Bill) error {
	_, err := r.tx.Exec(`
		UPDATE p_bills set 
		bill_number = ?,
		bill_date = ?,
		due_date = ?,
		vendor_id = ?,
		item_count = ?,
		sub_total = ?,
		discount_type = ?,
		discount_value = ?,
		tax_total = ?,
		shipping_fee = ?,
		total = ?,
		notes = ?,
		status = ?,
		updated = ?,
		updated_by =?
		WHERE bill_id = ?
	`, info.BillNumber, info.BillDate, info.DueDate, info.VendorID, info.ItemCount, info.Subtotal, info.DiscountType, info.DiscountValue, info.TaxTotal, info.ShippingFee, info.Total, info.Notes, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

// payment
func (r *purchaseorderRepository) CheckPaymentMadeNumberConfict(paymentMadeID, organizationID, paymentMadeNumber string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM p_payment_mades WHERE organization_id = ? AND payment_made_id != ? AND payment_made_number = ? AND status > 0 ", organizationID, paymentMadeID, paymentMadeNumber)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r purchaseorderRepository) CreatePaymentMade(info PaymentMade) error {
	_, err := r.tx.Exec(`
		INSERT INTO p_payment_mades 
		(
			organization_id,
			bill_id,
			vendor_id,
			payment_made_id,
			payment_made_number,
			payment_made_date,
			payment_method_id,
			amount,
			notes,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BillID, info.VendorID, info.PaymentMadeID, info.PaymentMadeNumber, info.PaymentMadeDate, info.PaymentMethodID, info.Amount, info.Notes, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *purchaseorderRepository) GetBillPaidCount(organizationID, billID string) (float64, error) {
	var sum float64
	row := r.tx.QueryRow("SELECT IFNULL(SUM(amount), 0) FROM p_payment_mades WHERE organization_id = ? AND bill_id = ? AND status > 0", organizationID, billID)
	err := row.Scan(&sum)
	return sum, err
}

func (r *purchaseorderRepository) UpdateBillStatus(id string, status int, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_bills SET
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE bill_id = ?
	`, status, time.Now(), byUser, id)
	return err
}

func (r *purchaseorderRepository) GetPaymentMadeByID(organizationID, id string) (*PaymentMadeResponse, error) {
	var res PaymentMadeResponse
	row := r.tx.QueryRow(`		
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
		WHERE p.organization_id = ? AND p.payment_made_id = ? AND p.status > 0  LIMIT 1
	`, organizationID, id)
	err := row.Scan(&res.OrganizationID, &res.BillID, &res.BillNumber, &res.VendorID, &res.VendorName, &res.PaymentMadeID, &res.PaymentMadeNumber, &res.PaymentMadeDate, &res.PaymentMethodID, &res.PaymentMethodName, &res.Amount, &res.Notes, &res.Status)
	return &res, err
}

func (r purchaseorderRepository) UpdatePaymentMade(id string, info PaymentMade) error {
	_, err := r.tx.Exec(`
		UPDATE p_payment_mades set 
		payment_made_number = ?,
		payment_made_date = ?,
		payment_method_id = ?,
		amount = ?,
		notes = ?,
		status = ?,
		updated = ?,
		updated_by =?
		WHERE payment_made_id = ?
	`, info.PaymentMadeNumber, info.PaymentMadeDate, info.PaymentMethodID, info.Amount, info.Notes, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *purchaseorderRepository) DeletePaymentMade(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update p_payment_mades SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE payment_made_id = ?
	`, time.Now(), byUser, id)
	return err
}
