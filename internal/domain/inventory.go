package domain

import "time"

type Inventory struct {
	ID          int       `json:"id"`
	TenantID    int       `json:"tenant_id"`
	WarehouseID int       `json:"warehouse_id"`
	ProductID   int       `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
