package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/warehouse/service"
)

type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(s *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: s}
}

func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "WarehouseHandler.GetWarehouses", "SUCCESS", "HTTP GET /warehouses", nil)
	data, err := h.service.ListWarehouses(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "WarehouseHandler.GetWarehouses", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}
