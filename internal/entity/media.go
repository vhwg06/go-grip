package entity

import "time"

type MediaAsset struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text,omitempty"`
	OwnerType string    `json:"owner_type,omitempty"`
	OwnerID   string    `json:"owner_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
