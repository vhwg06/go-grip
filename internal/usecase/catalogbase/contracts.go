package catalogbase

import "context"

// UseCase is the application port consumed by the REST adapter.  The adapter
// does not know the concrete service or its repository implementation.
type UseCase interface {
	ListCategories(context.Context) ([]map[string]any, error)
	CreateCategory(context.Context, map[string]any) (map[string]any, error)
	UpdateCategory(context.Context, string, map[string]any) (map[string]any, error)
	DeactivateCategory(context.Context, string) (map[string]any, error)
	DeleteCategory(context.Context, string) (map[string]any, error)

	ListDefinitions(context.Context) ([]map[string]any, error)
	CreateDefinition(context.Context, map[string]any) (map[string]any, error)
	UpdateDefinition(context.Context, string, map[string]any) (map[string]any, error)
	DeactivateDefinition(context.Context, string) (map[string]any, error)
	AddEnumValue(context.Context, string, map[string]any) (map[string]any, error)
	DeactivateEnumValue(context.Context, string, string) (map[string]any, error)

	ListMasters(context.Context, string) ([]map[string]any, error)
	CreateMaster(context.Context, string, map[string]any) (map[string]any, error)
	UpdateMaster(context.Context, string, string, map[string]any) (map[string]any, error)
	DeactivateMaster(context.Context, string, string) (map[string]any, error)

	CreateModel(context.Context, map[string]any) (map[string]any, error)
	ListModels(context.Context) ([]map[string]any, error)
	GetModel(context.Context, string) (map[string]any, error)
	UpdateModel(context.Context, string, map[string]any) (map[string]any, error)
	PublishModel(context.Context, string) (map[string]any, error)
	UnpublishModel(context.Context, string) (map[string]any, error)
	DiscontinueModel(context.Context, string) (map[string]any, error)
	DeleteModel(context.Context, string) (map[string]any, error)
	ReplaceMedia(context.Context, string, map[string]any) (map[string]any, error)

	CreateDimension(context.Context, string, map[string]any) (map[string]any, error)
	UpdateDimension(context.Context, string, string, map[string]any) (map[string]any, error)
	AddDimensionValue(context.Context, string, string, map[string]any) (map[string]any, error)
	DeactivateDimensionValue(context.Context, string, string, string) (map[string]any, error)

	CreateVariant(context.Context, string, map[string]any) (map[string]any, error)
	ListVariants(context.Context, string) ([]map[string]any, error)
	GetVariant(context.Context, string) (map[string]any, error)
	UpdateVariant(context.Context, string, map[string]any) (map[string]any, error)
	ActivateVariant(context.Context, string) (map[string]any, error)
	InactivateVariant(context.Context, string) (map[string]any, error)
	BulkSetPrice(context.Context, map[string]any) ([]map[string]any, error)

	ListPublicModels(context.Context, PublicFilter) (map[string]any, error)
	GetPublicModel(context.Context, string) (map[string]any, error)
	AvailableOptions(context.Context, string, map[string]string) (map[string]any, error)
	ResolvePublicVariant(context.Context, string, map[string]string) (map[string]any, error)
}

var _ UseCase = (*Service)(nil)
