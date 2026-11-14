package handler

import (
	"strconv"
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

func (h *UserHandler) GetUserByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.GetUserByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	start := time.Now()
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.CreateUser", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": data})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateUser(c.Request.Context(), id, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.UpdateUser", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "UserHandler.DeleteUser", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "User deleted successfully"})
}
