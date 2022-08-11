package warehouse

import (
	"errors"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type warehouseService struct {
}

func NewWarehouseService() *warehouseService {
	return &warehouseService{}
}

//Bay

func (s *warehouseService) GetBayList(filter BayFilter) (int, *[]BayResponse, error) {
	db := database.RDB()
	query := NewWarehouseQuery(db)
	count, err := query.GetBayCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBayList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *warehouseService) NewBay(info BayNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	isConflict, err := repo.CheckBayConfict("", info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "bay code conflict"
		return nil, errors.New(msg)
	}
	var bay Bay
	bay.BayID = "bay-" + xid.New().String()
	bay.OrganizationID = info.OrganizationID
	bay.Code = info.Code
	bay.Level = info.Level
	bay.Location = info.Location
	bay.Status = info.Status
	bay.Created = time.Now()
	bay.CreatedBy = info.User
	bay.Updated = time.Now()
	bay.UpdatedBy = info.User

	err = repo.CreateBay(bay)
	if err != nil {
		msg := "create bay error: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &bay.BayID, err
}

func (s *warehouseService) UpdateBay(bayID string, info BayNew) (*BayResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	isConflict, err := repo.CheckBayConfict(bayID, info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "bay conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetBayByID(bayID, info.OrganizationID)
	if err != nil {
		msg := "Bay not exist"
		return nil, errors.New(msg)
	}
	var bay Bay
	bay.Code = info.Code
	bay.Level = info.Level
	bay.Location = info.Location
	bay.Status = info.Status
	bay.Updated = time.Now()
	bay.UpdatedBy = info.User
	err = repo.UpdateBay(bayID, bay)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetBayByID(bayID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, err
}

func (s *warehouseService) GetBayByID(organizationID, id string) (*BayResponse, error) {
	db := database.RDB()
	query := NewWarehouseQuery(db)
	unit, err := query.GetBayByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *warehouseService) DeleteBay(bayID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	_, err = repo.GetBayByID(bayID, organizationID)
	if err != nil {
		msg := "Bay not exist"
		return errors.New(msg)
	}
	err = repo.DeleteBay(bayID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//Location

func (s *warehouseService) GetLocationList(filter LocationFilter) (int, *[]LocationResponse, error) {
	db := database.RDB()
	query := NewWarehouseQuery(db)
	count, err := query.GetLocationCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetLocationList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *warehouseService) NewLocation(info LocationNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	isConflict, err := repo.CheckLocationConfict("", info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "location code conflict"
		return nil, errors.New(msg)
	}
	var location Location
	location.LocationID = "location-" + xid.New().String()
	location.OrganizationID = info.OrganizationID
	location.Code = info.Code
	location.Level = info.Level
	location.Location = info.Location
	location.Status = info.Status
	location.Created = time.Now()
	location.CreatedBy = info.User
	location.Updated = time.Now()
	location.UpdatedBy = info.User

	err = repo.CreateLocation(location)
	if err != nil {
		msg := "create location error: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &location.LocationID, err
}

func (s *warehouseService) UpdateLocation(locationID string, info LocationNew) (*LocationResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	isConflict, err := repo.CheckLocationConfict(locationID, info.OrganizationID, info.Code)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "location conflict"
		return nil, errors.New(msg)
	}
	_, err = repo.GetLocationByID(locationID, info.OrganizationID)
	if err != nil {
		msg := "Location not exist"
		return nil, errors.New(msg)
	}
	var location Location
	location.Code = info.Code
	location.Level = info.Level
	location.Location = info.Location
	location.Status = info.Status
	location.Updated = time.Now()
	location.UpdatedBy = info.User
	err = repo.UpdateLocation(locationID, location)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetLocationByID(locationID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, err
}

func (s *warehouseService) GetLocationByID(organizationID, id string) (*LocationResponse, error) {
	db := database.RDB()
	query := NewWarehouseQuery(db)
	unit, err := query.GetLocationByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *warehouseService) DeleteLocation(locationID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewWarehouseRepository(tx)
	_, err = repo.GetLocationByID(locationID, organizationID)
	if err != nil {
		msg := "Location not exist"
		return errors.New(msg)
	}
	err = repo.DeleteLocation(locationID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}