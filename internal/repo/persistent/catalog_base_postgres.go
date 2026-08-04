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

// CatalogBaseRepo is the infrastructure adapter for the Catalog Base
// aggregate.  The application layer sees only domain entities and a single
// transaction boundary.
type CatalogBaseRepo struct {
	*postgres.Postgres
}

func NewCatalogBaseRepo(pg *postgres.Postgres) *CatalogBaseRepo {
	return &CatalogBaseRepo{Postgres: pg}
}

var _ repo.CatalogBaseRepository = (*CatalogBaseRepo)(nil)

func (r *CatalogBaseRepo) LoadCatalogBase(ctx context.Context) (entity.CatalogSnapshot, error) {
	var categories []models.CatalogBaseCategory
	if err := r.Gorm.WithContext(ctx).Order("position ASC, name ASC, id ASC").Find(&categories).Error; err != nil {
		return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load categories: %w", err)
	}
	var definitions []models.CatalogBaseDefinition
	if err := r.Gorm.WithContext(ctx).Order("ordering ASC, display_name ASC, id ASC").Find(&definitions).Error; err != nil {
		return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load definitions: %w", err)
	}
	var masters []models.CatalogBaseMaster
	if err := r.Gorm.WithContext(ctx).Order("kind ASC, name ASC, id ASC").Find(&masters).Error; err != nil {
		return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load masters: %w", err)
	}
	var productModels []models.CatalogBaseProductModel
	if err := r.Gorm.WithContext(ctx).Order("created_at ASC, id ASC").Find(&productModels).Error; err != nil {
		return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load product models: %w", err)
	}

	snapshot := entity.CatalogSnapshot{
		Categories:  make([]entity.CatalogCategory, 0, len(categories)),
		Definitions: make([]entity.CatalogAttributeDefinition, 0, len(definitions)),
		Masters:     make([]entity.CatalogMaster, 0, len(masters)),
		Models:      make([]entity.CatalogProductModel, 0, len(productModels)),
	}
	for _, row := range categories {
		snapshot.Categories = append(snapshot.Categories, models.CategoryToCatalogEntity(row))
	}
	for _, row := range definitions {
		snapshot.Definitions = append(snapshot.Definitions, models.DefinitionToCatalogEntity(row))
	}
	for _, row := range masters {
		snapshot.Masters = append(snapshot.Masters, models.MasterToCatalogEntity(row))
	}
	for _, row := range productModels {
		var images []models.CatalogBaseProductImage
		if err := r.Gorm.WithContext(ctx).Where("model_id = ?", row.ID).Order("ordering ASC, id ASC").Find(&images).Error; err != nil {
			return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load model images: %w", err)
		}
		var dimensions []models.CatalogBaseDimension
		if err := r.Gorm.WithContext(ctx).Where("model_id = ?", row.ID).Order("created_at ASC, id ASC").Find(&dimensions).Error; err != nil {
			return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load model dimensions: %w", err)
		}
		var variants []models.CatalogBaseVariant
		if err := r.Gorm.WithContext(ctx).Where("model_id = ?", row.ID).Order("created_at ASC, id ASC").Find(&variants).Error; err != nil {
			return entity.CatalogSnapshot{}, fmt.Errorf("catalog base repo: load model variants: %w", err)
		}
		snapshot.Models = append(snapshot.Models, models.ProductModelToCatalogEntity(row, images, dimensions, variants))
	}
	return snapshot, nil
}

func (r *CatalogBaseRepo) SaveCatalogBase(ctx context.Context, snapshot entity.CatalogSnapshot) error {
	return r.Gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := deleteCatalogBaseRows(tx); err != nil {
			return err
		}
		if err := createRows(tx, snapshot); err != nil {
			return fmt.Errorf("catalog base repo: save aggregate: %w", err)
		}
		return nil
	})
}

func deleteCatalogBaseRows(tx *gorm.DB) error {
	for _, value := range []any{
		&models.CatalogBaseVariant{},
		&models.CatalogBaseDimension{},
		&models.CatalogBaseProductImage{},
		&models.CatalogBaseProductModel{},
		&models.CatalogBaseMaster{},
		&models.CatalogBaseDefinition{},
		&models.CatalogBaseCategory{},
	} {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(value).Error; err != nil {
			return fmt.Errorf("catalog base repo: clear aggregate: %w", err)
		}
	}
	return nil
}

func createRows(tx *gorm.DB, snapshot entity.CatalogSnapshot) error {
	if err := createCategories(tx, snapshot.Categories); err != nil {
		return err
	}
	definitions := make([]models.CatalogBaseDefinition, 0, len(snapshot.Definitions))
	for _, value := range snapshot.Definitions {
		definitions = append(definitions, models.CatalogEntityToDefinition(value))
	}
	if len(definitions) > 0 {
		if err := tx.Create(&definitions).Error; err != nil {
			return fmt.Errorf("catalog base repo: create definitions: %w", err)
		}
	}
	masters := make([]models.CatalogBaseMaster, 0, len(snapshot.Masters))
	for _, value := range snapshot.Masters {
		masters = append(masters, models.CatalogEntityToMaster(value))
	}
	if len(masters) > 0 {
		if err := tx.Create(&masters).Error; err != nil {
			return fmt.Errorf("catalog base repo: create masters: %w", err)
		}
	}
	productModels := make([]models.CatalogBaseProductModel, 0, len(snapshot.Models))
	for _, value := range snapshot.Models {
		productModels = append(productModels, models.CatalogEntityToProductModel(value))
	}
	if len(productModels) > 0 {
		if err := tx.Create(&productModels).Error; err != nil {
			return fmt.Errorf("catalog base repo: create product models: %w", err)
		}
	}
	for _, model := range snapshot.Models {
		images := models.CatalogEntityToImages(model.ID, model.Images)
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return fmt.Errorf("catalog base repo: create model images: %w", err)
			}
		}
		dimensions := models.CatalogEntityToDimensions(model.ID, model.Dimensions)
		if len(dimensions) > 0 {
			if err := tx.Create(&dimensions).Error; err != nil {
				return fmt.Errorf("catalog base repo: create dimensions: %w", err)
			}
		}
		variants := models.CatalogEntityToVariants(model.ID, model.Variants)
		if len(variants) > 0 {
			if err := tx.Create(&variants).Error; err != nil {
				return fmt.Errorf("catalog base repo: create variants: %w", err)
			}
		}
	}
	return nil
}

// createCategories inserts parent categories before their children.  The
// schema intentionally keeps the hierarchy as a self-referencing foreign key,
// so inserting the complete slice as one batch can violate the immediate FK
// when a child happens to sort before its parent.
func createCategories(tx *gorm.DB, values []entity.CatalogCategory) error {
	remaining := append([]entity.CatalogCategory(nil), values...)
	for len(remaining) > 0 {
		progress := false
		next := make([]entity.CatalogCategory, 0, len(remaining))
		for _, value := range remaining {
			if value.ParentID != nil {
				found := false
				for _, candidate := range remaining {
					if candidate.ID == *value.ParentID {
						found = true
						break
					}
				}
				if found {
					// The parent may already have been inserted in an earlier
					// round; checking the database keeps this safe for chains.
					var count int64
					if err := tx.Model(&models.CatalogBaseCategory{}).Where("id = ?", *value.ParentID).Count(&count).Error; err != nil {
						return fmt.Errorf("catalog base repo: check category parent: %w", err)
					}
					if count == 0 {
						next = append(next, value)
						continue
					}
				}
			}
			row := models.CatalogEntityToCategory(value)
			if err := tx.Create(&row).Error; err != nil {
				return fmt.Errorf("catalog base repo: create category: %w", err)
			}
			progress = true
		}
		if !progress {
			return fmt.Errorf("catalog base repo: category hierarchy contains an unresolved parent")
		}
		remaining = next
	}
	return nil
}
