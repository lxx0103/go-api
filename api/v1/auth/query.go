package auth

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type authQuery struct {
	conn *sqlx.DB
}

func NewAuthQuery(connection *sqlx.DB) AuthQuery {
	return &authQuery{
		conn: connection,
	}
}

type AuthQuery interface {
	// //User Management
	// GetUserByID(int64, int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	// GetUserCount(UserFilter) (int, error)
	// GetUserList(UserFilter) (*[]UserResponse, error)
	// //Role Management
	GetRoleByID(id int64) (*RoleResponse, error)
	GetRoleCount(filter RoleFilter) (int, error)
	GetRoleList(filter RoleFilter) (*[]RoleResponse, error)
	// // //API Management
	// GetAPIByID(id int64) (*API, error)
	// GetAPICount(filter APIFilter) (int, error)
	// GetAPIList(filter APIFilter) (*[]API, error)
	// //Menu Management
	// GetMenuByID(id int64) (*Menu, error)
	// GetMenuCount(filter MenuFilter) (int, error)
	// GetMenuList(filter MenuFilter) (*[]Menu, error)
	// GetMenuAPIByID(int64) ([]int64, error)
	// GetRoleMenuByID(int64) ([]int64, error)
	// GetMyMenu(int64) ([]Menu, error)
}

// func (r *authQuery) GetUserByID(id int64, organizationID int64) (*User, error) {
// 	var user User
// 	var err error
// 	if organizationID != 0 {
// 		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? AND organization_id = ?", id, organizationID)
// 	} else {
// 		err = r.conn.Get(&user, "SELECT * FROM users WHERE id = ? ", id)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	user.Credential = ""
// 	return &user, nil
// }

func (r *authQuery) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.conn.Get(&user, `
		SELECT *
		FROM u_users
		WHERE email = ?
	`, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// func (r *authQuery) GetUserCount(filter UserFilter) (int, error) {
// 	where, args := []string{"status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Type; v == "wx" {
// 		where, args = append(where, "type = ?"), append(args, 2)
// 	}
// 	if v := filter.Type; v == "admin" {
// 		where, args = append(where, "type = ?"), append(args, 1)
// 	}
// 	if v := filter.OrganizationID; v != 0 {
// 		where, args = append(where, "organization_id = ?"), append(args, v)
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM users
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authQuery) GetUserList(filter UserFilter) (*[]UserResponse, error) {
// 	where, args := []string{"u.status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "u.name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Type; v == "wx" {
// 		where, args = append(where, "u.type = ?"), append(args, 2)
// 	}
// 	if v := filter.Type; v == "admin" {
// 		where, args = append(where, "u.type = ?"), append(args, 1)
// 	}
// 	if v := filter.OrganizationID; v != 0 {
// 		where, args = append(where, "u.organization_id = ?"), append(args, v)
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var users []UserResponse
// 	err := r.conn.Select(&users, `
// 		SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
// 		FROM users u
// 		LEFT JOIN organizations o
// 		ON u.organization_id = o.id
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &users, nil
// }

func (r *authQuery) GetRoleByID(id int64) (*RoleResponse, error) {
	var role RoleResponse
	err := r.conn.Get(&role, "SELECT id, name, organization_id, is_admin, is_default, priority, status FROM s_roles WHERE id = ? AND status > 0", id)
	return &role, err
}

func (r *authQuery) GetRoleCount(filter RoleFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_roles
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *authQuery) GetRoleList(filter RoleFilter) (*[]RoleResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != 0 {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var roles []RoleResponse
	err := r.conn.Select(&roles, `
		SELECT id, name, priority, organization_id, is_admin, is_default, status
		FROM s_roles
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &roles, err
}

// func (r *authQuery) GetAPIByID(id int64) (*API, error) {
// 	var api API
// 	err := r.conn.Get(&api, "SELECT * FROM apis WHERE id = ? ", id)
// 	return &api, err
// }

// func (r *authQuery) GetAPICount(filter APIFilter) (int, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Route; v != "" {
// 		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM apis
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authQuery) GetAPIList(filter APIFilter) (*[]API, error) {
// 	where, args := []string{"1 = 1"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.Route; v != "" {
// 		where, args = append(where, "route like ?"), append(args, "%"+v+"%")
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var apis []API
// 	err := r.conn.Select(&apis, `
// 		SELECT *
// 		FROM apis
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &apis, nil
// }

// func (r *authQuery) GetMenuByID(id int64) (*Menu, error) {
// 	var menu Menu
// 	err := r.conn.Get(&menu, "SELECT * FROM menus WHERE id = ? AND status > 0 ", id)
// 	return &menu, err
// }

// func (r *authQuery) GetMenuCount(filter MenuFilter) (int, error) {
// 	where, args := []string{"status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.OnlyTop; v {
// 		where, args = append(where, "parent_id = ?"), append(args, 0)
// 	}
// 	var count int
// 	err := r.conn.Get(&count, `
// 		SELECT count(1) as count
// 		FROM menus
// 		WHERE `+strings.Join(where, " AND "), args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (r *authQuery) GetMenuList(filter MenuFilter) (*[]Menu, error) {
// 	where, args := []string{"status > 0"}, []interface{}{}
// 	if v := filter.Name; v != "" {
// 		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
// 	}
// 	if v := filter.OnlyTop; v {
// 		where, args = append(where, "parent_id = ?"), append(args, 0)
// 	}
// 	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
// 	args = append(args, filter.PageSize)
// 	var menus []Menu
// 	err := r.conn.Select(&menus, `
// 		SELECT *
// 		FROM menus
// 		WHERE `+strings.Join(where, " AND ")+`
// 		LIMIT ?, ?
// 	`, args...)
// 	return &menus, err
// }

// func (r *authQuery) GetMenuAPIByID(menuID int64) ([]int64, error) {
// 	var apis []int64
// 	err := r.conn.Select(&apis, "SELECT api_id FROM menu_apis WHERE menu_id = ? and status > 0", menuID)
// 	return apis, err
// }

// func (r *authQuery) GetRoleMenuByID(roleID int64) ([]int64, error) {
// 	var menu []int64
// 	err := r.conn.Select(&menu, "SELECT menu_id FROM role_menus WHERE role_id = ? and status > 0", roleID)
// 	return menu, err
// }

// func (r *authQuery) GetMyMenu(roleID int64) ([]Menu, error) {
// 	var menu []Menu
// 	err := r.conn.Select(&menu, `
// 		SELECT m.* FROM role_menus rm
// 		LEFT JOIN menus m
// 		ON rm.menu_id = m.id
// 		WHERE rm.role_id = ?
// 		AND m.status > 0
// 		AND rm.status > 0
// 		ORDER BY parent_id ASC, ID ASC
// 	`, roleID)
// 	return menu, err
// }
