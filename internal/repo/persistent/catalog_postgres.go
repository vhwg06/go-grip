package persistent

import (
	"context"
	"fmt"

	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type GripCatalogRepo struct {
	*postgres.Postgres
}

func NewGripCatalogRepo(pg *postgres.Postgres) *GripCatalogRepo {
	return &GripCatalogRepo{Postgres: pg}
}

func (r *GripCatalogRepo) ListVisibleProducts(ctx context.Context, actor usermodule.Actor, filter catalogmodule.ProductRepoFilter) ([]catalogmodule.Product, int, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return nil, 0, nil
	}
	threshold := 0
	if actor.IsAdmin {
		threshold = 100
	}

	query := r.Gorm.WithContext(ctx).Model(&models.Product{}).
		Where("is_active = ?", true).
		Where("visibility_level <= ?", threshold)

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where(
			"COALESCE(title, name) ILIKE ? OR name ILIKE ? OR description ILIKE ?",
			keyword, keyword, keyword,
		)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("GripCatalogRepo.ListVisibleProducts: count: %w", err)
	}

	orderClause := "sort_order ASC, created_at DESC"
	switch filter.Sort {
	case "price_asc":
		orderClause = "price ASC, created_at DESC"
	case "price_desc":
		orderClause = "price DESC, created_at DESC"
	case "sales":
		orderClause = "sold_count DESC, created_at DESC"
	case "rating":
		orderClause = "rating DESC, created_at DESC"
	}

	var rows []models.Product
	if err := query.
		Limit(int(filter.Limit)).
		Offset(int(filter.Offset)).
		Order(orderClause).
		Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("GripCatalogRepo.ListVisibleProducts: find: %w", err)
	}

	products := make([]catalogmodule.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, models.ProductToModule(row))
	}

	return products, int(total), nil
}

func (r *GripCatalogRepo) GetVisibleProduct(ctx context.Context, actor usermodule.Actor, productID string) (catalogmodule.Product, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return catalogmodule.Product{}, catalogmodule.ErrNotFound
	}
	threshold := 0
	if actor.IsAdmin {
		threshold = 100
	}

	var row models.Product
	if err := r.Gorm.WithContext(ctx).
		Where("id = ?", productID).
		Where("is_active = ?", true).
		Where("visibility_level <= ?", threshold).
		First(&row).Error; err != nil {
		return catalogmodule.Product{}, fmt.Errorf("GripCatalogRepo.GetVisibleProduct: %w", err)
	}

	prod := models.ProductToModule(row)

	var detailRows []models.ProductDetail
	if err := r.Gorm.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&detailRows).Error; err == nil {
		prod.Specs = make([]catalogmodule.ProductSpecItem, 0, len(detailRows))
		for _, detail := range detailRows {
			switch detail.Key {
			case "sku":
				prod.SKU = detail.Value
			case "brand":
				prod.Brand = detail.Value
			default:
				prod.Specs = append(prod.Specs, catalogmodule.ProductSpecItem{Key: detail.Key, Value: detail.Value})
			}
		}
	}

	return prod, nil
}

func (r *GripCatalogRepo) ListCategories(ctx context.Context) ([]catalogmodule.Category, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return nil, nil
	}
	var rows []models.Category
	if err := r.Gorm.WithContext(ctx).
		Order("sort_order ASC, name ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("GripCatalogRepo.ListCategories: %w", err)
	}

	categories := make([]catalogmodule.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, models.CategoryToModule(row))
	}

	return categories, nil
}

func (r *GripCatalogRepo) ListSettings(ctx context.Context) ([]catalogmodule.Setting, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return nil, nil
	}
	var rows []models.Setting
	if err := r.Gorm.WithContext(ctx).
		Order("key ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("GripCatalogRepo.ListSettings: %w", err)
	}

	settings := make([]catalogmodule.Setting, 0, len(rows))
	for _, row := range rows {
		settings = append(settings, models.SettingToModule(row))
	}

	return settings, nil
}

func (r *GripCatalogRepo) GetSetting(ctx context.Context, key string) (catalogmodule.Setting, error) {
	if r.Postgres == nil || r.Gorm == nil {
		return catalogmodule.Setting{}, catalogmodule.ErrNotFound
	}
	var row models.Setting
	if err := r.Gorm.WithContext(ctx).
		Where("key = ?", key).
		First(&row).Error; err != nil {
		return catalogmodule.Setting{}, fmt.Errorf("GripCatalogRepo.GetSetting: %w", err)
	}

	return models.SettingToModule(row), nil
}
