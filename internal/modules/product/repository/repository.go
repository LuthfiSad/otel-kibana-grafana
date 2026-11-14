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
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, sku, price, stock, created_at, description, status FROM products")
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.CreatedAt, &p.Description, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	infrastructure.PushLog("repository", "ProductRepository.GetAll", "SUCCESS", "Fetched products", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	start := time.Now()
	var p domain.Product
	err := r.db.QueryRowContext(ctx, "SELECT id, tenant_id, name, sku, price, stock, created_at, description, status FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.TenantID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.CreatedAt, &p.Description, &p.Status)
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "ProductRepository.GetByID", "SUCCESS", "Fetched product by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	start := time.Now()
	query := "INSERT INTO products (tenant_id, name, sku, price, stock, created_at, description, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.TenantID, p.Name, p.SKU, p.Price, p.Stock, time.Now(), p.Description, p.Status).Scan(&p.ID)
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "ProductRepository.Create", "SUCCESS", "Created product", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	start := time.Now()
	query := "UPDATE products SET name = $1, sku = $2, price = $3, stock = $4, description = $5, status = $6 WHERE id = $7"
	_, err := r.db.ExecContext(ctx, query, p.Name, p.SKU, p.Price, p.Stock, p.Description, p.Status, p.ID)
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": p.ID})
		return err
	}
	infrastructure.PushLog("repository", "ProductRepository.Update", "SUCCESS", "Updated product", map[string]interface{}{"id": p.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "ProductRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "ProductRepository.Delete", "SUCCESS", "Deleted product", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
