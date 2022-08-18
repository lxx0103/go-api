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
	unit.Status = info.Status
	unit.Created = time.Now()
	unit.CreatedBy = info.Name
	unit.Updated = time.Now()
	unit.UpdatedBy = info.Name
	err = repo.CreateUnit(unit)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetUnitByID(unit.UnitID)
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
	oldUnit, err := repo.GetUnitByID(unitID)
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
	unit.UpdatedBy = info.Name
	unit.Updated = time.Now()
	unit.Status = info.Status
	err = repo.UpdateUnit(unitID, unit)
	if err != nil {
		msg := "update unit error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetUnitByID(unitID)
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
	oldUnit, err := repo.GetUnitByID(unitID)
	if err != nil {
		msg := "Unit not exist"
		return errors.New(msg)
	}
	if oldUnit.OrganizationID != organizationID {
		msg := "Unit not exist"
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
	manufacturer.CreatedBy = info.Name
	manufacturer.Updated = time.Now()
	manufacturer.UpdatedBy = info.Name
	err = repo.CreateManufacturer(manufacturer)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetManufacturerByID(manufacturer.ManufacturerID)
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
	oldManufacturer, err := repo.GetManufacturerByID(manufacturerID)
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
	manufacturer.UpdatedBy = info.Name
	manufacturer.Updated = time.Now()
	manufacturer.Status = info.Status
	err = repo.UpdateManufacturer(manufacturerID, manufacturer)
	if err != nil {
		msg := "update manufacturer error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetManufacturerByID(manufacturerID)
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
	oldManufacturer, err := repo.GetManufacturerByID(manufacturerID)
	if err != nil {
		msg := "Manufacturer not exist"
		return errors.New(msg)
	}
	if oldManufacturer.OrganizationID != organizationID {
		msg := "Manufacturer not exist"
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
	brand.CreatedBy = info.Name
	brand.Updated = time.Now()
	brand.UpdatedBy = info.Name
	err = repo.CreateBrand(brand)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetBrandByID(brand.BrandID)
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
	oldBrand, err := repo.GetBrandByID(brandID)
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
	brand.UpdatedBy = info.Name
	brand.Updated = time.Now()
	brand.Status = info.Status
	err = repo.UpdateBrand(brandID, brand)
	if err != nil {
		msg := "update brand error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetBrandByID(brandID)
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
	oldBrand, err := repo.GetBrandByID(brandID)
	if err != nil {
		msg := "Brand not exist"
		return errors.New(msg)
	}
	if oldBrand.OrganizationID != organizationID {
		msg := "Brand not exist"
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
	vendor.CreatedBy = info.Name
	vendor.Updated = time.Now()
	vendor.UpdatedBy = info.Name
	err = repo.CreateVendor(vendor)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetVendorByID(vendor.VendorID)
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
	oldVendor, err := repo.GetVendorByID(vendorID)
	if err != nil {
		msg := "Vendor not exist"
		return nil, errors.New(msg)
	}
	if oldVendor.OrganizationID != info.OrganizationID {
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
	vendor.UpdatedBy = info.Name
	vendor.Updated = time.Now()
	vendor.Status = info.Status
	err = repo.UpdateVendor(vendorID, vendor)
	if err != nil {
		msg := "update vendor error"
		return nil, errors.New(msg)
	}
	res, err := repo.GetVendorByID(vendorID)
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
	oldVendor, err := repo.GetVendorByID(vendorID)
	if err != nil {
		msg := "Vendor not exist"
		return errors.New(msg)
	}
	if oldVendor.OrganizationID != organizationID {
		msg := "Vendor not exist"
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
	err = repo.DeleteTax(taxID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
