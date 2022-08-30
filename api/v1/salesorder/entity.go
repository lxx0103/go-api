package salesorder

import "time"

type Salesorder struct {
	ID                   int64     `db:"id" json:"id"`
	OrganizationID       string    `db:"organization_id" json:"organization_id"`
	SalesorderID         string    `db:"salesorder_id" json:"salesorder_id"`
	SalesorderNumber     string    `db:"salesorder_number" json:"salesorder_number"`
	SalesorderDate       string    `db:"salesorder_date" json:"salesorder_date"`
	ExpectedShipmentDate string    `db:"expected_shipment_date" json:"expected_shipment_date"`
	CustomerID           string    `db:"customer_id" json:"customer_id"`
	ItemCount            int       `db:"item_count" json:"item_count"`
	Subtotal             float64   `db:"subtotal" json:"subtotal"`
	DiscountType         int       `db:"discount_type" json:"discount_type"`
	DiscountValue        float64   `db:"discount_value" json:"discount_value"`
	TaxTotal             float64   `db:"tax_total" json:"tax_total"`
	ShippingFee          float64   `db:"shipping_fee" json:"shipping_fee"`
	Total                float64   `db:"total" json:"total"`
	Notes                string    `db:"notes" json:"notes"`
	InvoiceStatus        int       `db:"invoice_status" json:"invoice_status"`
	PickingStatus        int       `db:"picking_status" json:"picking_status"`
	PackingStatus        int       `db:"packing_status" json:"packing_status"`
	ShippingStatus       int       `db:"shipping_status" json:"shipping_status"`
	Status               int       `db:"status" json:"status"`
	Created              time.Time `db:"created" json:"created"`
	CreatedBy            string    `db:"created_by" json:"created_by"`
	Updated              time.Time `db:"updated" json:"updated"`
	UpdatedBy            string    `db:"updated_by" json:"updated_by"`
}

type SalesorderItem struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   string    `db:"organization_id" json:"organization_id"`
	SalesorderID     string    `db:"salesorder_id" json:"salesorder_id"`
	SalesorderItemID string    `db:"salesorder_item_id" json:"salesorder_item_id"`
	ItemID           string    `db:"item_id" json:"item_id"`
	Quantity         int       `db:"quantity" json:"quantity"`
	Rate             float64   `db:"rate" json:"rate"`
	TaxID            string    `db:"tax_id" json:"tax_id"`
	TaxValue         float64   `db:"tax_value" json:"tax_value"`
	TaxAmount        float64   `db:"tax_amount" json:"tax_amount"`
	Amount           float64   `db:"amount" json:"amount"`
	QuantityInvoiced int       `db:"quantity_invoiced" json:"quantity_invoiced"`
	QuantityPicked   int       `db:"quantity_picked" json:"quantity_picked"`
	QuantityPacked   int       `db:"quantity_packed" json:"quantity_packed"`
	QuantityShipped  int       `db:"quantity_shipped" json:"quantity_shipped"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}

type Pickingorder struct {
	ID                 int64     `db:"id" json:"id"`
	OrganizationID     string    `db:"organization_id" json:"organization_id"`
	SalesorderID       string    `db:"salesorder_id" json:"salesorder_id"`
	PickingorderID     string    `db:"pickingorder_id" json:"pickingorder_id"`
	PickingorderNumber string    `db:"pickingorder_number" json:"pickingorder_number"`
	PickingorderDate   string    `db:"pickingorder_date" json:"pickingorder_date"`
	Notes              string    `db:"notes" json:"notes"`
	Status             int       `db:"status" json:"status"`
	Created            time.Time `db:"created" json:"created"`
	CreatedBy          string    `db:"created_by" json:"created_by"`
	Updated            time.Time `db:"updated" json:"updated"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
}

type PickingorderItem struct {
	ID                 int64     `db:"id" json:"id"`
	OrganizationID     string    `db:"organization_id" json:"organization_id"`
	PickingorderID     string    `db:"pickingorder_id" json:"pickingorder_id"`
	SalesorderItemID   string    `db:"salesorder_item_id" json:"salesorder_item_id"`
	PickingorderItemID string    `db:"pickingorder_item_id" json:"pickingorder_item_id"`
	ItemID             string    `db:"item_id" json:"item_id"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Status             int       `db:"status" json:"status"`
	Created            time.Time `db:"created" json:"created"`
	CreatedBy          string    `db:"created_by" json:"created_by"`
	Updated            time.Time `db:"updated" json:"updated"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
}

type PickingorderLog struct {
	ID                 int64     `db:"id" json:"id"`
	PickingorderLogID  string    `db:"pickingorder_log_id" json:"pickingorder_log_id"`
	OrganizationID     string    `db:"organization_id" json:"organization_id"`
	PickingorderID     string    `db:"pickingorder_id" json:"pickingorder_id"`
	SalesorderItemID   string    `db:"salesorder_item_id" json:"salesorder_item_id"`
	PickingorderItemID string    `db:"pickingorder_item_id" json:"pickingorder_item_id"`
	LocationID         string    `db:"location_id" json:"location_id"`
	ItemID             string    `db:"item_id" json:"item_id"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Status             int       `db:"status" json:"status"`
	Created            time.Time `db:"created" json:"created"`
	CreatedBy          string    `db:"created_by" json:"created_by"`
	Updated            time.Time `db:"updated" json:"updated"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
}

type PickingorderDetail struct {
	ID                   int64     `db:"id" json:"id"`
	PickingorderDetailID string    `db:"pickingorder_detail_id" json:"pickingorder_detail_id"`
	OrganizationID       string    `db:"organization_id" json:"organization_id"`
	PickingorderID       string    `db:"pickingorder_id" json:"pickingorder_id"`
	LocationID           string    `db:"location_id" json:"location_id"`
	ItemID               string    `db:"item_id" json:"item_id"`
	Quantity             int       `db:"quantity" json:"quantity"`
	QuantityPicked       int       `db:"quantity_picked" json:"quantity_picked"`
	Status               int       `db:"status" json:"status"`
	Created              time.Time `db:"created" json:"created"`
	CreatedBy            string    `db:"created_by" json:"created_by"`
	Updated              time.Time `db:"updated" json:"updated"`
	UpdatedBy            string    `db:"updated_by" json:"updated_by"`
}

type Package struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	SalesorderID   string    `db:"salesorder_id" json:"salesorder_id"`
	PackageID      string    `db:"package_id" json:"package_id"`
	PackageNumber  string    `db:"package_number" json:"package_number"`
	PackageDate    string    `db:"package_date" json:"package_date"`
	Notes          string    `db:"notes" json:"notes"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type PackageItem struct {
	ID               int64     `db:"id" json:"id"`
	OrganizationID   string    `db:"organization_id" json:"organization_id"`
	PackageID        string    `db:"package_id" json:"package_id"`
	SalesorderItemID string    `db:"salesorder_item_id" json:"salesorder_item_id"`
	PackageItemID    string    `db:"package_item_id" json:"package_item_id"`
	ItemID           string    `db:"item_id" json:"item_id"`
	Quantity         int       `db:"quantity" json:"quantity"`
	Status           int       `db:"status" json:"status"`
	Created          time.Time `db:"created" json:"created"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	Updated          time.Time `db:"updated" json:"updated"`
	UpdatedBy        string    `db:"updated_by" json:"updated_by"`
}
