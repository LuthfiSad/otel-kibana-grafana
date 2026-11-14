package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/modules/audit/service"
)

type AuditLogHandler struct {
	service *service.AuditLogService
}

func NewAuditLogHandler(service *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

func (h *AuditLogHandler) Record(c *gin.Context) {
	var req dto.AuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.service.RecordLog(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": log, "success": true})
}

func (h *AuditLogHandler) List(c *gin.Context) {
	logs, err := h.service.FetchLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs, "success": true})
}
