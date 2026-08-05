package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/importer"
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
	require.NoError(t, repo.StoreImportedProduct(context.Background(), importer.ImportProductDraft{ID: "p1", SKU: "sku", Title: "Phone"}))
	require.NoError(t, repo.StoreImportedPost(context.Background(), importer.ImportPostDraft{ID: "a1", Slug: "post"}))
}
