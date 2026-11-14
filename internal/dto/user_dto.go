package dto

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	TenantID int    `json:"tenant_id" binding:"required"`
}
