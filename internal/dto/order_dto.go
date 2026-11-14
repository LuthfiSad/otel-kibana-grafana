package dto

type OrderRequest struct {
	TenantID     int     `json:"tenant_id" binding:"required"`
	UserID       int     `json:"user_id" binding:"required"`
	TotalAmount  float64 `json:"total_amount" binding:"required"`
	ShippingAddr string  `json:"shipping_addr" binding:"required"`
}
