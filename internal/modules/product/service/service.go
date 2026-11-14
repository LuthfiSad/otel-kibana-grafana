package service

import (
	"context"
	"time"
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
	start := time.Now()
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		infrastructure.PushLog("service", "ProductService.ListProducts", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "ProductService.ListProducts", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	start := time.Now()
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		infrastructure.PushLog("service", "ProductService.GetProductByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("service", "ProductService.GetProductByID", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return data, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req dto.ProductRequest) (*domain.Product, error) {
	start := time.Now()
	p := domain.Product{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
		Status:      req.Status,
	}
	if err := s.repo.Create(ctx, &p); err != nil {
		infrastructure.PushLog("service", "ProductService.CreateProduct", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("service", "ProductService.CreateProduct", "SUCCESS", "Completed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return &p, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, req dto.ProductRequest) error {
	start := time.Now()
	p := domain.Product{
		ID:          id,
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
		Status:      req.Status,
	}
	if err := s.repo.Update(ctx, &p); err != nil {
		infrastructure.PushLog("service", "ProductService.UpdateProduct", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "ProductService.UpdateProduct", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	start := time.Now()
	if err := s.repo.Delete(ctx, id); err != nil {
		infrastructure.PushLog("service", "ProductService.DeleteProduct", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("service", "ProductService.DeleteProduct", "SUCCESS", "Completed", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
