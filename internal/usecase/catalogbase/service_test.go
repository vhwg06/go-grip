package catalogbase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newCatalogBaseTestService(t *testing.T) *Service {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&categoryRow{},
		&attributeDefinitionRow{},
		&masterRow{},
		&productModelRow{},
		&productImageRow{},
		&variantDimensionRow{},
		&variantRow{},
	))
	return New(db)
}

func catalogCategory(t *testing.T, service *Service) string {
	t.Helper()
	category, err := service.CreateCategory(context.Background(), map[string]any{
		"name": "Grip handles",
		"slug": "grip-handles-" + t.Name(),
	})
	require.NoError(t, err)
	return category["id"].(string)
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
	return dimension["id"].(string)
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

func TestCatalogBaseLifecycleAndCanonicalPublicProjection(t *testing.T) {
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
	modelID := model["id"].(string)
	require.Equal(t, "Draft", model["status"])

	dimensionID := catalogDimension(t, service, modelID, "Size", []map[string]any{
		{"id": "200-mm", "label": "200 mm", "active": true},
		{"id": "300-mm", "label": "300 mm", "active": true},
	})
	variant := catalogVariant(t, service, modelID, "200 mm", " SKU-001 ", 400000)
	variantID := variant["id"].(string)
	require.Equal(t, "sku-001", variant["sku"])
	require.Equal(t, "200 mm", jsonMap(jsonString(variant["selectedOptions"], "{}"))["Size"])
	require.True(t, variant["saleReady"].(bool))

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
	service := newCatalogBaseTestService(t)
	ctx := context.Background()
	categoryID := catalogCategory(t, service)
	pack, err := service.CreateMaster(ctx, "pack", map[string]any{
		"name": "Box 10", "sellingUnit": "Box", "quantity": 10, "baseUnit": "Piece",
	})
	require.NoError(t, err)
	packID := pack["id"].(string)

	model, err := service.CreateModel(ctx, map[string]any{"name": "Packed handle", "categoryId": categoryID, "fixedPackId": packID})
	require.NoError(t, err)
	modelID := model["id"].(string)
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
	readFirst, err := service.GetVariant(ctx, first["id"].(string))
	require.NoError(t, err)
	require.Equal(t, packID, readFirst["packId"])

	_, err = service.BulkSetPrice(ctx, map[string]any{
		"variantIds":   []any{first["id"], "missing-variant"},
		"sellingPrice": map[string]any{"amount": 250000, "currency": "VND"},
	})
	require.Error(t, err)
	unchanged, err := service.GetVariant(ctx, first["id"].(string))
	require.NoError(t, err)
	require.Equal(t, int64(100000), unchanged["sellingPrice"].(map[string]any)["amount"])

	updated, err := service.BulkSetPrice(ctx, map[string]any{
		"variantIds":   []any{first["id"], second["id"]},
		"sellingPrice": map[string]any{"amount": 250000, "currency": "VND"},
	})
	require.NoError(t, err)
	require.Len(t, updated, 2)
	for _, item := range updated {
		require.Equal(t, int64(250000), item["sellingPrice"].(map[string]any)["amount"])
	}
}

func ptrInt64(value int64) *int64 { return &value }
