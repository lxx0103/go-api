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
	if v := filter.UnitType; v != "" {
		where, args = append(where, "unit_type = ?"), append(args, v)
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
	if v := filter.UnitType; v != "" {
		where, args = append(where, "unit_type = ?"), append(args, v)
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

//Customer

func (r *settingQuery) GetCustomerByID(organizationID, id string) (*CustomerResponse, error) {
	var customer CustomerResponse
	err := r.conn.Get(&customer, "SELECT customer_id, organization_id, name, contact_salutation, contact_first_name, contact_last_name, contact_email, contact_phone, country, state, city, address1, address2, zip, phone, fax, status FROM s_customers WHERE organization_id = ? AND customer_id = ? AND status > 0", organizationID, id)
	return &customer, err
}

func (r *settingQuery) GetCustomerCount(filter CustomerFilter) (int, error) {
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
		FROM s_customers
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetCustomerList(filter CustomerFilter) (*[]CustomerResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var customers []CustomerResponse
	err := r.conn.Select(&customers, `
		SELECT customer_id, organization_id, name, contact_salutation, contact_first_name, contact_last_name, contact_email, contact_phone, country, state, city, address1, address2, zip, phone, fax, status
		FROM s_customers
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &customers, err
}

//Carrier

func (r *settingQuery) GetCarrierByID(organizationID, id string) (*CarrierResponse, error) {
	var carrier CarrierResponse
	err := r.conn.Get(&carrier, "SELECT carrier_id, organization_id, name, status FROM s_carriers WHERE organization_id = ? AND carrier_id = ? AND status > 0", organizationID, id)
	return &carrier, err
}

func (r *settingQuery) GetCarrierCount(filter CarrierFilter) (int, error) {
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
		FROM s_carriers
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetCarrierList(filter CarrierFilter) (*[]CarrierResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var carriers []CarrierResponse
	err := r.conn.Select(&carriers, `
		SELECT carrier_id, organization_id, name, status
		FROM s_carriers
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &carriers, err
}

//AdjustmentReason

func (r *settingQuery) GetAdjustmentReasonByID(organizationID, id string) (*AdjustmentReasonResponse, error) {
	var adjustmentReason AdjustmentReasonResponse
	err := r.conn.Get(&adjustmentReason, "SELECT adjustment_reason_id, organization_id, name, status FROM s_adjustment_reasons WHERE organization_id = ? AND adjustment_reason_id = ? AND status > 0", organizationID, id)
	return &adjustmentReason, err
}

func (r *settingQuery) GetAdjustmentReasonCount(filter AdjustmentReasonFilter) (int, error) {
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
		FROM s_adjustment_reasons
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetAdjustmentReasonList(filter AdjustmentReasonFilter) (*[]AdjustmentReasonResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var adjustmentReasons []AdjustmentReasonResponse
	err := r.conn.Select(&adjustmentReasons, `
		SELECT adjustment_reason_id, organization_id, name, status
		FROM s_adjustment_reasons
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &adjustmentReasons, err
}

//PaymentMethod

func (r *settingQuery) GetPaymentMethodByID(organizationID, id string) (*PaymentMethodResponse, error) {
	var paymentMethod PaymentMethodResponse
	err := r.conn.Get(&paymentMethod, "SELECT payment_method_id, organization_id, name, status FROM s_payment_methods WHERE organization_id = ? AND payment_method_id = ? AND status > 0", organizationID, id)
	return &paymentMethod, err
}

func (r *settingQuery) GetPaymentMethodCount(filter PaymentMethodFilter) (int, error) {
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
		FROM s_payment_methods
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *settingQuery) GetPaymentMethodList(filter PaymentMethodFilter) (*[]PaymentMethodResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var paymentMethods []PaymentMethodResponse
	err := r.conn.Select(&paymentMethods, `
		SELECT payment_method_id, organization_id, name, status
		FROM s_payment_methods
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &paymentMethods, err
}
