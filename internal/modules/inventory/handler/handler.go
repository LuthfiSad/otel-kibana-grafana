package handler

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/inventory/service"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(s *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: s}
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "InventoryHandler.GetInventory", "SUCCESS", "HTTP GET /inventory", nil)
	data, err := h.service.ListInventory(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "InventoryHandler.GetInventory", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *InventoryHandler) GetByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetInventoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Inventory item not found"})
		return
	}
	infrastructure.PushLog("controller", "InventoryHandler.GetByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	start := time.Now()
	var req dto.InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.AdjustStock(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "InventoryHandler.AdjustStock", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}
