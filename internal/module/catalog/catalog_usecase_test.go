package catalog_test

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestCatalogUseCaseCreateProduct(t *testing.T) {
	t.Parallel()
	repo := persistent.NewCatalogRepo(nil)
	uc := catalog.NewCatalogUseCase(repo, nil)

	p, err := uc.CreateProduct(context.Background(), catalog.Product{
		Title: "Test Keyboard",
		SKU:   "KB-001",
		Price: 100,
	})
	require.NoError(t, err)
	require.Equal(t, "Test Keyboard", p.Title)
	require.Equal(t, catalog.ProductStatusDraft, p.Status)
}
