package entity

const (
	MaxMediaUploadBytes   = 5 * 1024 * 1024
	MaxInitialImportItems = 25
)

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
