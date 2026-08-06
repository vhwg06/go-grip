package models

import (
	"encoding/json"
	"time"

	catalogbase "github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
)

type CatalogBaseCategory struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Slug      string    `gorm:"not null"`
	ParentID  *string   `gorm:"type:uuid"`
	Position  int       `gorm:"not null;default:0"`
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (CatalogBaseCategory) TableName() string { return "catalog_base_categories" }

type CatalogBaseDefinition struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	Key             string    `gorm:"not null"`
	DisplayName     string    `gorm:"column:display_name;not null;default:''"`
	Description     string    `gorm:"not null;default:''"`
	Ordering        int       `gorm:"not null;default:0"`
	ValueKind       string    `gorm:"column:value_kind;not null"`
	DataType        string    `gorm:"column:data_type;not null;default:''"`
	ReferenceTarget string    `gorm:"column:reference_target;not null;default:''"`
	UnitFamily      string    `gorm:"column:unit_family;not null;default:''"`
	Unit            string    `gorm:"not null;default:''"`
	Active          bool      `gorm:"not null;default:true"`
	EnumValues      string    `gorm:"column:enum_values;type:jsonb;not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (CatalogBaseDefinition) TableName() string { return "catalog_base_attribute_definitions" }

type CatalogBaseMaster struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Kind        string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null;default:''"`
	SwatchMedia string    `gorm:"column:swatch_media;type:jsonb;not null"`
	SellingUnit string    `gorm:"column:selling_unit;not null;default:''"`
	Quantity    *float64  `gorm:"type:numeric"`
	BaseUnit    string    `gorm:"column:base_unit;not null;default:''"`
	Active      bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (CatalogBaseMaster) TableName() string { return "catalog_base_masters" }

type CatalogBaseProductModel struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	Name            string    `gorm:"not null"`
	CategoryID      string    `gorm:"column:category_id;type:uuid;not null"`
	Description     string    `gorm:"not null;default:''"`
	WarrantySummary string    `gorm:"column:warranty_summary;type:jsonb"`
	FixedAttributes string    `gorm:"column:fixed_attributes;type:jsonb;not null"`
	FixedPackID     *string   `gorm:"column:fixed_pack_id;type:uuid"`
	Measurements    string    `gorm:"type:jsonb;not null"`
	Status          string    `gorm:"not null;default:'Draft'"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (CatalogBaseProductModel) TableName() string { return "catalog_base_product_models" }

type CatalogBaseProductImage struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	ModelID      string    `gorm:"column:model_id;type:uuid;not null"`
	URL          string    `gorm:"not null"`
	Ordering     int       `gorm:"not null;default:0"`
	PrimaryImage bool      `gorm:"column:primary_image;not null;default:false"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (CatalogBaseProductImage) TableName() string { return "catalog_base_product_images" }

type CatalogBaseDimension struct {
	ID            string    `gorm:"type:uuid;primaryKey"`
	ModelID       string    `gorm:"column:model_id;type:uuid;not null"`
	DefinitionID  string    `gorm:"column:definition_id;type:uuid;not null"`
	AllowedValues string    `gorm:"column:allowed_values;type:jsonb;not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (CatalogBaseDimension) TableName() string { return "catalog_base_variant_dimensions" }

type CatalogBaseVariant struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	ModelID              string    `gorm:"column:model_id;type:uuid;not null"`
	SelectedOptions      string    `gorm:"column:selected_options;type:jsonb;not null"`
	TechnicalValues      string    `gorm:"column:technical_values;type:jsonb;not null"`
	SKU                  string    `gorm:"not null;default:''"`
	SellingAmount        *int64    `gorm:"column:selling_amount"`
	SellingCurrency      string    `gorm:"column:selling_currency;not null;default:''"`
	PackID               *string   `gorm:"column:pack_id;type:uuid"`
	Status               string    `gorm:"not null;default:'Active'"`
	CanonicalCombination string    `gorm:"column:canonical_combination;not null"`
	History              string    `gorm:"type:jsonb;not null"`
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

func (CatalogBaseVariant) TableName() string { return "catalog_base_variants" }

func marshalJSON(value any, fallback string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fallback
	}
	return string(encoded)
}

func unmarshalJSON(raw, fallback string, target any) {
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), target); err == nil {
			return
		}
	}
	if err := json.Unmarshal([]byte(fallback), target); err != nil {
		return
	}
}

func CategoryToCatalogEntity(row CatalogBaseCategory) catalogbase.CatalogCategory {
	return catalogbase.CatalogCategory{ID: row.ID, Name: row.Name, Slug: row.Slug, ParentID: row.ParentID, Position: row.Position, Active: row.Active, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}

func CatalogEntityToCategory(value catalogbase.CatalogCategory) CatalogBaseCategory {
	return CatalogBaseCategory{ID: value.ID, Name: value.Name, Slug: value.Slug, ParentID: value.ParentID, Position: value.Position, Active: value.Active, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt}
}

func DefinitionToCatalogEntity(row CatalogBaseDefinition) catalogbase.CatalogAttributeDefinition {
	var values []catalogbase.CatalogEnumValue
	unmarshalJSON(row.EnumValues, "[]", &values)
	return catalogbase.CatalogAttributeDefinition{ID: row.ID, Key: row.Key, DisplayName: row.DisplayName, Description: row.Description, Ordering: row.Ordering, ValueKind: row.ValueKind, DataType: row.DataType, ReferenceTarget: row.ReferenceTarget, UnitFamily: row.UnitFamily, Unit: row.Unit, Active: row.Active, EnumValues: values, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}

func CatalogEntityToDefinition(value catalogbase.CatalogAttributeDefinition) CatalogBaseDefinition {
	return CatalogBaseDefinition{ID: value.ID, Key: value.Key, DisplayName: value.DisplayName, Description: value.Description, Ordering: value.Ordering, ValueKind: value.ValueKind, DataType: value.DataType, ReferenceTarget: value.ReferenceTarget, UnitFamily: value.UnitFamily, Unit: value.Unit, Active: value.Active, EnumValues: marshalJSON(value.EnumValues, "[]"), CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt}
}

func MasterToCatalogEntity(row CatalogBaseMaster) catalogbase.CatalogMaster {
	var media []string
	unmarshalJSON(row.SwatchMedia, "[]", &media)
	return catalogbase.CatalogMaster{ID: row.ID, Kind: row.Kind, Name: row.Name, Description: row.Description, SwatchMedia: media, SellingUnit: row.SellingUnit, Quantity: row.Quantity, BaseUnit: row.BaseUnit, Active: row.Active, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}

func CatalogEntityToMaster(value catalogbase.CatalogMaster) CatalogBaseMaster {
	return CatalogBaseMaster{ID: value.ID, Kind: value.Kind, Name: value.Name, Description: value.Description, SwatchMedia: marshalJSON(value.SwatchMedia, "[]"), SellingUnit: value.SellingUnit, Quantity: value.Quantity, BaseUnit: value.BaseUnit, Active: value.Active, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt}
}

func ProductModelToCatalogEntity(row CatalogBaseProductModel, images []CatalogBaseProductImage, dimensions []CatalogBaseDimension, variants []CatalogBaseVariant) catalogbase.CatalogProductModel {
	var warranty, fixed, measurements map[string]any
	unmarshalJSON(row.WarrantySummary, "{}", &warranty)
	unmarshalJSON(row.FixedAttributes, "{}", &fixed)
	unmarshalJSON(row.Measurements, "{}", &measurements)
	model := catalogbase.CatalogProductModel{ID: row.ID, Name: row.Name, CategoryID: row.CategoryID, Description: row.Description, WarrantySummary: warranty, FixedAttributes: fixed, FixedPackID: row.FixedPackID, Measurements: measurements, Status: row.Status, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
	model.Images = make([]catalogbase.CatalogProductImage, 0, len(images))
	for _, image := range images {
		model.Images = append(model.Images, catalogbase.CatalogProductImage{ID: image.ID, URL: image.URL, Ordering: image.Ordering, PrimaryImage: image.PrimaryImage, CreatedAt: image.CreatedAt})
	}
	model.Dimensions = make([]catalogbase.CatalogVariantDimension, 0, len(dimensions))
	for _, dimension := range dimensions {
		var values []catalogbase.CatalogDimensionValue
		unmarshalJSON(dimension.AllowedValues, "[]", &values)
		model.Dimensions = append(model.Dimensions, catalogbase.CatalogVariantDimension{ID: dimension.ID, DefinitionID: dimension.DefinitionID, AllowedValues: values, CreatedAt: dimension.CreatedAt, UpdatedAt: dimension.UpdatedAt})
	}
	model.Variants = make([]catalogbase.CatalogVariant, 0, len(variants))
	for _, variant := range variants {
		var selected map[string]string
		var technical map[string]any
		var history []catalogbase.CatalogHistoryEntry
		unmarshalJSON(variant.SelectedOptions, "{}", &selected)
		unmarshalJSON(variant.TechnicalValues, "{}", &technical)
		unmarshalJSON(variant.History, "[]", &history)
		var price *catalogbase.CatalogMoney
		if variant.SellingAmount != nil {
			price = &catalogbase.CatalogMoney{Amount: *variant.SellingAmount, Currency: variant.SellingCurrency}
		}
		model.Variants = append(model.Variants, catalogbase.CatalogVariant{ID: variant.ID, SelectedOptions: selected, TechnicalValues: technical, SKU: variant.SKU, SellingPrice: price, PackID: variant.PackID, Status: variant.Status, CanonicalCombination: variant.CanonicalCombination, History: history, CreatedAt: variant.CreatedAt, UpdatedAt: variant.UpdatedAt})
	}
	return model
}

func CatalogEntityToProductModel(value catalogbase.CatalogProductModel) CatalogBaseProductModel {
	return CatalogBaseProductModel{ID: value.ID, Name: value.Name, CategoryID: value.CategoryID, Description: value.Description, WarrantySummary: marshalJSON(value.WarrantySummary, "{}"), FixedAttributes: marshalJSON(value.FixedAttributes, "{}"), FixedPackID: value.FixedPackID, Measurements: marshalJSON(value.Measurements, "{}"), Status: value.Status, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt}
}

func CatalogEntityToImages(modelID string, values []catalogbase.CatalogProductImage) []CatalogBaseProductImage {
	rows := make([]CatalogBaseProductImage, 0, len(values))
	for _, value := range values {
		rows = append(rows, CatalogBaseProductImage{ID: value.ID, ModelID: modelID, URL: value.URL, Ordering: value.Ordering, PrimaryImage: value.PrimaryImage, CreatedAt: value.CreatedAt})
	}
	return rows
}

func CatalogEntityToDimensions(modelID string, values []catalogbase.CatalogVariantDimension) []CatalogBaseDimension {
	rows := make([]CatalogBaseDimension, 0, len(values))
	for _, value := range values {
		rows = append(rows, CatalogBaseDimension{ID: value.ID, ModelID: modelID, DefinitionID: value.DefinitionID, AllowedValues: marshalJSON(value.AllowedValues, "[]"), CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt})
	}
	return rows
}

func CatalogEntityToVariants(modelID string, values []catalogbase.CatalogVariant) []CatalogBaseVariant {
	rows := make([]CatalogBaseVariant, 0, len(values))
	for _, value := range values {
		var amount *int64
		currency := ""
		if value.SellingPrice != nil {
			amount, currency = &value.SellingPrice.Amount, value.SellingPrice.Currency
		}
		rows = append(rows, CatalogBaseVariant{ID: value.ID, ModelID: modelID, SelectedOptions: marshalJSON(value.SelectedOptions, "{}"), TechnicalValues: marshalJSON(value.TechnicalValues, "{}"), SKU: value.SKU, SellingAmount: amount, SellingCurrency: currency, PackID: value.PackID, Status: value.Status, CanonicalCombination: value.CanonicalCombination, History: marshalJSON(value.History, "[]"), CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt})
	}
	return rows
}
