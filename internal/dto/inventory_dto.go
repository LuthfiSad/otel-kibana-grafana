package dto

type InventoryRequest struct {
	TenantID    int `json:"tenant_id" binding:"required"`
	WarehouseID int `json:"warehouse_id" binding:"required"`
	ProductID   int `json:"product_id" binding:"required"`
	Quantity    int `json:"quantity" binding:"required"`
}
