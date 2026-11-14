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
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, domain, status, created_at FROM tenants")
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.Domain, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	infrastructure.PushLog("repository", "TenantRepository.GetAll", "SUCCESS", "Fetched tenants", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *TenantRepository) GetByID(ctx context.Context, id int) (*domain.Tenant, error) {
	start := time.Now()
	var t domain.Tenant
	err := r.db.QueryRowContext(ctx, "SELECT id, name, domain, status, created_at FROM tenants WHERE id = $1", id).
		Scan(&t.ID, &t.Name, &t.Domain, &t.Status, &t.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "TenantRepository.GetByID", "SUCCESS", "Fetched tenant by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &t, nil
}

func (r *TenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	start := time.Now()
	query := "INSERT INTO tenants (name, domain, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, t.Name, t.Domain, t.Status, time.Now()).Scan(&t.ID)
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "TenantRepository.Create", "SUCCESS", "Created tenant", map[string]interface{}{"id": t.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *TenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	start := time.Now()
	query := "UPDATE tenants SET name = $1, domain = $2, status = $3 WHERE id = $4"
	_, err := r.db.ExecContext(ctx, query, t.Name, t.Domain, t.Status, t.ID)
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": t.ID})
		return err
	}
	infrastructure.PushLog("repository", "TenantRepository.Update", "SUCCESS", "Updated tenant", map[string]interface{}{"id": t.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *TenantRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM tenants WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "TenantRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "TenantRepository.Delete", "SUCCESS", "Deleted tenant", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
