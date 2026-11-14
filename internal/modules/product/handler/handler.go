package handler

import (
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/product/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	start := time.Now()
	infrastructure.PushLog("controller", "ProductHandler.GetProducts", "SUCCESS", "HTTP GET /products", nil)
	data, err := h.service.ListProducts(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.GetProducts", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	data, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.GetProductByID", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "data": data})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	start := time.Now()
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.CreateProduct", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": data})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateProduct(c.Request.Context(), id, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.UpdateProduct", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Product updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	start := time.Now()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.DeleteProduct", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(200, gin.H{"success": true, "message": "Product deleted successfully"})
}
