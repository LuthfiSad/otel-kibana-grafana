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
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, user_id, total_amount, status, shipping_address, created_at FROM orders LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Order
	for rows.Next() {
		var o domain.Order
		rows.Scan(&o.ID, &o.TenantID, &o.UserID, &o.TotalAmount, &o.Status, &o.ShippingAddr, &o.CreatedAt)
		list = append(list, o)
	}
	infrastructure.PushLog("repository", "OrderRepo.GetAll", "SUCCESS", "Fetched orders", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *OrderRepository) Create(ctx context.Context, tenantID, userID int, totalAmount float64, addr string) (*domain.Order, error) {
	start := time.Now()
	var o domain.Order
	err := r.db.QueryRowContext(ctx, "INSERT INTO orders (tenant_id, user_id, total_amount, status, shipping_address) VALUES ($1, $2, $3, 'COMPLETED', $4) RETURNING id, tenant_id, user_id, total_amount, status, shipping_address, created_at", tenantID, userID, totalAmount, addr).
		Scan(&o.ID, &o.TenantID, &o.UserID, &o.TotalAmount, &o.Status, &o.ShippingAddr, &o.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "OrderRepo.Create", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("repository", "OrderRepo.Create", "SUCCESS", "Order created", map[string]interface{}{"id": o.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &o, nil
}
