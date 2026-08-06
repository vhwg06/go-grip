package catalogbase

import "context"

// UnitOfWork executes an application operation with one transaction-bound context.
type UnitOfWork interface {
	Within(ctx context.Context, fn func(context.Context) error) error
}

// CatalogCategoryRepository persists catalog category entities.
type CatalogCategoryRepository interface {
	List(ctx context.Context) ([]CatalogCategory, error)
	GetByID(ctx context.Context, id string) (CatalogCategory, error)
	Store(ctx context.Context, category CatalogCategory) error
	Update(ctx context.Context, category CatalogCategory) error
	Delete(ctx context.Context, id string) error
}

// CatalogAttributeDefinitionRepository persists attribute definitions and their embedded enum values.
type CatalogAttributeDefinitionRepository interface {
	List(ctx context.Context) ([]CatalogAttributeDefinition, error)
	GetByID(ctx context.Context, id string) (CatalogAttributeDefinition, error)
	Store(ctx context.Context, definition CatalogAttributeDefinition) error
	Update(ctx context.Context, definition CatalogAttributeDefinition) error
	Delete(ctx context.Context, id string) error
}

// CatalogMasterRepository persists Material, Finish, and Pack masters.
type CatalogMasterRepository interface {
	List(ctx context.Context, kind string) ([]CatalogMaster, error)
	GetByID(ctx context.Context, kind, id string) (CatalogMaster, error)
	Store(ctx context.Context, master CatalogMaster) error
	Update(ctx context.Context, master CatalogMaster) error
	Delete(ctx context.Context, kind, id string) error
}

// CatalogProductModelRepository persists ProductModel root records.
type CatalogProductModelRepository interface {
	List(ctx context.Context) ([]CatalogProductModel, error)
	GetByID(ctx context.Context, id string) (CatalogProductModel, error)
	Store(ctx context.Context, model CatalogProductModel) error
	Update(ctx context.Context, model CatalogProductModel) error
	Delete(ctx context.Context, id string) error
}

// CatalogProductImageRepository persists ProductModel image entities.
type CatalogProductImageRepository interface {
	ListByModelID(ctx context.Context, modelID string) ([]CatalogProductImage, error)
	Store(ctx context.Context, modelID string, image CatalogProductImage) error
	Update(ctx context.Context, modelID string, image CatalogProductImage) error
	Delete(ctx context.Context, modelID, imageID string) error
	DeleteByModelID(ctx context.Context, modelID string) error
}

// CatalogVariantDimensionRepository persists ProductModel variant dimension entities.
type CatalogVariantDimensionRepository interface {
	ListByModelID(ctx context.Context, modelID string) ([]CatalogVariantDimension, error)
	GetByID(ctx context.Context, modelID, id string) (CatalogVariantDimension, error)
	Store(ctx context.Context, modelID string, dimension CatalogVariantDimension) error
	Update(ctx context.Context, modelID string, dimension CatalogVariantDimension) error
	Delete(ctx context.Context, modelID, id string) error
	DeleteByModelID(ctx context.Context, modelID string) error
}

// CatalogVariantRepository persists ProductModel variant entities.
type CatalogVariantRepository interface {
	List(ctx context.Context) ([]CatalogVariant, error)
	ListByModelID(ctx context.Context, modelID string) ([]CatalogVariant, error)
	GetByID(ctx context.Context, modelID, id string) (CatalogVariant, error)
	Store(ctx context.Context, modelID string, variant CatalogVariant) error
	Update(ctx context.Context, modelID string, variant CatalogVariant) error
	Delete(ctx context.Context, modelID, id string) error
	DeleteByModelID(ctx context.Context, modelID string) error
}

// CatalogRepositories groups the repository ports required by Catalog Base service.
type CatalogRepositories struct {
	Categories        CatalogCategoryRepository
	Definitions       CatalogAttributeDefinitionRepository
	Masters           CatalogMasterRepository
	ProductModels     CatalogProductModelRepository
	ProductImages     CatalogProductImageRepository
	VariantDimensions CatalogVariantDimensionRepository
	Variants          CatalogVariantRepository
}
