package entity

import "time"

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
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
