package importer

import "context"

// ImportRepo defines the persistence port owned by the Importer module.
type ImportRepo interface {
	StoreImportedProduct(ctx context.Context, draft ImportProductDraft) error
	StoreImportedPost(ctx context.Context, draft ImportPostDraft) error
}
