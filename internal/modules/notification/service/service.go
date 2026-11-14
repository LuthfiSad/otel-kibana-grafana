package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/modules/notification/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Send(ctx context.Context, req dto.NotificationRequest) (*domain.Notification, error) {
	n := &domain.Notification{
		TenantID: req.TenantID,
		UserID:   req.UserID,
		Title:    req.Title,
		Message:  req.Message,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) List(ctx context.Context) ([]domain.Notification, error) {
	return s.repo.GetAll(ctx)
}
