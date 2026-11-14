package dto

type PaymentRequest struct {
	TenantID      int     `json:"tenant_id" binding:"required"`
	OrderID       int     `json:"order_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Status        string  `json:"status"`
}
