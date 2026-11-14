package domain

import "time"

type Order struct {
	ID           int       `json:"id"`
	TenantID     int       `json:"tenant_id"`
	UserID       int       `json:"user_id"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	ShippingAddr string    `json:"shipping_addr"`
	CreatedAt    time.Time `json:"created_at"`
}
