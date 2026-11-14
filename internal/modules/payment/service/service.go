package service

import (
	"context"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/payment/repository"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) GetPayments(ctx context.Context) ([]domain.Payment, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "PaymentService.GetPayments", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "PaymentService.GetPayments", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *PaymentService) ListPayments(ctx context.Context) ([]domain.Payment, error) {
	return s.GetPayments(ctx)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (*domain.Payment, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "PaymentService.GetPaymentByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "PaymentService.GetPaymentByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *PaymentService) CreatePayment(ctx context.Context, req dto.PaymentRequest) (*domain.Payment, error) {
	start := time.Now()
	p := domain.Payment{
		TenantID:      req.TenantID,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}
	if err := s.repo.Create(ctx, &p); err != nil {
		infrastructure.PushLog("service", "PaymentService.CreatePayment", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "PaymentService.CreatePayment", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}

func (s *PaymentService) UpdatePayment(ctx context.Context, id int, req dto.PaymentRequest) error {
	start := time.Now()
	p := domain.Payment{
		ID:            id,
		TenantID:      req.TenantID,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}
	if err := s.repo.Update(ctx, &p); err != nil {
		infrastructure.PushLog("service", "PaymentService.UpdatePayment", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "PaymentService.UpdatePayment", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "PaymentService.DeletePayment", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "PaymentService.DeletePayment", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
