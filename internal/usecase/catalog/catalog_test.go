package catalog

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestCatalogUseCaseCreateProduct(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewCatalogRepo(nil))
	product, err := uc.CreateProduct(context.Background(), entity.Product{Title: "Phone", SKU: "sku", Price: 100})
	require.NoError(t, err)
	require.NotEmpty(t, product.ID)
}
