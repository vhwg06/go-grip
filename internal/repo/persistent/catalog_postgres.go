package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type GripCatalogRepo struct {
	*postgres.Postgres
}

func NewGripCatalogRepo(pg *postgres.Postgres) *GripCatalogRepo {
	return &GripCatalogRepo{Postgres: pg}
}

var _ repo.CatalogRepository = (*GripCatalogRepo)(nil)

func (r *GripCatalogRepo) ListVisibleProducts(ctx context.Context, actor entity.Actor, filter repo.ProductFilter) ([]entity.Product, int, error) {
	threshold := actor.TrustLevel
	if threshold < 0 {
		threshold = -1
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

	products := make([]entity.Product, 0, len(rows))
	for _, row := range rows {
		products = append(products, models.ProductToEntity(row))
	}

	return products, int(total), nil
}

func (r *GripCatalogRepo) GetVisibleProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error) {
	threshold := actor.TrustLevel
	if threshold < 0 {
		threshold = -1
	}

	var row models.Product
	if err := r.Gorm.WithContext(ctx).
		Where("id = ?", productID).
		Where("is_active = ?", true).
		Where("visibility_level <= ?", threshold).
		First(&row).Error; err != nil {
		return entity.Product{}, fmt.Errorf("GripCatalogRepo.GetVisibleProduct: %w", err)
	}

	entityProduct := models.ProductToEntity(row)

	var detailRows []models.ProductDetail
	if err := r.Gorm.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&detailRows).Error; err == nil {
		entityProduct.Specs = make([]entity.ProductSpecItem, 0, len(detailRows))
		for _, detail := range detailRows {
			switch detail.Key {
			case "sku":
				entityProduct.SKU = detail.Value
			case "brand":
				entityProduct.Brand = detail.Value
			default:
				entityProduct.Specs = append(entityProduct.Specs, models.DetailToEntity(detail))
			}
		}
	}

	return entityProduct, nil
}

func (r *GripCatalogRepo) ListCategories(ctx context.Context) ([]entity.Category, error) {
	var rows []models.Category
	if err := r.Gorm.WithContext(ctx).
		Order("sort_order ASC, name ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("GripCatalogRepo.ListCategories: %w", err)
	}

	categories := make([]entity.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, models.CategoryToEntity(row))
	}

	return categories, nil
}

func (r *GripCatalogRepo) ListSettings(ctx context.Context) ([]entity.Setting, error) {
	var rows []models.Setting
	if err := r.Gorm.WithContext(ctx).
		Order("key ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("GripCatalogRepo.ListSettings: %w", err)
	}

	settings := make([]entity.Setting, 0, len(rows))
	for _, row := range rows {
		settings = append(settings, models.SettingToEntity(row))
	}

	return settings, nil
}

func (r *GripCatalogRepo) GetSetting(ctx context.Context, key string) (entity.Setting, error) {
	var row models.Setting
	if err := r.Gorm.WithContext(ctx).
		Where("key = ?", key).
		First(&row).Error; err != nil {
		return entity.Setting{}, fmt.Errorf("GripCatalogRepo.GetSetting: %w", err)
	}

	return models.SettingToEntity(row), nil
}
