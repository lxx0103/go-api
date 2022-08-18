package setting

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type settingQuery struct {
	conn *sqlx.DB
}

func NewSettingQuery(connection *sqlx.DB) *settingQuery {
	return &settingQuery{
		conn: connection,
	}
}

func (r *settingQuery) GetUnitByID(organizationID, id string) (*UnitResponse, error) {
	var unit UnitResponse
	err := r.conn.Get(&unit, "SELECT unit_id, organization_id, name, status FROM s_units WHERE organization_id = ? AND unit_id = ? AND status > 0", organizationID, id)
	return &unit, err
}

func (r *settingQuery) GetUnitCount(filter UnitFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_units
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetUnitList(filter UnitFilter) (*[]UnitResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var units []UnitResponse
	err := r.conn.Select(&units, `
		SELECT unit_id, organization_id, name, status
		FROM s_units
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &units, err
}

//Manufacturer

func (r *settingQuery) GetManufacturerByID(organizationID, id string) (*ManufacturerResponse, error) {
	var manufacturer ManufacturerResponse
	err := r.conn.Get(&manufacturer, "SELECT manufacturer_id, organization_id, name, status FROM s_manufacturers WHERE organization_id = ? AND manufacturer_id = ? AND status > 0", organizationID, id)
	return &manufacturer, err
}

func (r *settingQuery) GetManufacturerCount(filter ManufacturerFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_manufacturers
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetManufacturerList(filter ManufacturerFilter) (*[]ManufacturerResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var manufacturers []ManufacturerResponse
	err := r.conn.Select(&manufacturers, `
		SELECT manufacturer_id, organization_id, name, status
		FROM s_manufacturers
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &manufacturers, err
}

//Brand

func (r *settingQuery) GetBrandByID(organizationID, id string) (*BrandResponse, error) {
	var brand BrandResponse
	err := r.conn.Get(&brand, "SELECT brand_id, organization_id, name, status FROM s_brands WHERE organization_id = ? AND brand_id = ? AND status > 0", organizationID, id)
	return &brand, err
}

func (r *settingQuery) GetBrandCount(filter BrandFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_brands
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetBrandList(filter BrandFilter) (*[]BrandResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var brands []BrandResponse
	err := r.conn.Select(&brands, `
		SELECT brand_id, organization_id, name, status
		FROM s_brands
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &brands, err
}

//Vendor

func (r *settingQuery) GetVendorByID(organizationID, id string) (*VendorResponse, error) {
	var vendor VendorResponse
	err := r.conn.Get(&vendor, "SELECT vendor_id, organization_id, name, contact_salutation, contact_first_name, contact_last_name, contact_email, contact_phone, country, state, city, address1, address2, zip, phone, fax, status FROM s_vendors WHERE organization_id = ? AND vendor_id = ? AND status > 0", organizationID, id)
	return &vendor, err
}

func (r *settingQuery) GetVendorCount(filter VendorFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_vendors
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetVendorList(filter VendorFilter) (*[]VendorResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var vendors []VendorResponse
	err := r.conn.Select(&vendors, `
		SELECT vendor_id, organization_id, name, contact_salutation, contact_first_name, contact_last_name, contact_email, contact_phone, country, state, city, address1, address2, zip, phone, fax, status
		FROM s_vendors
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &vendors, err
}

//Tax

func (r *settingQuery) GetTaxByID(organizationID, id string) (*TaxResponse, error) {
	var tax TaxResponse
	err := r.conn.Get(&tax, "SELECT tax_id, organization_id, name, tax_value, status FROM s_taxes WHERE organization_id = ? AND tax_id = ? AND status > 0", organizationID, id)
	return &tax, err
}

func (r *settingQuery) GetTaxCount(filter TaxFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM s_taxes
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetTaxList(filter TaxFilter) (*[]TaxResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var taxes []TaxResponse
	err := r.conn.Select(&taxes, `
		SELECT tax_id, organization_id, name, tax_value, status
		FROM s_taxes
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &taxes, err
}
