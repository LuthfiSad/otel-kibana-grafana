package dto

type WarehouseRequest struct {
	TenantID int    `json:"tenant_id"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Status   string `json:"status"`
}
