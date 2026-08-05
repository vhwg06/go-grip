package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ImporterUseCase defines the application service interface for batch import operations.
type ImporterUseCase interface {
	Import(ctx context.Context, items []ImportItem) (ImportResult, error)
}

type importerUseCase struct {
	repo     ImportRepo
	maxItems int
}

// NewImporterUseCase constructs a new ImporterUseCase application service.
func NewImporterUseCase(r ImportRepo, maxItems int) ImporterUseCase {
	if maxItems <= 0 {
		maxItems = MaxInitialImportItems
	}
	return &importerUseCase{repo: r, maxItems: maxItems}
}

// Import processes a slice of ImportItem objects and records success/failure counts.
func (uc *importerUseCase) Import(ctx context.Context, items []ImportItem) (ImportResult, error) {
	if len(items) > uc.maxItems {
		return ImportResult{}, ErrInvalidInput
	}
	result := ImportResult{Failed: []ImportFailure{}}
	for i, item := range items {
		if err := uc.importOne(ctx, item); err != nil {
			result.Failed = append(result.Failed, ImportFailure{Index: i, Reason: err.Error()})
			continue
		}
		result.Imported++
	}
	return result, nil
}

func (uc *importerUseCase) importOne(ctx context.Context, item ImportItem) error {
	now := time.Now().UTC()
	switch item.Type {
	case ImportItemProduct:
		draft := ImportProductDraft{
			ID:        uuid.New().String(),
			Title:     fmt.Sprint(item.Data["title"]),
			SKU:       fmt.Sprint(item.Data["sku"]),
			CreatedAt: now,
			UpdatedAt: now,
		}
		return uc.repo.StoreImportedProduct(ctx, draft)
	case ImportItemPost:
		draft := ImportPostDraft{
			ID:        uuid.New().String(),
			Title:     fmt.Sprint(item.Data["title"]),
			Slug:      fmt.Sprint(item.Data["slug"]),
			CreatedAt: now,
			UpdatedAt: now,
		}
		return uc.repo.StoreImportedPost(ctx, draft)
	default:
		return ErrInvalidInput
	}
}
