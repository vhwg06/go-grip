package entity

type ImportItemType string

const (
	ImportItemProduct ImportItemType = "product"
	ImportItemPost    ImportItemType = "post"
)

type ImportItem struct {
	Type ImportItemType `json:"type"`
	Data map[string]any `json:"data"`
}

type ImportFailure struct {
	Index  int    `json:"index"`
	Reason string `json:"reason"`
}

type ImportResult struct {
	Imported int             `json:"imported"`
	Failed   []ImportFailure `json:"failed"`
}
