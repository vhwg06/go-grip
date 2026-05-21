package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestCatalogRepoTags(t *testing.T) {
	t.Parallel()
	repo := NewCatalogRepo(nil)
	require.NoError(t, repo.StoreTag(context.Background(), &entity.Tag{ID: "t1", Name: "hot"}))
	items, err := repo.ListTags(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
}
