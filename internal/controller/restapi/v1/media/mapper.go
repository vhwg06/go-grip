package media

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// toMediaUploadResponse maps media file upload result to openapi.MediaUploadResponse DTO.
func toMediaUploadResponse(id, url, filename string) openapi.MediaUploadResponse {
	fn := filename
	return openapi.MediaUploadResponse{
		Id:       id,
		Url:      url,
		Filename: &fn,
	}
}
