package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetAll(ctx context.Context) ([]domain.Payment, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, order_id, amount, payment_method, status, created_at FROM payments")
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.Amount, &p.PaymentMethod, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	infrastructure.PushLog("repository", "PaymentRepository.GetAll", "SUCCESS", "Fetched payments", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int) (*domain.Payment, error) {
	start := time.Now()
	var p domain.Payment
	err := r.db.QueryRowContext(ctx, "SELECT id, order_id, amount, payment_method, status, created_at FROM payments WHERE id = $1", id).
		Scan(&p.ID, &p.OrderID, &p.Amount, &p.PaymentMethod, &p.Status, &p.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "PaymentRepository.GetByID", "SUCCESS", "Fetched payment by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}

func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	start := time.Now()
	query := "INSERT INTO payments (order_id, amount, payment_method, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.OrderID, p.Amount, p.PaymentMethod, p.Status, time.Now()).Scan(&p.ID)
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "PaymentRepository.Create", "SUCCESS", "Created payment", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *PaymentRepository) Update(ctx context.Context, p *domain.Payment) error {
	start := time.Now()
	query := "UPDATE payments SET status = $1, amount = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, p.Status, p.Amount, p.ID)
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": p.ID})
		return err
	}
	infrastructure.PushLog("repository", "PaymentRepository.Update", "SUCCESS", "Updated payment", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM payments WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "PaymentRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "PaymentRepository.Delete", "SUCCESS", "Deleted payment", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
