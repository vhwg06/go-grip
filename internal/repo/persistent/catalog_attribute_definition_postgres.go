package persistent

import (
	"context"
	"fmt"

	catalogbase "github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

// CatalogAttributeDefinitionRepo persists attribute definitions and their
// embedded enum values.
type CatalogAttributeDefinitionRepo struct {
	db *gorm.DB
}

// NewCatalogAttributeDefinitionRepo creates an attribute definition
// repository backed by PostgreSQL.
func NewCatalogAttributeDefinitionRepo(pg *postgres.Postgres) *CatalogAttributeDefinitionRepo {
	if pg == nil {
		return &CatalogAttributeDefinitionRepo{}
	}
	return newCatalogAttributeDefinitionRepo(pg.Gorm)
}

func newCatalogAttributeDefinitionRepo(db *gorm.DB) *CatalogAttributeDefinitionRepo {
	return &CatalogAttributeDefinitionRepo{db: db}
}

var _ catalogbase.CatalogAttributeDefinitionRepository = (*CatalogAttributeDefinitionRepo)(nil)

// List returns definitions in their configured ordering.
func (r *CatalogAttributeDefinitionRepo) List(ctx context.Context) ([]catalogbase.CatalogAttributeDefinition, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	var rows []models.CatalogBaseDefinition
	if err := db.WithContext(ctx).Order("ordering ASC, display_name ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog definition repo: list: %w", err)
	}
	result := make([]catalogbase.CatalogAttributeDefinition, 0, len(rows))
	for _, row := range rows {
		result = append(result, models.DefinitionToCatalogEntity(row))
	}
	return result, nil
}

// GetByID returns one definition by identity.
func (r *CatalogAttributeDefinitionRepo) GetByID(ctx context.Context, id string) (catalogbase.CatalogAttributeDefinition, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return catalogbase.CatalogAttributeDefinition{}, err
	}
	var row models.CatalogBaseDefinition
	if err := db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return catalogbase.CatalogAttributeDefinition{}, fmt.Errorf("catalog definition repo: get %s: %w", id, err)
	}
	return models.DefinitionToCatalogEntity(row), nil
}

// Store creates one definition.
func (r *CatalogAttributeDefinitionRepo) Store(ctx context.Context, definition catalogbase.CatalogAttributeDefinition) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToDefinition(definition)
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog definition repo: store %s: %w", definition.ID, err)
	}
	return nil
}

// Update replaces one definition and its embedded enum values.
func (r *CatalogAttributeDefinitionRepo) Update(ctx context.Context, definition catalogbase.CatalogAttributeDefinition) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToDefinition(definition)
	if err := db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog definition repo: update %s: %w", definition.ID, err)
	}
	return nil
}

// Delete removes one definition by identity.
func (r *CatalogAttributeDefinitionRepo) Delete(ctx context.Context, id string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&models.CatalogBaseDefinition{}).Error; err != nil {
		return fmt.Errorf("catalog definition repo: delete %s: %w", id, err)
	}
	return nil
}
