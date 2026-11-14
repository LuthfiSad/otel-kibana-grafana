package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/shipment/service"
)

type ShipmentHandler struct {
	service *service.ShipmentService
}

func NewShipmentHandler(s *service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{service: s}
}

func (h *ShipmentHandler) GetShipments(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "ShipmentHandler.GetShipments", "SUCCESS", "HTTP GET /shipments", nil)
	data, err := h.service.ListShipments(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ShipmentHandler.GetShipments", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}
