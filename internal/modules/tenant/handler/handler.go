package handler

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
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

func (h *TenantHandler) GetTenantByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetTenantByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Tenant not found"})
		return
	}
	infrastructure.PushLog("controller", "TenantHandler.GetTenantByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	start := time.Now()
	var req dto.TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.CreateTenant(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "TenantHandler.CreateTenant", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": data})
}

func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateTenant(c.Request.Context(), id, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "TenantHandler.UpdateTenant", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Tenant updated successfully"})
}

func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteTenant(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "TenantHandler.DeleteTenant", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Tenant deleted successfully"})
}
