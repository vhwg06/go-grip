package content

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Content & Homepage capabilities.
type Handler struct {
	contentUC  contentmodule.ContentUseCase
	homepageUC contentmodule.HomepageUseCase
	logger     logger.Interface
}

// NewHandler constructs a new Content/Homepage vertical handler instance.
func NewHandler(contentUC contentmodule.ContentUseCase, homepageUC contentmodule.HomepageUseCase, l logger.Interface) *Handler {
	return &Handler{
		contentUC:  contentUC,
		homepageUC: homepageUC,
		logger:     l,
	}
}

// GetStaticPage handles GET /content/pages/{slug}
func (h *Handler) GetStaticPage(ctx context.Context, request openapi.GetStaticPageRequestObject) (openapi.GetStaticPageResponseObject, error) {
	page, err := h.contentUC.GetPage(ctx, request.Slug)
	if err != nil {
		status, errResp := mapContentError(err)
		switch status {
		case 404:
			return openapi.GetStaticPage404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetStaticPage500JSONResponse{}, nil
		}
	}

	pageDTO := toStaticPageResponse(page)
	return openapi.GetStaticPage200JSONResponse(pageDTO), nil
}

// GetHomepageConfig handles GET /homepage/config
func (h *Handler) GetHomepageConfig(ctx context.Context, request openapi.GetHomepageConfigRequestObject) (openapi.GetHomepageConfigResponseObject, error) {
	blocks, err := h.homepageUC.ListBlocks(ctx, true)
	if err != nil {
		return openapi.GetHomepageConfig500JSONResponse{}, nil
	}

	configDTO := toHomepageConfigResponse(blocks)
	return openapi.GetHomepageConfig200JSONResponse(configDTO), nil
}
