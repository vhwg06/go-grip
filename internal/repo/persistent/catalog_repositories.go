package persistent

import (
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

// NewCatalogRepositories wires the explicit Catalog Base CRUD repositories
// to the base application connection. UnitOfWork supplies a transaction
// through context when a use case needs atomic orchestration.
func NewCatalogRepositories(pg *postgres.Postgres) repo.CatalogRepositories {
	return repo.CatalogRepositories{
		Categories:        NewCatalogCategoryRepo(pg),
		Definitions:       NewCatalogAttributeDefinitionRepo(pg),
		Masters:           NewCatalogMasterRepo(pg),
		ProductModels:     NewCatalogProductModelRepo(pg),
		ProductImages:     NewCatalogProductImageRepo(pg),
		VariantDimensions: NewCatalogVariantDimensionRepo(pg),
		Variants:          NewCatalogVariantRepo(pg),
	}
}
