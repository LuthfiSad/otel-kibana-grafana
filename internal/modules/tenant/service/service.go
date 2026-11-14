package service

import (
	"context"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/tenant/repository"
)

type TenantService struct {
	repo *repository.TenantRepository
}

func NewTenantService(repo *repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "TenantService.ListTenants", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "TenantService.ListTenants", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *TenantService) GetTenantByID(ctx context.Context, id int) (*domain.Tenant, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "TenantService.GetTenantByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "TenantService.GetTenantByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *TenantService) CreateTenant(ctx context.Context, req dto.TenantRequest) (*domain.Tenant, error) {
	start := time.Now()
	t := domain.Tenant{
		Name:   req.Name,
		Domain: req.Domain,
		Status: req.Status,
	}
	if err := s.repo.Create(ctx, &t); err != nil {
		infrastructure.PushLog("service", "TenantService.CreateTenant", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "TenantService.CreateTenant", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &t, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, id int, req dto.TenantRequest) error {
	start := time.Now()
	t := domain.Tenant{
		ID:     id,
		Name:   req.Name,
		Domain: req.Domain,
		Status: req.Status,
	}
	if err := s.repo.Update(ctx, &t); err != nil {
		infrastructure.PushLog("service", "TenantService.UpdateTenant", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "TenantService.UpdateTenant", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "TenantService.DeleteTenant", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "TenantService.DeleteTenant", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
