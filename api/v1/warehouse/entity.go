package warehouse

import "time"

type Bay struct {
	ID             int64     `db:"id" json:"id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	BayID          string    `db:"bay_id" json:"bay_id"`
	Code           string    `db:"code" json:"code"`
	Level          int       `db:"level" json:"level"`
	Location       string    `db:"location" json:"location"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type Location struct {
	ID             int64     `db:"id" json:"id"`
	LocationID     string    `db:"location_id" json:"location_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Code           string    `db:"code" json:"code"`
	Level          string    `db:"level" json:"level"`
	BayID          string    `db:"bay_id" json:"bay_id"`
	ItemID         string    `db:"item_id" json:"item_id"`
	Capacity       int       `db:"capacity" json:"capacity"`
	Quantity       int       `db:"quantity" json:"quantity"`
	Available      int       `db:"available" json:"available"`
	CanPick        int       `db:"can_pick" json:"can_pick"`
	Alert          int       `db:"alert" json:"alert"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}

type Adjustment struct {
	ID                 int64     `db:"id" json:"id"`
	OrganizationID     string    `db:"organization_id" json:"organization_id"`
	LocationID         string    `db:"location_id" json:"location_id"`
	ItemID             string    `db:"item_id" json:"item_id"`
	AdjustmentID       string    `db:"adjustment_id" json:"adjustment_id"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Rate               float64   `db:"rate" json:"rate"`
	AdjustmentReasonID string    `db:"adjustment_reason_id" json:"adjustment_reason_id"`
	Remark             string    `db:"remark" json:"remark"`
	Status             int       `db:"status" json:"status"`
	Created            time.Time `db:"created" json:"created"`
	CreatedBy          string    `db:"created_by" json:"created_by"`
	Updated            time.Time `db:"updated" json:"updated"`
	UpdatedBy          string    `db:"updated_by" json:"updated_by"`
}
