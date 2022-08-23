package item

import "time"

type ItemGroup struct {
	ID             int64     `db:"id" json:"id"`
	ItemGroupID    string    `db:"item_group_id" json:"item_group_id"`
	Name           string    `db:"name" json:"name"`
	Unit           string    `db:"unit" json:"unit"`
	ManufacturerID string    `db:"manufacturer_id" json:"manufacturer_id"`
	BrandID        string    `db:"brand_id" json:"brand_id"`
	Description    string    `db:"description" json:"description"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type ItemGroupAttribute struct {
	ID          int64     `db:"id" json:"id"`
	ItemGroupID string    `db:"item_group_id" json:"item_group_id"`
	AttributeID string    `db:"attribute_id" json:"attribute_id"`
	Name        string    `db:"name" json:"name"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type ItemGroupAttributeOption struct {
	ID          int64     `db:"id" json:"id"`
	AttributeID string    `db:"attribute_id" json:"attribute_id"`
	OptionID    string    `db:"option_id" json:"option_id"`
	Value       string    `db:"value" json:"value"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type Item struct {
	ID              int64     `db:"id" json:"id"`
	OrganizationID  string    `db:"organization_id" json:"organization_id"`
	ItemID          string    `db:"item_id" json:"item_id"`
	SKU             string    `db:"sku" json:"sku"`
	Name            string    `db:"name" json:"name"`
	UnitID          string    `db:"unit_id" json:"unit_id"`
	ManufacturerID  string    `db:"manufacturer_id" json:"manufacturer_id"`
	BrandID         string    `db:"brand_id" json:"brand_id"`
	WeightUnit      string    `db:"weight_unit" json:"weight_unit"`
	Weight          float64   `db:"weight" json:"weight"`
	DimensionUnit   string    `db:"dimension_unit" json:"dimension_unit"`
	Length          float64   `db:"length" json:"length"`
	Width           float64   `db:"width" json:"width"`
	Height          float64   `db:"height" json:"height"`
	SellingPrice    float64   `db:"selling_price" json:"selling_price"`
	CostPrice       float64   `db:"cost_price" json:"cost_price"`
	ReorderStock    int       `db:"reorder_stock" json:"reorder_stock"`
	StockOnHand     int       `db:"stock_on_hand" json:"stock_on_hand"`
	StockAvailable  int       `db:"stock_available" json:"stock_available"`
	StockPicking    int       `db:"stock_picking" json:"stock_picking"`
	StockPacking    int       `db:"stock_packing" json:"stock_packing"`
	DefaultVendorID string    `db:"default_vendor_id" json:"default_vendor_id"`
	Description     string    `db:"description" json:"description"`
	TrackLocation   int       `db:"track_location" json:"track_location"`
	Status          int       `db:"status" json:"status"`
	Created         time.Time `db:"created" json:"created"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	Updated         time.Time `db:"updated" json:"updated"`
	UpdatedBy       string    `db:"updated_by" json:"updated_by"`
}

type ItemAttribute struct {
	ID          int64     `db:"id" json:"id"`
	ItemID      string    `db:"item_id" json:"item_id"`
	AttributeID string    `db:"item_attribute_id" json:"item_attribute_id"`
	OptionID    string    `db:"option_id" json:"option_id"`
	Status      int       `db:"status" json:"status"`
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	Updated     time.Time `db:"updated" json:"updated"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type Barcode struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	BarcodeID      string    `db:"barcode_id" json:"barcode_id"`
	Code           string    `db:"code" json:"code"`
	ItemID         string    `db:"item_id" json:"item_id"`
	Quantity       int       `db:"quantity" json:"quantity"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type ItemBatch struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	ItemID         string    `db:"item_id" json:"item_id"`
	BatchID        string    `db:"batch_id" json:"batch_id"`
	Type           string    `db:"type" json:"type"`
	ReferenceID    string    `db:"reference_id" json:"reference_id"`
	LocationID     string    `db:"location_id" json:"location_id"`
	Quantity       int       `db:"quantity" json:"quantity"`
	Rate           float64   `db:"rate" json:"rate"`
	Balance        int       `db:"balance" json:"balance"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
