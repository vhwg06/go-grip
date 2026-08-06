package catalog

// ProductTag maps a product to a tag.
type ProductTag struct {
	ProductID string `json:"product_id"`
	TagID     string `json:"tag_id"`
}
