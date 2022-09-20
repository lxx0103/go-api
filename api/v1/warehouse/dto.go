package warehouse

import (
	"go-api/core/request"
)

type BayFilter struct {
	Code           string `form:"code" binding:"omitempty,max=64,min=1"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type BayResponse struct {
	OrganizationID string `db:"organization_id" json:"organization_id"`
	BayID          string `db:"bay_id" json:"bay_id"`
	Code           string `db:"code" json:"code"`
	Level          int    `db:"level" json:"level"`
	Location       string `db:"location" json:"location"`
	Status         int    `db:"status" json:"status"`
}

type BayNew struct {
	Code           string `json:"code" binding:"required,min=1,max=64"`
	Level          int    `json:"level" binding:"required,min=1,max=64"`
	Location       string `json:"location" binding:"required,min=1"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}
type BayID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type LocationFilter struct {
	Code           string `form:"code" binding:"omitempty,max=64,min=1"`
	BayID          string `form:"bay_id" binding:"omitempty"`
	Level          string `form:"level" binding:"omitempty,min=1,max=64"`
	SKU            string `form:"sku" binding:"omitempty,max=64,min=1"`
	IsAlert        bool   `form:"is_alert" binding:"omitempty"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type LocationResponse struct {
	LocationID     string `db:"location_id" json:"location_id"`
	OrganizationID string `db:"organization_id" json:"organization_id"`
	Code           string `db:"code" json:"code"`
	Level          string `db:"level" json:"level"`
	BayID          string `db:"bay_id" json:"bay_id"`
	BayCode        string `db:"bay_code" json:"bay_code"`
	ItemID         string `db:"item_id" json:"item_id"`
	ItemName       string `db:"item_name" json:"item_name"`
	SKU            string `db:"sku" json:"sku"`
	Capacity       int    `db:"capacity" json:"capacity"`
	Quantity       int    `db:"quantity" json:"quantity"`
	Available      int    `db:"available" json:"available"`
	CanPick        int    `db:"can_pick" json:"can_pick"`
	Alert          int    `db:"alert" json:"alert"`
	Status         int    `db:"status" json:"status"`
}

type LocationNew struct {
	Code           string `json:"code" binding:"required,min=1,max=64"`
	Level          string `json:"level" binding:"required"`
	BayID          string `json:"bay_id" binding:"required"`
	ItemID         string `json:"item_id" binding:"required"`
	Capacity       int    `json:"capacity" binding:"required"`
	Alert          int    `json:"alert" binding:"omitempty"`
	Status         int    `json:"status" binding:"required,oneof=1 2"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}

type LocationID struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type LocationCode struct {
	Code string `uri:"code" binding:"required,min=1"`
}

type AdjustmentNew struct {
	LocationID         string  `json:"location_id" binding:"required"`
	AdjustmentReasonID string  `json:"adjustment_reason_id" binding:"required"`
	Quantity           int     `json:"quantity" binding:"required"`
	Rate               float64 `json:"rate" binding:"omitempty"`
	Remark             string  `json:"remark" binding:"required"`
	AdjustmentDate     string  `json:"adjustment_date" binding:"required,datetime=2006-01-02"`
	OrganizationID     string  `json:"organiztion_id" swaggerignore:"true"`
	User               string  `json:"user" swaggerignore:"true"`
	Email              string  `json:"email" swaggerignore:"true"`
}

type AdjustmentFilter struct {
	LocationID     string `form:"location_id" binding:"omitempty,max=64,min=1"`
	ItemID         string `form:"item_id" binding:"omitempty"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type AdjustmentResponse struct {
	OrganizationID       string  `db:"organization_id" json:"organization_id"`
	LocationID           string  `db:"location_id" json:"location_id"`
	LocationCode         string  `db:"location_code" json:"location_code"`
	ItemID               string  `db:"item_id" json:"item_id"`
	ItemName             string  `db:"item_name" json:"item_name"`
	SKU                  string  `db:"sku" json:"sku"`
	AdjustmentID         string  `db:"adjustment_id" json:"adjustment_id"`
	Quantity             int     `db:"quantity" json:"quantity"`
	OriginalQuantiy      int     `db:"original_quantity" json:"original_quantity"`
	NewQuantiy           int     `db:"new_quantity" json:"new_quantity"`
	Rate                 float64 `db:"rate" json:"rate"`
	AdjustmentDate       string  `db:"adjustment_date" json:"adjustment_date"`
	AdjustmentReasonID   string  `db:"adjustment_reason_id" json:"adjustment_reason_id"`
	AdjustmentReasonName string  `db:"adjustment_reason_name" json:"adjustment_reason_name"`
	Remark               string  `db:"remark" json:"remark"`
	Status               int     `db:"status" json:"status"`
}
