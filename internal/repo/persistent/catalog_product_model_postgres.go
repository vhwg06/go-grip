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

// CatalogProductModelRepo persists ProductModel root records. Images,
// dimensions, and variants are loaded by their own repositories and composed
// by the application orchestrator when a complete model is needed.
type CatalogProductModelRepo struct {
	db *gorm.DB
}

// NewCatalogProductModelRepo creates a ProductModel repository backed by
// PostgreSQL.
func NewCatalogProductModelRepo(pg *postgres.Postgres) *CatalogProductModelRepo {
	if pg == nil {
		return &CatalogProductModelRepo{}
	}
	return newCatalogProductModelRepo(pg.Gorm)
}

func newCatalogProductModelRepo(db *gorm.DB) *CatalogProductModelRepo {
	return &CatalogProductModelRepo{db: db}
}

var _ repo.CatalogProductModelRepository = (*CatalogProductModelRepo)(nil)

// List returns ProductModel roots in creation order.
func (r *CatalogProductModelRepo) List(ctx context.Context) ([]entity.CatalogProductModel, error) {
	db, err := catalogDB(r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseProductModel
	if err := db.WithContext(ctx).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog product model repo: list: %w", err)
	}
	result := make([]entity.CatalogProductModel, 0, len(rows))
	for _, row := range rows {
		result = append(result, models.ProductModelToCatalogEntity(row, nil, nil, nil))
	}
	return result, nil
}

// GetByID returns one ProductModel root by identity.
func (r *CatalogProductModelRepo) GetByID(ctx context.Context, id string) (entity.CatalogProductModel, error) {
	db, err := catalogDB(r.db)
	if err != nil {
		return entity.CatalogProductModel{}, err
	}
	var row models.CatalogBaseProductModel
	if err := db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return entity.CatalogProductModel{}, fmt.Errorf("catalog product model repo: get %s: %w", id, err)
	}
	return models.ProductModelToCatalogEntity(row, nil, nil, nil), nil
}

// Store creates one ProductModel root.
func (r *CatalogProductModelRepo) Store(ctx context.Context, model entity.CatalogProductModel) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToProductModel(model)
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog product model repo: store %s: %w", model.ID, err)
	}
	return nil
}

// Update replaces one ProductModel root.
func (r *CatalogProductModelRepo) Update(ctx context.Context, model entity.CatalogProductModel) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToProductModel(model)
	if err := db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog product model repo: update %s: %w", model.ID, err)
	}
	return nil
}

// Delete removes one ProductModel root by identity.
func (r *CatalogProductModelRepo) Delete(ctx context.Context, id string) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.CatalogBaseProductModel{}).Error; err != nil {
		return fmt.Errorf("catalog product model repo: delete %s: %w", id, err)
	}
	return nil
}
