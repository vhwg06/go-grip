package catalogbase

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type memoryCatalogStore struct {
	snapshot CatalogSnapshot
}

func newMemoryCatalogRepositories(store *memoryCatalogStore) CatalogRepositories {
	return CatalogRepositories{
		Categories:        &memoryCatalogCategoryRepo{store: store},
		Definitions:       &memoryCatalogDefinitionRepo{store: store},
		Masters:           &memoryCatalogMasterRepo{store: store},
		ProductModels:     &memoryCatalogProductModelRepo{store: store},
		ProductImages:     &memoryCatalogImageRepo{store: store},
		VariantDimensions: &memoryCatalogDimensionRepo{store: store},
		Variants:          &memoryCatalogVariantRepo{store: store},
	}
}

type memoryUnitOfWork struct {
	store *memoryCatalogStore
}

func (u *memoryUnitOfWork) Within(ctx context.Context, fn func(context.Context) error) error {
	original := cloneCatalogSnapshot(u.store.snapshot)
	u.store.snapshot = cloneCatalogSnapshot(u.store.snapshot)
	if err := fn(ctx); err != nil {
		u.store.snapshot = original
		return err
	}
	return nil
}

type memoryCatalogCategoryRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogCategoryRepo) List(context.Context) ([]CatalogCategory, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	return snapshot.Categories, nil
}

func (r *memoryCatalogCategoryRepo) GetByID(_ context.Context, id string) (CatalogCategory, error) {
	for _, category := range r.store.snapshot.Categories {
		if category.ID == id {
			return category, nil
		}
	}
	return CatalogCategory{}, fmt.Errorf("category %s not found", id)
}

func (r *memoryCatalogCategoryRepo) Store(_ context.Context, category CatalogCategory) error {
	r.store.snapshot.Categories = append(r.store.snapshot.Categories, category)
	return nil
}

func (r *memoryCatalogCategoryRepo) Update(_ context.Context, category CatalogCategory) error {
	for index := range r.store.snapshot.Categories {
		if r.store.snapshot.Categories[index].ID == category.ID {
			r.store.snapshot.Categories[index] = category
			return nil
		}
	}
	return fmt.Errorf("category %s not found", category.ID)
}

func (r *memoryCatalogCategoryRepo) Delete(_ context.Context, id string) error {
	for index := range r.store.snapshot.Categories {
		if r.store.snapshot.Categories[index].ID == id {
			r.store.snapshot.Categories = append(r.store.snapshot.Categories[:index], r.store.snapshot.Categories[index+1:]...)
			return nil
		}
	}
	return nil
}

type memoryCatalogDefinitionRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogDefinitionRepo) List(context.Context) ([]CatalogAttributeDefinition, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	return snapshot.Definitions, nil
}

func (r *memoryCatalogDefinitionRepo) GetByID(_ context.Context, id string) (CatalogAttributeDefinition, error) {
	for _, definition := range r.store.snapshot.Definitions {
		if definition.ID == id {
			return definition, nil
		}
	}
	return CatalogAttributeDefinition{}, fmt.Errorf("definition %s not found", id)
}

func (r *memoryCatalogDefinitionRepo) Store(_ context.Context, definition CatalogAttributeDefinition) error {
	r.store.snapshot.Definitions = append(r.store.snapshot.Definitions, definition)
	return nil
}

func (r *memoryCatalogDefinitionRepo) Update(_ context.Context, definition CatalogAttributeDefinition) error {
	for index := range r.store.snapshot.Definitions {
		if r.store.snapshot.Definitions[index].ID == definition.ID {
			r.store.snapshot.Definitions[index] = definition
			return nil
		}
	}
	return fmt.Errorf("definition %s not found", definition.ID)
}

func (r *memoryCatalogDefinitionRepo) Delete(_ context.Context, id string) error {
	for index := range r.store.snapshot.Definitions {
		if r.store.snapshot.Definitions[index].ID == id {
			r.store.snapshot.Definitions = append(r.store.snapshot.Definitions[:index], r.store.snapshot.Definitions[index+1:]...)
			return nil
		}
	}
	return nil
}

type memoryCatalogMasterRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogMasterRepo) List(_ context.Context, kind string) ([]CatalogMaster, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	result := make([]CatalogMaster, 0)
	for _, master := range snapshot.Masters {
		if kind == "" || master.Kind == kind {
			result = append(result, master)
		}
	}
	return result, nil
}

func (r *memoryCatalogMasterRepo) GetByID(_ context.Context, kind, id string) (CatalogMaster, error) {
	for _, master := range r.store.snapshot.Masters {
		if master.Kind == kind && master.ID == id {
			return master, nil
		}
	}
	return CatalogMaster{}, fmt.Errorf("master %s/%s not found", kind, id)
}

func (r *memoryCatalogMasterRepo) Store(_ context.Context, master CatalogMaster) error {
	r.store.snapshot.Masters = append(r.store.snapshot.Masters, master)
	return nil
}

func (r *memoryCatalogMasterRepo) Update(_ context.Context, master CatalogMaster) error {
	for index := range r.store.snapshot.Masters {
		if r.store.snapshot.Masters[index].ID == master.ID {
			r.store.snapshot.Masters[index] = master
			return nil
		}
	}
	return fmt.Errorf("master %s not found", master.ID)
}

func (r *memoryCatalogMasterRepo) Delete(_ context.Context, kind, id string) error {
	for index := range r.store.snapshot.Masters {
		if r.store.snapshot.Masters[index].Kind == kind && r.store.snapshot.Masters[index].ID == id {
			r.store.snapshot.Masters = append(r.store.snapshot.Masters[:index], r.store.snapshot.Masters[index+1:]...)
			return nil
		}
	}
	return nil
}

type memoryCatalogProductModelRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogProductModelRepo) List(context.Context) ([]CatalogProductModel, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	result := make([]CatalogProductModel, 0, len(snapshot.Models))
	for _, model := range snapshot.Models {
		model.Images = nil
		model.Dimensions = nil
		model.Variants = nil
		result = append(result, model)
	}
	return result, nil
}

func (r *memoryCatalogProductModelRepo) GetByID(_ context.Context, id string) (CatalogProductModel, error) {
	for _, model := range r.store.snapshot.Models {
		if model.ID == id {
			model.Images = nil
			model.Dimensions = nil
			model.Variants = nil
			return model, nil
		}
	}
	return CatalogProductModel{}, fmt.Errorf("model %s not found", id)
}

func (r *memoryCatalogProductModelRepo) Store(_ context.Context, model CatalogProductModel) error {
	model.Images = nil
	model.Dimensions = nil
	model.Variants = nil
	r.store.snapshot.Models = append(r.store.snapshot.Models, model)
	return nil
}

func (r *memoryCatalogProductModelRepo) Update(_ context.Context, model CatalogProductModel) error {
	for index := range r.store.snapshot.Models {
		if r.store.snapshot.Models[index].ID == model.ID {
			model.Images = nil
			model.Dimensions = nil
			model.Variants = nil
			r.store.snapshot.Models[index] = model
			return nil
		}
	}
	return fmt.Errorf("model %s not found", model.ID)
}

func (r *memoryCatalogProductModelRepo) Delete(_ context.Context, id string) error {
	for index := range r.store.snapshot.Models {
		if r.store.snapshot.Models[index].ID == id {
			r.store.snapshot.Models = append(r.store.snapshot.Models[:index], r.store.snapshot.Models[index+1:]...)
			return nil
		}
	}
	return nil
}

type memoryCatalogImageRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogImageRepo) ListByModelID(_ context.Context, modelID string) ([]CatalogProductImage, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	model := memoryModel(&memoryCatalogStore{snapshot: snapshot}, modelID)
	return append([]CatalogProductImage(nil), model.Images...), nil
}

func (r *memoryCatalogImageRepo) Store(_ context.Context, modelID string, image CatalogProductImage) error {
	memoryModel(r.store, modelID).Images = append(memoryModel(r.store, modelID).Images, image)
	return nil
}

func (r *memoryCatalogImageRepo) Update(_ context.Context, modelID string, image CatalogProductImage) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Images {
		if model.Images[index].ID == image.ID {
			model.Images[index] = image
			return nil
		}
	}
	return fmt.Errorf("image %s not found", image.ID)
}

func (r *memoryCatalogImageRepo) Delete(_ context.Context, modelID, imageID string) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Images {
		if model.Images[index].ID == imageID {
			model.Images = append(model.Images[:index], model.Images[index+1:]...)
			return nil
		}
	}
	return nil
}

func (r *memoryCatalogImageRepo) DeleteByModelID(_ context.Context, modelID string) error {
	memoryModel(r.store, modelID).Images = nil
	return nil
}

type memoryCatalogDimensionRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogDimensionRepo) ListByModelID(_ context.Context, modelID string) ([]CatalogVariantDimension, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	model := memoryModel(&memoryCatalogStore{snapshot: snapshot}, modelID)
	return append([]CatalogVariantDimension(nil), model.Dimensions...), nil
}

func (r *memoryCatalogDimensionRepo) GetByID(_ context.Context, modelID, id string) (CatalogVariantDimension, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	for _, dimension := range memoryModel(&memoryCatalogStore{snapshot: snapshot}, modelID).Dimensions {
		if dimension.ID == id {
			return dimension, nil
		}
	}
	return CatalogVariantDimension{}, fmt.Errorf("dimension %s not found", id)
}

func (r *memoryCatalogDimensionRepo) Store(_ context.Context, modelID string, dimension CatalogVariantDimension) error {
	memoryModel(r.store, modelID).Dimensions = append(memoryModel(r.store, modelID).Dimensions, dimension)
	return nil
}

func (r *memoryCatalogDimensionRepo) Update(_ context.Context, modelID string, dimension CatalogVariantDimension) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Dimensions {
		if model.Dimensions[index].ID == dimension.ID {
			model.Dimensions[index] = dimension
			return nil
		}
	}
	return fmt.Errorf("dimension %s not found", dimension.ID)
}

func (r *memoryCatalogDimensionRepo) Delete(_ context.Context, modelID, id string) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Dimensions {
		if model.Dimensions[index].ID == id {
			model.Dimensions = append(model.Dimensions[:index], model.Dimensions[index+1:]...)
			return nil
		}
	}
	return nil
}

func (r *memoryCatalogDimensionRepo) DeleteByModelID(_ context.Context, modelID string) error {
	memoryModel(r.store, modelID).Dimensions = nil
	return nil
}

type memoryCatalogVariantRepo struct{ store *memoryCatalogStore }

func (r *memoryCatalogVariantRepo) List(_ context.Context) ([]CatalogVariant, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	result := make([]CatalogVariant, 0)
	for _, model := range snapshot.Models {
		result = append(result, model.Variants...)
	}
	return result, nil
}

func (r *memoryCatalogVariantRepo) ListByModelID(_ context.Context, modelID string) ([]CatalogVariant, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	model := memoryModel(&memoryCatalogStore{snapshot: snapshot}, modelID)
	return append([]CatalogVariant(nil), model.Variants...), nil
}

func (r *memoryCatalogVariantRepo) GetByID(_ context.Context, modelID, id string) (CatalogVariant, error) {
	snapshot := cloneCatalogSnapshot(r.store.snapshot)
	for _, variant := range memoryModel(&memoryCatalogStore{snapshot: snapshot}, modelID).Variants {
		if variant.ID == id {
			return variant, nil
		}
	}
	return CatalogVariant{}, fmt.Errorf("variant %s not found", id)
}

func (r *memoryCatalogVariantRepo) Store(_ context.Context, modelID string, variant CatalogVariant) error {
	memoryModel(r.store, modelID).Variants = append(memoryModel(r.store, modelID).Variants, variant)
	return nil
}

func (r *memoryCatalogVariantRepo) Update(_ context.Context, modelID string, variant CatalogVariant) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Variants {
		if model.Variants[index].ID == variant.ID {
			model.Variants[index] = variant
			return nil
		}
	}
	return fmt.Errorf("variant %s not found", variant.ID)
}

func (r *memoryCatalogVariantRepo) Delete(_ context.Context, modelID, id string) error {
	model := memoryModel(r.store, modelID)
	for index := range model.Variants {
		if model.Variants[index].ID == id {
			model.Variants = append(model.Variants[:index], model.Variants[index+1:]...)
			return nil
		}
	}
	return nil
}

func (r *memoryCatalogVariantRepo) DeleteByModelID(_ context.Context, modelID string) error {
	memoryModel(r.store, modelID).Variants = nil
	return nil
}

func memoryModel(store *memoryCatalogStore, id string) *CatalogProductModel {
	for index := range store.snapshot.Models {
		if store.snapshot.Models[index].ID == id {
			return &store.snapshot.Models[index]
		}
	}
	panic(fmt.Sprintf("model %s not found", id))
}

func cloneCatalogSnapshot(snapshot CatalogSnapshot) CatalogSnapshot {
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}
	var clone CatalogSnapshot
	if err := json.Unmarshal(encoded, &clone); err != nil {
		panic(err)
	}
	return clone
}

func newCatalogBaseTestService(t *testing.T) *Service {
	t.Helper()
	store := &memoryCatalogStore{}
	return New(newMemoryCatalogRepositories(store), &memoryUnitOfWork{store: store})
}

func catalogCategory(t *testing.T, service *Service) string {
	t.Helper()
	category, err := service.CreateCategory(context.Background(), map[string]any{
		"name": "Grip handles",
		"slug": "grip-handles-" + t.Name(),
	})
	require.NoError(t, err)
	return mapString(t, category, "id")
}

func mapString(t *testing.T, object map[string]any, key string) string {
	t.Helper()
	value, ok := object[key].(string)
	require.True(t, ok)
	return value
}

func mapBool(t *testing.T, object map[string]any, key string) bool {
	t.Helper()
	value, ok := object[key].(bool)
	require.True(t, ok)
	return value
}

func catalogDimension(t *testing.T, service *Service, modelID, displayName string, values []map[string]any) string {
	t.Helper()
	definition, err := service.CreateDefinition(context.Background(), map[string]any{
		"key":         displayName + "-key-" + t.Name(),
		"displayName": displayName,
		"valueKind":   "Scalar",
		"dataType":    "Number",
		"unitFamily":  "length",
		"unit":        "mm",
	})
	require.NoError(t, err)
	allowedValues := make([]any, 0, len(values))
	for _, value := range values {
		allowedValues = append(allowedValues, value)
	}
	dimension, err := service.CreateDimension(context.Background(), modelID, map[string]any{
		"definitionId":  definition["id"],
		"allowedValues": allowedValues,
	})
	require.NoError(t, err)
	return mapString(t, dimension, "id")
}

func catalogVariant(t *testing.T, service *Service, modelID, selected, sku string, amount int64) map[string]any {
	t.Helper()
	variant, err := service.CreateVariant(context.Background(), modelID, map[string]any{
		"selectedOptions": map[string]any{"Size": selected},
		"sku":             sku,
		"sellingPrice":    map[string]any{"amount": amount, "currency": "VND"},
	})
	require.NoError(t, err)
	return variant
}

func TestCatalogBaseLifecycleAndCanonicalPublicProjection(t *testing.T) { //nolint:funlen // This scenario verifies the complete Catalog Base lifecycle contract.
	t.Parallel()
	service := newCatalogBaseTestService(t)
	ctx := context.Background()
	categoryID := catalogCategory(t, service)

	_, err := service.CreateDefinition(ctx, map[string]any{
		"key": "bad-reference", "valueKind": "Reference", "referenceTarget": "Material", "dataType": "Number",
	})
	require.Error(t, err)
	_, err = service.CreateDefinition(ctx, map[string]any{
		"key": "bad-unit", "valueKind": "Scalar", "dataType": "Number", "unitFamily": "length", "unit": "kg",
	})
	require.Error(t, err)

	model, err := service.CreateModel(ctx, map[string]any{
		"name": "Grip Handle A", "categoryId": categoryID,
		"description": "Catalog Base model", "warrantySummary": map[string]any{"term": "24 months"},
	})
	require.NoError(t, err)
	modelID := mapString(t, model, "id")
	require.Equal(t, "Draft", model["status"])

	dimensionID := catalogDimension(t, service, modelID, "Size", []map[string]any{
		{"id": "200-mm", "label": "200 mm", "active": true},
		{"id": "300-mm", "label": "300 mm", "active": true},
	})
	variant := catalogVariant(t, service, modelID, "200 mm", " SKU-001 ", 400000)
	variantID := mapString(t, variant, "id")
	require.Equal(t, "sku-001", variant["sku"])
	selectedOptions, ok := variant["selectedOptions"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "200 mm", selectedOptions["Size"])
	require.True(t, mapBool(t, variant, "saleReady"))

	_, err = service.CreateVariant(ctx, modelID, map[string]any{
		"selectedOptions": map[string]any{"Size": "20 cm"},
		"sku":             "SKU-002", "sellingPrice": map[string]any{"amount": 450000, "currency": "VND"},
	})
	require.Error(t, err, "20 cm and 200 mm must share one canonical combination")

	_, err = service.ReplaceMedia(ctx, modelID, map[string]any{"images": []any{
		map[string]any{"url": "https://cdn.example.test/primary.png", "ordering": 1, "primary": true},
	}})
	require.NoError(t, err)
	active, err := service.PublishModel(ctx, modelID)
	require.NoError(t, err)
	require.Equal(t, "Active", active["status"])

	public, err := service.GetPublicModel(ctx, modelID)
	require.NoError(t, err)
	require.Len(t, public["variants"], 1)
	publicList, err := service.ListPublicModels(ctx, PublicFilter{CategoryID: categoryID, MinPrice: ptrInt64(400000), MaxPrice: ptrInt64(400000)})
	require.NoError(t, err)
	require.Equal(t, 1, publicList["total"])

	resolved, err := service.ResolvePublicVariant(ctx, modelID, map[string]string{"Size": "20 cm"})
	require.NoError(t, err)
	require.Equal(t, variantID, resolved["id"])

	_, err = service.InactivateVariant(ctx, variantID)
	require.Error(t, err, "an active model must retain one sale-ready Variant")
	_, err = service.UpdateVariant(ctx, variantID, map[string]any{"sku": ""})
	require.Error(t, err, "an active model cannot lose its last sale-ready commercial identity")

	unpublished, err := service.UnpublishModel(ctx, modelID)
	require.NoError(t, err)
	require.Equal(t, "Inactive", unpublished["status"])
	_, err = service.InactivateVariant(ctx, variantID)
	require.NoError(t, err)
	_, err = service.GetPublicModel(ctx, modelID)
	require.Error(t, err, "inactive ProductModels are not publicly readable")

	// Adding a selectable value is explicitly allowed after a Variant exists;
	// replacing a dimension definition is not.
	_, err = service.AddDimensionValue(ctx, modelID, dimensionID, map[string]any{"id": "400-mm", "label": "400 mm"})
	require.NoError(t, err)
	replacement, err := service.CreateDefinition(ctx, map[string]any{"key": "replacement-" + t.Name(), "displayName": "Replacement", "valueKind": "Enum"})
	require.NoError(t, err)
	_, err = service.UpdateDimension(ctx, modelID, dimensionID, map[string]any{"definitionId": replacement["id"]})
	require.Error(t, err)
}

func TestCatalogBaseReferencePackLifecycleAndAtomicPrices(t *testing.T) {
	t.Parallel()
	service := newCatalogBaseTestService(t)
	ctx := context.Background()
	categoryID := catalogCategory(t, service)
	pack, err := service.CreateMaster(ctx, "pack", map[string]any{
		"name": "Box 10", "sellingUnit": "Box", "quantity": 10, "baseUnit": "Piece",
	})
	require.NoError(t, err)
	packID := mapString(t, pack, "id")

	model, err := service.CreateModel(ctx, map[string]any{"name": "Packed handle", "categoryId": categoryID, "fixedPackId": packID})
	require.NoError(t, err)
	modelID := mapString(t, model, "id")
	catalogDimension(t, service, modelID, "Size", []map[string]any{{"id": "200-mm", "label": "200 mm", "active": true}, {"id": "300-mm", "label": "300 mm", "active": true}})

	first := catalogVariant(t, service, modelID, "200 mm", "SKU-PACK-1", 100000)
	require.Equal(t, packID, first["packId"])
	second := catalogVariant(t, service, modelID, "300 mm", "SKU-PACK-2", 100000)

	_, err = service.DeactivateMaster(ctx, "pack", packID)
	require.NoError(t, err)
	_, err = service.CreateVariant(ctx, modelID, map[string]any{
		"selectedOptions": map[string]any{"Size": "200 mm"}, "sku": "SKU-PACK-3",
		"sellingPrice": map[string]any{"amount": 100000, "currency": "VND"},
	})
	require.Error(t, err, "new references to an inactive Pack must be rejected")
	readFirst, err := service.GetVariant(ctx, mapString(t, first, "id"))
	require.NoError(t, err)
	require.Equal(t, packID, readFirst["packId"])

	_, err = service.BulkSetPrice(ctx, map[string]any{
		"variantIds":   []any{first["id"], "missing-variant"},
		"sellingPrice": map[string]any{"amount": 250000, "currency": "VND"},
	})
	require.Error(t, err)
	unchanged, err := service.GetVariant(ctx, mapString(t, first, "id"))
	require.NoError(t, err)
	require.Equal(t, int64(100000), mapStringMap(t, unchanged, "sellingPrice")["amount"])

	updated, err := service.BulkSetPrice(ctx, map[string]any{
		"variantIds":   []any{first["id"], second["id"]},
		"sellingPrice": map[string]any{"amount": 250000, "currency": "VND"},
	})
	require.NoError(t, err)
	require.Len(t, updated, 2)
	for _, item := range updated {
		require.Equal(t, int64(250000), mapStringMap(t, item, "sellingPrice")["amount"])
	}
}

func TestCatalogBaseCanonicalizesMeasurementsAndReferenceIdentity(t *testing.T) {
	t.Parallel()
	service := newCatalogBaseTestService(t)
	ctx := context.Background()
	categoryID := catalogCategory(t, service)

	material, err := service.CreateMaster(ctx, "material", map[string]any{"name": "Inox 304"})
	require.NoError(t, err)
	materialID := mapString(t, material, "id")
	reference, err := service.CreateDefinition(ctx, map[string]any{
		"key":             "material-reference-" + t.Name(),
		"displayName":     "Material",
		"ordering":        20,
		"valueKind":       "Reference",
		"referenceTarget": "Material",
	})
	require.NoError(t, err)
	referenceID := mapString(t, reference, "id")
	surface, err := service.CreateDefinition(ctx, map[string]any{
		"key":         "surface-" + t.Name(),
		"displayName": "Surface",
		"ordering":    10,
		"valueKind":   "Enum",
	})
	require.NoError(t, err)
	surfaceID := mapString(t, surface, "id")
	surfaceValue, err := service.AddEnumValue(ctx, surfaceID, map[string]any{
		"key":   "brushed",
		"label": "Brushed",
	})
	require.NoError(t, err)
	variantReference, err := service.CreateDefinition(ctx, map[string]any{
		"key": "variant-material-reference-" + t.Name(), "displayName": "Variant Material", "valueKind": "Reference", "referenceTarget": "Material",
	})
	require.NoError(t, err)
	variantReferenceID := mapString(t, variantReference, "id")

	_, err = service.CreateDefinition(ctx, map[string]any{
		"key": "length-" + t.Name(), "displayName": "Overall length", "valueKind": "Scalar", "dataType": "Number", "unitFamily": "length", "unit": "mm",
	})
	require.NoError(t, err)

	model, err := service.CreateModel(ctx, map[string]any{
		"name": "Reference model", "categoryId": categoryID,
		"fixedAttributes": map[string]any{
			referenceID: "Inox 304",
			surfaceID:   surfaceValue["id"],
		},
		"measurements": map[string]any{"overallLength": map[string]any{"value": 20, "unit": "cm"}},
	})
	require.NoError(t, err)
	modelID := mapString(t, model, "id")
	fixed := mapStringMap(t, model, "fixedAttributes")
	require.Equal(t, materialID, fixed[referenceID])
	require.Equal(t, surfaceValue["id"], fixed[surfaceID])
	measurement := mapStringMap(t, model, "measurements")
	require.Equal(t, float64(200), measurementMapNumber(t, measurement["overallLength"], "value"))
	require.Equal(t, "mm", mapStringMap(t, measurement, "overallLength")["unit"])
	require.Equal(t, []map[string]any{
		{
			"key":   "Surface",
			"value": "Brushed",
		},
		{
			"key":   "Material",
			"value": "Inox 304",
		},
		{
			"key":   "Overall length",
			"value": "200 mm",
		},
	}, model["specs"])

	dimension, err := service.CreateDimension(ctx, modelID, map[string]any{
		"definitionId":  variantReferenceID,
		"allowedValues": []any{map[string]any{"id": materialID, "label": "Inox 304", "active": true}},
	})
	require.NoError(t, err)
	_, err = service.CreateVariant(ctx, modelID, map[string]any{
		"selectedOptions": map[string]any{"Variant Material": "Inox 304"},
		"sku":             "REF-001",
		"sellingPrice":    map[string]any{"amount": 100000, "currency": "VND"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, dimension["id"])
	variants, err := service.ListVariants(ctx, modelID)
	require.NoError(t, err)
	require.Len(t, variants, 1)
	require.Equal(t, materialID, mapStringMap(t, variants[0], "selectedOptions")["Variant Material"])
	_, err = service.DeactivateEnumValue(ctx, surfaceID, mapString(t, surfaceValue, "id"))
	require.NoError(t, err)
	_, err = service.DeactivateMaster(ctx, "material", materialID)
	require.NoError(t, err)

	_, err = service.ReplaceMedia(ctx, modelID, map[string]any{"images": []any{
		map[string]any{"url": "https://cdn.example.test/reference.png", "ordering": 1, "primary": true},
	}})
	require.NoError(t, err)
	_, err = service.PublishModel(ctx, modelID)
	require.NoError(t, err)
	publicModel, err := service.GetPublicModel(ctx, modelID)
	require.NoError(t, err)
	require.Equal(t, model["specs"], publicModel["specs"])
}

func TestCatalogBaseValidatesVariantTechnicalValuesAgainstDeclaredDefinitions(t *testing.T) {
	service := newCatalogBaseTestService(t)
	ctx := context.Background()
	categoryID := catalogCategory(t, service)
	model, err := service.CreateModel(ctx, map[string]any{"name": "Technical values", "categoryId": categoryID})
	require.NoError(t, err)
	modelID := mapString(t, model, "id")
	_, err = service.CreateDefinition(ctx, map[string]any{
		"key": "weight-" + t.Name(), "displayName": "Weight", "valueKind": "Scalar", "dataType": "Number", "unitFamily": "mass", "unit": "kg",
	})
	require.NoError(t, err)
	sizeDefinition, err := service.CreateDefinition(ctx, map[string]any{"key": "size-" + t.Name(), "displayName": "Size", "valueKind": "Enum"})
	require.NoError(t, err)
	dimension, err := service.CreateDimension(ctx, modelID, map[string]any{
		"definitionId":  sizeDefinition["id"],
		"allowedValues": []any{map[string]any{"id": "small", "label": "Small", "active": true}},
	})
	require.NoError(t, err)
	_, err = service.CreateVariant(ctx, modelID, map[string]any{
		"selectedOptions": map[string]any{"Size": "Small"},
		"technicalValues": map[string]any{"Weight": "1.2 kg"},
	})
	require.NoError(t, err)
	_, err = service.AddDimensionValue(ctx, modelID, mapString(t, dimension, "id"), map[string]any{"id": "large", "label": "Large"})
	require.NoError(t, err)

	_, err = service.CreateVariant(ctx, modelID, map[string]any{
		"selectedOptions": map[string]any{"Size": "Large"},
		"technicalValues": map[string]any{"Projection": "60 mm"},
	})
	require.Error(t, err, "undeclared technical attributes must not be accepted")
}

func measurementMapNumber(t *testing.T, value any, key string) float64 {
	t.Helper()
	measurement := mapStringMap(t, map[string]any{"measurement": value}, "measurement")
	parsed, ok := measurement[key].(float64)
	if ok {
		return parsed
	}
	integer, ok := measurement[key].(int)
	if ok {
		return float64(integer)
	}
	return 0
}

func mapStringMap(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := object[key].(map[string]any)
	require.True(t, ok)
	return value
}

func ptrInt64(value int64) *int64 { return &value }
