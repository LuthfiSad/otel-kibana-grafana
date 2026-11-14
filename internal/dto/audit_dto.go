package dto

type AuditLogRequest struct {
	TenantID int    `json:"tenant_id" binding:"required"`
	UserID   int    `json:"user_id" binding:"required"`
	Action   string `json:"action" binding:"required"`
	Details  string `json:"details" binding:"required"`
}
