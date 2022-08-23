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
