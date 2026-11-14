package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/tenant/service"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewTenantHandler(s *service.TenantService) *TenantHandler {
	return &TenantHandler{service: s}
}

func (h *TenantHandler) GetTenants(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "TenantHandler.GetTenants", "SUCCESS", "HTTP GET /tenants", nil)
	data, err := h.service.ListTenants(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "TenantHandler.GetTenants", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}
