package purchaseorder

import "time"

type Purchaseorder struct {
	ID                   int64     `db:"id" json:"id"`
	OrganizationID       string    `db:"organization_id" json:"organization_id"`
	PurchaseorderID      string    `db:"purchaseorder_id" json:"purchaseorder_id"`
	PurchaseorderNumber  string    `db:"purchaseorder_number" json:"purchaseorder_number"`
	PurchaseorderDate    string    `db:"purchaseorder_date" json:"purchaseorder_date"`
	ExpectedDeliveryDate string    `db:"expected_delivery_date" json:"expected_delivery_date"`
	VendorID             string    `db:"vendor_id" json:"vendor_id"`
	ItemCount            int       `db:"item_count" json:"item_count"`
	Subtotal             float64   `db:"subtotal" json:"subtotal"`
	DiscountType         int       `db:"discount_type" json:"discount_type"`
	DiscountValue        float64   `db:"discount_value" json:"discount_value"`
	ShippingFee          float64   `db:"shipping_fee" json:"shipping_fee"`
	Total                float64   `db:"total" json:"total"`
	Notes                string    `db:"notes" json:"notes"`
	Status               int       `db:"status" json:"status"`
	Created              time.Time `db:"created" json:"created"`
	CreatedBy            string    `db:"created_by" json:"created_by"`
	Updated              time.Time `db:"updated" json:"updated"`
	UpdatedBy            string    `db:"updated_by" json:"updated_by"`
}

type PurchaseorderItem struct {
	ID                  int64     `db:"id" json:"id"`
	OrganizationID      string    `db:"organization_id" json:"organization_id"`
	PurchaseorderID     string    `db:"purchaseorder_id" json:"purchaseorder_id"`
	PurchaseorderItemID string    `db:"purchaseorder_item_id" json:"purchaseorder_item_id"`
	ItemID              string    `db:"item_id" json:"item_id"`
	Quantity            int       `db:"quantity" json:"quantity"`
	Rate                float64   `db:"rate" json:"rate"`
	Amount              float64   `db:"amount" json:"amount"`
	QuantityReceived    int       `db:"quantity_received" json:"quantity_received"`
	Status              int       `db:"status" json:"status"`
	Created             time.Time `db:"created" json:"created"`
	CreatedBy           string    `db:"created_by" json:"created_by"`
	Updated             time.Time `db:"updated" json:"updated"`
	UpdatedBy           string    `db:"updated_by" json:"updated_by"`
}

type Barcode struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  string    `db:"organization_id" json:"organization_id"`
	BarcodeID       string    `db:"barcode_id" json:"barcode_id"`
	Code            string    `db:"code" json:"code"`
	PurchaseorderID string    `db:"purchaseorder_id" json:"purchaseorder_id"`
	Quantity        int       `db:"quantity" json:"quantity"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}
