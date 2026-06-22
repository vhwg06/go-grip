package persistent

import (
	"context"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestContentRepo(t *testing.T) {
	t.Parallel()
	repo := NewContentRepo(nil)
	if repo.Postgres == nil || repo.Pool == nil {
		t.Skip("Skipping TestContentRepo because PostgreSQL connection is nil")
	}
	ctx := context.Background()

	// 1. Create articles with different topics, tags, priorities, and publish dates
	now := time.Now().UTC()
	t1 := now.Add(-2 * time.Hour)
	t2 := now.Add(-1 * time.Hour)

	art1 := entity.ContentArticle{
		ID:          "a1",
		Title:       "Low Priority, Older",
		Slug:        "low-older",
		Status:      entity.ContentStatusPublished,
		Topic:       "tech",
		Tags:        []string{"go", "backend"},
		Priority:    1,
		PublishedAt: &t1,
		CreatedAt:   t1,
	}

	art2 := entity.ContentArticle{
		ID:          "a2",
		Title:       "Low Priority, Newer",
		Slug:        "low-newer",
		Status:      entity.ContentStatusPublished,
		Topic:       "tech",
		Tags:        []string{"go", "frontend"},
		Priority:    1,
		PublishedAt: &t2,
		CreatedAt:   t2,
	}

	art3 := entity.ContentArticle{
		ID:          "a3",
		Title:       "High Priority",
		Slug:        "high-priority",
		Status:      entity.ContentStatusPublished,
		Topic:       "design",
		Tags:        []string{"ui", "ux"},
		Priority:    10,
		PublishedAt: &t1,
		CreatedAt:   t1,
	}

	artDraft := entity.ContentArticle{
		ID:        "a4",
		Title:     "Draft Article",
		Slug:      "draft-slug",
		Status:    entity.ContentStatusDraft,
		Topic:     "tech",
		Priority:  5,
		CreatedAt: now,
	}

	require.NoError(t, repo.StoreArticle(ctx, &art1))
	require.NoError(t, repo.StoreArticle(ctx, &art2))
	require.NoError(t, repo.StoreArticle(ctx, &art3))
	require.NoError(t, repo.StoreArticle(ctx, &artDraft))

	// 2. Test ListArticles for PublicOnly with priority and publish date sorting
	// Expected order: a3 (Priority 10), then a2 (Priority 1, newer published_at), then a1 (Priority 1, older published_at)
	// Draft should be excluded.
	items, total, err := repo.ListArticles(ctx, entity.ArticleFilter{PublicOnly: true})
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Len(t, items, 3)
	require.Equal(t, "a3", items[0].ID)
	require.Equal(t, "a2", items[1].ID)
	require.Equal(t, "a1", items[2].ID)

	// 3. Test ListArticles filtering by Topic
	itemsTopic, totalTopic, err := repo.ListArticles(ctx, entity.ArticleFilter{Topic: "tech"})
	require.NoError(t, err)
	// Should return: draft (a4, priority 5), low-newer (a2, priority 1), low-older (a1, priority 1) (since PublicOnly is false by default)
	require.Equal(t, 3, totalTopic)
	require.Equal(t, "a4", itemsTopic[0].ID)
	require.Equal(t, "a2", itemsTopic[1].ID)
	require.Equal(t, "a1", itemsTopic[2].ID)

	// 4. Test ListArticles filtering by Tag
	itemsTag, totalTag, err := repo.ListArticles(ctx, entity.ArticleFilter{Tag: "go"})
	require.NoError(t, err)
	// Should return: low-newer (a2), low-older (a1)
	require.Equal(t, 2, totalTag)
	require.Equal(t, "a2", itemsTag[0].ID)
	require.Equal(t, "a1", itemsTag[1].ID)

	// 5. Test DeleteArticle
	require.NoError(t, repo.DeleteArticle(ctx, "a3"))
	_, err = repo.GetArticle(ctx, "a3")
	require.ErrorIs(t, err, entity.ErrNotFound)

	// After deletion, list public articles should not contain a3
	itemsAfterDelete, totalAfterDelete, err := repo.ListArticles(ctx, entity.ArticleFilter{PublicOnly: true})
	require.NoError(t, err)
	require.Equal(t, 2, totalAfterDelete)
	require.Equal(t, "a2", itemsAfterDelete[0].ID)
	require.Equal(t, "a1", itemsAfterDelete[1].ID)
}
