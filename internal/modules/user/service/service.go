package service

import (
	"context"
	"time"
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

func (s *UserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "UserService.GetUsers", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "UserService.GetUsers", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.GetUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "UserService.GetUserByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "UserService.GetUserByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *UserService) CreateUser(ctx context.Context, req dto.UserRequest) (*domain.User, error) {
	start := time.Now()
	u := domain.User{
		TenantID: req.TenantID,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
	}
	if err := s.repo.Create(ctx, &u); err != nil {
		infrastructure.PushLog("service", "UserService.CreateUser", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "UserService.CreateUser", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &u, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req dto.UserRequest) error {
	start := time.Now()
	u := domain.User{
		ID:       id,
		TenantID: req.TenantID,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
	}
	if err := s.repo.Update(ctx, &u); err != nil {
		infrastructure.PushLog("service", "UserService.UpdateUser", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "UserService.UpdateUser", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "UserService.DeleteUser", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "UserService.DeleteUser", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
