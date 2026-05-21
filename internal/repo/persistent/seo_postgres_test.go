package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestSEORepo(t *testing.T) {
	t.Parallel()
	repo := NewSEORepo(nil)
	require.NoError(t, repo.Store(context.Background(), &entity.SeoMetadata{ID: "s1", OwnerType: "product", OwnerID: "p1", Slug: "phone"}))
	meta, err := repo.GetByOwner(context.Background(), "product", "p1")
	require.NoError(t, err)
	require.Equal(t, "phone", meta.Slug)
}
