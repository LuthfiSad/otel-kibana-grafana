package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/user/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) RegisterUser(ctx context.Context, req dto.UserRequest) (*domain.User, error) {
	infrastructure.PushLog("service", "UserService.RegisterUser", "SUCCESS", "Processing user registration via DTO", map[string]interface{}{"email": req.Email})
	return s.repo.Create(ctx, req.TenantID, req.Name, req.Email)
}
