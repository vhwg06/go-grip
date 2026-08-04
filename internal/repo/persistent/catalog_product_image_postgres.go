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

// CatalogProductImageRepo persists ProductModel image entities.
type CatalogProductImageRepo struct {
	db *gorm.DB
}

// NewCatalogProductImageRepo creates a ProductModel image repository backed
// by PostgreSQL.
func NewCatalogProductImageRepo(pg *postgres.Postgres) *CatalogProductImageRepo {
	if pg == nil {
		return &CatalogProductImageRepo{}
	}
	return newCatalogProductImageRepo(pg.Gorm)
}

func newCatalogProductImageRepo(db *gorm.DB) *CatalogProductImageRepo {
	return &CatalogProductImageRepo{db: db}
}

var _ repo.CatalogProductImageRepository = (*CatalogProductImageRepo)(nil)

// ListByModelID returns images for one ProductModel in display order.
func (r *CatalogProductImageRepo) ListByModelID(ctx context.Context, modelID string) ([]entity.CatalogProductImage, error) {
	db, err := catalogDB(r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseProductImage
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Order("ordering ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog product image repo: list %s: %w", modelID, err)
	}
	result := make([]entity.CatalogProductImage, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.CatalogProductImage{ID: row.ID, URL: row.URL, Ordering: row.Ordering, PrimaryImage: row.PrimaryImage, CreatedAt: row.CreatedAt})
	}
	return result, nil
}

// Store creates one image for a ProductModel.
func (r *CatalogProductImageRepo) Store(ctx context.Context, modelID string, image entity.CatalogProductImage) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	row := models.CatalogBaseProductImage{ID: image.ID, ModelID: modelID, URL: image.URL, Ordering: image.Ordering, PrimaryImage: image.PrimaryImage, CreatedAt: image.CreatedAt}
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog product image repo: store %s/%s: %w", modelID, image.ID, err)
	}
	return nil
}

// Update replaces one image for a ProductModel.
func (r *CatalogProductImageRepo) Update(ctx context.Context, modelID string, image entity.CatalogProductImage) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	row := models.CatalogBaseProductImage{ID: image.ID, ModelID: modelID, URL: image.URL, Ordering: image.Ordering, PrimaryImage: image.PrimaryImage, CreatedAt: image.CreatedAt}
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog product image repo: update %s/%s: %w", modelID, image.ID, err)
	}
	return nil
}

// Delete removes one image from a ProductModel.
func (r *CatalogProductImageRepo) Delete(ctx context.Context, modelID, imageID string) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ? AND id = ?", modelID, imageID).Delete(&models.CatalogBaseProductImage{}).Error; err != nil {
		return fmt.Errorf("catalog product image repo: delete %s/%s: %w", modelID, imageID, err)
	}
	return nil
}

// DeleteByModelID removes all images belonging to one ProductModel.
func (r *CatalogProductImageRepo) DeleteByModelID(ctx context.Context, modelID string) error {
	db, err := catalogDB(r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Delete(&models.CatalogBaseProductImage{}).Error; err != nil {
		return fmt.Errorf("catalog product image repo: delete model %s: %w", modelID, err)
	}
	return nil
}
