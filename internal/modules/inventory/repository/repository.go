package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, warehouse_id, product_id, quantity, status, created_at FROM inventory_stocks LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Inventory
	for rows.Next() {
		var i domain.Inventory
		rows.Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.ProductID, &i.Quantity, &i.Status, &i.CreatedAt)
		list = append(list, i)
	}
	return list, nil
}

func (r *InventoryRepository) Create(ctx context.Context, i *domain.Inventory) error {
	query := `INSERT INTO inventory_stocks (tenant_id, warehouse_id, product_id, quantity, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRowContext(ctx, query, i.TenantID, i.WarehouseID, i.ProductID, i.Quantity, i.Status, time.Now()).Scan(&i.ID)
}

func (r *InventoryRepository) GetByID(ctx context.Context, id int) (*domain.Inventory, error) {
	var i domain.Inventory
	query := "SELECT id, tenant_id, warehouse_id, product_id, quantity, status, created_at FROM inventory_stocks WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.ProductID, &i.Quantity, &i.Status, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
