package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestImportRepo(t *testing.T) {
	t.Parallel()
	catalogRepo := NewCatalogRepo(nil)
	contentRepo := NewContentRepo(nil)
	if contentRepo.Postgres == nil || contentRepo.Pool == nil {
		t.Skip("Skipping TestImportRepo because PostgreSQL connection is nil")
	}
	repo := NewImportRepo(nil, catalogRepo, contentRepo)
	require.NoError(t, repo.StoreImportedProduct(context.Background(), &entity.Product{ID: "p1", SKU: "sku", Title: "Phone"}))
	require.NoError(t, repo.StoreImportedPost(context.Background(), &entity.ContentArticle{ID: "a1", Slug: "post"}))
}
