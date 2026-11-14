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

		// 2. Tenants
		v1.GET("/tenants", tenantHdl.GetTenants)
		v1.POST("/tenants", tenantHdl.CreateTenant)
		v1.GET("/tenants/:id", tenantHdl.GetTenantByID)
		v1.PUT("/tenants/:id", tenantHdl.UpdateTenant)
		v1.DELETE("/tenants/:id", tenantHdl.DeleteTenant)

		// 3. Users
		v1.GET("/users", userHdl.GetUsers)
		v1.POST("/users", userHdl.CreateUser)
		v1.GET("/users/:id", userHdl.GetUserByID)
		v1.PUT("/users/:id", userHdl.UpdateUser)
		v1.DELETE("/users/:id", userHdl.DeleteUser)

		// 4. Orders
		v1.GET("/orders", orderHdl.GetOrders)
		v1.POST("/orders", orderHdl.CreateOrder)
		v1.GET("/orders/:id", orderHdl.GetOrderByID)
		v1.PUT("/orders/:id", orderHdl.UpdateOrder)
		v1.DELETE("/orders/:id", orderHdl.DeleteOrder)

		// 5. Products
		v1.GET("/products", productHdl.GetProducts)
		v1.POST("/products", productHdl.CreateProduct)
		v1.GET("/products/:id", productHdl.GetProductByID)
		v1.PUT("/products/:id", productHdl.UpdateProduct)
		v1.DELETE("/products/:id", productHdl.DeleteProduct)

		// 6. Payments
		v1.GET("/payments", paymentHdl.GetPayments)
		v1.POST("/payments", paymentHdl.CreatePayment)
		v1.GET("/payments/:id", paymentHdl.GetPaymentByID)
		v1.PUT("/payments/:id", paymentHdl.UpdatePayment)
		v1.DELETE("/payments/:id", paymentHdl.DeletePayment)

		// 7. Shipments
		v1.GET("/shipments", shipmentHdl.GetShipments)
		v1.POST("/shipments", shipmentHdl.CreateShipment)
		v1.GET("/shipments/:id", shipmentHdl.GetShipmentByID)
		v1.PUT("/shipments/:id", shipmentHdl.UpdateShipment)
		v1.DELETE("/shipments/:id", shipmentHdl.DeleteShipment)

		// 8. Warehouses
		v1.GET("/warehouses", warehouseHdl.GetWarehouses)
		v1.POST("/warehouses", warehouseHdl.CreateWarehouse)
		v1.GET("/warehouses/:id", warehouseHdl.GetWarehouseByID)
		v1.PUT("/warehouses/:id", warehouseHdl.UpdateWarehouse)
		v1.DELETE("/warehouses/:id", warehouseHdl.DeleteWarehouse)

		// 9. Inventory
		v1.GET("/inventory", inventoryHdl.GetInventory)
		v1.POST("/inventory", inventoryHdl.AdjustStock)
		v1.GET("/inventory/:id", inventoryHdl.GetByID)
	}

	r.Run(":8080")
}
