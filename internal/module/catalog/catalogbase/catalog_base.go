package catalogbase

import (
	"strings"
	"time"
)

const (
	CatalogDraft        = "Draft"
	CatalogActive       = "Active"
	CatalogInactive     = "Inactive"
	CatalogDiscontinued = "Discontinued"

	CatalogVariantActive   = "Active"
	CatalogVariantInactive = "Inactive"
)

// CatalogCategory represents a category entity inside catalog base.
type CatalogCategory struct {
	ID        string
	Name      string
	Slug      string
	ParentID  *string
	Position  int
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CatalogEnumValue represents a selectable enum value option.
type CatalogEnumValue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

// CatalogAttributeDefinition defines custom technical specifications and attributes.
type CatalogAttributeDefinition struct {
	ID              string
	Key             string
	DisplayName     string
	Description     string
	Ordering        int
	ValueKind       string
	DataType        string
	ReferenceTarget string
	UnitFamily      string
	Unit            string
	Active          bool
	EnumValues      []CatalogEnumValue
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CatalogMaster defines master product metadata.
type CatalogMaster struct {
	ID          string
	Kind        string
	Name        string
	Description string
	SwatchMedia []string
	SellingUnit string
	Quantity    *float64
	BaseUnit    string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CatalogProductImage represents an image asset linked to a catalog product model.
type CatalogProductImage struct {
	ID           string
	URL          string
	Ordering     int
	PrimaryImage bool
	CreatedAt    time.Time
}

// CatalogDimensionValue represents an allowed dimension value.
type CatalogDimensionValue struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

// CatalogVariantDimension defines variant dimensions for a product model.
type CatalogVariantDimension struct {
	ID            string
	DefinitionID  string
	AllowedValues []CatalogDimensionValue
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CatalogHistoryEntry records historical actions performed on catalog entities.
type CatalogHistoryEntry struct {
	Action string    `json:"action"`
	At     time.Time `json:"at"`
}

// CatalogMoney represents price monetary amounts in minor units.
type CatalogMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// CatalogVariant represents a specific SKU variant under a product model.
type CatalogVariant struct {
	ID                   string
	SelectedOptions      map[string]string
	TechnicalValues      map[string]any
	SKU                  string
	SellingPrice         *CatalogMoney
	PackID               *string
	Status               string
	CanonicalCombination string
	History              []CatalogHistoryEntry
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// CatalogProductModel represents a product model grouping variants, dimensions, and images.
type CatalogProductModel struct {
	ID              string
	Name            string
	CategoryID      string
	Description     string
	WarrantySummary map[string]any
	FixedAttributes map[string]any
	FixedPackID     *string
	Measurements    map[string]any
	Status          string
	Images          []CatalogProductImage
	Dimensions      []CatalogVariantDimension
	Variants        []CatalogVariant
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CatalogSnapshot is the application-level composition used to evaluate
// Catalog Base invariants across separately persisted entities.
type CatalogSnapshot struct {
	Categories  []CatalogCategory
	Definitions []CatalogAttributeDefinition
	Masters     []CatalogMaster
	Models      []CatalogProductModel
}

// HasPrimaryImage reports whether the model has exactly one primary image.
func (m CatalogProductModel) HasPrimaryImage() bool {
	count := 0
	for _, image := range m.Images {
		if image.PrimaryImage {
			count++
		}
	}
	return count == 1
}

// SaleReadyVariantCount returns the number of active, sale-ready variants for the model.
func (m CatalogProductModel) SaleReadyVariantCount() int {
	count := 0
	for _, variant := range m.Variants {
		if variant.SaleReady() {
			count++
		}
	}
	return count
}

// SaleReady reports whether a variant is active with valid SKU and price.
func (v CatalogVariant) SaleReady() bool {
	return v.Status == CatalogVariantActive && v.SKU != "" && v.SellingPrice != nil && v.SellingPrice.Amount > 0 && strings.EqualFold(v.SellingPrice.Currency, "VND")
}
