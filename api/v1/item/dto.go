package item

import (
	"go-api/core/request"
)

type ItemNew struct {
	SKU             string  `json:"sku" binding:"required,min=4,max=64"`
	Name            string  `json:"name" binding:"required,min=4,max=255"`
	UnitID          string  `json:"unit_id" binding:"required,min=6,max=64"`
	ManufacturerID  string  `json:"manufacturer_id" binding:"omitempty,min=6,max=64"`
	BrandID         string  `json:"brand_id" binding:"omitempty,min=6,max=64"`
	WeightUnit      string  `json:"weight_unit" binding:"omitempty,min=6,max=64"`
	Weight          float64 `json:"weight" binding:"omitempty"`
	DimensionUnit   string  `json:"dimension_unit" binding:"omitempty,min=6,max=64"`
	Length          float64 `json:"length" binding:"omitempty"`
	Width           float64 `json:"width" binding:"omitempty"`
	Height          float64 `json:"height" binding:"omitempty"`
	SellingPrice    float64 `json:"selling_price" binding:"omitempty"`
	CostPrice       float64 `json:"cost_price" binding:"omitempty"`
	ReorderStock    int     `json:"reorder_stock" binding:"omitempty"`
	DefaultVendorID string  `json:"default_vendor_id" binding:"omitempty"`
	Description     string  `json:"description" binding:"omitempty"`
	TrackLocation   int     `json:"track_location" binding:"required,oneof=1 2"`
	Status          int     `json:"status" binding:"required,oneof=1 2"`
	OrganizationID  string  `json:"organiztion_id" swaggerignore:"true"`
	User            string  `json:"user" swaggerignore:"true"`
	Email           string  `json:"email" swaggerignore:"true"`
}

type ItemFilter struct {
	Name           string `form:"name" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type ItemResponse struct {
	ItemID            string  `db:"item_id" json:"item_id"`
	OrganizationID    string  `db:"organization_id" json:"organization_id"`
	SKU               string  `db:"sku" json:"sku"`
	Name              string  `db:"name" json:"name"`
	UnitID            string  `db:"unit_id" json:"unit_id"`
	UnitName          string  `db:"unit_name" json:"unit_name"`
	ManufacturerID    string  `db:"manufacturer_id" json:"manufacturer_id"`
	ManufacturerName  string  `db:"manufacturer_name" json:"manufacturer_name"`
	BrandID           string  `db:"brand_id" json:"brand_id"`
	BrandName         string  `db:"brand_name" json:"brand_name"`
	WeightUnit        string  `db:"weight_unit" json:"weight_unit"`
	WeightUnitName    string  `db:"weight_unit_name" json:"weight_unit_name"`
	Weight            float64 `db:"weight" json:"weight"`
	DimensionUnit     string  `db:"dimension_unit" json:"dimension_unit"`
	DimensionUnitName string  `db:"dimension_unit_name" json:"dimension_unit_name"`
	Length            float64 `db:"length" json:"length"`
	Width             float64 `db:"width" json:"width"`
	Height            float64 `db:"height" json:"height"`
	SellingPrice      float64 `db:"selling_price" json:"selling_price"`
	CostPrice         float64 `db:"cost_price" json:"cost_price"`
	ReorderStock      int     `db:"reorder_stock" json:"reorder_stock"`
	StockOnHand       int     `db:"stock_on_hand" json:"stock_on_hand"`
	StockAvailable    int     `db:"stock_available" json:"stock_available"`
	StockPicking      int     `db:"stock_picking" json:"stock_picking"`
	StockPacking      int     `db:"stock_packing" json:"stock_packing"`
	DefaultVendorID   string  `db:"default_vendor_id" json:"default_vendor_id"`
	Description       string  `db:"description" json:"description"`
	TrackLocation     int     `db:"track_location" json:"track_location"`
	Status            int     `db:"status" json:"status"`
}

type ItemID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type BarcodeFilter struct {
	Code           string `form:"code" binding:"omitempty,max=64,min=1"`
	ItemID         string `form:"item_id" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type BarcodeResponse struct {
	BarcodeID      string `db:"barcode_id" json:"barcode_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	ItemID         string `db:"item_id" json:"item_id"`
	ItemName       string `db:"item_name" json:"item_name"`
	Code           string `db:"code" json:"code"`
	SKU            string `db:"sku" json:"sku"`
	Unit           string `db:"unit" json:"unit"`
	Quantity       int    `db:"quantity" json:"quantity"`
	Status         int    `db:"status" json:"status"`
}

type BarcodeNew struct {
	Code           string `json:"code" binding:"required,min=1,max=64"`
	ItemID         string `json:"item_id" binding:"required,min=1,max=64"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}
type BarcodeID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type ItemBatchResponse struct {
	OrganizationID string `db:"organization_id" json:"organization_id"`
	ItemID         string `db:"item_id" json:"item_id"`
	SKU            string `db:"sku" json:"sku"`
	ItemName       string `db:"item_name" json:"item_name"`
	BatchID        string `db:"batch_id" json:"batch_id"`
	Type           string `db:"type" json:"type"`
	ReferenceID    string `db:"reference_id" json:"reference_id"`
	LocationID     string `db:"location_id" json:"location_id"`
	Quantity       int    `db:"quantity" json:"quantity"`
	Balance        int    `db:"balance" json:"balance"`
	Status         int    `db:"status" json:"status"`
}

type BarcodeCode struct {
	Code string `uri:"code" binding:"required,min=1"`
}
