package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]domain.Inventory, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, warehouse_id, product_id, quantity, status, created_at FROM inventory_stocks LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "InventoryRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Inventory
	for rows.Next() {
		var i domain.Inventory
		rows.Scan(&i.ID, &i.TenantID, &i.WarehouseID, &i.ProductID, &i.Quantity, &i.Status, &i.CreatedAt)
		list = append(list, i)
	}
	infrastructure.PushLog("repository", "InventoryRepo.GetAll", "SUCCESS", "Fetched inventories", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}
