package persistent

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
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

func (r *ImportRepo) StoreImportedProduct(ctx context.Context, product *entity.Product) error {
	return r.catalog.StoreProduct(ctx, product)
}

func (r *ImportRepo) StoreImportedPost(ctx context.Context, article *entity.ContentArticle) error {
	return r.content.StoreArticle(ctx, article)
}
