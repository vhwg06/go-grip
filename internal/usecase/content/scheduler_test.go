package content

import (
	"context"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestContentSchedulerCatchUp(t *testing.T) {
	t.Parallel()
	repo := persistent.NewContentRepo(nil)
	uc := New(repo)
	past := time.Now().Add(-time.Hour)
	_, err := uc.CreateArticle(context.Background(), entity.ContentArticle{Title: "Post", Slug: "post", Status: entity.ContentStatusScheduled, ScheduledAt: &past})
	require.NoError(t, err)
	count, err := uc.CatchUpScheduled(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
