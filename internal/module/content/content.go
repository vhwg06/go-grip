package content

import (
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
