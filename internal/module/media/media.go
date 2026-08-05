package media

import "time"

const (
	// MaxMediaUploadBytes defines the default 5MB maximum file upload size.
	MaxMediaUploadBytes = 5 * 1024 * 1024
)

// AllowedMediaMimeTypes lists the supported image MIME types.
var AllowedMediaMimeTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

// IsAllowedMediaMimeType checks whether a media upload MIME type is supported.
func IsAllowedMediaMimeType(mimeType string) bool {
	_, ok := AllowedMediaMimeTypes[mimeType]
	return ok
}

// MediaAsset represents a file/image asset stored in media storage and metadata database.
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
