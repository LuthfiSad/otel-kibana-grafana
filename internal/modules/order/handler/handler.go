package handler

import (
	"strconv"
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

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.GetOrderByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	start := time.Now()
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.CreateOrder", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": data})
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateOrder(c.Request.Context(), id, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.UpdateOrder", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Order updated successfully"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteOrder(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "OrderHandler.DeleteOrder", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Order deleted successfully"})
}
