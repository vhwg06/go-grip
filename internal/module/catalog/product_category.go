package catalog

// ProductCategory maps a product to a category.
type ProductCategory struct {
	ProductID  string `json:"product_id"`
	CategoryID string `json:"category_id"`
}
