# Content Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Content` business module.

---

## 1. Owned Symbols
- **Entities**: `ContentArticle`, `ContentStatus`, `ArticleFilter`, `HomepageConfig`, `HomepageBlock`, `StaticPage`, `SupportChannel`, `SEOMetadataVO`
- **Errors**: `ErrNotFound`, `ErrInvalidInput`, `ErrUnauthorized`
- **Use Cases**: `ContentUseCase`, `HomepageUseCase`, `PublishSchedulerUseCase`
- **Repository Ports**: `ContentRepo`, `HomepageRepo`, `SupportRepo`

---

## 2. Ports & Interfaces
```go
package content

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusScheduled ContentStatus = "scheduled"
	ContentStatusPublished ContentStatus = "published"
)

type ContentArticle struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Body        string        `json:"body"`
	Status      ContentStatus `json:"status"`
	ScheduledAt *time.Time    `json:"scheduled_at,omitempty"`
	PublishedAt *time.Time    `json:"published_at,omitempty"`
	AuthorID    string        `json:"author_id,omitempty"`
	ImageURL    string        `json:"image_url,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Topic       string        `json:"topic,omitempty"`
	Priority    int           `json:"priority"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ArticleFilter struct {
	PublicOnly bool
	Topic      string
	Tag        string
	Pagination pagination.Pagination
}

type ContentRepo interface {
	Store(ctx context.Context, article ContentArticle) (ContentArticle, error)
	GetByID(ctx context.Context, id string) (ContentArticle, error)
	GetBySlug(ctx context.Context, slug string) (ContentArticle, error)
	List(ctx context.Context, filter ArticleFilter) ([]ContentArticle, int, error)
	Update(ctx context.Context, article ContentArticle) (ContentArticle, error)
	PublishDueScheduledArticles(ctx context.Context, now time.Time) (int, error)
}
```

---

## 3. Infrastructure & Delivery Consumers
- `internal/repo/persistent/content_postgres.go`
- `internal/repo/persistent/homepage_postgres.go`
- `internal/repo/persistent/support_postgres.go`
- `internal/controller/restapi/v1/content/`
- `internal/app/app.go`
