package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	TenantID  int       `json:"tenant_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
