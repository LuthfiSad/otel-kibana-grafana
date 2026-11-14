package service

import (
	"context"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/dto"
	"otel-kibana-grafana-test/internal/infrastructure"
	"otel-kibana-grafana-test/internal/modules/product/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, tenantID int, req dto.ProductRequest) (*domain.Product, error) {
	infrastructure.PushLog("service", "ProductService.CreateProduct", "SUCCESS", "Creating product via DTO", map[string]interface{}{"name": req.Name})
	return s.repo.Create(ctx, tenantID, req.Name, req.Description, req.Price)
}
