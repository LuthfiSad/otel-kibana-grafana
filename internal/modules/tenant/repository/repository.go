package repository

import (
	"context"
	"database/sql"
	"time"

	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) GetAll(ctx context.Context) ([]domain.Tenant, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, subdomain, status, created_at FROM tenants")
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		rows.Scan(&t.ID, &t.Name, &t.Subdomain, &t.Status, &t.CreatedAt)
		list = append(list, t)
	}
	infrastructure.PushLog("repository", "TenantRepo.GetAll", "SUCCESS", "Fetched tenants", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}
