package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

// CatalogCategoryRepo persists catalog category entities.
type CatalogCategoryRepo struct {
	db *gorm.DB
}

// NewCatalogCategoryRepo creates a category repository backed by PostgreSQL.
func NewCatalogCategoryRepo(pg *postgres.Postgres) *CatalogCategoryRepo {
	if pg == nil {
		return &CatalogCategoryRepo{}
	}
	return newCatalogCategoryRepo(pg.Gorm)
}

func newCatalogCategoryRepo(db *gorm.DB) *CatalogCategoryRepo {
	return &CatalogCategoryRepo{db: db}
}

var _ repo.CatalogCategoryRepository = (*CatalogCategoryRepo)(nil)

// List returns categories in their stable display order.
func (r *CatalogCategoryRepo) List(ctx context.Context) ([]entity.CatalogCategory, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseCategory
	if err := db.WithContext(ctx).Order("position ASC, name ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog category repo: list: %w", err)
	}
	result := make([]entity.CatalogCategory, 0, len(rows))
	for _, row := range rows {
		result = append(result, models.CategoryToCatalogEntity(row))
	}
	return result, nil
}

// GetByID returns one category by identity.
func (r *CatalogCategoryRepo) GetByID(ctx context.Context, id string) (entity.CatalogCategory, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return entity.CatalogCategory{}, err
	}
	var row models.CatalogBaseCategory
	if err := db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return entity.CatalogCategory{}, fmt.Errorf("catalog category repo: get %s: %w", id, err)
	}
	return models.CategoryToCatalogEntity(row), nil
}

// Store creates one category.
func (r *CatalogCategoryRepo) Store(ctx context.Context, category entity.CatalogCategory) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToCategory(category)
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog category repo: store %s: %w", category.ID, err)
	}
	return nil
}

// Update replaces the persisted fields of one category.
func (r *CatalogCategoryRepo) Update(ctx context.Context, category entity.CatalogCategory) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToCategory(category)
	if err := db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog category repo: update %s: %w", category.ID, err)
	}
	return nil
}

// Delete removes one category by identity.
func (r *CatalogCategoryRepo) Delete(ctx context.Context, id string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.CatalogBaseCategory{}).Error; err != nil {
		return fmt.Errorf("catalog category repo: delete %s: %w", id, err)
	}
	return nil
}
