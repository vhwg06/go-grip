package catalog

// Category defines a product taxonomy grouping.
type Category struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
	Position int     `json:"position"`
	IsActive bool    `json:"is_active"`
}
