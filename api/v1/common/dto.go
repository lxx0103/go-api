package common

import (
	"go-api/core/request"
)

type HistoryNew struct {
	HistoryType    string `json:"history_type" binding:"required,min=6,max=64"`
	HistoryTime    string `json:"history_time" binding:"required,datetime=2006-01-02 15:04:05"`
	HistoryBy      string `json:"history_by" binding:"required"`
	ReferenceID    string `json:"reference_id" binding:"required"`
	Description    string `json:"description" binding:"required,max=255"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	User           string `json:"user" swaggerignore:"true"`
}
type HistoryFilter struct {
	ReferenceID    string `form:"reference_id" binding:"required,max=64"`
	HistoryType    string `form:"history_type" binding:"required,max=64"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
	request.PageInfo
}

type HistoryResponse struct {
	OrganizationID string `db:"organization_id" json:"organization_id"`
	HistoryID      string `db:"history_id" json:"history_id"`
	HistoryType    string `db:"history_type" json:"history_type"`
	ReferenceID    string `db:"reference_id" json:"reference_id"`
	HistoryTime    string `db:"history_time" json:"history_time"`
	HistoryBy      string `db:"history_by" json:"history_by"`
	Description    string `db:"description" json:"description"`
}

type NumberFilter struct {
	NumberType     string `form:"number_type" binding:"required,oneof=purchaseorder salesorder purchasereceive pickingorder package shippingorder"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
}
