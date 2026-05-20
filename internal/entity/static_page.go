package entity

import "time"

type StaticPage struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Body        string        `json:"body"`
	TemplateKey string        `json:"template_key"`
	Status      ContentStatus `json:"status"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
