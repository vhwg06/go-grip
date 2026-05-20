package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestMediaRepo(t *testing.T) {
	t.Parallel()
	repo := NewMediaRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &entity.MediaAsset{ID: "m1", FileName: "a.jpg"}))
	items, total, err := repo.List(context.Background(), entity.Pagination{})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
	require.NoError(t, repo.Delete(context.Background(), "m1"))
}
