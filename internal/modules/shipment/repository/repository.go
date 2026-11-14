package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type ShipmentRepository struct {
	db *sql.DB
}

func NewShipmentRepository(db *sql.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) GetAll(ctx context.Context) ([]domain.Shipment, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, order_id, status, tracking_number, created_at FROM shipments")
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.Shipment
	for rows.Next() {
		var s domain.Shipment
		if err := rows.Scan(&s.ID, &s.OrderID, &s.Status, &s.TrackingNumber, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	infrastructure.PushLog("repository", "ShipmentRepository.GetAll", "SUCCESS", "Fetched shipments", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *ShipmentRepository) GetByID(ctx context.Context, id int) (*domain.Shipment, error) {
	start := time.Now()
	var s domain.Shipment
	err := r.db.QueryRowContext(ctx, "SELECT id, order_id, status, tracking_number, created_at FROM shipments WHERE id = $1", id).
		Scan(&s.ID, &s.OrderID, &s.Status, &s.TrackingNumber, &s.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "ShipmentRepository.GetByID", "SUCCESS", "Fetched shipment by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &s, nil
}

func (r *ShipmentRepository) Create(ctx context.Context, s *domain.Shipment) error {
	start := time.Now()
	query := "INSERT INTO shipments (order_id, status, tracking_number, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, s.OrderID, s.Status, s.TrackingNumber, time.Now()).Scan(&s.ID)
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "ShipmentRepository.Create", "SUCCESS", "Created shipment", map[string]interface{}{"id": s.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *ShipmentRepository) Update(ctx context.Context, s *domain.Shipment) error {
	start := time.Now()
	query := "UPDATE shipments SET status = $1, tracking_number = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, s.Status, s.TrackingNumber, s.ID)
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": s.ID})
		return err
	}
	infrastructure.PushLog("repository", "ShipmentRepository.Update", "SUCCESS", "Updated shipment", map[string]interface{}{"id": s.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *ShipmentRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM shipments WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "ShipmentRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "ShipmentRepository.Delete", "SUCCESS", "Deleted shipment", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
