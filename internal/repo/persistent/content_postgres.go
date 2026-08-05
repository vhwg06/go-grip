package persistent

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ContentRepo struct {
	*postgres.Postgres
	mu       sync.RWMutex
	articles map[string]contentmodule.ContentArticle
	pages    map[string]contentmodule.StaticPage
}

func NewContentRepo(pg *postgres.Postgres) *ContentRepo {
	return &ContentRepo{
		Postgres: pg,
		articles: map[string]contentmodule.ContentArticle{},
		pages:    map[string]contentmodule.StaticPage{},
	}
}

var _ contentmodule.ContentRepo = (*ContentRepo)(nil)

func (r *ContentRepo) useMemory() bool {
	return r == nil || r.Postgres == nil || r.Pool == nil
}

func uuidOrNil(id string) *string {
	if id == "" {
		return nil
	}
	if len(id) != 36 {
		return nil
	}
	return &id
}

func (r *ContentRepo) StoreArticle(ctx context.Context, article *contentmodule.ContentArticle) error {
	if r.useMemory() {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		r.articles[article.ID] = *article
		return nil
	}

	sql, args, err := r.Builder.
		Insert("content_articles").
		Columns("id, title, slug, body, status, scheduled_at, published_at, author_id, image_url, tags, topic, priority, created_at, updated_at").
		Values(
			article.ID, article.Title, article.Slug, article.Body, string(article.Status),
			article.ScheduledAt, article.PublishedAt, uuidOrNil(article.AuthorID),
			article.ImageURL, article.Tags, article.Topic, article.Priority,
			article.CreatedAt, article.UpdatedAt,
		).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
			title = EXCLUDED.title, 
			slug = EXCLUDED.slug, 
			body = EXCLUDED.body, 
			status = EXCLUDED.status, 
			scheduled_at = EXCLUDED.scheduled_at, 
			published_at = EXCLUDED.published_at, 
			author_id = EXCLUDED.author_id, 
			image_url = EXCLUDED.image_url, 
			tags = EXCLUDED.tags, 
			topic = EXCLUDED.topic, 
			priority = EXCLUDED.priority, 
			updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return fmt.Errorf("ContentRepo - StoreArticle - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContentRepo - StoreArticle - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *ContentRepo) UpdateArticle(ctx context.Context, article *contentmodule.ContentArticle) error {
	return r.StoreArticle(ctx, article)
}

func (r *ContentRepo) ListArticles(ctx context.Context, filter contentmodule.ArticleFilter) ([]contentmodule.ContentArticle, int, error) {
	if r.useMemory() {
		_ = ctx
		r.mu.RLock()
		defer r.mu.RUnlock()

		items := make([]contentmodule.ContentArticle, 0, len(r.articles))
		for _, article := range r.articles {
			if filter.PublicOnly && article.Status != contentmodule.ContentStatusPublished {
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
			return []contentmodule.ContentArticle{}, total, nil
		}

		end := page.Offset + page.Limit
		if end > total {
			end = total
		}

		return append([]contentmodule.ContentArticle(nil), items[page.Offset:end]...), total, nil
	}

	page := filter.Pagination.Normalize()

	countBuilder := r.Builder.Select("COUNT(*)").From("content_articles")
	if filter.PublicOnly {
		countBuilder = countBuilder.Where(sq.Eq{"status": string(contentmodule.ContentStatusPublished)})
	}
	if filter.Topic != "" {
		countBuilder = countBuilder.Where(sq.Eq{"topic": filter.Topic})
	}
	if filter.Tag != "" {
		countBuilder = countBuilder.Where("? = ANY(tags)", filter.Tag)
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("ContentRepo - ListArticles - count builder: %w", err)
	}

	var total int
	if err = r.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ContentRepo - ListArticles - count query: %w", err)
	}

	dataBuilder := r.Builder.
		Select("id, title, slug, body, status, scheduled_at, published_at, author_id, image_url, tags, topic, priority, created_at, updated_at").
		From("content_articles")

	if filter.PublicOnly {
		dataBuilder = dataBuilder.Where(sq.Eq{"status": string(contentmodule.ContentStatusPublished)})
	}
	if filter.Topic != "" {
		dataBuilder = dataBuilder.Where(sq.Eq{"topic": filter.Topic})
	}
	if filter.Tag != "" {
		dataBuilder = dataBuilder.Where("? = ANY(tags)", filter.Tag)
	}

	dataBuilder = dataBuilder.OrderBy("priority DESC, COALESCE(published_at, created_at) DESC").
		Limit(uint64(page.Limit)).
		Offset(uint64(page.Offset))

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("ContentRepo - ListArticles - data builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("ContentRepo - ListArticles - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	articles := make([]contentmodule.ContentArticle, 0, page.Limit)
	for rows.Next() {
		var art contentmodule.ContentArticle
		var authorPtr *string
		var statusStr string
		var imagePtr *string
		var topicPtr *string
		var tags []string

		err = rows.Scan(
			&art.ID, &art.Title, &art.Slug, &art.Body, &statusStr,
			&art.ScheduledAt, &art.PublishedAt, &authorPtr, &imagePtr,
			&tags, &topicPtr, &art.Priority, &art.CreatedAt, &art.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ContentRepo - ListArticles - rows.Scan: %w", err)
		}

		art.Status = contentmodule.ContentStatus(statusStr)
		if authorPtr != nil {
			art.AuthorID = *authorPtr
		}
		if imagePtr != nil {
			art.ImageURL = *imagePtr
		}
		if topicPtr != nil {
			art.Topic = *topicPtr
		}
		art.Tags = tags

		articles = append(articles, art)
	}

	return articles, total, nil
}

func (r *ContentRepo) GetArticle(ctx context.Context, idOrSlug string) (contentmodule.ContentArticle, error) {
	if r.useMemory() {
		_ = ctx
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, article := range r.articles {
			if article.ID == idOrSlug || article.Slug == idOrSlug {
				return article, nil
			}
		}
		return contentmodule.ContentArticle{}, contentmodule.ErrNotFound
	}

	sql, args, err := r.Builder.
		Select("id, title, slug, body, status, scheduled_at, published_at, author_id, image_url, tags, topic, priority, created_at, updated_at").
		From("content_articles").
		Where("id = ? OR slug = ?", idOrSlug, idOrSlug).
		ToSql()
	if err != nil {
		return contentmodule.ContentArticle{}, fmt.Errorf("ContentRepo - GetArticle - r.Builder: %w", err)
	}

	var art contentmodule.ContentArticle
	var authorPtr *string
	var statusStr string
	var imagePtr *string
	var topicPtr *string
	var tags []string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&art.ID, &art.Title, &art.Slug, &art.Body, &statusStr,
		&art.ScheduledAt, &art.PublishedAt, &authorPtr, &imagePtr,
		&tags, &topicPtr, &art.Priority, &art.CreatedAt, &art.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return contentmodule.ContentArticle{}, contentmodule.ErrNotFound
		}
		return contentmodule.ContentArticle{}, fmt.Errorf("ContentRepo - GetArticle - r.Pool.QueryRow: %w", err)
	}

	art.Status = contentmodule.ContentStatus(statusStr)
	if authorPtr != nil {
		art.AuthorID = *authorPtr
	}
	if imagePtr != nil {
		art.ImageURL = *imagePtr
	}
	if topicPtr != nil {
		art.Topic = *topicPtr
	}
	art.Tags = tags

	return art, nil
}

func (r *ContentRepo) DeleteArticle(ctx context.Context, id string) error {
	if r.useMemory() {
		_ = ctx
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.articles[id]; !ok {
			return contentmodule.ErrNotFound
		}
		delete(r.articles, id)
		return nil
	}

	sql, args, err := r.Builder.
		Delete("content_articles").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("ContentRepo - DeleteArticle - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContentRepo - DeleteArticle - r.Pool.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return contentmodule.ErrNotFound
	}
	return nil
}

func (r *ContentRepo) StorePage(ctx context.Context, page *contentmodule.StaticPage) error {
	if r.useMemory() {
		_ = ctx
		if page.ID == "" {
			page.ID = uuid.New().String()
		}
		r.mu.Lock()
		defer r.mu.Unlock()
		r.pages[page.Slug] = *page
		return nil
	}

	if page.ID == "" {
		page.ID = uuid.New().String()
	}

	sql, args, err := r.Builder.
		Insert("static_pages").
		Columns("id, title, slug, body, template_key, status, gallery, updated_at").
		Values(page.ID, page.Title, page.Slug, page.Body, page.TemplateKey, string(page.Status), page.Gallery, page.UpdatedAt).
		Suffix(`ON CONFLICT (slug) DO UPDATE SET 
			title = EXCLUDED.title, 
			body = EXCLUDED.body, 
			template_key = EXCLUDED.template_key, 
			status = EXCLUDED.status, 
			gallery = EXCLUDED.gallery, 
			updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return fmt.Errorf("ContentRepo - StorePage - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContentRepo - StorePage - r.Pool.Exec: %w", err)
	}
	return nil
}

func (r *ContentRepo) UpdatePage(ctx context.Context, page *contentmodule.StaticPage) error {
	existing, err := r.GetPageBySlug(ctx, page.Slug)
	if err == nil {
		page.ID = existing.ID
	}
	return r.StorePage(ctx, page)
}

func (r *ContentRepo) GetPageBySlug(ctx context.Context, slug string) (contentmodule.StaticPage, error) {
	if r.useMemory() {
		_ = ctx
		r.mu.RLock()
		defer r.mu.RUnlock()
		page, ok := r.pages[slug]
		if !ok {
			return contentmodule.StaticPage{}, contentmodule.ErrNotFound
		}
		return page, nil
	}

	sql, args, err := r.Builder.
		Select("id, title, slug, body, template_key, status, gallery, updated_at").
		From("static_pages").
		Where("slug = ?", slug).
		ToSql()
	if err != nil {
		return contentmodule.StaticPage{}, fmt.Errorf("ContentRepo - GetPageBySlug - r.Builder: %w", err)
	}

	var pg contentmodule.StaticPage
	var statusStr string
	var gallery []string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&pg.ID, &pg.Title, &pg.Slug, &pg.Body, &pg.TemplateKey, &statusStr, &gallery, &pg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return contentmodule.StaticPage{}, contentmodule.ErrNotFound
		}
		return contentmodule.StaticPage{}, fmt.Errorf("ContentRepo - GetPageBySlug - r.Pool.QueryRow: %w", err)
	}

	pg.Status = contentmodule.ContentStatus(statusStr)
	pg.Gallery = gallery

	return pg, nil
}

func (r *ContentRepo) PublishDue(ctx context.Context) (int, error) {
	if r.useMemory() {
		_ = ctx
		now := time.Now().UTC()
		r.mu.Lock()
		defer r.mu.Unlock()

		published := 0
		for id, article := range r.articles {
			if article.Status == contentmodule.ContentStatusScheduled && article.ScheduledAt != nil && !article.ScheduledAt.After(now) {
				article.Status = contentmodule.ContentStatusPublished
				article.PublishedAt = &now
				r.articles[id] = article
				published++
			}
		}
		return published, nil
	}

	now := time.Now().UTC()

	sql, args, err := r.Builder.
		Update("content_articles").
		Set("status", string(contentmodule.ContentStatusPublished)).
		Set("published_at", now).
		Where("status = ? AND scheduled_at <= ?", string(contentmodule.ContentStatusScheduled), now).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ContentRepo - PublishDue - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("ContentRepo - PublishDue - r.Pool.Exec: %w", err)
	}

	return int(result.RowsAffected()), nil
}
