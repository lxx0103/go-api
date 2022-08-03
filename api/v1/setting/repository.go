package setting

import (
	"database/sql"
	"time"
)

type settingRepository struct {
	tx *sql.Tx
}

func NewSettingRepository(transaction *sql.Tx) *settingRepository {
	return &settingRepository{
		tx: transaction,
	}
}

func (r *settingRepository) GetUnitByID(unitID string) (*UnitResponse, error) {
	var res UnitResponse
	row := r.tx.QueryRow(`SELECT unit_id, organization_id, name, status FROM s_units WHERE unit_id = ? AND status > 0 LIMIT 1`, unitID)
	err := row.Scan(&res.UnitID, &res.OrganizationID, &res.Name, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckUnitConfict(unitID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_units WHERE organization_id = ? AND unit_id != ? AND name = ? AND status > 0", organizationID, unitID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateUnit(info Unit) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_units
		(
			unit_id,
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.UnitID, info.OrganizationID, info.Name, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateUnit(id string, info Unit) error {
	_, err := r.tx.Exec(`
		Update s_units SET
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE unit_id = ?
	`, info.Name, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteUnit(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_units SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE unit_id = ?
	`, time.Now(), byUser, id)
	return err
}

// Manufacturer

func (r *settingRepository) GetManufacturerByID(manufacturerID string) (*ManufacturerResponse, error) {
	var res ManufacturerResponse
	row := r.tx.QueryRow(`SELECT manufacturer_id, organization_id, name, status FROM s_manufacturers WHERE manufacturer_id = ? AND status > 0 LIMIT 1`, manufacturerID)
	err := row.Scan(&res.ManufacturerID, &res.OrganizationID, &res.Name, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckManufacturerConfict(manufacturerID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_manufacturers WHERE organization_id = ? AND manufacturer_id != ? AND name = ? AND status > 0", organizationID, manufacturerID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateManufacturer(info Manufacturer) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_manufacturers
		(
			manufacturer_id,
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.ManufacturerID, info.OrganizationID, info.Name, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateManufacturer(id string, info Manufacturer) error {
	_, err := r.tx.Exec(`
		Update s_manufacturers SET
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE manufacturer_id = ?
	`, info.Name, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteManufacturer(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_manufacturers SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE manufacturer_id = ?
	`, time.Now(), byUser, id)
	return err
}

// Brand

func (r *settingRepository) GetBrandByID(brandID string) (*BrandResponse, error) {
	var res BrandResponse
	row := r.tx.QueryRow(`SELECT brand_id, organization_id, name, status FROM s_brands WHERE brand_id = ? AND status > 0 LIMIT 1`, brandID)
	err := row.Scan(&res.BrandID, &res.OrganizationID, &res.Name, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckBrandConfict(brandID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_brands WHERE organization_id = ? AND brand_id != ? AND name = ? AND status > 0", organizationID, brandID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateBrand(info Brand) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_brands
		(
			brand_id,
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, info.BrandID, info.OrganizationID, info.Name, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateBrand(id string, info Brand) error {
	_, err := r.tx.Exec(`
		Update s_brands SET
		name = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE brand_id = ?
	`, info.Name, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteBrand(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_brands SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE brand_id = ?
	`, time.Now(), byUser, id)
	return err
}
