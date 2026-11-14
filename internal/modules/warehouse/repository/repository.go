package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, location, status, created_at FROM warehouses LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		rows.Scan(&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Status, &w.CreatedAt)
		list = append(list, w)
	}
	infrastructure.PushLog("repository", "WarehouseRepo.GetAll", "SUCCESS", "Fetched warehouses", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}
