package entity

import (
	"strings"
	"time"
)

// Catalog Base is kept as a small set of aggregates.  Persistence models are
// deliberately not part of this package; the repository maps these values to
// PostgreSQL rows at the infrastructure boundary.

const (
	CatalogDraft        = "Draft"
	CatalogActive       = "Active"
	CatalogInactive     = "Inactive"
	CatalogDiscontinued = "Discontinued"

	CatalogVariantActive   = "Active"
	CatalogVariantInactive = "Inactive"
)

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

type CatalogEnumValue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

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

type CatalogProductImage struct {
	ID           string
	URL          string
	Ordering     int
	PrimaryImage bool
	CreatedAt    time.Time
}

type CatalogDimensionValue struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type CatalogVariantDimension struct {
	ID            string
	DefinitionID  string
	AllowedValues []CatalogDimensionValue
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CatalogHistoryEntry struct {
	Action string    `json:"action"`
	At     time.Time `json:"at"`
}

type CatalogMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

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
// Catalog Base invariants across separately persisted entities. It is not a
// persistence contract; repository CRUD and transaction coordination remain
// outside the domain package.
type CatalogSnapshot struct {
	Categories  []CatalogCategory
	Definitions []CatalogAttributeDefinition
	Masters     []CatalogMaster
	Models      []CatalogProductModel
}

func (m CatalogProductModel) HasPrimaryImage() bool {
	count := 0
	for _, image := range m.Images {
		if image.PrimaryImage {
			count++
		}
	}
	return count == 1
}

func (m CatalogProductModel) SaleReadyVariantCount() int {
	count := 0
	for _, variant := range m.Variants {
		if variant.SaleReady() {
			count++
		}
	}
	return count
}

func (v CatalogVariant) SaleReady() bool {
	return v.Status == CatalogVariantActive && v.SKU != "" && v.SellingPrice != nil && v.SellingPrice.Amount > 0 && strings.EqualFold(v.SellingPrice.Currency, "VND")
}
