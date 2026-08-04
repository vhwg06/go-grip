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

// CatalogVariantRepo persists ProductModel variants.
type CatalogVariantRepo struct {
	db *gorm.DB
}

// NewCatalogVariantRepo creates a variant repository backed by PostgreSQL.
func NewCatalogVariantRepo(pg *postgres.Postgres) *CatalogVariantRepo {
	if pg == nil {
		return &CatalogVariantRepo{}
	}
	return newCatalogVariantRepo(pg.Gorm)
}

func newCatalogVariantRepo(db *gorm.DB) *CatalogVariantRepo {
	return &CatalogVariantRepo{db: db}
}

var _ repo.CatalogVariantRepository = (*CatalogVariantRepo)(nil)

// List returns variants across ProductModels in stable creation order.
func (r *CatalogVariantRepo) List(ctx context.Context) ([]entity.CatalogVariant, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseVariant
	if err := db.WithContext(ctx).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog variant repo: list: %w", err)
	}
	result := make([]entity.CatalogVariant, 0, len(rows))
	for _, row := range rows {
		result = append(result, variantToEntity(row))
	}
	return result, nil
}

// ListByModelID returns variants for one ProductModel.
func (r *CatalogVariantRepo) ListByModelID(ctx context.Context, modelID string) ([]entity.CatalogVariant, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseVariant
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog variant repo: list %s: %w", modelID, err)
	}
	result := make([]entity.CatalogVariant, 0, len(rows))
	for _, row := range rows {
		result = append(result, variantToEntity(row))
	}
	return result, nil
}

// GetByID returns one variant for its owning ProductModel.
func (r *CatalogVariantRepo) GetByID(ctx context.Context, modelID, id string) (entity.CatalogVariant, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return entity.CatalogVariant{}, err
	}
	var row models.CatalogBaseVariant
	if err := db.WithContext(ctx).Where("model_id = ? AND id = ?", modelID, id).First(&row).Error; err != nil {
		return entity.CatalogVariant{}, fmt.Errorf("catalog variant repo: get %s/%s: %w", modelID, id, err)
	}
	return variantToEntity(row), nil
}

func variantToEntity(row models.CatalogBaseVariant) entity.CatalogVariant {
	model := models.ProductModelToCatalogEntity(models.CatalogBaseProductModel{}, nil, nil, []models.CatalogBaseVariant{row})
	return model.Variants[0]
}

// Store creates one variant for a ProductModel.
func (r *CatalogVariantRepo) Store(ctx context.Context, modelID string, variant entity.CatalogVariant) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToVariants(modelID, []entity.CatalogVariant{variant})[0]
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog variant repo: store %s/%s: %w", modelID, variant.ID, err)
	}
	return nil
}

// Update replaces one variant for a ProductModel.
func (r *CatalogVariantRepo) Update(ctx context.Context, modelID string, variant entity.CatalogVariant) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToVariants(modelID, []entity.CatalogVariant{variant})[0]
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog variant repo: update %s/%s: %w", modelID, variant.ID, err)
	}
	return nil
}

// Delete removes one variant from a ProductModel.
func (r *CatalogVariantRepo) Delete(ctx context.Context, modelID, id string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ? AND id = ?", modelID, id).Delete(&models.CatalogBaseVariant{}).Error; err != nil {
		return fmt.Errorf("catalog variant repo: delete %s/%s: %w", modelID, id, err)
	}
	return nil
}

// DeleteByModelID removes all variants belonging to one ProductModel.
func (r *CatalogVariantRepo) DeleteByModelID(ctx context.Context, modelID string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Delete(&models.CatalogBaseVariant{}).Error; err != nil {
		return fmt.Errorf("catalog variant repo: delete model %s: %w", modelID, err)
	}
	return nil
}
