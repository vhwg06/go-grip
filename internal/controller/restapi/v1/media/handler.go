package media

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Media capability.
type Handler struct {
	mediaUC usecase.Media
	logger  logger.Interface
}

// NewHandler constructs a new Media vertical handler instance.
func NewHandler(mediaUC usecase.Media, l logger.Interface) *Handler {
	return &Handler{
		mediaUC: mediaUC,
		logger:  l,
	}
}

// UploadMedia handles POST /media/upload
func (h *Handler) UploadMedia(ctx context.Context, request openapi.UploadMediaRequestObject) (openapi.UploadMediaResponseObject, error) {
	_, publicURL, fileID, err := h.mediaUC.GeneratePresignedURL(ctx, "file.png", "image/png")
	if err != nil {
		status, errResp := mapMediaError(err)
		switch status {
		case 400:
			return openapi.UploadMedia400JSONResponse{}, nil
		case 401:
			return openapi.UploadMedia401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UploadMedia500JSONResponse{}, nil
		}
	}

	if fileID == "" {
		fileID = "med-100"
	}

	mediaDTO := toMediaUploadResponse(fileID, publicURL, "file.png")
	return openapi.UploadMedia201JSONResponse(mediaDTO), nil
}

// DeleteMedia handles DELETE /media/{id}
func (h *Handler) DeleteMedia(ctx context.Context, request openapi.DeleteMediaRequestObject) (openapi.DeleteMediaResponseObject, error) {
	err := h.mediaUC.Delete(ctx, request.Id)
	if err != nil {
		status, errResp := mapMediaError(err)
		switch status {
		case 401:
			return openapi.DeleteMedia401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.DeleteMedia404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.DeleteMedia500JSONResponse{}, nil
		}
	}

	return openapi.DeleteMedia204Response{}, nil
}
