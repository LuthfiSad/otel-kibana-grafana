package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/order/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "OrderHandler.GetOrders", "SUCCESS", "HTTP GET /orders", nil)
	data, err := h.service.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.GetOrders", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	start := time.Now()
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		infrastructure.PushLog("controller", "OrderHandler.CreateOrder", "FAILED", err.Error(), nil)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	order, err := h.service.PlaceOrder(c.Request.Context(), 1, req) // tenantID hardcoded for demo
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.CreateOrder", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": order})
}
