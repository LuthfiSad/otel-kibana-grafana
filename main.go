package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/middleware"
	"otel-kibana-grafana-test/internal/modules/audit/handler"
	"otel-kibana-grafana-test/internal/modules/audit/repository"
	"otel-kibana-grafana-test/internal/modules/audit/service"
)

func main() {
	db := infrastructure.NewDB()

	r := gin.Default()
	r.Use(otelgin.Middleware("enterprise-service"))

	auditRepo := repository.NewAuditLogRepository(db)
	auditSvc := service.NewAuditLogService(auditRepo)
	auditHdl := handler.NewAuditLogHandler(auditSvc)

	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// Audit Log
		v1.POST("/audit", auditHdl.Record)
		v1.GET("/audit", auditHdl.List)

		// Placeholder groupings for expanded enterprise modules (Tenant, User, Order, Product, Payment, Shipment, Inventory, Warehouse)
		v1.GET("/tenants", func(c *gin.Context) { c.JSON(200, gin.H{"module": "tenants", "status": "active"}) })
		v1.POST("/tenants", func(c *gin.Context) { c.JSON(201, gin.H{"module": "tenants", "status": "created"}) })

		v1.GET("/users", func(c *gin.Context) { c.JSON(200, gin.H{"module": "users", "status": "active"}) })
		v1.POST("/users", func(c *gin.Context) { c.JSON(201, gin.H{"module": "users", "status": "created"}) })

		v1.GET("/orders", func(c *gin.Context) { c.JSON(200, gin.H{"module": "orders", "status": "active"}) })
		v1.POST("/orders", func(c *gin.Context) { c.JSON(201, gin.H{"module": "orders", "status": "created"}) })

		v1.GET("/products", func(c *gin.Context) { c.JSON(200, gin.H{"module": "products", "status": "active"}) })
		v1.POST("/products", func(c *gin.Context) { c.JSON(201, gin.H{"module": "products", "status": "created"}) })

		v1.GET("/payments", func(c *gin.Context) { c.JSON(200, gin.H{"module": "payments", "status": "active"}) })
		v1.POST("/payments", func(c *gin.Context) { c.JSON(201, gin.H{"module": "payments", "status": "created"}) })

		v1.GET("/shipments", func(c *gin.Context) { c.JSON(200, gin.H{"module": "shipments", "status": "active"}) })
		v1.POST("/shipments", func(c *gin.Context) { c.JSON(201, gin.H{"module": "shipments", "status": "created"}) })

		v1.GET("/inventory", func(c *gin.Context) { c.JSON(200, gin.H{"module": "inventory", "status": "active"}) })
		v1.GET("/warehouses", func(c *gin.Context) { c.JSON(200, gin.H{"module": "warehouses", "status": "active"}) })
	}

	r.Run(":8080")
}
