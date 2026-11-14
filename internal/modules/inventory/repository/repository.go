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
	rows, err := r.db.QueryContext(ctx, "SELECT id, product_id, warehouse_id, quantity, last_updated FROM inventory")
	if err != nil {
		infrastructure.PushLog("repository", "InventoryRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Inventory
	for rows.Next() {
		var i domain.Inventory
		if err := rows.Scan(&i.ID, &i.ProductID, &i.WarehouseID, &i.Quantity, &i.LastUpdated); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	infrastructure.PushLog("repository", "InventoryRepository.GetAll", "SUCCESS", "Fetched inventory", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id int) (*domain.Inventory, error) {
	start := time.Now()
	var i domain.Inventory
	err := r.db.QueryRowContext(ctx, "SELECT id, product_id, warehouse_id, quantity, last_updated FROM inventory WHERE id = $1", id).
		Scan(&i.ID, &i.ProductID, &i.WarehouseID, &i.Quantity, &i.LastUpdated)
	if err != nil {
		infrastructure.PushLog("repository", "InventoryRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "InventoryRepository.GetByID", "SUCCESS", "Fetched inventory by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &i, nil
}

func (r *InventoryRepository) Upsert(ctx context.Context, i *domain.Inventory) error {
	start := time.Now()
	query := `
		INSERT INTO inventory (product_id, warehouse_id, quantity, last_updated)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (product_id, warehouse_id)
		DO UPDATE SET quantity = inventory.quantity + $3, last_updated = $4
		RETURNING id`
	err := r.db.QueryRowContext(ctx, query, i.ProductID, i.WarehouseID, i.Quantity, time.Now()).Scan(&i.ID)
	if err != nil {
		infrastructure.PushLog("repository", "InventoryRepository.Upsert", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "InventoryRepository.Upsert", "SUCCESS", "Upserted inventory", map[string]interface{}{"id": i.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
