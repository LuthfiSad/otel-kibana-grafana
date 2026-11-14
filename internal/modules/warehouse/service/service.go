package service

import (
	"context"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/warehouse/repository"
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(repo *repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) GetWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "WarehouseService.GetWarehouses", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "WarehouseService.GetWarehouses", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *WarehouseService) ListWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	return s.GetWarehouses(ctx)
}

func (s *WarehouseService) GetWarehouseByID(ctx context.Context, id int) (*domain.Warehouse, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "WarehouseService.GetWarehouseByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "WarehouseService.GetWarehouseByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, req dto.WarehouseRequest) (*domain.Warehouse, error) {
	start := time.Now()
	w := domain.Warehouse{
		TenantID: req.TenantID,
		Name:     req.Name,
		Location: req.Location,
		Status:   req.Status,
	}
	if err := s.repo.Create(ctx, &w); err != nil {
		infrastructure.PushLog("service", "WarehouseService.CreateWarehouse", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "WarehouseService.CreateWarehouse", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &w, nil
}

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, id int, req dto.WarehouseRequest) error {
	start := time.Now()
	w := domain.Warehouse{
		ID:       id,
		TenantID: req.TenantID,
		Name:     req.Name,
		Location: req.Location,
		Status:   req.Status,
	}
	if err := s.repo.Update(ctx, &w); err != nil {
		infrastructure.PushLog("service", "WarehouseService.UpdateWarehouse", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "WarehouseService.UpdateWarehouse", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "WarehouseService.DeleteWarehouse", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "WarehouseService.DeleteWarehouse", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
