package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, description, price, sku, status, created_at FROM products LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Product
	for rows.Next() {
		var p domain.Product
		rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price, &p.SKU, &p.Status, &p.CreatedAt)
		list = append(list, p)
	}
	infrastructure.PushLog("repository", "ProductRepo.GetAll", "SUCCESS", "Fetched products", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *ProductRepository) Create(ctx context.Context, tenantID int, name, desc string, price float64) (*domain.Product, error) {
	start := time.Now()
	var p domain.Product
	sku := "SKU-" + time.Now().Format("20060102150405")
	err := r.db.QueryRowContext(ctx, "INSERT INTO products (tenant_id, name, description, price, sku, status) VALUES ($1, $2, $3, $4, $5, 'ACTIVE') RETURNING id, tenant_id, name, description, price, sku, status, created_at", tenantID, name, desc, price, sku).
		Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price, &p.SKU, &p.Status, &p.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepo.Create", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("repository", "ProductRepo.Create", "SUCCESS", "Product created", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}
