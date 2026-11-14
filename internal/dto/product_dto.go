package dto

type ProductRequest struct {
	TenantID    int     `json:"tenant_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	SKU         string  `json:"sku"`
	Status      string  `json:"status"`
}
