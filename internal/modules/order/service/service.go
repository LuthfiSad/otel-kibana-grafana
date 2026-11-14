package service

import (
	"context"
	"time"
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

func (s *OrderService) GetOrders(ctx context.Context) ([]domain.Order, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "OrderService.GetOrders", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "OrderService.GetOrders", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	return s.GetOrders(ctx)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (*domain.Order, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "OrderService.GetOrderByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "OrderService.GetOrderByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.OrderRequest) (*domain.Order, error) {
	start := time.Now()
	o := domain.Order{
		TenantID:     req.TenantID,
		UserID:       req.UserID,
		TotalAmount:  req.TotalAmount,
		ShippingAddr: req.ShippingAddr,
		Status:       req.Status,
	}
	if err := s.repo.Create(ctx, &o); err != nil {
		infrastructure.PushLog("service", "OrderService.CreateOrder", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "OrderService.CreateOrder", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &o, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, req dto.OrderRequest) error {
	start := time.Now()
	o := domain.Order{
		ID:           id,
		TenantID:     req.TenantID,
		UserID:       req.UserID,
		TotalAmount:  req.TotalAmount,
		ShippingAddr: req.ShippingAddr,
		Status:       req.Status,
	}
	if err := s.repo.Update(ctx, &o); err != nil {
		infrastructure.PushLog("service", "OrderService.UpdateOrder", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "OrderService.UpdateOrder", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "OrderService.DeleteOrder", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "OrderService.DeleteOrder", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
