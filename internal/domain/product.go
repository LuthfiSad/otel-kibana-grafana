package domain

import "time"

type Product struct {
	ID          int       `json:"id"`
	TenantID    int       `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	SKU         string    `json:"sku"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
