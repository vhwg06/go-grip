package content

import (
	"context"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestContentUseCaseArticle(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewContentRepo(nil))
	ctx := context.Background()

	// Test CreateArticle
	article, err := uc.CreateArticle(ctx, entity.ContentArticle{Title: "Post", Slug: "post"})
	require.NoError(t, err)
	require.NotEmpty(t, article.ID)
	require.Equal(t, entity.ContentStatusDraft, article.Status)

	// Test UpdateArticle
	article.Title = "Updated Post"
	updated, err := uc.UpdateArticle(ctx, article)
	require.NoError(t, err)
	require.Equal(t, "Updated Post", updated.Title)

	// Test GetArticle
	fetched, err := uc.GetArticle(ctx, article.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Post", fetched.Title)

	// Test DeleteArticle
	err = uc.DeleteArticle(ctx, article.ID)
	require.NoError(t, err)

	_, err = uc.GetArticle(ctx, article.ID)
	require.ErrorIs(t, err, entity.ErrNotFound)
}

func TestContentUseCaseSortingAndFiltering(t *testing.T) {
	t.Parallel()
	uc := New(persistent.NewContentRepo(nil))
	ctx := context.Background()

	now := time.Now().UTC()
	t1 := now.Add(-2 * time.Hour)
	t2 := now.Add(-1 * time.Hour)

	_, err := uc.CreateArticle(ctx, entity.ContentArticle{
		Title:       "Low Priority, Older",
		Slug:        "low-older",
		Status:      entity.ContentStatusPublished,
		Topic:       "news",
		Tags:        []string{"announcement"},
		Priority:    1,
		PublishedAt: &t1,
		CreatedAt:   t1,
	})
	require.NoError(t, err)

	_, err = uc.CreateArticle(ctx, entity.ContentArticle{
		Title:       "Low Priority, Newer",
		Slug:        "low-newer",
		Status:      entity.ContentStatusPublished,
		Topic:       "news",
		Tags:        []string{"announcement", "featured"},
		Priority:    1,
		PublishedAt: &t2,
		CreatedAt:   t2,
	})
	require.NoError(t, err)

	_, err = uc.CreateArticle(ctx, entity.ContentArticle{
		Title:       "High Priority",
		Slug:        "high",
		Status:      entity.ContentStatusPublished,
		Topic:       "promo",
		Tags:        []string{"featured"},
		Priority:    10,
		PublishedAt: &t1,
		CreatedAt:   t1,
	})
	require.NoError(t, err)

	// Test ListArticles PublicOnly
	items, total, err := uc.ListArticles(ctx, entity.ArticleFilter{PublicOnly: true})
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Equal(t, "high", items[0].Slug)
	require.Equal(t, "low-newer", items[1].Slug)
	require.Equal(t, "low-older", items[2].Slug)

	// Test ListArticles Filtering by Topic
	itemsTopic, totalTopic, err := uc.ListArticles(ctx, entity.ArticleFilter{Topic: "news"})
	require.NoError(t, err)
	require.Equal(t, 2, totalTopic)
	require.Equal(t, "low-newer", itemsTopic[0].Slug)

	// Test ListArticles Filtering by Tag
	itemsTag, totalTag, err := uc.ListArticles(ctx, entity.ArticleFilter{Tag: "featured"})
	require.NoError(t, err)
	require.Equal(t, 2, totalTag)
	require.Equal(t, "high", itemsTag[0].Slug)
	require.Equal(t, "low-newer", itemsTag[1].Slug)
}
