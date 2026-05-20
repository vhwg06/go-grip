package persistent

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestContentRepo(t *testing.T) {
	t.Parallel()
	repo := NewContentRepo(nil)
	require.NoError(t, repo.StoreArticle(context.Background(), &entity.ContentArticle{ID: "a1", Slug: "post", Status: entity.ContentStatusPublished}))
	items, total, err := repo.ListArticles(context.Background(), true, entity.Pagination{})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
}
