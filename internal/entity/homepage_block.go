package entity

type HomepageBlock struct {
	ID        string         `json:"id"`
	BlockType string         `json:"block_type"`
	Config    map[string]any `json:"config"`
	Position  int            `json:"position"`
	IsActive  bool           `json:"is_active"`
}
