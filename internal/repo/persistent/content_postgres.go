package persistent

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type ContentRepo struct {
	*postgres.Postgres
	mu       sync.RWMutex
	articles map[string]entity.ContentArticle
	pages    map[string]entity.StaticPage
}

func NewContentRepo(pg *postgres.Postgres) *ContentRepo {
	return &ContentRepo{Postgres: pg, articles: map[string]entity.ContentArticle{}, pages: map[string]entity.StaticPage{}}
}

func (r *ContentRepo) StoreArticle(ctx context.Context, article *entity.ContentArticle) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.articles[article.ID] = *article
	return nil
}

func (r *ContentRepo) UpdateArticle(ctx context.Context, article *entity.ContentArticle) error {
	return r.StoreArticle(ctx, article)
}

func (r *ContentRepo) ListArticles(ctx context.Context, filter entity.ArticleFilter) ([]entity.ContentArticle, int, error) {
	_ = ctx
	page := filter.Pagination.Normalize()
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]entity.ContentArticle, 0, len(r.articles))
	for _, article := range r.articles {
		if filter.PublicOnly && article.Status != entity.ContentStatusPublished {
			continue
		}
		if filter.Topic != "" && article.Topic != filter.Topic {
			continue
		}
		if filter.Tag != "" && !slices.Contains(article.Tags, filter.Tag) {
			continue
		}
		items = append(items, article)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Priority != items[j].Priority {
			return items[i].Priority > items[j].Priority
		}
		timeI := items[i].CreatedAt
		if items[i].PublishedAt != nil {
			timeI = *items[i].PublishedAt
		}
		timeJ := items[j].CreatedAt
		if items[j].PublishedAt != nil {
			timeJ = *items[j].PublishedAt
		}
		return timeI.After(timeJ)
	})

	total := len(items)
	if page.Offset > total {
		return []entity.ContentArticle{}, total, nil
	}
	end := min(page.Offset+page.Limit, total)
	return items[page.Offset:end], total, nil
}

func (r *ContentRepo) GetArticle(ctx context.Context, idOrSlug string) (entity.ContentArticle, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, article := range r.articles {
		if article.ID == idOrSlug || article.Slug == idOrSlug {
			return article, nil
		}
	}
	return entity.ContentArticle{}, entity.ErrNotFound
}

func (r *ContentRepo) DeleteArticle(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.articles[id]; !ok {
		return entity.ErrNotFound
	}
	delete(r.articles, id)
	return nil
}

func (r *ContentRepo) StorePage(ctx context.Context, page *entity.StaticPage) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pages[page.Slug] = *page
	return nil
}

func (r *ContentRepo) UpdatePage(ctx context.Context, page *entity.StaticPage) error {
	return r.StorePage(ctx, page)
}

func (r *ContentRepo) GetPageBySlug(ctx context.Context, slug string) (entity.StaticPage, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	page, ok := r.pages[slug]
	if !ok {
		return entity.StaticPage{}, entity.ErrNotFound
	}
	return page, nil
}

func (r *ContentRepo) PublishDue(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	count := 0
	for id, article := range r.articles {
		if article.Status == entity.ContentStatusScheduled && article.ScheduledAt != nil && !article.ScheduledAt.After(now) {
			article.Status = entity.ContentStatusPublished
			article.PublishedAt = &now
			r.articles[id] = article
			count++
		}
	}
	return count, nil
}
