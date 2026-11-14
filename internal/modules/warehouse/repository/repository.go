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
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, location, status, created_at FROM warehouses")
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		if err := rows.Scan(&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Status, &w.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	infrastructure.PushLog("repository", "WarehouseRepository.GetAll", "SUCCESS", "Fetched warehouses", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *WarehouseRepository) GetByID(ctx context.Context, id int) (*domain.Warehouse, error) {
	start := time.Now()
	var w domain.Warehouse
	err := r.db.QueryRowContext(ctx, "SELECT id, tenant_id, name, location, status, created_at FROM warehouses WHERE id = $1", id).
		Scan(&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Status, &w.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "WarehouseRepository.GetByID", "SUCCESS", "Fetched warehouse by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &w, nil
}

func (r *WarehouseRepository) Create(ctx context.Context, w *domain.Warehouse) error {
	start := time.Now()
	if w.TenantID == 0 {
		w.TenantID = 1
	}
	query := "INSERT INTO warehouses (tenant_id, name, location, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, w.TenantID, w.Name, w.Location, w.Status, time.Now()).Scan(&w.ID)
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "WarehouseRepository.Create", "SUCCESS", "Created warehouse", map[string]interface{}{"id": w.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *WarehouseRepository) Update(ctx context.Context, w *domain.Warehouse) error {
	start := time.Now()
	query := "UPDATE warehouses SET name = $1, location = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, w.Name, w.Location, w.ID)
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": w.ID})
		return err
	}
	infrastructure.PushLog("repository", "WarehouseRepository.Update", "SUCCESS", "Updated warehouse", map[string]interface{}{"id": w.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *WarehouseRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM warehouses WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "WarehouseRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "WarehouseRepository.Delete", "SUCCESS", "Deleted warehouse", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
