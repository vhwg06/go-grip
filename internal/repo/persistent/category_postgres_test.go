package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestCatalogRepoCategories(t *testing.T) {
	t.Parallel()
	repo := NewCatalogRepo(nil)
	require.NoError(t, repo.StoreCategory(context.Background(), &entity.Category{ID: "c1", Name: "Phones"}))
	items, err := repo.ListCategories(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
}
