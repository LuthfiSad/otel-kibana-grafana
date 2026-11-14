package service

import (
	"context"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/shipment/repository"
)

type ShipmentService struct {
	repo *repository.ShipmentRepository
}

func NewShipmentService(repo *repository.ShipmentRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

func (s *ShipmentService) ListShipments(ctx context.Context) ([]domain.Shipment, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "ShipmentService.ListShipments", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "ShipmentService.ListShipments", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *ShipmentService) GetShipmentByID(ctx context.Context, id int) (*domain.Shipment, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "ShipmentService.GetShipmentByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "ShipmentService.GetShipmentByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *ShipmentService) CreateShipment(ctx context.Context, req dto.ShipmentRequest) (*domain.Shipment, error) {
	start := time.Now()
	sh := domain.Shipment{
		TenantID:       req.TenantID,
		OrderID:        req.OrderID,
		Courier:        req.Courier,
		TrackingNumber: req.TrackingNumber,
		Status:         req.Status,
	}
	if err := s.repo.Create(ctx, &sh); err != nil {
		infrastructure.PushLog("service", "ShipmentService.CreateShipment", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "ShipmentService.CreateShipment", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &sh, nil
}

func (s *ShipmentService) UpdateShipment(ctx context.Context, id int, req dto.ShipmentRequest) error {
	start := time.Now()
	sh := domain.Shipment{
		ID:             id,
		TenantID:       req.TenantID,
		OrderID:        req.OrderID,
		Courier:        req.Courier,
		TrackingNumber: req.TrackingNumber,
		Status:         req.Status,
	}
	if err := s.repo.Update(ctx, &sh); err != nil {
		infrastructure.PushLog("service", "ShipmentService.UpdateShipment", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "ShipmentService.UpdateShipment", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *ShipmentService) DeleteShipment(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "ShipmentService.DeleteShipment", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "ShipmentService.DeleteShipment", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
