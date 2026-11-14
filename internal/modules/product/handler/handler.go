package handler

import (
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

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	start := time.Now()
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		infrastructure.PushLog("controller", "ProductHandler.CreateProduct", "FAILED", err.Error(), nil)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	product, err := h.service.CreateProduct(c.Request.Context(), 1, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	infrastructure.PushLog("controller", "ProductHandler.CreateProduct", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	c.JSON(201, gin.H{"success": true, "data": product})
}
