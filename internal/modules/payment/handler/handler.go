package handler

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
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

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Payment not found"})
		return
	}
	infrastructure.PushLog("controller", "PaymentHandler.GetPaymentByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	start := time.Now()
	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "PaymentHandler.CreatePayment", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": data})
}

func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdatePayment(c.Request.Context(), id, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "PaymentHandler.UpdatePayment", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Payment updated successfully"})
}

func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeletePayment(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "PaymentHandler.DeletePayment", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Payment deleted successfully"})
}
