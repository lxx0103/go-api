package purchaseorder

import (
	"go-api/core/request"
)

type PurchaseorderNew struct {
	PurchaseorderNumber  string  `json:"purchaseorder_number" binding:"required,min=6,max=64"`
	PurchaseorderDate    string  `json:"purchaseorder_date" binding:"required,datetime=2006-01-02"`
	ExpectedDeliveryDate string  `json:"expected_delivery_date" binding:"required,datetime=2006-01-02"`
	VendorID             string  `json:"vendor_id" binding:"required"`
	DiscountType         int     `json:"discount_type" binding:"omitempty,oneof=1 2"`
	DiscountValue        float64 `json:"discount_value" binding:"omitempty"`
	ShippingFee          float64 `json:"shipping_fee" binding:"omitempty"`
	Notes                string  `json:"notes" binding:"omitempty"`
	Status               int     `json:"status" binding:"required,oneof=1 2"`
	Items                []PurchaseorderItemNew
	OrganizationID       string `json:"organiztion_id" swaggerignore:"true"`
	User                 string `json:"user" swaggerignore:"true"`
}

type PurchaseorderItemNew struct {
	ItemID   string  `json:"item_id" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Rate     float64 `json:"rate" binding:"required"`
}

type PurchaseorderFilter struct {
	PurchaseorderNumber string `form:"purchaseorder_number" binding:"omitempty,max=64,min=1"`
	OrganizationID      string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type PurchaseorderResponse struct {
	PurchaseorderID      string  `db:"purchaseorder_id" json:"purchaseorder_id"`
	OrganizationID       string  `db:"organization_id" json:"organization_id"`
	PurchaseorderNumber  string  `db:"purchaseorder_number" json:"purchaseorder_number"`
	PurchaseorderDate    string  `db:"purchaseorder_date" json:"purchaseorder_date"`
	ExpectedDeliveryDate string  `db:"expected_delivery_date" json:"expected_delivery_date"`
	VendorID             string  `db:"vendor_id" json:"vendor_id"`
	VendorName           string  `db:"vendor_name" json:"vendor_name"`
	ItemCount            int     `db:"item_count" json:"item_count"`
	Subtotal             float64 `db:"subtotal" json:"subtotal"`
	DiscountType         int     `db:"discount_type" json:"discount_type"`
	DiscountValue        float64 `db:"discount_value" json:"discount_value"`
	ShippingFee          float64 `db:"shipping_fee" json:"shipping_fee"`
	Total                float64 `db:"total" json:"total"`
	Notes                string  `db:"notes" json:"notes"`
	Status               int     `db:"status" json:"status"`
}

type PurchaseorderID struct {
	ID string `uri:"id" binding:"required,min=1"`
}
