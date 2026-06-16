package webapi

import (
	"mime"
	"path/filepath"
)

func guessExtension(contentType string, fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return ext
	}
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".bin"
	}
}
