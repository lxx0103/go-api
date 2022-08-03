package item

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type itemQuery struct {
	conn *sqlx.DB
}

func NewItemQuery(connection *sqlx.DB) *itemQuery {
	return &itemQuery{
		conn: connection,
	}
}

func (r *itemQuery) GetItemByID(organizationID, id string) (*ItemResponse, error) {
	var item ItemResponse
	err := r.conn.Get(&item, `
	SELECT 
	i.sku, 
	i.name, 
	i.unit_id, 
	u.name as unit_name, 
	i.manufacturer_id, 
	IFNULL(m.name, "") as manufacturer_name, 
	i.brand_id,
	IFNULL(b.name, "") as brand_name,
	i.weight_unit,
	i.weight,
	i.dimension_unit,
	i.length,
	i.width,
	i.height,
	i.selling_price,
	i.cost_price,
	i.openning_stock,
	i.openning_stock_rate,
	i.reorder_stock,
	i.default_vendor_id,
	i.description,
	i.status
	FROM i_items i
	LEFT JOIN s_units u
	ON i.unit_id = u.unit_id
	LEFT JOIN s_manufacturers m
	ON i.manufacturer_id = m.manufacturer_id
	LEFT JOIN s_brands b
	ON i.brand_id = b.brand_id
	WHERE i.organization_id = ? AND i.item_id = ? AND i.status > 0
	`, organizationID, id)
	return &item, err
}

func (r *itemQuery) GetItemCount(filter ItemFilter) (int, error) {
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
		FROM i_items
		WHERE `+strings.Join(where, " AND "), args...)
	return count, err
}

func (r *itemQuery) GetItemList(filter ItemFilter) (*[]ItemResponse, error) {
	where, args := []string{"i.status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "i.name like ?"), append(args, "%"+v+"%")
	}
	if v := filter.OrganizationID; v != "" {
		where, args = append(where, "i.organization_id = ?"), append(args, v)
	}
	args = append(args, filter.PageID*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var items []ItemResponse
	err := r.conn.Select(&items, `
		SELECT 
		i.sku, 
		i.name, 
		i.unit_id, 
		u.name as unit_name, 
		i.manufacturer_id, 
		IFNULL(m.name, "") as manufacturer_name, 
		i.brand_id,
		IFNULL(b.name, "") as brand_name,
		i.weight_unit,
		i.weight,
		i.dimension_unit,
		i.length,
		i.width,
		i.height,
		i.selling_price,
		i.cost_price,
		i.openning_stock,
		i.openning_stock_rate,
		i.reorder_stock,
		i.default_vendor_id,
		i.description,
		i.status
		FROM i_items i
		LEFT JOIN s_units u
		ON i.unit_id = u.unit_id
		LEFT JOIN s_manufacturers m
		ON i.manufacturer_id = m.manufacturer_id
		LEFT JOIN s_brands b
		ON i.brand_id = b.brand_id
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	return &items, err
}
