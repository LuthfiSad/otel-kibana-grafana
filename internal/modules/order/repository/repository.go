package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, user_id, total_amount, status, created_at FROM orders")
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.TenantID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	infrastructure.PushLog("repository", "OrderRepository.GetAll", "SUCCESS", "Fetched orders", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*domain.Order, error) {
	start := time.Now()
	var o domain.Order
	err := r.db.QueryRowContext(ctx, "SELECT id, tenant_id, user_id, total_amount, status, created_at FROM orders WHERE id = $1", id).
		Scan(&o.ID, &o.TenantID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "OrderRepository.GetByID", "SUCCESS", "Fetched order by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &o, nil
}

func (r *OrderRepository) Create(ctx context.Context, o *domain.Order) error {
	start := time.Now()
	query := "INSERT INTO orders (tenant_id, user_id, total_amount, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, o.TenantID, o.UserID, o.TotalAmount, o.Status, time.Now()).Scan(&o.ID)
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "OrderRepository.Create", "SUCCESS", "Created order", map[string]interface{}{"id": o.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *OrderRepository) Update(ctx context.Context, o *domain.Order) error {
	start := time.Now()
	query := "UPDATE orders SET status = $1, total_amount = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, o.Status, o.TotalAmount, o.ID)
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": o.ID})
		return err
	}
	infrastructure.PushLog("repository", "OrderRepository.Update", "SUCCESS", "Updated order", map[string]interface{}{"id": o.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "OrderRepository.Delete", "SUCCESS", "Deleted order", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
