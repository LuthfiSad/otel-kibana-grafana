package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/modules/shipment/repository"
)

type ShipmentService struct {
	repo *repository.ShipmentRepository
}

func NewShipmentService(repo *repository.ShipmentRepository) *ShipmentService {
	return &ShipmentService{repo: repo}
}

func (s *ShipmentService) ListShipments(ctx context.Context) ([]domain.Shipment, error) {
	return s.repo.GetAll(ctx)
}
