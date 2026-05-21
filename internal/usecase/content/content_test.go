package content

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestContentUseCaseArticle(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewContentRepo(nil))
	article, err := uc.CreateArticle(context.Background(), entity.ContentArticle{Title: "Post", Slug: "post"})
	require.NoError(t, err)
	require.NotEmpty(t, article.ID)
}
