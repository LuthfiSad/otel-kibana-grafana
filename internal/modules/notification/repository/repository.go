package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `INSERT INTO notifications (tenant_id, user_id, title, message, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	n.CreatedAt = time.Now()
	n.Status = "SENT"
	return r.db.QueryRowContext(ctx, query, n.TenantID, n.UserID, n.Title, n.Message, n.Status, n.CreatedAt).Scan(&n.ID)
}

func (r *NotificationRepository) GetAll(ctx context.Context) ([]domain.Notification, error) {
	query := `SELECT id, tenant_id, user_id, title, message, status, created_at FROM notifications`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.TenantID, &n.UserID, &n.Title, &n.Message, &n.Status, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
