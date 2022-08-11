package warehouse

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type warehouseQuery struct {
	conn *sqlx.DB
}

func NewWarehouseQuery(connection *sqlx.DB) *warehouseQuery {
	return &warehouseQuery{
		conn: connection,
	}
}

//Bay
func (r *warehouseQuery) GetBayCount(filter BayFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Code; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM w_bays
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *warehouseQuery) GetBayList(filter BayFilter) (*[]BayResponse, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.Code; v != "" {
		where, args = append(where, "code like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var bays []BayResponse
	err := r.conn.Select(&bays, `
		SELECT 
		bay_id, 
		organization_id,
		code,
		level, 
		location, 
		status
		FROM w_bays 
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &bays, err
}

func (r *warehouseQuery) GetBayByID(organizationID, bayID string) (*BayResponse, error) {
	var bay BayResponse
	err := r.conn.Get(&bay, `
		SELECT 
		bay_id, 
		organization_id,
		code,
		level, 
		location, 
		status
		FROM w_bays 
		WHERE organization_id = ? AND bay_id = ? AND status > 0
	`, organizationID, bayID)
	return &bay, err
}

//Location
func (r *warehouseQuery) GetLocationCount(filter LocationFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "organization_id = ?"), append(args, v)
	}
	if v := filter.BayID; v != "" {
		where, args = append(where, "bay_id = ?"), append(args, v)
	}
	if v := filter.Code; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM w_locations
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *warehouseQuery) GetLocationList(filter LocationFilter) (*[]LocationResponse, error) {
	where, args := []string{"l.status > 0"}, []interface{}{}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "l.organization_id = ?"), append(args, v)
	}
	if v := filter.BayID; v != "" {
		where, args = append(where, "l.bay_id = ?"), append(args, v)
	}
	if v := filter.Code; v != "" {
		where, args = append(where, "l.code like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var locations []LocationResponse
	err := r.conn.Select(&locations, `
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
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &locations, err
}

func (r *warehouseQuery) GetLocationByID(organizationID, locationID string) (*LocationResponse, error) {
	var location LocationResponse
	err := r.conn.Get(&location, `
		SELECT 
		l.location_id, 
		l.organization_id,
		l.code,
		l.level, 
		l.bay_id,
		b.name as bay_name,
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
		WHERE l.organization_id = ? AND l.location_id = ? AND l.status > 0
	`, organizationID, locationID)
	return &location, err
}
