package persistent

import (
	"context"
	"fmt"

	catalogbase "github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

// CatalogVariantDimensionRepo persists ProductModel variant dimensions.
type CatalogVariantDimensionRepo struct {
	db *gorm.DB
}

// NewCatalogVariantDimensionRepo creates a variant dimension repository
// backed by PostgreSQL.
func NewCatalogVariantDimensionRepo(pg *postgres.Postgres) *CatalogVariantDimensionRepo {
	if pg == nil {
		return &CatalogVariantDimensionRepo{}
	}
	return newCatalogVariantDimensionRepo(pg.Gorm)
}

func newCatalogVariantDimensionRepo(db *gorm.DB) *CatalogVariantDimensionRepo {
	return &CatalogVariantDimensionRepo{db: db}
}

var _ catalogbase.CatalogVariantDimensionRepository = (*CatalogVariantDimensionRepo)(nil)

// ListByModelID returns dimensions belonging to one ProductModel.
func (r *CatalogVariantDimensionRepo) ListByModelID(ctx context.Context, modelID string) ([]catalogbase.CatalogVariantDimension, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseDimension
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog variant dimension repo: list %s: %w", modelID, err)
	}
	result := make([]catalogbase.CatalogVariantDimension, 0, len(rows))
	for _, row := range rows {
		result = append(result, dimensionToEntity(row))
	}
	return result, nil
}

func dimensionToEntity(row models.CatalogBaseDimension) catalogbase.CatalogVariantDimension {
	model := models.ProductModelToCatalogEntity(models.CatalogBaseProductModel{}, nil, []models.CatalogBaseDimension{row}, nil)
	return model.Dimensions[0]
}

// GetByID returns one variant dimension for its owning ProductModel.
func (r *CatalogVariantDimensionRepo) GetByID(ctx context.Context, modelID, id string) (catalogbase.CatalogVariantDimension, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return catalogbase.CatalogVariantDimension{}, err
	}
	var row models.CatalogBaseDimension
	if err := db.WithContext(ctx).Where("model_id = ? AND id = ?", modelID, id).First(&row).Error; err != nil {
		return catalogbase.CatalogVariantDimension{}, fmt.Errorf("catalog variant dimension repo: get %s/%s: %w", modelID, id, err)
	}
	return dimensionToEntity(row), nil
}

// Store creates one variant dimension.
func (r *CatalogVariantDimensionRepo) Store(ctx context.Context, modelID string, dimension catalogbase.CatalogVariantDimension) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToDimensions(modelID, []catalogbase.CatalogVariantDimension{dimension})[0]
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog variant dimension repo: store %s/%s: %w", modelID, dimension.ID, err)
	}
	return nil
}

// Update replaces one variant dimension.
func (r *CatalogVariantDimensionRepo) Update(ctx context.Context, modelID string, dimension catalogbase.CatalogVariantDimension) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToDimensions(modelID, []catalogbase.CatalogVariantDimension{dimension})[0]
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog variant dimension repo: update %s/%s: %w", modelID, dimension.ID, err)
	}
	return nil
}

// Delete removes one variant dimension.
func (r *CatalogVariantDimensionRepo) Delete(ctx context.Context, modelID, id string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ? AND id = ?", modelID, id).Delete(&models.CatalogBaseDimension{}).Error; err != nil {
		return fmt.Errorf("catalog variant dimension repo: delete %s/%s: %w", modelID, id, err)
	}
	return nil
}

// DeleteByModelID removes all dimensions belonging to one ProductModel.
func (r *CatalogVariantDimensionRepo) DeleteByModelID(ctx context.Context, modelID string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Delete(&models.CatalogBaseDimension{}).Error; err != nil {
		return fmt.Errorf("catalog variant dimension repo: delete model %s: %w", modelID, err)
	}
	return nil
}
