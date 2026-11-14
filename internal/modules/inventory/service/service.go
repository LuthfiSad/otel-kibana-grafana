package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/modules/inventory/repository"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) ListInventories(ctx context.Context) ([]domain.Inventory, error) {
	return s.repo.GetAll(ctx)
}
