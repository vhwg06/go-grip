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

// CatalogMasterRepo persists Material, Finish, and Pack masters.
type CatalogMasterRepo struct {
	db *gorm.DB
}

// NewCatalogMasterRepo creates a catalog master repository backed by
// PostgreSQL.
func NewCatalogMasterRepo(pg *postgres.Postgres) *CatalogMasterRepo {
	if pg == nil {
		return &CatalogMasterRepo{}
	}
	return newCatalogMasterRepo(pg.Gorm)
}

func newCatalogMasterRepo(db *gorm.DB) *CatalogMasterRepo {
	return &CatalogMasterRepo{db: db}
}

var _ repo.CatalogMasterRepository = (*CatalogMasterRepo)(nil)

// List returns masters for a kind. An empty kind returns all supported rows.
func (r *CatalogMasterRepo) List(ctx context.Context, kind string) ([]entity.CatalogMaster, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return nil, err
	}
	query := db.WithContext(ctx)
	if kind != "" {
		query = query.Where("kind = ?", kind)
	}
	var rows []models.CatalogBaseMaster
	if err := query.Order("kind ASC, name ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("catalog master repo: list %s: %w", kind, err)
	}
	result := make([]entity.CatalogMaster, 0, len(rows))
	for _, row := range rows {
		result = append(result, models.MasterToCatalogEntity(row))
	}
	return result, nil
}

// GetByID returns one master by kind and identity.
func (r *CatalogMasterRepo) GetByID(ctx context.Context, kind, id string) (entity.CatalogMaster, error) {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return entity.CatalogMaster{}, err
	}
	var row models.CatalogBaseMaster
	if err := db.WithContext(ctx).Where("kind = ? AND id = ?", kind, id).First(&row).Error; err != nil {
		return entity.CatalogMaster{}, fmt.Errorf("catalog master repo: get %s/%s: %w", kind, id, err)
	}
	return models.MasterToCatalogEntity(row), nil
}

// Store creates one master.
func (r *CatalogMasterRepo) Store(ctx context.Context, master entity.CatalogMaster) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToMaster(master)
	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("catalog master repo: store %s: %w", master.ID, err)
	}
	return nil
}

// Update replaces one master.
func (r *CatalogMasterRepo) Update(ctx context.Context, master entity.CatalogMaster) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	row := models.CatalogEntityToMaster(master)
	if err := db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("catalog master repo: update %s: %w", master.ID, err)
	}
	return nil
}

// Delete removes one master by kind and identity.
func (r *CatalogMasterRepo) Delete(ctx context.Context, kind, id string) error {
	db, err := catalogDBForContext(ctx, r.db)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("kind = ? AND id = ?", kind, id).Delete(&models.CatalogBaseMaster{}).Error; err != nil {
		return fmt.Errorf("catalog master repo: delete %s/%s: %w", kind, id, err)
	}
	return nil
}
