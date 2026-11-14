package domain

import "time"

type Tenant struct {
	ID        int
	Name      string
	Subdomain string
	Status    string
	CreatedAt time.Time
}

type User struct {
	ID        int
	TenantID  int
	Name      string
	Email     string
	Status    string
	CreatedAt time.Time
}

type Product struct {
	ID          int
	TenantID    int
	Name        string
	Description string
	Price       float64
	SKU         string
	Status      string
	CreatedAt   time.Time
}

type Order struct {
	ID           int
	TenantID     int
	UserID       int
	TotalAmount  float64
	Status       string
	ShippingAddr string
	CreatedAt    time.Time
}
