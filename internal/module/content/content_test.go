package content

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type contentRepoStub struct {
	mu       sync.RWMutex
	articles map[string]ContentArticle
	pages    map[string]StaticPage
}

func newContentRepoStub() *contentRepoStub {
	return &contentRepoStub{
		articles: map[string]ContentArticle{},
		pages:    map[string]StaticPage{},
	}
}

func (r *contentRepoStub) StoreArticle(ctx context.Context, article *ContentArticle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.articles[article.ID] = *article
	return nil
}

func (r *contentRepoStub) UpdateArticle(ctx context.Context, article *ContentArticle) error {
	return r.StoreArticle(ctx, article)
}

func (r *contentRepoStub) ListArticles(ctx context.Context, filter ArticleFilter) ([]ContentArticle, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]ContentArticle, 0, len(r.articles))
	for _, article := range r.articles {
		if filter.PublicOnly && article.Status != ContentStatusPublished {
			continue
		}
		if filter.Topic != "" && article.Topic != filter.Topic {
			continue
		}
		if filter.Tag != "" {
			found := false
			for _, tag := range article.Tags {
				if tag == filter.Tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		items = append(items, article)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Priority != items[j].Priority {
			return items[i].Priority > items[j].Priority
		}
		leftTime := items[i].CreatedAt
		if items[i].PublishedAt != nil {
			leftTime = *items[i].PublishedAt
		}
		rightTime := items[j].CreatedAt
		if items[j].PublishedAt != nil {
			rightTime = *items[j].PublishedAt
		}
		return leftTime.After(rightTime)
	})

	total := len(items)
	page := filter.Pagination.Normalize()
	if page.Offset >= total {
		return []ContentArticle{}, total, nil
	}

	end := page.Offset + page.Limit
	if end > total {
		end = total
	}

	return append([]ContentArticle(nil), items[page.Offset:end]...), total, nil
}

func (r *contentRepoStub) GetArticle(ctx context.Context, idOrSlug string) (ContentArticle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, article := range r.articles {
		if article.ID == idOrSlug || article.Slug == idOrSlug {
			return article, nil
		}
	}
	return ContentArticle{}, ErrNotFound
}

func (r *contentRepoStub) DeleteArticle(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.articles[id]; !ok {
		return ErrNotFound
	}
	delete(r.articles, id)
	return nil
}

func (r *contentRepoStub) StorePage(ctx context.Context, page *StaticPage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pages[page.Slug] = *page
	return nil
}

func (r *contentRepoStub) UpdatePage(ctx context.Context, page *StaticPage) error {
	return r.StorePage(ctx, page)
}

func (r *contentRepoStub) GetPageBySlug(ctx context.Context, slug string) (StaticPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	page, ok := r.pages[slug]
	if !ok {
		return StaticPage{}, ErrNotFound
	}
	return page, nil
}

func (r *contentRepoStub) PublishDue(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()

	published := 0
	for id, article := range r.articles {
		if article.Status == ContentStatusScheduled && article.ScheduledAt != nil && !article.ScheduledAt.After(now) {
			article.Status = ContentStatusPublished
			article.PublishedAt = &now
			r.articles[id] = article
			published++
		}
	}
	return published, nil
}

type homepageRepoStub struct {
	mu    sync.RWMutex
	items map[string]HomepageBlock
}

func newHomepageRepoStub() *homepageRepoStub {
	return &homepageRepoStub{items: map[string]HomepageBlock{}}
}

func (r *homepageRepoStub) Store(ctx context.Context, block *HomepageBlock) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[block.ID] = *block
	return nil
}

func (r *homepageRepoStub) List(ctx context.Context, activeOnly bool) ([]HomepageBlock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]HomepageBlock, 0, len(r.items))
	for _, item := range r.items {
		if activeOnly && !item.IsActive {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *homepageRepoStub) Update(ctx context.Context, block *HomepageBlock) error {
	return r.Store(ctx, block)
}

func (r *homepageRepoStub) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

type supportRepoStub struct {
	mu    sync.RWMutex
	items map[string]SupportChannel
}

func newSupportRepoStub() *supportRepoStub {
	return &supportRepoStub{items: map[string]SupportChannel{}}
}

func (r *supportRepoStub) List(ctx context.Context, enabledOnly bool) ([]SupportChannel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]SupportChannel, 0, len(r.items))
	for _, item := range r.items {
		if enabledOnly && !item.IsEnabled {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *supportRepoStub) Update(ctx context.Context, channel *SupportChannel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[channel.ID] = *channel
	return nil
}

func TestContentUseCaseArticle(t *testing.T) {
	t.Parallel()
	uc := NewContentUseCase(newContentRepoStub())
	ctx := context.Background()

	// Test CreateArticle
	article, err := uc.CreateArticle(ctx, ContentArticle{Title: "Post", Slug: "post"})
	require.NoError(t, err)
	require.NotEmpty(t, article.ID)
	require.Equal(t, ContentStatusDraft, article.Status)

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
	require.ErrorIs(t, err, ErrNotFound)
}

func TestContentUseCaseSortingAndFiltering(t *testing.T) {
	t.Parallel()
	uc := NewContentUseCase(newContentRepoStub())
	ctx := context.Background()

	now := time.Now().UTC()
	t1 := now.Add(-2 * time.Hour)
	t2 := now.Add(-1 * time.Hour)

	_, err := uc.CreateArticle(ctx, ContentArticle{
		Title:       "Low Priority, Older",
		Slug:        "low-older",
		Status:      ContentStatusPublished,
		Topic:       "news",
		Tags:        []string{"announcement"},
		Priority:    1,
		PublishedAt: &t1,
		CreatedAt:   t1,
	})
	require.NoError(t, err)

	_, err = uc.CreateArticle(ctx, ContentArticle{
		Title:       "Low Priority, Newer",
		Slug:        "low-newer",
		Status:      ContentStatusPublished,
		Topic:       "news",
		Tags:        []string{"announcement", "featured"},
		Priority:    1,
		PublishedAt: &t2,
		CreatedAt:   t2,
	})
	require.NoError(t, err)

	_, err = uc.CreateArticle(ctx, ContentArticle{
		Title:       "High Priority",
		Slug:        "high",
		Status:      ContentStatusPublished,
		Topic:       "promo",
		Tags:        []string{"featured"},
		Priority:    10,
		PublishedAt: &t1,
		CreatedAt:   t1,
	})
	require.NoError(t, err)

	// Test ListArticles PublicOnly
	items, total, err := uc.ListArticles(ctx, ArticleFilter{PublicOnly: true})
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Equal(t, "high", items[0].Slug)
	require.Equal(t, "low-newer", items[1].Slug)
	require.Equal(t, "low-older", items[2].Slug)

	// Test ListArticles Filtering by Topic
	itemsTopic, totalTopic, err := uc.ListArticles(ctx, ArticleFilter{Topic: "news"})
	require.NoError(t, err)
	require.Equal(t, 2, totalTopic)
	require.Equal(t, "low-newer", itemsTopic[0].Slug)

	// Test ListArticles Filtering by Tag
	itemsTag, totalTag, err := uc.ListArticles(ctx, ArticleFilter{Tag: "featured"})
	require.NoError(t, err)
	require.Equal(t, 2, totalTag)
	require.Equal(t, "high", itemsTag[0].Slug)
	require.Equal(t, "low-newer", itemsTag[1].Slug)
}

func TestHomepageUseCase(t *testing.T) {
	t.Parallel()
	uc := NewHomepageUseCase(newHomepageRepoStub(), newSupportRepoStub())
	block, err := uc.StoreBlock(context.Background(), HomepageBlock{BlockType: "banner", IsActive: true})
	require.NoError(t, err)
	require.NotEmpty(t, block.ID)
	items, err := uc.ListBlocks(context.Background(), true)
	require.NoError(t, err)
	require.Len(t, items, 1)
}

func TestContentSchedulerCatchUp(t *testing.T) {
	t.Parallel()
	repo := newContentRepoStub()
	uc := NewContentUseCase(repo)
	past := time.Now().Add(-time.Hour)
	_, err := uc.CreateArticle(context.Background(), ContentArticle{Title: "Post", Slug: "post", Status: ContentStatusScheduled, ScheduledAt: &past})
	require.NoError(t, err)
	count, err := uc.CatchUpScheduled(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
