package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/modules/inventory/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) List(ctx context.Context) ([]domain.Inventory, error) {
	return s.repo.GetAll(ctx)
}

func (s *InventoryService) Adjust(ctx context.Context, req dto.InventoryRequest) (*domain.Inventory, error) {
	inv := domain.Inventory{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Quantity:    req.Quantity,
	}
	if err := s.repo.Create(ctx, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

func (s *InventoryService) GetByID(ctx context.Context, id int) (*domain.Inventory, error) {
	return s.repo.GetByID(ctx, id)
}
