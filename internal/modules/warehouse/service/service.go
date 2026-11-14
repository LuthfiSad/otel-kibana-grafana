package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/modules/warehouse/repository"
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(repo *repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) ListWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	return s.repo.GetAll(ctx)
}
