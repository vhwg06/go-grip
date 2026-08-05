package catalog

import (
	"context"

	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// ProductRepoFilter defines low-level filter criteria for repository queries.
type ProductRepoFilter struct {
	Keyword  string
	Category string
	Brand    string
	MinPrice *int64
	MaxPrice *int64
	Sort     string
	Limit    uint64
	Offset   uint64
}

// CatalogRepo defines persistence port for core catalog CRUD operations.
type CatalogRepo interface {
	StoreProduct(ctx context.Context, product *Product) error
	ListProducts(ctx context.Context, filter ProductRepoFilter) ([]Product, int, error)
	GetProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id string) error
	StoreCategory(ctx context.Context, category *Category) error
	ListCategories(ctx context.Context) ([]Category, error)
	StoreTag(ctx context.Context, tag *Tag) error
	ListTags(ctx context.Context) ([]Tag, error)
}

// GripCatalogRepo defines persistence port for user-visible storefront catalog operations.
type GripCatalogRepo interface {
	ListCategories(ctx context.Context) ([]Category, error)
	ListVisibleProducts(ctx context.Context, actor usermodule.Actor, filter ProductRepoFilter) ([]Product, int, error)
	GetVisibleProduct(ctx context.Context, actor usermodule.Actor, productID string) (Product, error)
	ListSettings(ctx context.Context) ([]Setting, error)
	GetSetting(ctx context.Context, key string) (Setting, error)
}
