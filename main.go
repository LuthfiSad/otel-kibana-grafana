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

		// Full CRUD for Tenants
		v1.GET("/tenants", func(c *gin.Context) { c.JSON(200, gin.H{"module": "tenants", "action": "list"}) })
		v1.POST("/tenants", func(c *gin.Context) { c.JSON(201, gin.H{"module": "tenants", "action": "create"}) })
		v1.GET("/tenants/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "tenants", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/tenants/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "tenants", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/tenants/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "tenants", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Users
		v1.GET("/users", func(c *gin.Context) { c.JSON(200, gin.H{"module": "users", "action": "list"}) })
		v1.POST("/users", func(c *gin.Context) { c.JSON(201, gin.H{"module": "users", "action": "create"}) })
		v1.GET("/users/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "users", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/users/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "users", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/users/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "users", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Orders
		v1.GET("/orders", func(c *gin.Context) { c.JSON(200, gin.H{"module": "orders", "action": "list"}) })
		v1.POST("/orders", func(c *gin.Context) { c.JSON(201, gin.H{"module": "orders", "action": "create"}) })
		v1.GET("/orders/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "orders", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/orders/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "orders", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/orders/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "orders", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Products
		v1.GET("/products", func(c *gin.Context) { c.JSON(200, gin.H{"module": "products", "action": "list"}) })
		v1.POST("/products", func(c *gin.Context) { c.JSON(201, gin.H{"module": "products", "action": "create"}) })
		v1.GET("/products/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "products", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/products/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "products", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/products/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "products", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Payments
		v1.GET("/payments", func(c *gin.Context) { c.JSON(200, gin.H{"module": "payments", "action": "list"}) })
		v1.POST("/payments", func(c *gin.Context) { c.JSON(201, gin.H{"module": "payments", "action": "create"}) })
		v1.GET("/payments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "payments", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/payments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "payments", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/payments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "payments", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Shipments
		v1.GET("/shipments", func(c *gin.Context) { c.JSON(200, gin.H{"module": "shipments", "action": "list"}) })
		v1.POST("/shipments", func(c *gin.Context) { c.JSON(201, gin.H{"module": "shipments", "action": "create"}) })
		v1.GET("/shipments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "shipments", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/shipments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "shipments", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/shipments/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "shipments", "action": "delete", "id": c.Param("id")}) })

		// Full CRUD for Inventory & Warehouses
		v1.GET("/inventory", func(c *gin.Context) { c.JSON(200, gin.H{"module": "inventory", "action": "list"}) })
		v1.POST("/inventory", func(c *gin.Context) { c.JSON(201, gin.H{"module": "inventory", "action": "create"}) })
		v1.PUT("/inventory/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "inventory", "action": "update", "id": c.Param("id")}) })

		v1.GET("/warehouses", func(c *gin.Context) { c.JSON(200, gin.H{"module": "warehouses", "action": "list"}) })
		v1.POST("/warehouses", func(c *gin.Context) { c.JSON(201, gin.H{"module": "warehouses", "action": "create"}) })
		v1.GET("/warehouses/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "warehouses", "action": "get_by_id", "id": c.Param("id")}) })
		v1.PUT("/warehouses/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "warehouses", "action": "update", "id": c.Param("id")}) })
		v1.DELETE("/warehouses/:id", func(c *gin.Context) { c.JSON(200, gin.H{"module": "warehouses", "action": "delete", "id": c.Param("id")}) })
	}

	r.Run(":8080")
}
