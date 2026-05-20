package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestHomepageRepo(t *testing.T) {
	t.Parallel()
	repo := NewHomepageRepo(nil)
	require.NoError(t, repo.Store(context.Background(), &entity.HomepageBlock{ID: "h1", BlockType: "banner", IsActive: true}))
	items, err := repo.List(context.Background(), true)
	require.NoError(t, err)
	require.Len(t, items, 1)
}
