package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Catalog capability.
type Handler struct {
	catalogUC catalogmodule.CatalogUseCase
	logger    logger.Interface
}

// NewHandler constructs a new Catalog vertical handler instance.
func NewHandler(catalogUC catalogmodule.CatalogUseCase, l logger.Interface) *Handler {
	return &Handler{
		catalogUC: catalogUC,
		logger:    l,
	}
}

// ListProducts handles GET /catalog/products
func (h *Handler) ListProducts(ctx context.Context, request openapi.ListProductsRequestObject) (openapi.ListProductsResponseObject, error) {
	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	filter := catalogmodule.ProductFilter{
		Pagination: pagination.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	}
	if request.Params.CategoryId != nil {
		filter.CategoryID = *request.Params.CategoryId
	}
	if request.Params.Search != nil {
		filter.Keyword = *request.Params.Search
	}

	products, total, err := h.catalogUC.ListProducts(ctx, filter)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		default:
			return openapi.ListProducts500JSONResponse{
				InternalErrorResponseJSONResponse: openapi.InternalErrorResponseJSONResponse(errResp),
			}, nil
		}
	}

	listDTO := toProductListResponse(products, total)
	return openapi.ListProducts200JSONResponse(listDTO), nil
}

// CreateProduct handles POST /catalog/products
func (h *Handler) CreateProduct(ctx context.Context, request openapi.CreateProductRequestObject) (openapi.CreateProductResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateProduct400JSONResponse{}, nil
	}

	p := catalogmodule.Product{
		Title: request.Body.Title,
		Price: int64(request.Body.Price),
	}
	if request.Body.Description != nil {
		p.Description = *request.Body.Description
	}
	if request.Body.Stock != nil {
		p.StockCount = *request.Body.Stock
	}
	if request.Body.Sku != nil {
		p.SKU = *request.Body.Sku
	}
	if request.Body.CategoryId != nil {
		p.CategoryID = *request.Body.CategoryId
	}
	if request.Body.ImageUrl != nil {
		p.ImageURL = *request.Body.ImageUrl
	}

	createdProduct, err := h.catalogUC.CreateProduct(ctx, p)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 400:
			return openapi.CreateProduct400JSONResponse{}, nil
		case 401:
			return openapi.CreateProduct401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.CreateProduct403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 409:
			return openapi.CreateProduct409JSONResponse{
				ConflictResponseJSONResponse: openapi.ConflictResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateProduct500JSONResponse{}, nil
		}
	}

	productDTO := toProductResponse(createdProduct)
	return openapi.CreateProduct201JSONResponse(productDTO), nil
}

// GetProductByID handles GET /catalog/products/{id}
func (h *Handler) GetProductByID(ctx context.Context, request openapi.GetProductByIDRequestObject) (openapi.GetProductByIDResponseObject, error) {
	product, err := h.catalogUC.GetProduct(ctx, request.Id)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 404:
			return openapi.GetProductByID404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetProductByID500JSONResponse{}, nil
		}
	}

	productDTO := toProductResponse(product)
	return openapi.GetProductByID200JSONResponse(productDTO), nil
}

// UpdateProduct handles PUT /catalog/products/{id}
func (h *Handler) UpdateProduct(ctx context.Context, request openapi.UpdateProductRequestObject) (openapi.UpdateProductResponseObject, error) {
	if request.Body == nil {
		return openapi.UpdateProduct400JSONResponse{}, nil
	}

	p := catalogmodule.Product{
		ID:    request.Id,
		Title: request.Body.Title,
		Price: int64(request.Body.Price),
	}
	if request.Body.Description != nil {
		p.Description = *request.Body.Description
	}
	if request.Body.Stock != nil {
		p.StockCount = *request.Body.Stock
	}
	if request.Body.Sku != nil {
		p.SKU = *request.Body.Sku
	}
	if request.Body.CategoryId != nil {
		p.CategoryID = *request.Body.CategoryId
	}
	if request.Body.ImageUrl != nil {
		p.ImageURL = *request.Body.ImageUrl
	}

	updatedProduct, err := h.catalogUC.UpdateProduct(ctx, p)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 400:
			return openapi.UpdateProduct400JSONResponse{}, nil
		case 401:
			return openapi.UpdateProduct401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.UpdateProduct403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.UpdateProduct404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UpdateProduct500JSONResponse{}, nil
		}
	}

	productDTO := toProductResponse(updatedProduct)
	return openapi.UpdateProduct200JSONResponse(productDTO), nil
}

// DeleteProduct handles DELETE /catalog/products/{id}
func (h *Handler) DeleteProduct(ctx context.Context, request openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	err := h.catalogUC.DeleteProduct(ctx, request.Id)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 401:
			return openapi.DeleteProduct401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.DeleteProduct403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.DeleteProduct404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.DeleteProduct500JSONResponse{}, nil
		}
	}

	return openapi.DeleteProduct204Response{}, nil
}

// ListCategories handles GET /catalog/categories
func (h *Handler) ListCategories(ctx context.Context, request openapi.ListCategoriesRequestObject) (openapi.ListCategoriesResponseObject, error) {
	categories, err := h.catalogUC.ListCategories(ctx)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		default:
			return openapi.ListCategories500JSONResponse{
				InternalErrorResponseJSONResponse: openapi.InternalErrorResponseJSONResponse(errResp),
			}, nil
		}
	}

	catDTO := toCategoryListResponse(categories)
	return openapi.ListCategories200JSONResponse(catDTO), nil
}

// CreateCategory handles POST /catalog/categories
func (h *Handler) CreateCategory(ctx context.Context, request openapi.CreateCategoryRequestObject) (openapi.CreateCategoryResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateCategory400JSONResponse{}, nil
	}

	cat := catalogmodule.Category{
		Name: request.Body.Name,
	}

	createdCat, err := h.catalogUC.CreateCategory(ctx, cat)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 400:
			return openapi.CreateCategory400JSONResponse{}, nil
		case 401:
			return openapi.CreateCategory401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.CreateCategory403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateCategory500JSONResponse{}, nil
		}
	}

	catDTO := toCategoryResponse(createdCat)
	return openapi.CreateCategory201JSONResponse(catDTO), nil
}

// ListTags handles GET /catalog/tags
func (h *Handler) ListTags(ctx context.Context, request openapi.ListTagsRequestObject) (openapi.ListTagsResponseObject, error) {
	tags, err := h.catalogUC.ListTags(ctx)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		default:
			return openapi.ListTags500JSONResponse{
				InternalErrorResponseJSONResponse: openapi.InternalErrorResponseJSONResponse(errResp),
			}, nil
		}
	}

	tagDTO := toTagListResponse(tags)
	return openapi.ListTags200JSONResponse(tagDTO), nil
}

// CreateTag handles POST /catalog/tags
func (h *Handler) CreateTag(ctx context.Context, request openapi.CreateTagRequestObject) (openapi.CreateTagResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateTag400JSONResponse{}, nil
	}

	t := catalogmodule.Tag{
		Name: request.Body.Name,
	}

	createdTag, err := h.catalogUC.CreateTag(ctx, t)
	if err != nil {
		status, errResp := mapCatalogError(err)
		switch status {
		case 400:
			return openapi.CreateTag400JSONResponse{}, nil
		case 401:
			return openapi.CreateTag401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.CreateTag403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateTag500JSONResponse{}, nil
		}
	}

	tagDTO := toTagResponse(createdTag)
	return openapi.CreateTag201JSONResponse(tagDTO), nil
}
