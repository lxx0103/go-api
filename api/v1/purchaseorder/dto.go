package purchaseorder

import (
	"go-api/core/request"
)

type PurchaseorderNew struct {
	PurchaseorderNumber  string                 `json:"purchaseorder_number" binding:"required,min=6,max=64"`
	PurchaseorderDate    string                 `json:"purchaseorder_date" binding:"required,datetime=2006-01-02"`
	ExpectedDeliveryDate string                 `json:"expected_delivery_date" binding:"required,datetime=2006-01-02"`
	VendorID             string                 `json:"vendor_id" binding:"required"`
	DiscountType         int                    `json:"discount_type" binding:"omitempty,oneof=1 2"`
	DiscountValue        float64                `json:"discount_value" binding:"omitempty"`
	ShippingFee          float64                `json:"shipping_fee" binding:"omitempty"`
	Notes                string                 `json:"notes" binding:"omitempty"`
	Items                []PurchaseorderItemNew `json:"items" binding:"required"`
	OrganizationID       string                 `json:"organiztion_id" swaggerignore:"true"`
	User                 string                 `json:"user" swaggerignore:"true"`
	Email                string                 `json:"email" swaggerignore:"true"`
}

type PurchaseorderItemNew struct {
	PurchaseorderItemID string  `json:"purchaseorder_item_id" binding:"omitempty"`
	ItemID              string  `json:"item_id" binding:"required"`
	Quantity            int     `json:"quantity" binding:"required"`
	Rate                float64 `json:"rate" binding:"required"`
	TaxID               string  `json:"tax_id" binding:"omitempty"`
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
	ItemCount            float64 `db:"item_count" json:"item_count"`
	Subtotal             float64 `db:"sub_total" json:"sub_total"`
	TaxTotal             float64 `db:"tax_total" json:"tax_total"`
	DiscountType         int     `db:"discount_type" json:"discount_type"`
	DiscountValue        float64 `db:"discount_value" json:"discount_value"`
	ShippingFee          float64 `db:"shipping_fee" json:"shipping_fee"`
	Total                float64 `db:"total" json:"total"`
	Notes                string  `db:"notes" json:"notes"`
	BillingStatus        int     `db:"billing_status" json:"billing_status"`
	ReceiveStatus        int     `db:"receive_status" json:"receive_status"`
	Status               int     `db:"status" json:"status"`
}

type PurchaseorderItemResponse struct {
	OrganizationID      string  `db:"organization_id" json:"organization_id"`
	PurchaseorderID     string  `db:"purchaseorder_id" json:"purchaseorder_id"`
	PurchaseorderItemID string  `db:"purchaseorder_item_id" json:"purchaseorder_item_id"`
	ItemID              string  `db:"item_id" json:"item_id"`
	ItemName            string  `db:"item_name" json:"item_name"`
	SKU                 string  `db:"sku" json:"sku"`
	Quantity            int     `db:"quantity" json:"quantity"`
	Rate                float64 `db:"rate" json:"rate"`
	TaxID               string  `db:"tax_id" json:"tax_id"`
	TaxValue            float64 `db:"tax_value" json:"tax_value"`
	TaxAmount           float64 `db:"tax_amount" json:"tax_amount"`
	Amount              float64 `db:"amount" json:"amount"`
	QuantityReceived    int     `db:"quantity_received" json:"quantity_received"`
	QuantityBilled      int     `db:"quantity_billed" json:"quantity_billed"`
	Status              int     `db:"status" json:"status"`
}

type PurchaseorderID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type PurchasereceiveNew struct {
	PurchasereceiveNumber string                   `json:"purchasereceive_number" binding:"required,min=6,max=64"`
	PurchasereceiveDate   string                   `json:"purchasereceive_date" binding:"required,datetime=2006-01-02"`
	Notes                 string                   `json:"notes" binding:"omitempty"`
	Items                 []PurchasereceiveItemNew `json:"items" binding:"required"`
	OrganizationID        string                   `json:"organiztion_id" swaggerignore:"true"`
	User                  string                   `json:"user" swaggerignore:"true"`
	Email                 string                   `json:"email" swaggerignore:"true"`
}

type PurchasereceiveItemNew struct {
	ItemID   string `json:"item_id" binding:"omitempty"`
	Quantity int    `json:"quantity" binding:"required"`
}

type PurchasereceiveResponse struct {
	ID                    int64  `db:"id" json:"id"`
	OrganizationID        string `db:"organization_id" json:"organization_id"`
	PurchaseorderID       string `db:"purchaseorder_id" json:"purchaseorder_id"`
	PurchaseorderNumber   string `db:"purchaseorder_number" json:"purchaseorder_number"`
	PurchasereceiveID     string `db:"purchasereceive_id" json:"purchasereceive_id"`
	PurchasereceiveNumber string `db:"purchasereceive_number" json:"purchasereceive_number"`
	PurchasereceiveDate   string `db:"purchasereceive_date" json:"purchasereceive_date"`
	Notes                 string `db:"notes" json:"notes"`
	Status                int    `db:"status" json:"status"`
}

type PurchasereceiveFilter struct {
	PurchasereceiveNumber string `form:"purchasereceive_number" binding:"omitempty,max=64,min=1"`
	OrganizationID        string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type PurchasereceiveID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type PurchasereceiveItemResponse struct {
	OrganizationID        string `db:"organization_id" json:"organization_id"`
	PurchasereceiveID     string `db:"purchasereceive_id" json:"purchasereceive_id"`
	PurchaseorderItemID   string `db:"purchaseorder_item_id" json:"purchaseorder_item_id"`
	PurchasereceiveItemID string `db:"purchasereceive_item_id" json:"purchasereceive_item_id"`
	ItemID                string `db:"item_id" json:"item_id"`
	ItemName              string `db:"item_name" json:"item_name"`
	SKU                   string `db:"sku" json:"sku"`
	Quantity              int    `db:"quantity" json:"quantity"`
	Status                int    `db:"status" json:"status"`
}

type PurchasereceiveDetailResponse struct {
	PurchasereceiveDetailID string `db:"purchasereceive_detail_id" json:"purchasereceive_detail_id"`
	OrganizationID          string `db:"organization_id" json:"organization_id"`
	PurchasereceiveID       string `db:"purchasereceive_id" json:"purchasereceive_id"`
	PurchaseorderItemID     string `db:"purchaseorder_item_id" json:"purchaseorder_item_id"`
	PurchasereceiveItemID   string `db:"purchasereceive_item_id" json:"purchasereceive_item_id"`
	LocationID              string `db:"location_id" json:"location_id"`
	LocationCode            string `db:"location_code" json:"location_code"`
	ItemID                  string `db:"item_id" json:"item_id"`
	ItemName                string `db:"item_name" json:"item_name"`
	SKU                     string `db:"sku" json:"sku"`
	Quantity                int    `db:"quantity" json:"quantity"`
	Status                  int    `db:"status" json:"status"`
}

type BillNew struct {
	BillNumber     string        `json:"bill_number" binding:"required,min=6,max=64"`
	BillDate       string        `json:"bill_date" binding:"required,datetime=2006-01-02"`
	DueDate        string        `json:"due_date" binding:"required,datetime=2006-01-02"`
	VendorID       string        `json:"vendor_id" binding:"required"`
	DiscountType   int           `json:"discount_type" binding:"omitempty,oneof=1 2"`
	DiscountValue  float64       `json:"discount_value" binding:"omitempty"`
	ShippingFee    float64       `json:"shipping_fee" binding:"omitempty"`
	Notes          string        `json:"notes" binding:"omitempty"`
	Items          []BillItemNew `json:"items" binding:"required"`
	OrganizationID string        `json:"organiztion_id" swaggerignore:"true"`
	User           string        `json:"user" swaggerignore:"true"`
	Email          string        `json:"email" swaggerignore:"true"`
}

type BillItemNew struct {
	PurchaseorderItemID string  `json:"purchaseorder_item_id" binding:"required"`
	ItemID              string  `json:"item_id" binding:"required"`
	Quantity            int     `json:"quantity" binding:"required"`
	Rate                float64 `json:"rate" binding:"required"`
	TaxID               string  `json:"tax_id" binding:"omitempty"`
}

type BillResponse struct {
	OrganizationID      string  `db:"organization_id" json:"organization_id"`
	PurchaseorderID     string  `db:"purchaseorder_id" json:"purchaseorder_id"`
	PurchaseorderNumber string  `db:"purchaseorder_number" json:"purchaseorder_number"`
	BillID              string  `db:"bill_id" json:"bill_id"`
	BillNumber          string  `db:"bill_number" json:"bill_number"`
	BillDate            string  `db:"bill_date" json:"bill_date"`
	DueDate             string  `db:"due_date" json:"due_date"`
	VendorID            string  `db:"vendor_id" json:"vendor_id"`
	VendorName          string  `db:"vendor_name" json:"vendor_name"`
	ItemCount           float64 `db:"item_count" json:"item_count"`
	Subtotal            float64 `db:"sub_total" json:"sub_total"`
	DiscountType        int     `db:"discount_type" json:"discount_type"`
	DiscountValue       float64 `db:"discount_value" json:"discount_value"`
	TaxTotal            float64 `db:"tax_total" json:"tax_total"`
	ShippingFee         float64 `db:"shipping_fee" json:"shipping_fee"`
	Total               float64 `db:"total" json:"total"`
	Notes               string  `db:"notes" json:"notes"`
	Status              int     `db:"status" json:"status"`
}

type BillFilter struct {
	PurchaseorderID string `form:"purchaseorder_id" binding:"omitempty,max=64"`
	BillNumber      string `form:"bill_number" binding:"omitempty,max=64,min=1"`
	OrganizationID  string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type BillID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type BillItemResponse struct {
	OrganizationID      string  `db:"organization_id" json:"organization_id"`
	BillID              string  `db:"bill_id" json:"bill_id"`
	PurchaseorderItemID string  `db:"purchaseorder_item_id" json:"purchaseorder_item_id"`
	BillItemID          string  `db:"bill_item_id" json:"bill_item_id"`
	ItemID              string  `db:"item_id" json:"item_id"`
	ItemName            string  `db:"item_name" json:"item_name"`
	SKU                 string  `db:"sku" json:"sku"`
	Quantity            int     `db:"quantity" json:"quantity"`
	Rate                float64 `db:"rate" json:"rate"`
	TaxID               string  `db:"tax_id" json:"tax_id"`
	TaxValue            float64 `db:"tax_value" json:"tax_value"`
	TaxAmount           float64 `db:"tax_amount" json:"tax_amount"`
	Amount              float64 `db:"amount" json:"amount"`
	Status              int     `db:"status" json:"status"`
}

type PaymentMadeNew struct {
	PaymentMadeNumber string  `json:"payment_made_number" binding:"required,min=6,max=64"`
	PaymentMadeDate   string  `json:"payment_made_date" binding:"required,datetime=2006-01-02"`
	PaymentMethodID   string  `json:"payment_method_id" binding:"required,min=6,max=64"`
	Amount            float64 `json:"amount" binding:"required"`
	Notes             string  `json:"notes" binding:"omitempty"`
	OrganizationID    string  `json:"organiztion_id" swaggerignore:"true"`
	User              string  `json:"user" swaggerignore:"true"`
	Email             string  `json:"email" swaggerignore:"true"`
}

type PaymentMadeFilter struct {
	BillID            string `form:"bill_id" binding:"omitempty,max=64"`
	PaymentMadeNumber string `form:"payment_made_number" binding:"omitempty,max=64,min=1"`
	PaymentMethodID   string `form:"payment_method_id" binding:"omitempty,max=64"`
	OrganizationID    string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type PaymentMadeResponse struct {
	OrganizationID    string  `db:"organization_id" json:"organization_id"`
	BillID            string  `db:"bill_id" json:"bill_id"`
	BillNumber        string  `db:"bill_number" json:"bill_number"`
	VendorID          string  `db:"vendor_id" json:"vendor_id"`
	VendorName        string  `db:"vendor_name" json:"vendor_name"`
	PaymentMadeID     string  `db:"payment_made_id" json:"payment_made_id"`
	PaymentMadeNumber string  `db:"payment_made_number" json:"payment_made_number"`
	PaymentMadeDate   string  `db:"payment_made_date" json:"payment_made_date"`
	PaymentMethodID   string  `db:"payment_method_id" json:"payment_method_id"`
	PaymentMethodName string  `db:"payment_method_name" json:"payment_method_name"`
	Amount            float64 `db:"amount" json:"amount"`
	Notes             string  `db:"notes" json:"notes"`
	Status            int     `db:"status" json:"status"`
}

type PaymentMadeID struct {
	ID string `uri:"id" binding:"required,min=1"`
}
