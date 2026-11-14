package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/modules/tenant/repository"
)

type TenantService struct {
	repo *repository.TenantRepository
}

func NewTenantService(repo *repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	return s.repo.GetAll(ctx)
}
