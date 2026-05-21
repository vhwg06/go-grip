package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo     repo.ImportRepo
	maxItems int
}

func New(r repo.ImportRepo, maxItems int) *UseCase {
	if maxItems <= 0 {
		maxItems = entity.MaxInitialImportItems
	}
	return &UseCase{repo: r, maxItems: maxItems}
}

func (uc *UseCase) Import(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error) {
	if len(items) > uc.maxItems {
		return entity.ImportResult{}, entity.ErrInvalidInput
	}
	result := entity.ImportResult{Failed: []entity.ImportFailure{}}
	for i, item := range items {
		if err := uc.importOne(ctx, item); err != nil {
			result.Failed = append(result.Failed, entity.ImportFailure{Index: i, Reason: err.Error()})
			continue
		}
		result.Imported++
	}
	return result, nil
}

func (uc *UseCase) importOne(ctx context.Context, item entity.ImportItem) error {
	now := time.Now().UTC()
	switch item.Type {
	case entity.ImportItemProduct:
		product := entity.Product{ID: uuid.New().String(), Title: fmt.Sprint(item.Data["title"]), SKU: fmt.Sprint(item.Data["sku"]), Status: entity.ProductStatusDraft, CreatedAt: now, UpdatedAt: now}
		return uc.repo.StoreImportedProduct(ctx, &product)
	case entity.ImportItemPost:
		article := entity.ContentArticle{ID: uuid.New().String(), Title: fmt.Sprint(item.Data["title"]), Slug: fmt.Sprint(item.Data["slug"]), Status: entity.ContentStatusDraft, CreatedAt: now, UpdatedAt: now}
		return uc.repo.StoreImportedPost(ctx, &article)
	default:
		return entity.ErrInvalidInput
	}
}
