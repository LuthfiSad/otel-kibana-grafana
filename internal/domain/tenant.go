package domain

import "time"

type Tenant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
