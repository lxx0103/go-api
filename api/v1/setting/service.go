package setting

import (
	"errors"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type settingService struct {
}

func NewSettingService() *settingService {
	return &settingService{}
}

func (s *settingService) GetUnitByID(organizationID, id string) (*UnitResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	unit, err := query.GetUnitByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *settingService) NewUnit(info UnitNew) (*UnitResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckUnitConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Unit name conflict"
		return nil, errors.New(msg)
	}
	var unit Unit
	unit.UnitID = "unit-" + xid.New().String()
	unit.OrganizationID = info.OrganizationID
	unit.Name = info.Name
	unit.UnitType = "custom"
	unit.Status = info.Status
	unit.Created = time.Now()
	unit.CreatedBy = info.User
	unit.Updated = time.Now()
	unit.UpdatedBy = info.User
	err = repo.CreateUnit(unit)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetUnitByID(unit.UnitID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetUnitList(filter UnitFilter) (int, *[]UnitResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetUnitCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetUnitList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateUnit(unitID string, info UnitNew) (*UnitResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckUnitConfict(unitID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "unit name conflict"
		return nil, errors.New(msg)
	}
	oldUnit, err := repo.GetUnitByID(unitID, info.OrganizationID)
	if err != nil {
		msg := "Unit not exist"
		return nil, errors.New(msg)
	}
	if oldUnit.OrganizationID != info.OrganizationID {
		msg := "Unit not exist"
		return nil, errors.New(msg)
	}
	var unit Unit
	unit.Name = info.Name
	unit.UpdatedBy = info.User
	unit.Updated = time.Now()
	unit.Status = info.Status
	err = repo.UpdateUnit(unitID, unit)
	if err != nil {
		msg := "update unit error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetUnitByID(unitID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteUnit(unitID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetUnitByID(unitID, organizationID)
	if err != nil {
		msg := "Unit not exist"
		return errors.New(msg)
	}
	itemUnitCount, err := repo.GetItemUnitCount(unitID, organizationID)
	if err != nil {
		msg := "get item unit count error"
		return errors.New(msg)
	}
	if itemUnitCount > 0 {
		msg := "Unit used by item can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteUnit(unitID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *settingService) GetManufacturerByID(organizationID, id string) (*ManufacturerResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	manufacturer, err := query.GetManufacturerByID(organizationID, id)
	if err != nil {
		msg := "get manufacturer error: " + err.Error()
		return nil, errors.New(msg)
	}
	return manufacturer, nil
}

func (s *settingService) NewManufacturer(info ManufacturerNew) (*ManufacturerResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckManufacturerConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Manufacturer name conflict"
		return nil, errors.New(msg)
	}
	var manufacturer Manufacturer
	manufacturer.ManufacturerID = "man-" + xid.New().String()
	manufacturer.OrganizationID = info.OrganizationID
	manufacturer.Name = info.Name
	manufacturer.Status = info.Status
	manufacturer.Created = time.Now()
	manufacturer.CreatedBy = info.User
	manufacturer.Updated = time.Now()
	manufacturer.UpdatedBy = info.User
	err = repo.CreateManufacturer(manufacturer)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetManufacturerByID(manufacturer.ManufacturerID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetManufacturerList(filter ManufacturerFilter) (int, *[]ManufacturerResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetManufacturerCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetManufacturerList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateManufacturer(manufacturerID string, info ManufacturerNew) (*ManufacturerResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckManufacturerConfict(manufacturerID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "manufacturer name conflict"
		return nil, errors.New(msg)
	}
	oldManufacturer, err := repo.GetManufacturerByID(manufacturerID, info.OrganizationID)
	if err != nil {
		msg := "Manufacturer not exist"
		return nil, errors.New(msg)
	}
	if oldManufacturer.OrganizationID != info.OrganizationID {
		msg := "Manufacturer not exist"
		return nil, errors.New(msg)
	}
	var manufacturer Manufacturer
	manufacturer.Name = info.Name
	manufacturer.UpdatedBy = info.User
	manufacturer.Updated = time.Now()
	manufacturer.Status = info.Status
	err = repo.UpdateManufacturer(manufacturerID, manufacturer)
	if err != nil {
		msg := "update manufacturer error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetManufacturerByID(manufacturerID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteManufacturer(manufacturerID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetManufacturerByID(manufacturerID, organizationID)
	if err != nil {
		msg := "Manufacturer not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetItemManufacturerCount(manufacturerID, organizationID)
	if err != nil {
		msg := "get item unit count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Manufacturer used by item can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteManufacturer(manufacturerID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//brand

func (s *settingService) GetBrandByID(organizationID, id string) (*BrandResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	brand, err := query.GetBrandByID(organizationID, id)
	if err != nil {
		msg := "get brand error: " + err.Error()
		return nil, errors.New(msg)
	}
	return brand, nil
}

func (s *settingService) NewBrand(info BrandNew) (*BrandResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckBrandConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Brand name conflict"
		return nil, errors.New(msg)
	}
	var brand Brand
	brand.BrandID = "brn-" + xid.New().String()
	brand.OrganizationID = info.OrganizationID
	brand.Name = info.Name
	brand.Status = info.Status
	brand.Created = time.Now()
	brand.CreatedBy = info.User
	brand.Updated = time.Now()
	brand.UpdatedBy = info.User
	err = repo.CreateBrand(brand)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetBrandByID(brand.BrandID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetBrandList(filter BrandFilter) (int, *[]BrandResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetBrandCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBrandList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateBrand(brandID string, info BrandNew) (*BrandResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckBrandConfict(brandID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "brand name conflict"
		return nil, errors.New(msg)
	}
	oldBrand, err := repo.GetBrandByID(brandID, info.OrganizationID)
	if err != nil {
		msg := "Brand not exist"
		return nil, errors.New(msg)
	}
	if oldBrand.OrganizationID != info.OrganizationID {
		msg := "Brand not exist"
		return nil, errors.New(msg)
	}
	var brand Brand
	brand.Name = info.Name
	brand.UpdatedBy = info.User
	brand.Updated = time.Now()
	brand.Status = info.Status
	err = repo.UpdateBrand(brandID, brand)
	if err != nil {
		msg := "update brand error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetBrandByID(brandID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteBrand(brandID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetBrandByID(brandID, organizationID)
	if err != nil {
		msg := "Brand not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetItemBrandCount(brandID, organizationID)
	if err != nil {
		msg := "get item unit count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Brand used by item can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteBrand(brandID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//vendor

func (s *settingService) GetVendorByID(organizationID, id string) (*VendorResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	vendor, err := query.GetVendorByID(organizationID, id)
	if err != nil {
		msg := "get vendor error: " + err.Error()
		return nil, errors.New(msg)
	}
	return vendor, nil
}

func (s *settingService) NewVendor(info VendorNew) (*VendorResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckVendorConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Vendor name conflict"
		return nil, errors.New(msg)
	}
	var vendor Vendor
	vendor.VendorID = "ven-" + xid.New().String()
	vendor.OrganizationID = info.OrganizationID
	vendor.Name = info.Name
	vendor.ContactSalutation = info.ContactSalutation
	vendor.ContactFirstName = info.ContactFirstName
	vendor.ContactLastName = info.ContactLastName
	vendor.ContactEmail = info.ContactEmail
	vendor.ContactPhone = info.ContactPhone
	vendor.Country = info.Country
	vendor.State = info.State
	vendor.City = info.City
	vendor.Address1 = info.Address1
	vendor.Address2 = info.Address2
	vendor.Zip = info.Zip
	vendor.Phone = info.Phone
	vendor.Fax = info.Fax
	vendor.Status = info.Status
	vendor.Created = time.Now()
	vendor.CreatedBy = info.User
	vendor.Updated = time.Now()
	vendor.UpdatedBy = info.User
	err = repo.CreateVendor(vendor)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetVendorByID(vendor.VendorID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetVendorList(filter VendorFilter) (int, *[]VendorResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetVendorCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetVendorList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateVendor(vendorID string, info VendorNew) (*VendorResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckVendorConfict(vendorID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "vendor name conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetVendorByID(vendorID, info.OrganizationID)
	if err != nil {
		msg := "Vendor not exist"
		return nil, errors.New(msg)
	}
	var vendor Vendor
	vendor.Name = info.Name
	vendor.ContactSalutation = info.ContactSalutation
	vendor.ContactFirstName = info.ContactFirstName
	vendor.ContactLastName = info.ContactLastName
	vendor.ContactEmail = info.ContactEmail
	vendor.ContactPhone = info.ContactPhone
	vendor.Country = info.Country
	vendor.State = info.State
	vendor.City = info.City
	vendor.Address1 = info.Address1
	vendor.Address2 = info.Address2
	vendor.Zip = info.Zip
	vendor.Phone = info.Phone
	vendor.Fax = info.Fax
	vendor.UpdatedBy = info.User
	vendor.Updated = time.Now()
	vendor.Status = info.Status
	err = repo.UpdateVendor(vendorID, vendor)
	if err != nil {
		msg := "update vendor error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetVendorByID(vendorID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteVendor(vendorID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetVendorByID(vendorID, organizationID)
	if err != nil {
		msg := "Vendor not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetPOVendorCount(vendorID, organizationID)
	if err != nil {
		msg := "get po vendor count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Vendor used in Purchase order can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteVendor(vendorID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//tax

func (s *settingService) GetTaxByID(organizationID, id string) (*TaxResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	tax, err := query.GetTaxByID(organizationID, id)
	if err != nil {
		msg := "get tax error: " + err.Error()
		return nil, errors.New(msg)
	}
	return tax, nil
}

func (s *settingService) NewTax(info TaxNew) (*TaxResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckTaxConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Tax name conflict"
		return nil, errors.New(msg)
	}
	var tax Tax
	tax.TaxID = "tax-" + xid.New().String()
	tax.OrganizationID = info.OrganizationID
	tax.Name = info.Name
	tax.TaxValue = info.TaxValue
	tax.Status = info.Status
	tax.Created = time.Now()
	tax.CreatedBy = info.User
	tax.Updated = time.Now()
	tax.UpdatedBy = info.User
	err = repo.CreateTax(tax)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetTaxByID(tax.TaxID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetTaxList(filter TaxFilter) (int, *[]TaxResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetTaxCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetTaxList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateTax(taxID string, info TaxNew) (*TaxResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckTaxConfict(taxID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "tax name conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetTaxByID(taxID, info.OrganizationID)
	if err != nil {
		msg := "Tax not exist"
		return nil, errors.New(msg)
	}
	var tax Tax
	tax.Name = info.Name
	tax.TaxValue = info.TaxValue
	tax.UpdatedBy = info.User
	tax.Updated = time.Now()
	tax.Status = info.Status
	err = repo.UpdateTax(taxID, tax)
	if err != nil {
		msg := "update tax error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetTaxByID(taxID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteTax(taxID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetTaxByID(taxID, organizationID)
	if err != nil {
		msg := "Tax not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetPOTaxCount(taxID, organizationID)
	if err != nil {
		msg := "get po tax count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Tax used in Purchase order can not be deleted"
		return errors.New(msg)
	}
	usedCount, err = repo.GetSOTaxCount(taxID, organizationID)
	if err != nil {
		msg := "get so tax count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Tax used in Sales order can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteTax(taxID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//customer

func (s *settingService) GetCustomerByID(organizationID, id string) (*CustomerResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	customer, err := query.GetCustomerByID(organizationID, id)
	if err != nil {
		msg := "get customer error: " + err.Error()
		return nil, errors.New(msg)
	}
	return customer, nil
}

func (s *settingService) NewCustomer(info CustomerNew) (*CustomerResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckCustomerConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Customer name conflict"
		return nil, errors.New(msg)
	}
	var customer Customer
	customer.CustomerID = "cus-" + xid.New().String()
	customer.OrganizationID = info.OrganizationID
	customer.Name = info.Name
	customer.ContactSalutation = info.ContactSalutation
	customer.ContactFirstName = info.ContactFirstName
	customer.ContactLastName = info.ContactLastName
	customer.ContactEmail = info.ContactEmail
	customer.ContactPhone = info.ContactPhone
	customer.Country = info.Country
	customer.State = info.State
	customer.City = info.City
	customer.Address1 = info.Address1
	customer.Address2 = info.Address2
	customer.Zip = info.Zip
	customer.Phone = info.Phone
	customer.Fax = info.Fax
	customer.Status = info.Status
	customer.Created = time.Now()
	customer.CreatedBy = info.User
	customer.Updated = time.Now()
	customer.UpdatedBy = info.User
	err = repo.CreateCustomer(customer)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetCustomerByID(customer.CustomerID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetCustomerList(filter CustomerFilter) (int, *[]CustomerResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetCustomerCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetCustomerList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateCustomer(customerID string, info CustomerNew) (*CustomerResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckCustomerConfict(customerID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "customer name conflict"
		return nil, errors.New(msg)
	}
	oldCustomer, err := repo.GetCustomerByID(customerID, info.OrganizationID)
	if err != nil {
		msg := "Customer not exist"
		return nil, errors.New(msg)
	}
	if oldCustomer.OrganizationID != info.OrganizationID {
		msg := "Customer not exist"
		return nil, errors.New(msg)
	}
	var customer Customer
	customer.Name = info.Name
	customer.ContactSalutation = info.ContactSalutation
	customer.ContactFirstName = info.ContactFirstName
	customer.ContactLastName = info.ContactLastName
	customer.ContactEmail = info.ContactEmail
	customer.ContactPhone = info.ContactPhone
	customer.Country = info.Country
	customer.State = info.State
	customer.City = info.City
	customer.Address1 = info.Address1
	customer.Address2 = info.Address2
	customer.Zip = info.Zip
	customer.Phone = info.Phone
	customer.Fax = info.Fax
	customer.UpdatedBy = info.User
	customer.Updated = time.Now()
	customer.Status = info.Status
	err = repo.UpdateCustomer(customerID, customer)
	if err != nil {
		msg := "update customer error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetCustomerByID(customerID, info.OrganizationID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteCustomer(customerID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetCustomerByID(customerID, organizationID)
	if err != nil {
		msg := "Customer not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetSOCustomerCount(customerID, organizationID)
	if err != nil {
		msg := "get so customer count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Customer used in Sales order can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteCustomer(customerID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//carrier

func (s *settingService) GetCarrierByID(organizationID, id string) (*CarrierResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	carrier, err := query.GetCarrierByID(organizationID, id)
	if err != nil {
		msg := "get carrier error: " + err.Error()
		return nil, errors.New(msg)
	}
	return carrier, nil
}

func (s *settingService) NewCarrier(info CarrierNew) (*CarrierResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckCarrierConfict("", info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "Carrier name conflict"
		return nil, errors.New(msg)
	}
	var carrier Carrier
	carrier.CarrierID = "car-" + xid.New().String()
	carrier.OrganizationID = info.OrganizationID
	carrier.Name = info.Name
	carrier.Status = info.Status
	carrier.Created = time.Now()
	carrier.CreatedBy = info.User
	carrier.Updated = time.Now()
	carrier.UpdatedBy = info.User
	err = repo.CreateCarrier(carrier)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetCarrierByID(info.OrganizationID, carrier.CarrierID)
	tx.Commit()
	return res, err
}

func (s *settingService) GetCarrierList(filter CarrierFilter) (int, *[]CarrierResponse, error) {
	db := database.RDB()
	query := NewSettingQuery(db)
	count, err := query.GetCarrierCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetCarrierList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *settingService) UpdateCarrier(carrierID string, info CarrierNew) (*CarrierResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	isConflict, err := repo.CheckCarrierConfict(carrierID, info.OrganizationID, info.Name)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "carrier name conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetCarrierByID(info.OrganizationID, carrierID)
	if err != nil {
		msg := "Carrier not exist"
		return nil, errors.New(msg)
	}
	var carrier Carrier
	carrier.Name = info.Name
	carrier.UpdatedBy = info.User
	carrier.Updated = time.Now()
	carrier.Status = info.Status
	err = repo.UpdateCarrier(carrierID, carrier)
	if err != nil {
		msg := "update carrier error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetCarrierByID(info.OrganizationID, carrierID)
	tx.Commit()
	return res, err
}

func (s *settingService) DeleteCarrier(carrierID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewSettingRepository(tx)
	_, err = repo.GetCarrierByID(organizationID, carrierID)
	if err != nil {
		msg := "Carrier not exist"
		return errors.New(msg)
	}
	usedCount, err := repo.GetShppingCarrierCount(carrierID, organizationID)
	if err != nil {
		msg := "get so customer count error"
		return errors.New(msg)
	}
	if usedCount > 0 {
		msg := "Carrier used in Shipping order can not be deleted"
		return errors.New(msg)
	}
	err = repo.DeleteCarrier(carrierID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
