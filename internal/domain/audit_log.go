package domain

import "time"

type AuditLog struct {
	ID        int       `json:"id"`
	TenantID  int       `json:"tenant_id"`
	UserID    int       `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
