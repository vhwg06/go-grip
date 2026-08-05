package importer

import "time"

const (
	// MaxInitialImportItems defines default batch import limit.
	MaxInitialImportItems = 25
)

// ImportItemType represents supported bulk import entity categories.
type ImportItemType string

const (
	ImportItemProduct ImportItemType = "product"
	ImportItemPost    ImportItemType = "post"
)

// ImportItem represents a single item to import.
type ImportItem struct {
	Type ImportItemType `json:"type"`
	Data map[string]any `json:"data"`
}

// ImportFailure records a failed item index and error reason.
type ImportFailure struct {
	Index  int    `json:"index"`
	Reason string `json:"reason"`
}

// ImportResult summarizes the batch import outcome.
type ImportResult struct {
	Imported int             `json:"imported"`
	Failed   []ImportFailure `json:"failed"`
}

// ImportProductDraft contains product data to be stored by repo adapter.
type ImportProductDraft struct {
	ID        string
	Title     string
	SKU       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ImportPostDraft contains post data to be stored by repo adapter.
type ImportPostDraft struct {
	ID        string
	Title     string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
