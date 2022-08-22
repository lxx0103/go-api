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
	row := r.tx.QueryRow("SELECT count(1) FROM s_units WHERE organization_id = ? AND unit_id != ? AND name = ? AND status > 0 AND unit_type = ?", organizationID, unitID, name, "custom")
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
			unit_type,
			organization_id,
			name,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.UnitID, info.UnitType, info.OrganizationID, info.Name, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
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
		WHERE unit_id = ? AND unit_type = ?
	`, time.Now(), byUser, id, "custom")
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

// Vendor

func (r *settingRepository) GetVendorByID(vendorID string) (*VendorResponse, error) {
	var res VendorResponse
	row := r.tx.QueryRow(`
	SELECT 
	vendor_id,
	organization_id,
	name,
	contact_salutation,
	contact_first_name,
	contact_last_name,
	contact_email,
	contact_phone,
	country,
	state,
	city,
	address1,
	address2,
	zip,
	phone,
	fax,
	status
	FROM s_vendors 
	WHERE vendor_id = ? AND status > 0 LIMIT 1`, vendorID)
	err := row.Scan(&res.VendorID, &res.OrganizationID, &res.Name, &res.ContactSalutation, &res.ContactFirstName, &res.ContactLastName, &res.ContactEmail, &res.ContactPhone, &res.Country, &res.State, &res.City, &res.Address1, &res.Address2, &res.Zip, &res.Phone, &res.Fax, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckVendorConfict(vendorID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_vendors WHERE organization_id = ? AND vendor_id != ? AND name = ? AND status > 0", organizationID, vendorID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateVendor(info Vendor) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_vendors
		(
			vendor_id,
			organization_id,
			name,
			contact_salutation,
			contact_first_name,
			contact_last_name,
			contact_email,
			contact_phone,
			country,
			state,
			city,
			address1,
			address2,
			zip,
			phone,
			fax,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.VendorID, info.OrganizationID, info.Name, info.ContactSalutation, info.ContactFirstName, info.ContactLastName, info.ContactEmail, info.ContactPhone, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Phone, info.Fax, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateVendor(id string, info Vendor) error {
	_, err := r.tx.Exec(`
		Update s_vendors SET
		name = ?,
		contact_salutation = ?,
		contact_first_name = ?,
		contact_last_name = ?,
		contact_email = ?,
		contact_phone = ?,
		country = ?,
		state = ?,
		city = ?,
		address1 = ?,
		address2 = ?,
		zip = ?,
		phone = ?,
		fax = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE vendor_id = ?
	`, info.Name, info.ContactSalutation, info.ContactFirstName, info.ContactLastName, info.ContactEmail, info.ContactPhone, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Phone, info.Fax, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteVendor(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_vendors SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE vendor_id = ?
	`, time.Now(), byUser, id)
	return err
}

// Tax

func (r *settingRepository) GetTaxByID(taxID, organizationID string) (*TaxResponse, error) {
	var res TaxResponse
	row := r.tx.QueryRow(`
	SELECT 
	tax_id,
	organization_id,
	name,
	tax_value,
	status
	FROM s_taxes 
	WHERE tax_id = ? AND organization_id = ? AND status > 0 LIMIT 1`, taxID, organizationID)
	err := row.Scan(&res.TaxID, &res.OrganizationID, &res.Name, &res.TaxValue, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckTaxConfict(taxID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_taxes WHERE organization_id = ? AND tax_id != ? AND name = ? AND status > 0", organizationID, taxID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateTax(info Tax) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_taxes
		(
			tax_id,
			organization_id,
			name,
			tax_value,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.TaxID, info.OrganizationID, info.Name, info.TaxValue, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateTax(id string, info Tax) error {
	_, err := r.tx.Exec(`
		Update s_taxes SET
		name = ?,
		tax_value = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE tax_id = ?
	`, info.Name, info.TaxValue, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteTax(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_taxes SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE tax_id = ?
	`, time.Now(), byUser, id)
	return err
}

// Customer

func (r *settingRepository) GetCustomerByID(customerID string) (*CustomerResponse, error) {
	var res CustomerResponse
	row := r.tx.QueryRow(`
	SELECT 
	customer_id,
	organization_id,
	name,
	contact_salutation,
	contact_first_name,
	contact_last_name,
	contact_email,
	contact_phone,
	country,
	state,
	city,
	address1,
	address2,
	zip,
	phone,
	fax,
	status
	FROM s_customers 
	WHERE customer_id = ? AND status > 0 LIMIT 1`, customerID)
	err := row.Scan(&res.CustomerID, &res.OrganizationID, &res.Name, &res.ContactSalutation, &res.ContactFirstName, &res.ContactLastName, &res.ContactEmail, &res.ContactPhone, &res.Country, &res.State, &res.City, &res.Address1, &res.Address2, &res.Zip, &res.Phone, &res.Fax, &res.Status)
	return &res, err
}

func (r *settingRepository) CheckCustomerConfict(customerID, organizationID, name string) (bool, error) {
	var existed int
	row := r.tx.QueryRow("SELECT count(1) FROM s_customers WHERE organization_id = ? AND customer_id != ? AND name = ? AND status > 0", organizationID, customerID, name)
	err := row.Scan(&existed)
	if err != nil {
		return true, err
	}
	return existed != 0, nil
}

func (r *settingRepository) CreateCustomer(info Customer) error {
	_, err := r.tx.Exec(`
		INSERT INTO s_customers
		(
			customer_id,
			organization_id,
			name,
			contact_salutation,
			contact_first_name,
			contact_last_name,
			contact_email,
			contact_phone,
			country,
			state,
			city,
			address1,
			address2,
			zip,
			phone,
			fax,
			status,
			created,
			created_by,
			updated,
			updated_by
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, info.CustomerID, info.OrganizationID, info.Name, info.ContactSalutation, info.ContactFirstName, info.ContactLastName, info.ContactEmail, info.ContactPhone, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Phone, info.Fax, info.Status, info.Created, info.CreatedBy, info.Updated, info.UpdatedBy)
	return err
}

func (r *settingRepository) UpdateCustomer(id string, info Customer) error {
	_, err := r.tx.Exec(`
		Update s_customers SET
		name = ?,
		contact_salutation = ?,
		contact_first_name = ?,
		contact_last_name = ?,
		contact_email = ?,
		contact_phone = ?,
		country = ?,
		state = ?,
		city = ?,
		address1 = ?,
		address2 = ?,
		zip = ?,
		phone = ?,
		fax = ?,
		status = ?,
		updated = ?,
		updated_by = ?
		WHERE customer_id = ?
	`, info.Name, info.ContactSalutation, info.ContactFirstName, info.ContactLastName, info.ContactEmail, info.ContactPhone, info.Country, info.State, info.City, info.Address1, info.Address2, info.Zip, info.Phone, info.Fax, info.Status, info.Updated, info.UpdatedBy, id)
	return err
}

func (r *settingRepository) DeleteCustomer(id, byUser string) error {
	_, err := r.tx.Exec(`
		Update s_customers SET
		status = -1,
		updated = ?,
		updated_by = ?
		WHERE customer_id = ?
	`, time.Now(), byUser, id)
	return err
}
