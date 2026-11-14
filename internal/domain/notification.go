package domain

import "time"

type Notification struct {
	ID        int       `json:"id"`
	TenantID  int       `json:"tenant_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"` // PENDING, SENT, FAILED
	CreatedAt time.Time `json:"created_at"`
}
