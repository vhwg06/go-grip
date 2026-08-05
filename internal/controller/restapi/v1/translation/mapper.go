package translation

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toTranslationResponse maps domain entity.Translation to openapi.TranslationResponse.
func toTranslationResponse(t entity.Translation) openapi.TranslationResponse {
	return openapi.TranslationResponse{
		Source:      t.Source,
		Target:      t.Destination,
		Text:        t.Original,
		Translation: t.Translation,
	}
}

// toTranslationHistoryResponse maps domain entity.TranslationHistory to openapi.TranslationHistoryResponse.
func toTranslationHistoryResponse(th entity.TranslationHistory) openapi.TranslationHistoryResponse {
	items := make([]openapi.TranslationResponse, len(th.History))
	for i, item := range th.History {
		items[i] = toTranslationResponse(item)
	}
	return openapi.TranslationHistoryResponse{
		History: items,
	}
}
