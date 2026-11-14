package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/audit/handler"
	"otel-kibana-grafana-test/internal/modules/audit/repository"
	"otel-kibana-grafana-test/internal/modules/audit/service"
)

func main() {
	db := infrastructure.NewDB()

	r := gin.Default()
	r.Use(otelgin.Middleware("audit-service"))

	auditRepo := repository.NewAuditLogRepository(db)
	auditSvc := service.NewAuditLogService(auditRepo)
	auditHdl := handler.NewAuditLogHandler(auditSvc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/audit", auditHdl.Record)
		v1.GET("/audit", auditHdl.List)
	}

	r.Run(":8080")
}
