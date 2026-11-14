package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (tenant_id, user_id, action, details, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, log.TenantID, log.UserID, log.Action, log.Details, time.Now()).Scan(&log.ID)
}

func (r *AuditLogRepository) GetAll(ctx context.Context) ([]domain.AuditLog, error) {
	query := `SELECT id, tenant_id, user_id, action, details, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 100`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var l domain.AuditLog
		if err := rows.Scan(&l.ID, &l.TenantID, &l.UserID, &l.Action, &l.Details, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
