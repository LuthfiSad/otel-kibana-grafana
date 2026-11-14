package dto

type ShipmentRequest struct {
	TenantID       int    `json:"tenant_id" binding:"required"`
	OrderID        int    `json:"order_id" binding:"required"`
	Courier        string `json:"courier" binding:"required"`
	TrackingNumber string `json:"tracking_number" binding:"required"`
}
