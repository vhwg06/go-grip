package importer

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Importer capability.
type Handler struct {
	importerUC usecase.Importer
	logger     logger.Interface
}

// NewHandler constructs a new Importer vertical handler instance.
func NewHandler(importerUC usecase.Importer, l logger.Interface) *Handler {
	return &Handler{
		importerUC: importerUC,
		logger:     l,
	}
}

// ExecuteImport handles POST /importer/execute
func (h *Handler) ExecuteImport(ctx context.Context, request openapi.ExecuteImportRequestObject) (openapi.ExecuteImportResponseObject, error) {
	if request.Body == nil {
		return openapi.ExecuteImport400JSONResponse{}, nil
	}

	items := make([]entity.ImportItem, len(request.Body.Items))
	for i, item := range request.Body.Items {
		var dataMap map[string]any
		if item.Data != nil {
			dataMap = *item.Data
		}
		items[i] = entity.ImportItem{
			Type: entity.ImportItemType(item.Type),
			Data: dataMap,
		}
	}

	res, err := h.importerUC.Import(ctx, items)
	if err != nil {
		status, errResp := mapImporterError(err)
		switch status {
		case 400:
			return openapi.ExecuteImport400JSONResponse{}, nil
		case 401:
			return openapi.ExecuteImport401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.ExecuteImport403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ExecuteImport500JSONResponse{}, nil
		}
	}

	importDTO := toImportExecuteResponse(res)
	return openapi.ExecuteImport200JSONResponse(importDTO), nil
}
