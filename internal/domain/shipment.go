package domain

import "time"

type Shipment struct {
	ID             int
	TenantID       int
	OrderID        int
	Courier        string
	TrackingNumber string
	Status         string
	CreatedAt      time.Time
}
