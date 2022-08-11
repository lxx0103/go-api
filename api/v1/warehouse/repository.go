package warehouse

import (
	"database/sql"
	"fmt"
	"time"
)

type warehouseRepository struct {
	tx *sql.Tx
}

func NewWarehouseRepository(tx *sql.Tx) *warehouseRepository {
	return &warehouseRepository{tx: tx}
}

//Bay

func (r *warehouseRepository) CheckBayConfict(bayID, organizationID, code string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM w_bays WHERE organization_id = ? AND bay_id != ? AND code = ? AND status > 0 ", organizationID, bayID, code)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r warehouseRepository) CreateBay(info Bay) error {
	_, err := r.tx.Exec(`
		INSERT INTO w_bays 
		(
			organization_id,
			bay_id,
			code,
			level,
			location,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.BayID, info.Code, info.Level, info.Location, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *warehouseRepository) GetBayByID(bayID, organizationID string) (*BayResponse, error) {
	var res BayResponse
	row := r.tx.QueryRow(`
		SELECT 
		bay_id, 
		organization_id,
		code,
		level, 
		location, 
		status
		FROM w_bays 
		WHERE bay_id = ? AND organization_id = ? LIMIT 1
	`, bayID, organizationID)
	err := row.Scan(&res.BayID, &res.OrganizationID, &res.Code, &res.Level, &res.Location, &res.Status)
	return &res, err
}

func (r *warehouseRepository) UpdateBay(id string, info Bay) error {
	_, err := r.tx.Exec(`
		Update w_bays SET
		code = ?,
		level = ?,
		location = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE bay_id = ?
	`, info.Code, info.Level, info.Location, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *warehouseRepository) DeleteBay(id, byUser string) error {
	fmt.Println(id)
	_, err := r.tx.Exec(`
		Update w_bays SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE bay_id = ?
	`, time.Now(), byUser, id)
	return err
}

//Location

func (r *warehouseRepository) CheckLocationConfict(locationID, organizationID, code string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM w_locations WHERE organization_id = ? AND location_id != ? AND code = ? AND status > 0 ", organizationID, locationID, code)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r warehouseRepository) CreateLocation(info Location) error {
	_, err := r.tx.Exec(`
		INSERT INTO w_locations 
		(
			organization_id,
			location_id,
			code,
			level,
			location,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.LocationID, info.Code, info.Level, info.Location, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *warehouseRepository) GetLocationByID(locationID, organizationID string) (*LocationResponse, error) {
	var res LocationResponse
	row := r.tx.QueryRow(`
		SELECT 
		location_id, 
		organization_id,
		code,
		level, 
		location, 
		status
		FROM w_locations 
		WHERE location_id = ? AND organization_id = ? LIMIT 1
	`, locationID, organizationID)
	err := row.Scan(&res.LocationID, &res.OrganizationID, &res.Code, &res.Level, &res.Location, &res.Status)
	return &res, err
}

func (r *warehouseRepository) UpdateLocation(id string, info Location) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		code = ?,
		level = ?,
		location = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, info.Code, info.Level, info.Location, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *warehouseRepository) DeleteLocation(id, byUser string) error {
	fmt.Println(id)
	_, err := r.tx.Exec(`
		Update w_locations SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, time.Now(), byUser, id)
	return err
}
