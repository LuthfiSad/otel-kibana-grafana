package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) GetAll(ctx context.Context) ([]domain.Shipment, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, order_id, courier, tracking_number, status, created_at FROM shipments LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Shipment
	for rows.Next() {
		var s domain.Shipment
		rows.Scan(&s.ID, &s.TenantID, &s.OrderID, &s.Courier, &s.TrackingNumber, &s.Status, &s.CreatedAt)
		list = append(list, s)
	}
	infrastructure.PushLog("repository", "ShipmentRepo.GetAll", "SUCCESS", "Fetched shipments", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}
