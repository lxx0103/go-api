package item

import (
	"errors"
	"go-api/api/v1/setting"
	"go-api/core/database"
	"time"

	"github.com/rs/xid"
)

type itemService struct {
}

func NewItemService() *itemService {
	return &itemService{}
}

func (s *itemService) NewItem(info ItemNew) (*string, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		msg := "begin transaction error"
		return nil, errors.New(msg)
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckSKUConfict("", info.OrganizationID, info.SKU)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "item SKU exists"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetUnitByID(info.OrganizationID, info.UnitID)
	if err != nil {
		return nil, err
	}
	// _, err = settingService.GetBrandByID(info.BrandID)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = settingService.GetManufacturerByID(info.ManufacturerID)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = settingService.GetVendorByID(info.DefaultVendorID)
	// if err != nil {
	// 	return nil, err
	// }
	var item Item
	item.ItemID = "item-" + xid.New().String()
	item.OrganizationID = info.OrganizationID
	item.SKU = info.SKU
	item.Name = info.Name
	item.UnitID = info.UnitID
	item.ManufacturerID = info.ManufacturerID
	item.BrandID = info.BrandID
	item.WeightUnit = info.WeightUnit
	item.Weight = info.Weight
	item.DimensionUnit = info.DimensionUnit
	item.Length = info.Length
	item.Width = info.Width
	item.Height = info.Height
	item.SellingPrice = info.SellingPrice
	item.CostPrice = info.CostPrice
	item.OpenningStock = info.OpenningStock
	item.OpenningStockRate = info.OpenningStockRate
	item.ReorderStock = info.ReorderStock
	item.DefaultVendorID = info.DefaultVendorID
	item.Description = info.Description
	item.Status = info.Status
	item.Created = time.Now()
	item.CreatedBy = info.User
	item.Updated = time.Now()
	item.UpdatedBy = info.User

	err = repo.CreateItem(item)
	if err != nil {
		msg := "create itemerror: " + err.Error()
		return nil, errors.New(msg)
	}
	tx.Commit()
	return &item.ItemID, err
}

func (s *itemService) GetItemList(filter ItemFilter) (int, *[]ItemResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	count, err := query.GetItemCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetItemList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *itemService) UpdateItem(itemID string, info ItemNew) (*ItemResponse, error) {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	isConflict, err := repo.CheckSKUConfict(itemID, info.OrganizationID, info.SKU)
	if err != nil {
		msg := "check conflict error: " + err.Error()
		return nil, errors.New(msg)
	}
	if isConflict {
		msg := "item SKU conflict"
		return nil, errors.New(msg)
	}
	settingService := setting.NewSettingService()
	_, err = settingService.GetUnitByID(info.OrganizationID, info.UnitID)
	if err != nil {
		return nil, err
	}
	_, err = settingService.GetBrandByID(info.OrganizationID, info.BrandID)
	if err != nil {
		return nil, err
	}
	_, err = settingService.GetManufacturerByID(info.OrganizationID, info.ManufacturerID)
	if err != nil {
		return nil, err
	}
	_, err = settingService.GetVendorByID(info.OrganizationID, info.DefaultVendorID)
	if err != nil {
		return nil, err
	}
	oldItem, err := repo.GetItemByID(itemID)
	if err != nil {
		msg := "Item not exist"
		return nil, errors.New(msg)
	}
	if oldItem.OrganizationID != info.OrganizationID {
		msg := "Item not exist"
		return nil, errors.New(msg)
	}
	var item Item
	item.SKU = info.SKU
	item.Name = info.Name
	item.UnitID = info.UnitID
	item.ManufacturerID = info.ManufacturerID
	item.BrandID = info.BrandID
	item.WeightUnit = info.WeightUnit
	item.Weight = info.Weight
	item.DimensionUnit = info.DimensionUnit
	item.Length = info.Length
	item.Width = info.Width
	item.Height = info.Height
	item.SellingPrice = info.SellingPrice
	item.CostPrice = info.CostPrice
	item.OpenningStock = info.OpenningStock
	item.OpenningStockRate = info.OpenningStockRate
	item.ReorderStock = info.ReorderStock
	item.DefaultVendorID = info.DefaultVendorID
	item.Description = info.Description
	item.Status = info.Status
	item.Updated = time.Now()
	item.UpdatedBy = info.User
	err = repo.UpdateItem(itemID, item)
	if err != nil {
		return nil, err
	}
	res, err := repo.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, err
}

func (s *itemService) GetItemByID(organizationID, id string) (*ItemResponse, error) {
	db := database.RDB()
	query := NewItemQuery(db)
	unit, err := query.GetItemByID(organizationID, id)
	if err != nil {
		msg := "get unit error: " + err.Error()
		return nil, errors.New(msg)
	}
	return unit, nil
}

func (s *itemService) DeleteItem(itemID, organizationID, user string) error {
	db := database.WDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewItemRepository(tx)
	oldItem, err := repo.GetItemByID(itemID)
	if err != nil {
		msg := "Item not exist"
		return errors.New(msg)
	}
	if oldItem.OrganizationID != organizationID {
		msg := "Item not exist"
		return errors.New(msg)
	}
	err = repo.DeleteItem(itemID, user)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
