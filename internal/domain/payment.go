package domain

import "time"

type Payment struct {
	ID            int
	TenantID      int
	OrderID       int
	Amount        float64
	PaymentMethod string
	Status        string
	CreatedAt     time.Time
}
