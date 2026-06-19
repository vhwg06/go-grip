package persistent

import (
	"context"
	"slices"
	"sort"
	"strings"
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
	cloned := cloneContentArticle(*article)
	r.articles[cloned.ID] = cloned
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
		items = append(items, cloneContentArticle(article))
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
			return cloneContentArticle(article), nil
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
	cloned := cloneStaticPage(*page)
	r.pages[cloned.Slug] = cloned
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
	return cloneStaticPage(page), nil
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

func cloneContentArticle(article entity.ContentArticle) entity.ContentArticle {
	cloned := article
	cloned.ID = strings.Clone(article.ID)
	cloned.Title = strings.Clone(article.Title)
	cloned.Slug = strings.Clone(article.Slug)
	cloned.Body = strings.Clone(article.Body)
	cloned.AuthorID = strings.Clone(article.AuthorID)
	cloned.ImageURL = strings.Clone(article.ImageURL)
	cloned.Topic = strings.Clone(article.Topic)
	if article.ScheduledAt != nil {
		scheduledAt := *article.ScheduledAt
		cloned.ScheduledAt = &scheduledAt
	}
	if article.PublishedAt != nil {
		publishedAt := *article.PublishedAt
		cloned.PublishedAt = &publishedAt
	}
	if len(article.Tags) > 0 {
		cloned.Tags = make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			cloned.Tags[i] = strings.Clone(tag)
		}
	}
	return cloned
}

func cloneStaticPage(page entity.StaticPage) entity.StaticPage {
	cloned := page
	cloned.ID = strings.Clone(page.ID)
	cloned.Title = strings.Clone(page.Title)
	cloned.Slug = strings.Clone(page.Slug)
	cloned.Body = strings.Clone(page.Body)
	cloned.TemplateKey = strings.Clone(page.TemplateKey)
	if len(page.Gallery) > 0 {
		cloned.Gallery = make([]string, len(page.Gallery))
		for i, item := range page.Gallery {
			cloned.Gallery[i] = strings.Clone(item)
		}
	}
	return cloned
}
