package translation

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Translation capability.
type Handler struct {
	translationUC usecase.Translation
	logger        logger.Interface
}

// NewHandler constructs a new Translation vertical handler instance.
func NewHandler(translationUC usecase.Translation, l logger.Interface) *Handler {
	return &Handler{
		translationUC: translationUC,
		logger:        l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// TranslateText handles POST /translation/do
func (h *Handler) TranslateText(ctx context.Context, request openapi.TranslateTextRequestObject) (openapi.TranslateTextResponseObject, error) {
	if request.Body == nil {
		return openapi.TranslateText400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	inputTranslation := entity.Translation{
		Source:      request.Body.Source,
		Destination: request.Body.Target,
		Original:    request.Body.Text,
	}

	res, err := h.translationUC.Translate(ctx, actor.UserID, inputTranslation)
	if err != nil {
		status, errResp := mapTranslationError(err)
		switch status {
		case 400:
			return openapi.TranslateText400JSONResponse{}, nil
		case 401:
			return openapi.TranslateText401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.TranslateText500JSONResponse{}, nil
		}
	}

	resDTO := toTranslationResponse(res)
	return openapi.TranslateText200JSONResponse(resDTO), nil
}

// GetTranslationHistory handles GET /translation/history
func (h *Handler) GetTranslationHistory(ctx context.Context, request openapi.GetTranslationHistoryRequestObject) (openapi.GetTranslationHistoryResponseObject, error) {
	actor := getActor(ctx)
	history, err := h.translationUC.History(ctx, actor.UserID)
	if err != nil {
		status, errResp := mapTranslationError(err)
		switch status {
		case 401:
			return openapi.GetTranslationHistory401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetTranslationHistory500JSONResponse{}, nil
		}
	}

	historyDTO := toTranslationHistoryResponse(history)
	return openapi.GetTranslationHistory200JSONResponse(historyDTO), nil
}
