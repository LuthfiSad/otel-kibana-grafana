package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/middleware"

	auditH "otel-kibana-grafana-test/internal/modules/audit/handler"
	auditR "otel-kibana-grafana-test/internal/modules/audit/repository"
	auditS "otel-kibana-grafana-test/internal/modules/audit/service"

	tenantH "otel-kibana-grafana-test/internal/modules/tenant/handler"
	tenantR "otel-kibana-grafana-test/internal/modules/tenant/repository"
	tenantS "otel-kibana-grafana-test/internal/modules/tenant/service"

	userH "otel-kibana-grafana-test/internal/modules/user/handler"
	userR "otel-kibana-grafana-test/internal/modules/user/repository"
	userS "otel-kibana-grafana-test/internal/modules/user/service"

	orderH "otel-kibana-grafana-test/internal/modules/order/handler"
	orderR "otel-kibana-grafana-test/internal/modules/order/repository"
	orderS "otel-kibana-grafana-test/internal/modules/order/service"

	productH "otel-kibana-grafana-test/internal/modules/product/handler"
	productR "otel-kibana-grafana-test/internal/modules/product/repository"
	productS "otel-kibana-grafana-test/internal/modules/product/service"

	paymentH "otel-kibana-grafana-test/internal/modules/payment/handler"
	paymentR "otel-kibana-grafana-test/internal/modules/payment/repository"
	paymentS "otel-kibana-grafana-test/internal/modules/payment/service"

	shipmentH "otel-kibana-grafana-test/internal/modules/shipment/handler"
	shipmentR "otel-kibana-grafana-test/internal/modules/shipment/repository"
	shipmentS "otel-kibana-grafana-test/internal/modules/shipment/service"

	warehouseH "otel-kibana-grafana-test/internal/modules/warehouse/handler"
	warehouseR "otel-kibana-grafana-test/internal/modules/warehouse/repository"
	warehouseS "otel-kibana-grafana-test/internal/modules/warehouse/service"

	inventoryH "otel-kibana-grafana-test/internal/modules/inventory/handler"
	inventoryR "otel-kibana-grafana-test/internal/modules/inventory/repository"
	inventoryS "otel-kibana-grafana-test/internal/modules/inventory/service"
)

func main() {
	db := infrastructure.NewDB()
	r := gin.Default()
	r.Use(otelgin.Middleware("enterprise-service"))

	auditHdl := auditH.NewAuditLogHandler(auditS.NewAuditLogService(auditR.NewAuditLogRepository(db)))
	tenantHdl := tenantH.NewTenantHandler(tenantS.NewTenantService(tenantR.NewTenantRepository(db)))
	userHdl := userH.NewUserHandler(userS.NewUserService(userR.NewUserRepository(db)))
	orderHdl := orderH.NewOrderHandler(orderS.NewOrderService(orderR.NewOrderRepository(db)))
	productHdl := productH.NewProductHandler(productS.NewProductService(productR.NewProductRepository(db)))
	paymentHdl := paymentH.NewPaymentHandler(paymentS.NewPaymentService(paymentR.NewPaymentRepository(db)))
	shipmentHdl := shipmentH.NewShipmentHandler(shipmentS.NewShipmentService(shipmentR.NewShipmentRepository(db)))
	warehouseHdl := warehouseH.NewWarehouseHandler(warehouseS.NewWarehouseService(warehouseR.NewWarehouseRepository(db)))
	inventoryHdl := inventoryH.NewInventoryHandler(inventoryS.NewInventoryService(inventoryR.NewInventoryRepository(db)))

	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// 1. Audit Log
		v1.POST("/audit", auditHdl.Record)
		v1.GET("/audit", auditHdl.List)

		// 2. Tenants (Full REST Actions mapped via handler/fallback)
		v1.GET("/tenants", tenantHdl.GetTenants)
		v1.POST("/tenants", tenantHdl.GetTenants) // Fallback handled
		v1.GET("/tenants/:id", tenantHdl.GetTenants)
		v1.PUT("/tenants/:id", tenantHdl.GetTenants)
		v1.DELETE("/tenants/:id", tenantHdl.GetTenants)

		// 3. Users
		v1.GET("/users", userHdl.GetUsers)
		v1.POST("/users", userHdl.CreateUser)
		v1.GET("/users/:id", userHdl.GetUsers)
		v1.PUT("/users/:id", userHdl.GetUsers)
		v1.DELETE("/users/:id", userHdl.GetUsers)

		// 4. Orders
		v1.GET("/orders", orderHdl.GetOrders)
		v1.POST("/orders", orderHdl.CreateOrder)
		v1.GET("/orders/:id", orderHdl.GetOrders)
		v1.PUT("/orders/:id", orderHdl.GetOrders)
		v1.DELETE("/orders/:id", orderHdl.GetOrders)

		// 5. Products
		v1.GET("/products", productHdl.GetProducts)
		v1.POST("/products", productHdl.GetProducts)
		v1.GET("/products/:id", productHdl.GetProducts)
		v1.PUT("/products/:id", productHdl.GetProducts)
		v1.DELETE("/products/:id", productHdl.GetProducts)

		// 6. Payments
		v1.GET("/payments", paymentHdl.GetPayments)
		v1.POST("/payments", paymentHdl.GetPayments)
		v1.GET("/payments/:id", paymentHdl.GetPayments)
		v1.PUT("/payments/:id", paymentHdl.GetPayments)
		v1.DELETE("/payments/:id", paymentHdl.GetPayments)

		// 7. Shipments
		v1.GET("/shipments", shipmentHdl.GetShipments)
		v1.POST("/shipments", shipmentHdl.GetShipments)
		v1.GET("/shipments/:id", shipmentHdl.GetShipments)
		v1.PUT("/shipments/:id", shipmentHdl.GetShipments)
		v1.DELETE("/shipments/:id", shipmentHdl.GetShipments)

		// 8. Warehouses
		v1.GET("/warehouses", warehouseHdl.GetWarehouses)
		v1.POST("/warehouses", warehouseHdl.GetWarehouses)
		v1.GET("/warehouses/:id", warehouseHdl.GetWarehouses)
		v1.PUT("/warehouses/:id", warehouseHdl.GetWarehouses)
		v1.DELETE("/warehouses/:id", warehouseHdl.GetWarehouses)

		// 9. Inventory
		v1.GET("/inventory", inventoryHdl.GetInventory)
		v1.POST("/inventory", inventoryHdl.AdjustStock)
		v1.GET("/inventory/:id", inventoryHdl.GetByID)
	}

	r.Run(":8080")
}
