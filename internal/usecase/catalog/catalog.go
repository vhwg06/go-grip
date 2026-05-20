package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repo.CatalogRepo
}

func New(r repo.CatalogRepo) *UseCase { return &UseCase{repo: r} }

func (uc *UseCase) CreateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	now := time.Now().UTC()
	product.ID = uuid.New().String()
	product.Status = defaultProductStatus(product.Status)
	product.CreatedAt = now
	product.UpdatedAt = now
	if product.Price < 0 || product.SKU == "" || product.Title == "" {
		return entity.Product{}, entity.ErrInvalidInput
	}
	if err := uc.repo.StoreProduct(ctx, &product); err != nil {
		return entity.Product{}, fmt.Errorf("CatalogUseCase - CreateProduct - repo.StoreProduct: %w", err)
	}
	return product, nil
}

func (uc *UseCase) ListProducts(ctx context.Context, filter entity.ProductFilter) ([]entity.Product, int, error) {
	page := filter.Pagination.Normalize()
	items, total, err := uc.repo.ListProducts(ctx, repo.ProductFilter{
		Keyword: filter.Keyword, Brand: filter.Brand, MinPrice: filter.MinPrice, MaxPrice: filter.MaxPrice,
		Sort: filter.Sort, Limit: uint64(page.Limit), Offset: uint64(page.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("CatalogUseCase - ListProducts - repo.ListProducts: %w", err)
	}
	return items, total, nil
}

func (uc *UseCase) GetProduct(ctx context.Context, id string) (entity.Product, error) {
	return uc.repo.GetProduct(ctx, id)
}

func (uc *UseCase) UpdateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	product.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdateProduct(ctx, &product); err != nil {
		return entity.Product{}, fmt.Errorf("CatalogUseCase - UpdateProduct - repo.UpdateProduct: %w", err)
	}
	return product, nil
}

func (uc *UseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.DeleteProduct(ctx, id)
}

func (uc *UseCase) CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	category.ID = uuid.New().String()
	category.IsActive = true
	if err := uc.repo.StoreCategory(ctx, &category); err != nil {
		return entity.Category{}, err
	}
	return category, nil
}

func (uc *UseCase) ListCategories(ctx context.Context) ([]entity.Category, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) CreateTag(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	tag.ID = uuid.New().String()
	if err := uc.repo.StoreTag(ctx, &tag); err != nil {
		return entity.Tag{}, err
	}
	return tag, nil
}

func (uc *UseCase) ListTags(ctx context.Context) ([]entity.Tag, error) {
	return uc.repo.ListTags(ctx)
}

func defaultProductStatus(status entity.ProductStatus) entity.ProductStatus {
	if status == "" {
		return entity.ProductStatusDraft
	}
	return status
}
