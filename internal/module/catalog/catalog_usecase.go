package catalog

import (
	"context"
	"fmt"
	"time"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/google/uuid"
)

// CatalogUseCase defines application service interface for Catalog management.
type CatalogUseCase interface {
	CreateProduct(ctx context.Context, product Product) (Product, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]Product, int, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, product Product) (Product, error)
	DeleteProduct(ctx context.Context, id string) error
	CreateCategory(ctx context.Context, category Category) (Category, error)
	ListCategories(ctx context.Context) ([]Category, error)
	CreateTag(ctx context.Context, tag Tag) (Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	ListVisibleProducts(ctx context.Context, actor usermodule.Actor, filter ProductFilter) ([]Product, int, error)
	GetVisibleProduct(ctx context.Context, actor usermodule.Actor, productID string) (Product, error)
	ListPublicSettings(ctx context.Context) ([]Setting, error)
	GetPublicSetting(ctx context.Context, key string) (Setting, error)
}

type catalogUseCase struct {
	repo     CatalogRepo
	gripRepo GripCatalogRepo
}

// NewCatalogUseCase constructs a new CatalogUseCase instance.
func NewCatalogUseCase(r CatalogRepo, grip GripCatalogRepo) CatalogUseCase {
	return &catalogUseCase{repo: r, gripRepo: grip}
}

func (uc *catalogUseCase) CreateProduct(ctx context.Context, product Product) (Product, error) {
	now := time.Now().UTC()
	product.ID = uuid.New().String()
	product.Status = defaultProductStatus(product.Status)
	product.CreatedAt = now
	product.UpdatedAt = now
	if product.Price < 0 || product.SKU == "" || product.Title == "" {
		return Product{}, ErrInvalidInput
	}
	if err := uc.repo.StoreProduct(ctx, &product); err != nil {
		return Product{}, fmt.Errorf("CatalogUseCase - CreateProduct - repo.StoreProduct: %w", err)
	}
	return product, nil
}

func (uc *catalogUseCase) ListProducts(ctx context.Context, filter ProductFilter) ([]Product, int, error) {
	page := filter.Pagination.Normalize()
	items, total, err := uc.repo.ListProducts(ctx, ProductRepoFilter{
		Keyword: filter.Keyword, Category: filter.CategoryID, Brand: filter.Brand, MinPrice: filter.MinPrice, MaxPrice: filter.MaxPrice,
		Sort: filter.Sort, Limit: uint64(page.Limit), Offset: uint64(page.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("CatalogUseCase - ListProducts - repo.ListProducts: %w", err)
	}
	return items, total, nil
}

func (uc *catalogUseCase) GetProduct(ctx context.Context, id string) (Product, error) {
	return uc.repo.GetProduct(ctx, id)
}

func (uc *catalogUseCase) UpdateProduct(ctx context.Context, product Product) (Product, error) {
	product.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdateProduct(ctx, &product); err != nil {
		return Product{}, fmt.Errorf("CatalogUseCase - UpdateProduct - repo.UpdateProduct: %w", err)
	}
	return product, nil
}

func (uc *catalogUseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.DeleteProduct(ctx, id)
}

func (uc *catalogUseCase) CreateCategory(ctx context.Context, category Category) (Category, error) {
	category.ID = uuid.New().String()
	category.IsActive = true
	if err := uc.repo.StoreCategory(ctx, &category); err != nil {
		return Category{}, err
	}
	return category, nil
}

func (uc *catalogUseCase) ListCategories(ctx context.Context) ([]Category, error) {
	if uc.gripRepo != nil {
		return uc.gripRepo.ListCategories(ctx)
	}
	return uc.repo.ListCategories(ctx)
}

func (uc *catalogUseCase) CreateTag(ctx context.Context, tag Tag) (Tag, error) {
	tag.ID = uuid.New().String()
	if err := uc.repo.StoreTag(ctx, &tag); err != nil {
		return Tag{}, err
	}
	return tag, nil
}

func (uc *catalogUseCase) ListTags(ctx context.Context) ([]Tag, error) {
	return uc.repo.ListTags(ctx)
}

func (uc *catalogUseCase) ListVisibleProducts(ctx context.Context, actor usermodule.Actor, filter ProductFilter) ([]Product, int, error) {
	if uc.gripRepo == nil {
		return nil, 0, ErrNotFound
	}

	page := filter.Pagination.Normalize()
	items, total, err := uc.gripRepo.ListVisibleProducts(ctx, actor, ProductRepoFilter{
		Keyword:  filter.Keyword,
		Category: filter.CategoryID,
		Brand:    filter.Brand,
		MinPrice: filter.MinPrice,
		MaxPrice: filter.MaxPrice,
		Sort:     filter.Sort,
		Limit:    uint64(page.Limit),
		Offset:   uint64(page.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("CatalogUseCase - ListVisibleProducts - gripRepo.ListVisibleProducts: %w", err)
	}

	for i := range items {
		items[i].MaxPurchaseableQuantity = resolveMaxPurchasable(items[i])
	}

	return items, total, nil
}

func (uc *catalogUseCase) GetVisibleProduct(ctx context.Context, actor usermodule.Actor, productID string) (Product, error) {
	if uc.gripRepo == nil {
		return Product{}, ErrNotFound
	}

	product, err := uc.gripRepo.GetVisibleProduct(ctx, actor, productID)
	if err != nil {
		return Product{}, fmt.Errorf("CatalogUseCase - GetVisibleProduct - gripRepo.GetVisibleProduct: %w", err)
	}

	product.MaxPurchaseableQuantity = resolveMaxPurchasable(product)

	return product, nil
}

func (uc *catalogUseCase) ListPublicSettings(ctx context.Context) ([]Setting, error) {
	if uc.gripRepo == nil {
		return nil, ErrNotFound
	}
	settings, err := uc.gripRepo.ListSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("CatalogUseCase - ListPublicSettings - gripRepo.ListSettings: %w", err)
	}
	return settings, nil
}

func (uc *catalogUseCase) GetPublicSetting(ctx context.Context, key string) (Setting, error) {
	if uc.gripRepo == nil {
		return Setting{}, ErrNotFound
	}
	setting, err := uc.gripRepo.GetSetting(ctx, key)
	if err != nil {
		return Setting{}, fmt.Errorf("CatalogUseCase - GetPublicSetting - gripRepo.GetSetting: %w", err)
	}
	return setting, nil
}

func defaultProductStatus(status ProductStatus) ProductStatus {
	if status == "" {
		return ProductStatusDraft
	}
	return status
}

func resolveMaxPurchasable(product Product) int {
	displayStock := product.StockCount - product.LockedCount

	if product.PurchaseLimit > 0 && product.PurchaseLimit < displayStock {
		return product.PurchaseLimit
	}

	if displayStock < 0 {
		return 0
	}

	return displayStock
}
