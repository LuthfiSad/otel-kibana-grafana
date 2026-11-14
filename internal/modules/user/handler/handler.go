package handler

import (
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/user/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "UserHandler.GetUsers", "SUCCESS", "HTTP GET /users", nil)
	data, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.GetUsers", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	start := time.Now()
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		infrastructure.PushLog("controller", "UserHandler.CreateUser", "FAILED", err.Error(), nil)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.RegisterUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.CreateUser", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": user})
}
