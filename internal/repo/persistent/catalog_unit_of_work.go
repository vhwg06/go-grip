package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

// CatalogUnitOfWork coordinates catalog repositories on one database
// transaction. The callback receives repository adapters bound to the
// transaction, so a child operation cannot silently fall back to the base
// connection.
type CatalogUnitOfWork struct {
	db *gorm.DB
}

// NewCatalogUnitOfWork creates a catalog transaction coordinator for the
// configured PostgreSQL connection.
func NewCatalogUnitOfWork(pg *postgres.Postgres) *CatalogUnitOfWork {
	if pg == nil {
		return &CatalogUnitOfWork{}
	}
	return &CatalogUnitOfWork{db: pg.Gorm}
}

var _ repo.CatalogUnitOfWork = (*CatalogUnitOfWork)(nil)

// Within executes fn atomically and rolls back when any repository operation
// returns an error.
func (u *CatalogUnitOfWork) Within(ctx context.Context, fn func(repo.CatalogRepositories) error) error {
	if u == nil || u.db == nil {
		return fmt.Errorf("catalog unit of work: database is not configured")
	}
	if fn == nil {
		return fmt.Errorf("catalog unit of work: callback is nil")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(catalogRepositories(tx))
	})
}

// NewCatalogRepositories wires one set of catalog repositories to the base
// application connection. Transaction callbacks use the same factory with a
// transaction-bound *gorm.DB and are kept private to this infrastructure
// package.
func NewCatalogRepositories(pg *postgres.Postgres) repo.CatalogRepositories {
	if pg == nil {
		return repo.CatalogRepositories{}
	}
	return catalogRepositories(pg.Gorm)
}

func catalogRepositories(db *gorm.DB) repo.CatalogRepositories {
	return repo.CatalogRepositories{
		Categories:        newCatalogCategoryRepo(db),
		Definitions:       newCatalogAttributeDefinitionRepo(db),
		Masters:           newCatalogMasterRepo(db),
		ProductModels:     newCatalogProductModelRepo(db),
		ProductImages:     newCatalogProductImageRepo(db),
		VariantDimensions: newCatalogVariantDimensionRepo(db),
		Variants:          newCatalogVariantRepo(db),
	}
}
