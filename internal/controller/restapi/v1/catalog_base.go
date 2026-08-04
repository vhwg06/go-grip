package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/usecase/catalogbase"
	"github.com/gofiber/fiber/v2"
)

func (r *V1) catalogBaseError(ctx *fiber.Ctx, err error) error {
	status, body := catalogbase.ErrorStatus(err)
	return ctx.Status(status).JSON(body)
}

func (r *V1) requireCatalogBase(ctx *fiber.Ctx) (*catalogbase.Service, error) {
	if r.catalogBase == nil {
		return nil, &catalogbase.APIError{Status: http.StatusInternalServerError, Code: "catalog_base_not_configured", Message: "Catalog Base service is not configured"}
	}
	return r.catalogBase, nil
}

func catalogBody(ctx *fiber.Ctx) (map[string]any, error) {
	body := strings.TrimSpace(string(ctx.Body()))
	if body == "" {
		return map[string]any{}, nil
	}
	input := map[string]any{}
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		return nil, &catalogbase.APIError{Status: http.StatusBadRequest, Code: "invalid_json", Message: "request body must be valid JSON"}
	}
	return input, nil
}

func (r *V1) catalogBaseListCategories(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.ListCategories(ctx.UserContext())
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func (r *V1) catalogBaseCreateCategory(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateCategory(ctx.UserContext(), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseUpdateCategory(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateCategory(ctx.UserContext(), ctx.Params("categoryId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDeactivateCategory(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeactivateCategory(ctx.UserContext(), ctx.Params("categoryId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDeleteCategory(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeleteCategory(ctx.UserContext(), ctx.Params("categoryId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseListDefinitions(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.ListDefinitions(ctx.UserContext())
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func (r *V1) catalogBaseCreateDefinition(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateDefinition(ctx.UserContext(), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseUpdateDefinition(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateDefinition(ctx.UserContext(), ctx.Params("definitionId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDeactivateDefinition(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeactivateDefinition(ctx.UserContext(), ctx.Params("definitionId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseAddEnumValue(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.AddEnumValue(ctx.UserContext(), ctx.Params("definitionId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseDeactivateEnumValue(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeactivateEnumValue(ctx.UserContext(), ctx.Params("definitionId"), ctx.Params("enumValueId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseListMasters(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.ListMasters(ctx.UserContext(), ctx.Params("masterKind"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func (r *V1) catalogBaseCreateMaster(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateMaster(ctx.UserContext(), ctx.Params("masterKind"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseUpdateMaster(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateMaster(ctx.UserContext(), ctx.Params("masterKind"), ctx.Params("masterId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDeactivateMaster(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeactivateMaster(ctx.UserContext(), ctx.Params("masterKind"), ctx.Params("masterId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseListModels(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.ListModels(ctx.UserContext())
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func (r *V1) catalogBaseCreateModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateModel(ctx.UserContext(), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseGetModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.GetModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseUpdateModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateModel(ctx.UserContext(), ctx.Params("modelId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDeleteModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeleteModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBasePublishModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.PublishModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseUnpublishModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UnpublishModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseDiscontinueModel(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DiscontinueModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseReplaceMedia(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.ReplaceMedia(ctx.UserContext(), ctx.Params("modelId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseCreateDimension(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateDimension(ctx.UserContext(), ctx.Params("modelId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseUpdateDimension(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateDimension(ctx.UserContext(), ctx.Params("modelId"), ctx.Params("dimensionId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseAddDimensionValue(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.AddDimensionValue(ctx.UserContext(), ctx.Params("modelId"), ctx.Params("dimensionId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseDeactivateDimensionValue(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.DeactivateDimensionValue(ctx.UserContext(), ctx.Params("modelId"), ctx.Params("dimensionId"), ctx.Params("valueId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseListVariants(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.ListVariants(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func (r *V1) catalogBaseCreateVariant(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.CreateVariant(ctx.UserContext(), ctx.Params("modelId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.Status(http.StatusCreated).JSON(item)
}

func (r *V1) catalogBaseGetVariant(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.GetVariant(ctx.UserContext(), ctx.Params("variantId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseUpdateVariant(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.UpdateVariant(ctx.UserContext(), ctx.Params("variantId"), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseActivateVariant(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.ActivateVariant(ctx.UserContext(), ctx.Params("variantId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseInactivateVariant(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	item, err := service.InactivateVariant(ctx.UserContext(), ctx.Params("variantId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(item)
}

func (r *V1) catalogBaseBulkPrice(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	items, err := service.BulkSetPrice(ctx.UserContext(), input)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(items)
}

func parsePublicPrice(ctx *fiber.Ctx, name string) (*int64, error) {
	value := strings.TrimSpace(ctx.Query(name))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return nil, &catalogbase.APIError{Status: http.StatusBadRequest, Code: "invalid_price_filter", Message: "price filters must be non-negative integers"}
	}
	return &parsed, nil
}

func (r *V1) catalogBasePublicList(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	minPrice, err := parsePublicPrice(ctx, "minPrice")
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	maxPrice, err := parsePublicPrice(ctx, "maxPrice")
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	result, err := service.ListPublicModels(ctx.UserContext(), catalogbase.PublicFilter{Page: ctx.QueryInt("page", 1), Limit: ctx.QueryInt("limit", 20), CategoryID: ctx.Query("categoryId"), MaterialID: ctx.Query("materialId"), FinishID: ctx.Query("finishId"), MinPrice: minPrice, MaxPrice: maxPrice, Search: ctx.Query("search"), Sort: ctx.Query("sort")})
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(result)
}

func (r *V1) catalogBasePublicDetail(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	result, err := service.GetPublicModel(ctx.UserContext(), ctx.Params("modelId"))
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(result)
}

func selectedOptionsQuery(ctx *fiber.Ctx) (map[string]string, error) {
	raw := strings.TrimSpace(ctx.Query("selected"))
	if raw == "" {
		return map[string]string{}, nil
	}
	selected := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &selected); err != nil {
		return nil, &catalogbase.APIError{Status: http.StatusBadRequest, Code: "invalid_selected_options", Message: "selected must be a JSON object"}
	}
	result := map[string]string{}
	for key, value := range selected {
		result[key] = strings.TrimSpace(strings.Trim(strings.TrimSpace(toString(value)), "\""))
	}
	return result, nil
}

func toString(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(stringify(value)), "\\\"", "\""))
}

func stringify(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

func (r *V1) catalogBasePublicOptions(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	selected, err := selectedOptionsQuery(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	result, err := service.AvailableOptions(ctx.UserContext(), ctx.Params("modelId"), selected)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(result)
}

func (r *V1) catalogBasePublicResolve(ctx *fiber.Ctx) error {
	service, err := r.requireCatalogBase(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	input, err := catalogBody(ctx)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	selectedValue, ok := input["selectedOptions"]
	if !ok {
		return r.catalogBaseError(ctx, &catalogbase.APIError{Status: http.StatusBadRequest, Code: "invalid_selected_options", Message: "selectedOptions is required"})
	}
	selected := map[string]string{}
	if object, ok := selectedValue.(map[string]any); ok {
		for key, value := range object {
			selected[key] = strings.TrimSpace(toString(value))
		}
	} else {
		return r.catalogBaseError(ctx, &catalogbase.APIError{Status: http.StatusBadRequest, Code: "invalid_selected_options", Message: "selectedOptions must be an object"})
	}
	result, err := service.ResolvePublicVariant(ctx.UserContext(), ctx.Params("modelId"), selected)
	if err != nil {
		return r.catalogBaseError(ctx, err)
	}
	return ctx.JSON(result)
}
