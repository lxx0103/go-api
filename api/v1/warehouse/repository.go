package warehouse

import (
	"database/sql"
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
			bay_id,
			item_id,
			capacity,
			quantity,
			available,
			can_pick,
			alert,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.LocationID, info.Code, info.Level, info.BayID, info.ItemID, info.Capacity, info.Quantity, info.Available, info.CanPick, info.Alert, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *warehouseRepository) GetLocationByID(locationID, organizationID string) (*LocationResponse, error) {
	var res LocationResponse
	row := r.tx.QueryRow(`
		SELECT 
		l.location_id, 
		l.organization_id,
		l.code,
		l.level, 
		l.bay_id,
		b.code as bay_code,
		l.item_id,
		i.name as item_name,
		i.sku,
		l.capacity,
		l.quantity,
		l.available,
		l.can_pick,
		l.alert, 
		l.status
		FROM w_locations l
		LEFT JOIN w_bays b
		ON l.bay_id = b.bay_id
		LEFT JOIN i_items i
		ON l.item_id = i.item_id
		WHERE l.location_id = ? AND l.organization_id = ? AND l.status > 0
	`, locationID, organizationID)
	err := row.Scan(&res.LocationID, &res.OrganizationID, &res.Code, &res.Level, &res.BayID, &res.BayCode, &res.ItemID, &res.ItemName, &res.SKU, &res.Capacity, &res.Quantity, &res.Available, &res.CanPick, &res.Alert, &res.Status)
	return &res, err
}

func (r *warehouseRepository) UpdateLocation(id string, info Location) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		code = ?,
		level = ?,
		bay_id = ?,
		item_id = ?,
		capacity = ?,
		quantity = ?,
		available = ?,
		can_pick = ?,
		alert = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, info.Code, info.Level, info.BayID, info.ItemID, info.Capacity, info.Quantity, info.Available, info.CanPick, info.Alert, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *warehouseRepository) DeleteLocation(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, time.Now(), byUser, id)
	return err
}

func (r *warehouseRepository) GetItemAvailable(itemID, organizationID string) (int, error) {
	var available int
	row := r.tx.QueryRow("SELECT sum(available) FROM w_locations WHERE organization_id = ? AND item_id = ? AND status > 0 ", organizationID, itemID)
	err := row.Scan(&available)
	return available, err
}

func (r *warehouseRepository) GetNextLocation(itemID, organizationID string) (*LocationResponse, error) {
	var location LocationResponse
	row := r.tx.QueryRow("SELECT location_id,available FROM w_locations WHERE organization_id = ? AND item_id = ? AND available > 0  AND status > 0  limit 1", organizationID, itemID)
	err := row.Scan(&location.LocationID, &location.Available)
	return &location, err
}

func (r *warehouseRepository) ReceiveItem(locationID string, quantity int, byUser string) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		available = available - ?,
		quantity = quantity + ?,
		can_pick = can_pick + ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, quantity, quantity, quantity, time.Now(), byUser, locationID)
	return err
}

func (r *warehouseRepository) UpdateLocationCanPick(locationID string, quantity int, byUser string) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		can_pick = can_pick - ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, quantity, time.Now(), byUser, locationID)
	return err
}

func (r *warehouseRepository) UpdateLocationPicked(locationID string, quantity int, byUser string) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		available = available + ?,
		quantity = quantity - ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, quantity, quantity, time.Now(), byUser, locationID)
	return err
}

func (r *warehouseRepository) UpdateLocationQuantity(info Location) error {
	_, err := r.tx.Exec(`
		Update w_locations SET
		available = ?,
		quantity = ?,
		can_pick = ?,
		updated = ?,
		updated_by = ?
		WHERE location_id = ?
	`, info.Available, info.Quantity, info.CanPick, info.Updated, info.UpdatedBy, info.LocationID)
	return err
}

func (r warehouseRepository) CreateAdjustment(info Adjustment) error {
	_, err := r.tx.Exec(`
		INSERT INTO i_adjustments 
		(
			organization_id,
			location_id,
			item_id,
			adjustment_id,
			quantity,
			rate,
			reason,
			remark,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.OrganizationID, info.LocationID, info.ItemID, info.AdjustmentID, info.Quantity, info.Rate, info.Reason, info.Remark, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}
