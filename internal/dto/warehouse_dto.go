package dto

type WarehouseRequest struct {
	TenantID int    `json:"tenant_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}
