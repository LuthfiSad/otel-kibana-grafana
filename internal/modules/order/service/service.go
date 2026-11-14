package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/order/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *OrderService) PlaceOrder(ctx context.Context, tenantID int, req dto.OrderRequest) (*domain.Order, error) {
	infrastructure.PushLog("service", "OrderService.PlaceOrder", "SUCCESS", "Processing order via DTO", map[string]interface{}{"user_id": req.UserID})
	return s.repo.Create(ctx, tenantID, req.UserID, req.TotalAmount, req.ShippingAddr)
}
