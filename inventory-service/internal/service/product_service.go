package service

import (
	"context"

	"github.com/Assignment/inventory-service/internal/domain"
	"github.com/Assignment/inventory-service/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	GetProductByID(ctx context.Context, id string) (domain.Product, error)
	UpdateProduct(ctx context.Context, product domain.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	return s.productRepo.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	return s.productRepo.List(ctx, filter)
}
