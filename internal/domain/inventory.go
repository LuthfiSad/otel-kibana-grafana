package domain

import "time"

type Inventory struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	WarehouseID int       `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	LastUpdated time.Time `json:"last_updated"`
}
