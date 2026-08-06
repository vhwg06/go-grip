package catalog

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
)

// toProductResponse maps catalogmodule.Product to openapi.ProductResponse.
func toProductResponse(p catalogmodule.Product) openapi.ProductResponse {
	priceInt := int(p.Price)
	stock := p.StockCount
	desc := p.Description
	sku := p.SKU
	catID := p.CategoryID
	img := p.ImageURL
	slug := p.Title // fallback slug

	return openapi.ProductResponse{
		Id:          p.ID,
		Title:       p.Title,
		Slug:        &slug,
		Description: &desc,
		Price:       priceInt,
		Stock:       &stock,
		Sku:         &sku,
		CategoryId:  &catID,
		ImageUrl:    &img,
		CreatedAt:   &p.CreatedAt,
		UpdatedAt:   &p.UpdatedAt,
	}
}

// toProductListResponse maps []catalogmodule.Product and total count to openapi.ProductListResponse.
func toProductListResponse(products []catalogmodule.Product, total int) openapi.ProductListResponse {
	items := make([]openapi.ProductResponse, len(products))
	for i, p := range products {
		items[i] = toProductResponse(p)
	}
	return openapi.ProductListResponse{
		Items: items,
		Total: total,
	}
}

// toCategoryResponse maps catalogmodule.Category to openapi.CategoryResponse.
func toCategoryResponse(c catalogmodule.Category) openapi.CategoryResponse {
	slug := c.Name
	return openapi.CategoryResponse{
		Id:   c.ID,
		Name: c.Name,
		Slug: &slug,
	}
}

// toCategoryListResponse maps []catalogmodule.Category to []openapi.CategoryResponse.
func toCategoryListResponse(categories []catalogmodule.Category) []openapi.CategoryResponse {
	res := make([]openapi.CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = toCategoryResponse(c)
	}
	return res
}

// toTagResponse maps catalogmodule.Tag to openapi.TagResponse.
func toTagResponse(t catalogmodule.Tag) openapi.TagResponse {
	slug := t.Name
	return openapi.TagResponse{
		Id:   t.ID,
		Name: t.Name,
		Slug: &slug,
	}
}

// toTagListResponse maps []catalogmodule.Tag to []openapi.TagResponse.
func toTagListResponse(tags []catalogmodule.Tag) []openapi.TagResponse {
	res := make([]openapi.TagResponse, len(tags))
	for i, t := range tags {
		res[i] = toTagResponse(t)
	}
	return res
}

// toPublicProductModelListResponse maps catalogbase use-case output map to openapi.PublicProductModelListResponse.
func toPublicProductModelListResponse(resMap map[string]any) openapi.PublicProductModelListResponse {
	items := make([]map[string]interface{}, 0)
	if rawItems, ok := resMap["items"].([]map[string]any); ok {
		items = rawItems
	}
	total := 0
	if rawTotal, ok := resMap["total"].(int); ok {
		total = rawTotal
	}
	page := 1
	if rawPage, ok := resMap["page"].(int); ok {
		page = rawPage
	}
	limit := 20
	if rawLimit, ok := resMap["limit"].(int); ok {
		limit = rawLimit
	}
	return openapi.PublicProductModelListResponse{
		Items: &items,
		Total: &total,
		Page:  &page,
		Limit: &limit,
	}
}
