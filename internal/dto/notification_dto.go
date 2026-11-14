package dto

type NotificationRequest struct {
	TenantID int    `json:"tenant_id" binding:"required"`
	UserID   int    `json:"user_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Message  string `json:"message" binding:"required"`
}
