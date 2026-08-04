package catalogbase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/require"
)

type memoryCatalogBaseRepository struct {
	snapshot entity.CatalogSnapshot
}

func (r *memoryCatalogBaseRepository) LoadCatalogBase(context.Context) (entity.CatalogSnapshot, error) {
	return cloneCatalogSnapshot(r.snapshot), nil
}

func (r *memoryCatalogBaseRepository) SaveCatalogBase(_ context.Context, snapshot entity.CatalogSnapshot) error {
	r.snapshot = cloneCatalogSnapshot(snapshot)
	return nil
}

func cloneCatalogSnapshot(snapshot entity.CatalogSnapshot) entity.CatalogSnapshot {
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		panic(err)
	}
	var clone entity.CatalogSnapshot
	if err := json.Unmarshal(encoded, &clone); err != nil {
		panic(err)
	}
	return clone
}

func newCatalogBaseTestService(t *testing.T) *Service {
	t.Helper()
	return New(&memoryCatalogBaseRepository{})
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
		"key": "material-reference-" + t.Name(), "displayName": "Material", "valueKind": "Reference", "referenceTarget": "Material",
	})
	require.NoError(t, err)
	referenceID := mapString(t, reference, "id")
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
		"fixedAttributes": map[string]any{referenceID: "Inox 304"},
		"measurements":    map[string]any{"overallLength": map[string]any{"value": 20, "unit": "cm"}},
	})
	require.NoError(t, err)
	modelID := mapString(t, model, "id")
	fixed := mapStringMap(t, model, "fixedAttributes")
	require.Equal(t, materialID, fixed[referenceID])
	measurement := mapStringMap(t, model, "measurements")
	require.Equal(t, float64(200), measurementMapNumber(t, measurement["overallLength"], "value"))
	require.Equal(t, "mm", mapStringMap(t, measurement, "overallLength")["unit"])

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
