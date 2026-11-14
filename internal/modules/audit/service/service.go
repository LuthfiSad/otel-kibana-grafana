package service

import (
	"context"
	"log"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/modules/audit/repository"
)

type AuditLogService struct {
	repo *repository.AuditLogRepository
}

func NewAuditLogService(repo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) RecordLog(ctx context.Context, req dto.AuditLogRequest) (*domain.AuditLog, error) {
	logEntry := &domain.AuditLog{
		TenantID: req.TenantID,
		UserID:   req.UserID,
		Action:   req.Action,
		Details:  req.Details,
	}

	if err := s.repo.Create(ctx, logEntry); err != nil {
		log.Printf("ERROR: Failed to create audit log: %v", err)
		return nil, err
	}

	log.Printf("INFO: Created audit log ID %d for action %s", logEntry.ID, logEntry.Action)
	return logEntry, nil
}

func (s *AuditLogService) FetchLogs(ctx context.Context) ([]domain.AuditLog, error) {
	return s.repo.GetAll(ctx)
}
