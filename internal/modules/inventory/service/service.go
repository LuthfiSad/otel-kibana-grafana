package service

import (
	"context"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/inventory/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetInventory(ctx context.Context) ([]domain.Inventory, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "InventoryService.GetInventory", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "InventoryService.GetInventory", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *InventoryService) ListInventory(ctx context.Context) ([]domain.Inventory, error) {
	return s.GetInventory(ctx)
}

func (s *InventoryService) GetByID(ctx context.Context, id int) (*domain.Inventory, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "InventoryService.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "InventoryService.GetByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *InventoryService) GetInventoryByID(ctx context.Context, id int) (*domain.Inventory, error) {
	return s.GetByID(ctx, id)
}

func (s *InventoryService) AdjustStock(ctx context.Context, req dto.InventoryRequest) (*domain.Inventory, error) {
	start := time.Now()
	i := &domain.Inventory{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Quantity:    req.Quantity,
	}
	if err := s.repo.Upsert(ctx, i); err != nil {
		infrastructure.PushLog("service", "InventoryService.AdjustStock", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "InventoryService.AdjustStock", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return i, nil
}
