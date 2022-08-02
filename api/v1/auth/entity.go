package auth

import "time"

type User struct {
	ID             int64     `db:"id" json:"id"`
	UserID         string    `db:"user_id" json:"user_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	RoleID         string    `db:"role_id" json:"role_id"`
	Email          string    `db:"email" json:"email"`
	Password       string    `db:"password" json:"password"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
type Role struct {
	ID             int64     `db:"id" json:"id"`
	RoleID         string    `db:"role_id" json:"role_id"`
	OrganizationID string    `db:"organization_id" json:"organization_id"`
	Name           string    `db:"name" json:"name"`
	Priority       int       `db:"priority" json:"priority"`
	IsDefault      int       `db:"is_default" json:"is_default"`
	IsAdmin        int       `db:"is_admin" json:"is_admin"`
	Status         int       `db:"status" json:"status"`
	Created        time.Time `db:"created" json:"created"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	Updated        time.Time `db:"updated" json:"updated"`
	UpdatedBy      string    `db:"updated_by" json:"updated_by"`
}
type Menu struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Action    string    `db:"action" json:"action"`
	Title     string    `db:"title" json:"title"`
	Path      string    `db:"path" json:"path"`
	Component string    `db:"component" json:"component"`
	IsHidden  int64     `db:"is_hidden" json:"is_hidden"`
	ParentID  int64     `db:"parent_id" json:"parent_id"`
	Status    int       `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	Updated   time.Time `db:"updated" json:"updated"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
