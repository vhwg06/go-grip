package catalog

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestCatalogUseCaseSearchFilter(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewCatalogRepo(nil))
	_, err := uc.CreateProduct(context.Background(), entity.Product{Title: "Phone", SKU: "sku", Price: 100, Brand: "Acme"})
	require.NoError(t, err)
	items, total, err := uc.ListProducts(context.Background(), entity.ProductFilter{Brand: "Acme"})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
}
