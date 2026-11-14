package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/payment/repository"
	"time"
)

type PaymentService struct { repo *repository.PaymentRepository }
func NewPaymentService(repo *repository.PaymentRepository) *PaymentService { return &PaymentService{repo: repo} }

func (s *PaymentService) ListPayments(ctx context.Context) ([]domain.Payment, error) {
	return s.repo.GetAll(ctx)
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req dto.PaymentRequest) (*domain.Payment, error) {
	start := time.Now()
	res, err := s.repo.Create(ctx, req.TenantID, req.OrderID, req.Amount, req.PaymentMethod)
	if err != nil {
		infrastructure.PushLog("service", "PaymentService.ProcessPayment", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "PaymentService.ProcessPayment", "SUCCESS", "Processed", map[string]interface{}{"id": res.ID, "duration_ms": time.Since(start).Milliseconds()})
	return res, nil
}
