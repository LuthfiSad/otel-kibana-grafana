package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type PaymentRepository struct { db *sql.DB }
func NewPaymentRepository(db *sql.DB) *PaymentRepository { return &PaymentRepository{db: db} }

func (r *PaymentRepository) GetAll(ctx context.Context) ([]domain.Payment, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, order_id, amount, payment_method, status, created_at FROM payments LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()
	var list []domain.Payment
	for rows.Next() {
		var p domain.Payment
		rows.Scan(&p.ID, &p.TenantID, &p.OrderID, &p.Amount, &p.PaymentMethod, &p.Status, &p.CreatedAt)
		list = append(list, p)
	}
	infrastructure.PushLog("repository", "PaymentRepo.GetAll", "SUCCESS", "Fetched", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *PaymentRepository) Create(ctx context.Context, tenantID, orderID int, amount float64, method string) (*domain.Payment, error) {
	start := time.Now()
	var p domain.Payment
	err := r.db.QueryRowContext(ctx, "INSERT INTO payments (tenant_id, order_id, amount, payment_method, status) VALUES ($1, $2, $3, $4, 'SUCCESS') RETURNING id, tenant_id, order_id, amount, payment_method, status, created_at", tenantID, orderID, amount, method).
		Scan(&p.ID, &p.TenantID, &p.OrderID, &p.Amount, &p.PaymentMethod, &p.Status, &p.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepo.Create", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("repository", "PaymentRepo.Create", "SUCCESS", "Created", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}
