package persistent

import (
	"context"

	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/internal/module/importer"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type ImportRepo struct {
	*postgres.Postgres
	catalog *CatalogRepo
	content *ContentRepo
}

func NewImportRepo(pg *postgres.Postgres, catalog *CatalogRepo, content *ContentRepo) *ImportRepo {
	return &ImportRepo{Postgres: pg, catalog: catalog, content: content}
}

func (r *ImportRepo) StoreImportedProduct(ctx context.Context, draft importer.ImportProductDraft) error {
	p := catalogmodule.Product{
		ID:        draft.ID,
		Title:     draft.Title,
		SKU:       draft.SKU,
		Status:    catalogmodule.ProductStatusDraft,
		CreatedAt: draft.CreatedAt,
		UpdatedAt: draft.UpdatedAt,
	}
	return r.catalog.StoreProduct(ctx, &p)
}

func (r *ImportRepo) StoreImportedPost(ctx context.Context, draft importer.ImportPostDraft) error {
	a := contentmodule.ContentArticle{
		ID:        draft.ID,
		Title:     draft.Title,
		Slug:      draft.Slug,
		Status:    contentmodule.ContentStatusDraft,
		CreatedAt: draft.CreatedAt,
		UpdatedAt: draft.UpdatedAt,
	}
	return r.content.StoreArticle(ctx, &a)
}
