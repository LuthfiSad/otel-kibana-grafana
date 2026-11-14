package domain

import "time"

type Warehouse struct {
	ID        int
	TenantID  int
	Name      string
	Location  string
	Status    string
	CreatedAt time.Time
}
