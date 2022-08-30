package salesorder

import (
	"go-api/core/request"
)

type SalesorderNew struct {
	SalesorderNumber     string              `json:"salesorder_number" binding:"required,min=6,max=64"`
	SalesorderDate       string              `json:"salesorder_date" binding:"required,datetime=2006-01-02"`
	ExpectedShipmentDate string              `json:"expected_shipment_date" binding:"required,datetime=2006-01-02"`
	CustomerID           string              `json:"customer_id" binding:"required"`
	DiscountType         int                 `json:"discount_type" binding:"omitempty,oneof=1 2"`
	DiscountValue        float64             `json:"discount_value" binding:"omitempty"`
	ShippingFee          float64             `json:"shipping_fee" binding:"omitempty"`
	Notes                string              `json:"notes" binding:"omitempty"`
	Items                []SalesorderItemNew `json:"items" binding:"required"`
	OrganizationID       string              `json:"organiztion_id" swaggerignore:"true"`
	User                 string              `json:"user" swaggerignore:"true"`
	Email                string              `json:"email" swaggerignore:"true"`
}

type SalesorderItemNew struct {
	SalesorderItemID string  `json:"salesorder_item_id" binding:"omitempty"`
	ItemID           string  `json:"item_id" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required"`
	Rate             float64 `json:"rate" binding:"required"`
	TaxID            string  `json:"tax_id" binding:"omitempty"`
}

type SalesorderFilter struct {
	SalesorderNumber string `form:"salesorder_number" binding:"omitempty,max=64,min=1"`
	OrganizationID   string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type SalesorderResponse struct {
	SalesorderID         string  `db:"salesorder_id" json:"salesorder_id"`
	OrganizationID       string  `db:"organization_id" json:"organization_id"`
	SalesorderNumber     string  `db:"salesorder_number" json:"salesorder_number"`
	SalesorderDate       string  `db:"salesorder_date" json:"salesorder_date"`
	ExpectedShipmentDate string  `db:"expected_shipment_date" json:"expected_shipment_date"`
	CustomerID           string  `db:"customer_id" json:"customer_id"`
	CustomerName         string  `db:"customer_name" json:"customer_name"`
	ItemCount            float64 `db:"item_count" json:"item_count"`
	Subtotal             float64 `db:"sub_total" json:"sub_total"`
	TaxTotal             float64 `db:"tax_total" json:"tax_total"`
	DiscountType         int     `db:"discount_type" json:"discount_type"`
	DiscountValue        float64 `db:"discount_value" json:"discount_value"`
	ShippingFee          float64 `db:"shipping_fee" json:"shipping_fee"`
	Total                float64 `db:"total" json:"total"`
	Notes                string  `db:"notes" json:"notes"`
	InvoiceStatus        int     `db:"invoice_status" json:"invoice_status"`
	PickingStatus        int     `db:"picking_status" json:"picking_status"`
	PackingStatus        int     `db:"packing_status" json:"packing_status"`
	ShippingStatus       int     `db:"shipping_status" json:"shipping_status"`
	Status               int     `db:"status" json:"status"`
}

type SalesorderItemResponse struct {
	OrganizationID   string  `db:"organization_id" json:"organization_id"`
	SalesorderID     string  `db:"salesorder_id" json:"salesorder_id"`
	SalesorderItemID string  `db:"salesorder_item_id" json:"salesorder_item_id"`
	ItemID           string  `db:"item_id" json:"item_id"`
	ItemName         string  `db:"item_name" json:"item_name"`
	SKU              string  `db:"sku" json:"sku"`
	Quantity         int     `db:"quantity" json:"quantity"`
	Rate             float64 `db:"rate" json:"rate"`
	TaxID            string  `db:"tax_id" json:"tax_id"`
	TaxValue         float64 `db:"tax_value" json:"tax_value"`
	TaxAmount        float64 `db:"tax_amount" json:"tax_amount"`
	Amount           float64 `db:"amount" json:"amount"`
	QuantityInvoiced int     `db:"quantity_invoiced" json:"quantity_invoiced"`
	QuantityPicked   int     `db:"quantity_picked" json:"quantity_picked"`
	QuantityPacked   int     `db:"quantity_packed" json:"quantity_packed"`
	QuantityShipped  int     `db:"quantity_shipped" json:"quantity_shipped"`
	Status           int     `db:"status" json:"status"`
}

type SalesorderID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type PickingorderNew struct {
	PickingorderNumber string                `json:"pickingorder_number" binding:"required,min=6,max=64"`
	PickingorderDate   string                `json:"pickingorder_date" binding:"required,datetime=2006-01-02"`
	Notes              string                `json:"notes" binding:"omitempty"`
	Items              []PickingorderItemNew `json:"items" binding:"required"`
	OrganizationID     string                `json:"organiztion_id" swaggerignore:"true"`
	User               string                `json:"user" swaggerignore:"true"`
	Email              string                `json:"email" swaggerignore:"true"`
}

type PickingorderBatch struct {
	SOID               []string `json:"so_id" binding:"required,min=1"`
	PickingorderNumber string   `json:"pickingorder_number" binding:"required,min=6,max=64"`
	PickingorderDate   string   `json:"pickingorder_date" binding:"required,datetime=2006-01-02"`
	Notes              string   `json:"notes" binding:"omitempty"`
	Assigned           string   `json:"assigned" binding:"omitempty"`
	OrganizationID     string   `json:"organiztion_id" swaggerignore:"true"`
	User               string   `json:"user" swaggerignore:"true"`
	Email              string   `json:"email" swaggerignore:"true"`
}

type PickingorderItemNew struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type PickingorderResponse struct {
	OrganizationID     string `db:"organization_id" json:"organization_id"`
	SalesorderID       string `db:"salesorder_id" json:"salesorder_id"`
	SalesorderNumber   string `db:"salesorder_number" json:"salesorder_number"`
	PickingorderID     string `db:"pickingorder_id" json:"pickingorder_id"`
	PickingorderNumber string `db:"pickingorder_number" json:"pickingorder_number"`
	PickingorderDate   string `db:"pickingorder_date" json:"pickingorder_date"`
	Notes              string `db:"notes" json:"notes"`
	Status             int    `db:"status" json:"status"`
}

type PickingorderFilter struct {
	SalesorderID       string `form:"salesorder_id" binding:"omitempty,max=64"`
	PickingorderNumber string `form:"pickingorder_number" binding:"omitempty,max=64,min=1"`
	OrganizationID     string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type PickingorderID struct {
	ID string `uri:"id" binding:"required,min=1"`
}
type PickingorderItemResponse struct {
	OrganizationID     string `db:"organization_id" json:"organization_id"`
	PickingorderID     string `db:"pickingorder_id" json:"pickingorder_id"`
	SalesorderItemID   string `db:"salesorder_item_id" json:"salesorder_item_id"`
	PickingorderItemID string `db:"pickingorder_item_id" json:"pickingorder_item_id"`
	ItemID             string `db:"item_id" json:"item_id"`
	ItemName           string `db:"item_name" json:"item_name"`
	SKU                string `db:"sku" json:"sku"`
	Quantity           int    `db:"quantity" json:"quantity"`
	Status             int    `db:"status" json:"status"`
}

type PickingorderLogResponse struct {
	PickingorderLogID  string `db:"pickingorder_log_id" json:"pickingorder_log_id"`
	OrganizationID     string `db:"organization_id" json:"organization_id"`
	PickingorderID     string `db:"pickingorder_id" json:"pickingorder_id"`
	SalesorderItemID   string `db:"salesorder_item_id" json:"salesorder_item_id"`
	PickingorderItemID string `db:"pickingorder_item_id" json:"pickingorder_item_id"`
	LocationID         string `db:"location_id" json:"location_id"`
	LocationCode       string `db:"location_code" json:"location_code"`
	ItemID             string `db:"item_id" json:"item_id"`
	ItemName           string `db:"item_name" json:"item_name"`
	SKU                string `db:"sku" json:"sku"`
	Quantity           int    `db:"quantity" json:"quantity"`
	Status             int    `db:"status" json:"status"`
}

type PickingorderDetailResponse struct {
	PickingorderDetailID string `db:"pickingorder_detail_id" json:"pickingorder_detail_id"`
	OrganizationID       string `db:"organization_id" json:"organization_id"`
	PickingorderID       string `db:"pickingorder_id" json:"pickingorder_id"`
	LocationID           string `db:"location_id" json:"location_id"`
	LocationCode         string `db:"location_code" json:"location_code"`
	ItemID               string `db:"item_id" json:"item_id"`
	ItemName             string `db:"item_name" json:"item_name"`
	SKU                  string `db:"sku" json:"sku"`
	Quantity             int    `db:"quantity" json:"quantity"`
	QuantityPicked       int    `db:"quantity_picked" json:"quantity_picked"`
	Status               int    `db:"status" json:"status"`
}

type PickingFromLocationNew struct {
	LocationID     string `json:"location_id" binding:"required,min=1"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
	Email          string `json:"email" swaggerignore:"true"`
}

type PackageNew struct {
	PackageNumber  string           `json:"package_number" binding:"required,min=6,max=64"`
	PackageDate    string           `json:"package_date" binding:"required,datetime=2006-01-02"`
	Notes          string           `json:"notes" binding:"omitempty"`
	Items          []PackageItemNew `json:"items" binding:"required"`
	OrganizationID string           `json:"organiztion_id" swaggerignore:"true"`
	User           string           `json:"user" swaggerignore:"true"`
	Email          string           `json:"email" swaggerignore:"true"`
}
type PackageItemNew struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}
