package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/payment/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) GetPayments(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "PaymentHandler.GetPayments", "SUCCESS", "HTTP GET /payments", nil)
	data, err := h.service.ListPayments(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "PaymentHandler.GetPayments", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}
