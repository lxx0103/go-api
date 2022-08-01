package auth

import (
	"database/sql"
	"time"
)

type authRepository struct {
	tx *sql.Tx
}

func NewAuthRepository(transaction *sql.Tx) AuthRepository {
	return &authRepository{
		tx: transaction,
	}
}

type AuthRepository interface {
	// // GetCredential(SigninRequest) (UserAuth, error)
	CreateUser(User) (int64, error)
	// GetUserByID(int64) (*UserResponse, error)
	// CheckConfict(int, string) (bool, error)
	// UpdateUser(int64, UserResponse, string) error
	// UpdatePassword(int64, string, string) error
	// // GetAuthCount(filter AuthFilter) (int, error)
	// // GetAuthList(filter AuthFilter) ([]Auth, error)

	//  Role Management
	GetRoleByID(int64) (*RoleResponse, error)
	CheckRoleConfict(int64, string) (bool, error)
	CreateRole(info Role) (int64, error)
	UpdateRole(int64, Role) error
	DeleteRole(int64, string) error
	// // API Management
	// CreateAPI(APINew) (int64, error)
	// UpdateAPI(int64, APINew) (int64, error)
	// GetAPIByID(int64) (*API, error)
	// // Menu Management
	// GetMenuByID(id int64) (*Menu, error)
	// CreateMenu(info MenuNew) (int64, error)
	// UpdateMenu(int64, Menu, string) error
	// DeleteMenu(int64, string) error
	// // Privilege Management
	// NewMenuAPI(int64, MenuAPINew) error
	// NewRoleMenu(int64, RoleMenuNew) error
}

func (r *authRepository) CreateUser(newUser User) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO u_users
		(
			organization_id,
			role_id,
			email,
			password,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, newUser.OrganizationID, newUser.RoleID, newUser.Email, newUser.Password, newUser.Status, time.Now(), newUser.CreatedBy, time.Now(), newUser.UpdatedBy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// // func (r *authRepository) GetUserByID(id int64) (*UserResponse, error) {
// // 	var res UserResponse
// // 	row := r.tx.QueryRow(`
// // 	SELECT u.id as id, u.type as type, u.identifier as identifier, u.organization_id as organization_id, u.position_id as position_id, u.role_id as role_id, u.name as name, u.email as email, u.gender as gender, u.phone as phone, u.birthday as birthday, u.address as address, u.status as status, IFNULL(o.name, "ADMIN") as organization_name
// // 	FROM users u
// // 	LEFT JOIN organizations o
// // 	ON u.organization_id = o.id
// // 	WHERE u.id = ?
// // 	`, id)
// // 	err := row.Scan(&res.ID, &res.Type, &res.Identifier, &res.OrganizationID, &res.PositionID, &res.RoleID, &res.Name, &res.Email, &res.Gender, &res.Phone, &res.Birthday, &res.Address, &res.Status, &res.OrganizationName)
// // 	if err != nil {
// // 		msg := "用户不存在:" + err.Error()
// // 		return nil, errors.New(msg)
// // 	}
// // 	return &res, nil
// // }

// func (r *authRepository) CheckConfict(authType int, identifier string) (bool, error) {
// 	var existed int
// 	row := r.tx.QueryRow("SELECT count(1) FROM users WHERE type = ? AND identifier = ?", authType, identifier)
// 	err := row.Scan(&existed)
// 	if err != nil {
// 		return true, err
// 	}
// 	return existed != 0, nil
// }
// func (r *authRepository) UpdateUser(id int64, info UserResponse, by string) error {
// 	_, err := r.tx.Exec(`
// 		Update users SET
// 		name = ?,
// 		email = ?,
// 		role_id = ?,
// 		position_id = ?,
// 		gender = ?,
// 		phone = ?,
// 		birthday = ?,
// 		address = ?,
// 		status = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Email, info.RoleID, info.PositionID, info.Gender, info.Phone, info.Birthday, info.Address, info.Status, time.Now(), by, id)
// 	if err != nil {
// 		msg := "更新失败:" + err.Error()
// 		return errors.New(msg)
// 	}
// 	return nil
// }

func (r *authRepository) GetRoleByID(id int64) (*RoleResponse, error) {
	var res RoleResponse
	row := r.tx.QueryRow(`SELECT id, organization_id, priority, name, is_admin, is_default, status FROM s_roles WHERE id = ? AND status > 0 LIMIT 1`, id)
	err := row.Scan(&res.ID, &res.OrganizationID, &res.Priority, &res.Name, &res.IsAdmin, &res.IsDefault, &res.Status)
	return &res, err
}

func (r *authRepository) CheckRoleConfict(roleID int64, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_roles WHERE id != ? AND name = ?", roleID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *authRepository) CreateRole(info Role) (int64, error) {
	result, err := r.tx.Exec(`
		INSERT INTO s_roles
		(
			organization_id,
			name,
			priority,
			is_default,
			is_admin,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.Name, info.Priority, info.IsDefault, info.IsAdmin, info.Status, time.Now(), info.CreatedBy, time.Now(), info.UpdatedBy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *authRepository) UpdateRole(id int64, info Role) error {
	_, err := r.tx.Exec(`
		Update s_roles SET
		name = ?,
		priority = ?,
		is_admin = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, info.Name, info.Priority, info.IsAdmin, info.Status, time.Now(), info.UpdatedBy, id)
	return err
}

func (r *authRepository) DeleteRole(id int64, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_roles SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE id = ?
	`, time.Now(), byUser, id)
	return err
}

// func (r *authRepository) CreateAPI(info APINew) (int64, error) {
// 	result, err := r.tx.Exec(`
// 		INSERT INTO apis
// 		(
// 			name,
// 			route,
// 			method,
// 			status,
// 			created,
// 			created_by,
// 			updated,
// 			updated_by
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
// 	`, info.Name, info.Route, info.Method, info.Status, time.Now(), info.User, time.Now(), info.User)
// 	if err != nil {
// 		return 0, err
// 	}
// 	id, err := result.LastInsertId()
// 	return id, err
// }

// func (r *authRepository) UpdateAPI(id int64, info APINew) (int64, error) {
// 	result, err := r.tx.Exec(`
// 		Update apis SET
// 		name = ?,
// 		route = ?,
// 		method = ?,
// 		status = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Route, info.Method, info.Status, time.Now(), info.User, id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	affected, err := result.RowsAffected()
// 	return affected, err
// }

// func (r *authRepository) GetAPIByID(id int64) (*API, error) {
// 	var res API
// 	row := r.tx.QueryRow(`SELECT id, name, route, method, status, created, created_by, updated, updated_by FROM apis WHERE id = ? LIMIT 1`, id)
// 	err := row.Scan(&res.ID, &res.Name, &res.Route, &res.Method, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
// 	if err != nil {
// 		msg := "API不存在:" + err.Error()
// 		return nil, errors.New(msg)
// 	}
// 	return &res, nil
// }

// func (r *authRepository) GetMenuByID(id int64) (*Menu, error) {
// 	var res Menu
// 	row := r.tx.QueryRow(`SELECT id, name, action, title, path, component, is_hidden, parent_id, status, created, created_by, updated, updated_by FROM menus WHERE id = ? LIMIT 1`, id)
// 	err := row.Scan(&res.ID, &res.Name, &res.Action, &res.Title, &res.Path, &res.Component, &res.IsHidden, &res.ParentID, &res.Status, &res.Created, &res.CreatedBy, &res.Updated, &res.UpdatedBy)
// 	if err != nil {
// 		msg := "菜单不存在:" + err.Error()
// 		return nil, errors.New(msg)
// 	}
// 	return &res, nil
// }

// func (r *authRepository) CreateMenu(info MenuNew) (int64, error) {
// 	result, err := r.tx.Exec(`
// 		INSERT INTO menus
// 		(
// 			name,
// 			action,
// 			title,
// 			path,
// 			component,
// 			is_hidden,
// 			parent_id,
// 			status,
// 			created,
// 			created_by,
// 			updated,
// 			updated_by
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, 1, time.Now(), info.User, time.Now(), info.User)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result.LastInsertId()
// }

// func (r *authRepository) UpdateMenu(id int64, info Menu, byUser string) error {
// 	fmt.Println(info.Component)
// 	_, err := r.tx.Exec(`
// 		Update menus SET
// 		name = ?,
// 		action = ?,
// 		title = ?,
// 		path = ?,
// 		component = ?,
// 		is_hidden = ?,
// 		parent_id = ?,
// 		status = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, info.Name, info.Action, info.Title, info.Path, info.Component, info.IsHidden, info.ParentID, info.Status, time.Now(), byUser, id)
// 	return err
// }

// func (r *authRepository) DeleteMenu(id int64, byUser string) error {
// 	_, err := r.tx.Exec(`
// 		Update menus SET
// 		status = -1,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, time.Now(), byUser, id)
// 	return err
// }

// func (r *authRepository) NewRoleMenu(role_id int64, info RoleMenuNew) error {
// 	_, err := r.tx.Exec(`
// 		Update role_menus SET
// 		status = -1,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE role_id = ?
// 	`, time.Now(), info.User, role_id)
// 	if err != nil {
// 		return err
// 	}
// 	sql := `
// 	INSERT INTO role_menus
// 	(
// 		role_id,
// 		menu_id,
// 		status,
// 		created,
// 		created_by,
// 		updated,
// 		updated_by
// 	)
// 	VALUES
// 	`
// 	for i := 0; i < len(info.IDS); i++ {
// 		sql += "(" + fmt.Sprint(role_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
// 	}
// 	sql = sql[:len(sql)-1]
// 	_, err = r.tx.Exec(sql)
// 	return err
// }

// func (r *authRepository) NewMenuAPI(menu_id int64, info MenuAPINew) error {
// 	_, err := r.tx.Exec(`
// 		Update menu_apis SET
// 		status = -1,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE menu_id = ?
// 	`, time.Now(), info.User, menu_id)
// 	if err != nil {
// 		return err
// 	}
// 	sql := `
// 	INSERT INTO menu_apis
// 	(
// 		menu_id,
// 		api_id,
// 		status,
// 		created,
// 		created_by,
// 		updated,
// 		updated_by
// 	)
// 	VALUES
// 	`
// 	for i := 0; i < len(info.IDS); i++ {
// 		sql += "(" + fmt.Sprint(menu_id) + "," + fmt.Sprint(info.IDS[i]) + ",1,\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\",\"" + time.Now().Format("2006-01-02 15:01:01") + "\",\"" + info.User + "\"),"
// 	}
// 	sql = sql[:len(sql)-1]
// 	_, err = r.tx.Exec(sql)
// 	return err
// }

// func (r *authRepository) UpdatePassword(id int64, password, by string) error {
// 	_, err := r.tx.Exec(`
// 		Update users SET
// 		credential = ?,
// 		updated = ?,
// 		updated_by = ?
// 		WHERE id = ?
// 	`, password, time.Now(), by, id)
// 	if err != nil {
// 		msg := "更新失败:" + err.Error()
// 		return errors.New(msg)
// 	}
// 	return nil
// }
