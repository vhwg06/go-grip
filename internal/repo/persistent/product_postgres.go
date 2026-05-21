package persistent

import (
	"context"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type CatalogRepo struct {
	*postgres.Postgres
	mu         sync.RWMutex
	products   map[string]entity.Product
	categories map[string]entity.Category
	tags       map[string]entity.Tag
}

func NewCatalogRepo(pg *postgres.Postgres) *CatalogRepo {
	return &CatalogRepo{Postgres: pg, products: map[string]entity.Product{}, categories: map[string]entity.Category{}, tags: map[string]entity.Tag{}}
}

func (r *CatalogRepo) StoreProduct(ctx context.Context, product *entity.Product) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.products {
		if existing.SKU == product.SKU && existing.ID != product.ID {
			return entity.ErrDuplicateSKU
		}
	}
	r.products[product.ID] = *product
	return nil
}

func (r *CatalogRepo) GetProduct(ctx context.Context, id string) (entity.Product, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, ok := r.products[id]
	if !ok {
		return entity.Product{}, entity.ErrNotFound
	}
	return product, nil
}

func (r *CatalogRepo) ListProducts(ctx context.Context, filter repo.ProductFilter) ([]entity.Product, int, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.Product, 0, len(r.products))
	for _, product := range r.products {
		if filter.Brand != "" && product.Brand != filter.Brand {
			continue
		}
		if filter.MinPrice != nil && product.Price < *filter.MinPrice {
			continue
		}
		if filter.MaxPrice != nil && product.Price > *filter.MaxPrice {
			continue
		}
		items = append(items, product)
	}
	return items, len(items), nil
}

func (r *CatalogRepo) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[product.ID]; !ok {
		return entity.ErrNotFound
	}
	r.products[product.ID] = *product
	return nil
}

func (r *CatalogRepo) DeleteProduct(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.products, id)
	return nil
}

func (r *CatalogRepo) StoreCategory(ctx context.Context, category *entity.Category) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.categories[category.ID] = *category
	return nil
}

func (r *CatalogRepo) ListCategories(ctx context.Context) ([]entity.Category, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.Category, 0, len(r.categories))
	for _, category := range r.categories {
		items = append(items, category)
	}
	return items, nil
}

func (r *CatalogRepo) StoreTag(ctx context.Context, tag *entity.Tag) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[tag.ID] = *tag
	return nil
}

func (r *CatalogRepo) ListTags(ctx context.Context) ([]entity.Tag, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.Tag, 0, len(r.tags))
	for _, tag := range r.tags {
		items = append(items, tag)
	}
	return items, nil
}
