package persistent

import (
	"context"
	"testing"

	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/stretchr/testify/require"
)

func TestCatalogRepoProducts(t *testing.T) {
	t.Parallel()

	catalogRepo := NewCatalogRepo(nil)
	product := catalogmodule.Product{ID: "p1", Title: "Phone", SKU: "sku-1", Price: 100, Brand: "Acme"}

	require.NoError(t, catalogRepo.StoreProduct(context.Background(), &product))
	got, err := catalogRepo.GetProduct(context.Background(), "p1")
	require.NoError(t, err)
	require.Equal(t, product.SKU, got.SKU)

	items, total, err := catalogRepo.ListProducts(context.Background(), catalogmodule.ProductRepoFilter{Brand: "Acme"})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
}
