// Package catalogbase contains the Catalog Base application service.
//
// The service owns the catalog invariants and only talks to the repository
// port.  GORM rows and SQL transactions live in internal/repo/persistent.
package catalogbase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string { return e.Message }

func bad(message string) error {
	return &APIError{Status: http.StatusBadRequest, Code: "invalid_catalog_command", Message: message}
}

func conflict(message string) error {
	return &APIError{Status: http.StatusConflict, Code: "catalog_conflict", Message: message}
}

func notFound(message string) error {
	return &APIError{Status: http.StatusNotFound, Code: "catalog_not_found", Message: message}
}

func methodNotAllowed(message string) error {
	return &APIError{Status: http.StatusMethodNotAllowed, Code: "catalog_method_not_allowed", Message: message}
}

func ErrorStatus(err error) (int, map[string]any) {
	if err == nil {
		return http.StatusOK, nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status, map[string]any{"code": apiErr.Code, "message": apiErr.Message}
	}
	return http.StatusInternalServerError, map[string]any{"code": "catalog_internal_error", "message": "catalog operation failed"}
}

type Service struct {
	repository repo.CatalogBaseRepository
}

func New(repository repo.CatalogBaseRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) load(ctx context.Context) (entity.CatalogSnapshot, error) {
	if s == nil || s.repository == nil {
		return entity.CatalogSnapshot{}, errors.New("catalog base repository is not configured")
	}
	return s.repository.LoadCatalogBase(ctx)
}

func (s *Service) save(ctx context.Context, snapshot entity.CatalogSnapshot) error {
	if s == nil || s.repository == nil {
		return errors.New("catalog base repository is not configured")
	}
	return s.repository.SaveCatalogBase(ctx, snapshot)
}

func findCategory(snapshot *entity.CatalogSnapshot, id string) (*entity.CatalogCategory, error) {
	for index := range snapshot.Categories {
		if snapshot.Categories[index].ID == id {
			return &snapshot.Categories[index], nil
		}
	}
	return nil, notFound("category not found")
}

func findDefinition(snapshot *entity.CatalogSnapshot, id string) (*entity.CatalogAttributeDefinition, error) {
	for index := range snapshot.Definitions {
		if snapshot.Definitions[index].ID == id {
			return &snapshot.Definitions[index], nil
		}
	}
	return nil, notFound("attribute definition not found")
}

func findMaster(snapshot *entity.CatalogSnapshot, kind, id string) (*entity.CatalogMaster, error) {
	for index := range snapshot.Masters {
		if snapshot.Masters[index].Kind == kind && snapshot.Masters[index].ID == id {
			return &snapshot.Masters[index], nil
		}
	}
	return nil, notFound("catalog master not found")
}

func findModel(snapshot *entity.CatalogSnapshot, id string) (*entity.CatalogProductModel, error) {
	for index := range snapshot.Models {
		if snapshot.Models[index].ID == id {
			return &snapshot.Models[index], nil
		}
	}
	return nil, notFound("ProductModel not found")
}

func findVariant(snapshot *entity.CatalogSnapshot, id string) (*entity.CatalogVariant, *entity.CatalogProductModel, error) {
	for modelIndex := range snapshot.Models {
		for variantIndex := range snapshot.Models[modelIndex].Variants {
			if snapshot.Models[modelIndex].Variants[variantIndex].ID == id {
				return &snapshot.Models[modelIndex].Variants[variantIndex], &snapshot.Models[modelIndex], nil
			}
		}
	}
	return nil, nil, notFound("Variant not found")
}

func mapValue(input map[string]any, names ...string) (any, bool) {
	for _, name := range names {
		if value, ok := input[name]; ok {
			return value, true
		}
	}
	return nil, false
}

func stringValue(input map[string]any, names ...string) string {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func boolValue(input map[string]any, fallback bool, names ...string) bool {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return fallback
	}
	if parsed, ok := value.(bool); ok {
		return parsed
	}
	return strings.EqualFold(strings.TrimSpace(fmt.Sprint(value)), "true")
}

func intValue(input map[string]any, fallback int, names ...string) int {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return fallback
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func floatValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func recordMap(value any) map[string]any {
	if typed, ok := value.(map[string]any); ok && typed != nil {
		return typed
	}
	return map[string]any{}
}

func parseStringMap(value any) map[string]string {
	result := map[string]string{}
	switch typed := value.(type) {
	case map[string]string:
		return typed
	case map[string]any:
		for key, raw := range typed {
			if raw != nil {
				result[key] = fmt.Sprint(raw)
			}
		}
	}
	return result
}

func jsonMap(value map[string]any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return value
}

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F-]{20,}$`)

func now() time.Time { return time.Now().UTC() }

func normalizeText(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(value))), " ")
}

func canonicalUnit(value string) (float64, string, bool) {
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(value)))
	if len(parts) != 2 {
		return 0, "", false
	}
	number, err := strconv.ParseFloat(parts[0], 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, "", false
	}
	type unit struct {
		factor float64
		family string
	}
	units := map[string]unit{"mm": {1, "length"}, "cm": {10, "length"}, "m": {1000, "length"}, "g": {1, "mass"}, "kg": {1000, "mass"}}
	valueUnit, ok := units[parts[1]]
	if !ok {
		return 0, "", false
	}
	return number * valueUnit.factor, valueUnit.family, true
}

func canonicalValue(value string) string {
	if number, family, ok := canonicalUnit(value); ok {
		return fmt.Sprintf("%s:%0.6f", family, number)
	}
	return normalizeText(value)
}

func compatibleUnit(family, unit string) bool {
	if family == "" && unit == "" {
		return true
	}
	unit = strings.ToLower(strings.TrimSpace(unit))
	switch strings.ToLower(strings.TrimSpace(family)) {
	case "length":
		return unit == "mm" || unit == "cm" || unit == "m"
	case "mass", "weight":
		return unit == "g" || unit == "kg"
	case "area":
		return unit == "mm2" || unit == "cm2" || unit == "m2"
	case "volume":
		return unit == "ml" || unit == "l"
	default:
		return false
	}
}

func validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit string) error {
	valueKind, dataType = strings.TrimSpace(valueKind), strings.TrimSpace(dataType)
	referenceTarget = strings.TrimSpace(referenceTarget)
	switch valueKind {
	case "Scalar":
		if dataType == "" || (dataType != "Text" && dataType != "Number" && dataType != "Boolean") {
			return bad("Scalar requires a supported dataType")
		}
		if dataType == "Number" && !compatibleUnit(unitFamily, unit) {
			return bad("numeric Scalar has an incompatible unit")
		}
		if dataType != "Number" && (unitFamily != "" || unit != "") {
			return bad("only numeric Scalar definitions may declare a unit")
		}
	case "Enum":
		if dataType != "" || referenceTarget != "" || unitFamily != "" || unit != "" {
			return bad("Enum definitions cannot carry Scalar or Reference semantics")
		}
	case "Reference":
		if referenceTarget != "Material" && referenceTarget != "Finish" && referenceTarget != "Pack" {
			return bad("Reference requires one supported referenceTarget")
		}
		if dataType != "" || unitFamily != "" || unit != "" {
			return bad("Reference definitions cannot carry Scalar semantics")
		}
	default:
		return bad("unsupported attribute valueKind")
	}
	return nil
}

func warrantyValue(input map[string]any, names ...string) (map[string]any, error) {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return nil, nil
	}
	object, ok := value.(map[string]any)
	if !ok || stringValue(object, "term") == "" {
		return nil, bad("warrantySummary.term is required")
	}
	return map[string]any{"term": stringValue(object, "term"), "note": stringValue(object, "note")}, nil
}

func validateMeasurements(value any) error {
	for name, raw := range recordMap(value) {
		measurement := recordMap(raw)
		unit := stringValue(measurement, "unit")
		if unit == "" {
			continue
		}
		family := ""
		if strings.Contains(strings.ToLower(name), "length") {
			family = "length"
		}
		if family != "" && !compatibleUnit(family, unit) {
			return bad("length measurement has an incompatible unit")
		}
		if rawValue, ok := mapValue(measurement, "value"); ok {
			number, valid := floatValue(rawValue)
			if !valid || math.IsNaN(number) || math.IsInf(number, 0) {
				return bad("measurement value must be numeric")
			}
			if canonicalNumber, canonicalFamily, canonicalUnitName, canonical := canonicalMeasurement(number, unit); canonical {
				if family != "" && !strings.EqualFold(family, canonicalFamily) {
					return bad("measurement has an incompatible unit family")
				}
				measurement["value"], measurement["unit"] = canonicalNumber, canonicalUnitName
			}
		}
	}
	return nil
}

func canonicalMeasurement(number float64, unit string) (float64, string, string, bool) {
	value, family, valid := canonicalUnit(fmt.Sprintf("%v %s", number, unit))
	if !valid {
		return 0, "", "", false
	}
	baseUnit := map[string]string{"length": "mm", "mass": "g", "area": "mm2", "volume": "ml"}[family]
	if baseUnit == "" {
		return 0, "", "", false
	}
	return value, family, baseUnit, true
}

func categoryOutput(value entity.CatalogCategory) map[string]any {
	return map[string]any{"id": value.ID, "name": value.Name, "slug": value.Slug, "parentId": value.ParentID, "position": value.Position, "active": value.Active}
}

func definitionOutput(value entity.CatalogAttributeDefinition) map[string]any {
	result := map[string]any{"id": value.ID, "key": value.Key, "displayName": value.DisplayName, "description": value.Description, "ordering": value.Ordering, "valueKind": value.ValueKind, "active": value.Active, "enumValues": value.EnumValues}
	if value.DataType != "" {
		result["dataType"] = value.DataType
	}
	if value.ReferenceTarget != "" {
		result["referenceTarget"] = value.ReferenceTarget
	}
	if value.UnitFamily != "" {
		result["unitFamily"] = value.UnitFamily
	}
	if value.Unit != "" {
		result["unit"] = value.Unit
	}
	return result
}

func masterOutput(value entity.CatalogMaster) map[string]any {
	result := map[string]any{"id": value.ID, "kind": value.Kind, "name": value.Name, "description": value.Description, "swatchMedia": value.SwatchMedia, "sellingUnit": value.SellingUnit, "baseUnit": value.BaseUnit, "active": value.Active}
	if value.Quantity != nil {
		result["quantity"] = *value.Quantity
	}
	return result
}

func dimensionOutput(snapshot *entity.CatalogSnapshot, modelID string, dimension entity.CatalogVariantDimension) (map[string]any, error) {
	definition, err := findDefinition(snapshot, dimension.DefinitionID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"id": dimension.ID, "definitionId": dimension.DefinitionID, "key": definition.DisplayName, "allowedValues": dimension.AllowedValues, "modelId": modelID}, nil
}

func variantOutput(value entity.CatalogVariant) map[string]any {
	selectedOptions := make(map[string]any, len(value.SelectedOptions))
	for key, option := range value.SelectedOptions {
		selectedOptions[key] = option
	}
	result := map[string]any{"id": value.ID, "selectedOptions": selectedOptions, "technicalValues": jsonMap(value.TechnicalValues), "sku": value.SKU, "status": value.Status, "saleReady": value.SaleReady(), "canonicalCombination": value.CanonicalCombination, "history": value.History, "packId": nil, "sellingPrice": nil}
	if value.PackID != nil {
		result["packId"] = *value.PackID
	}
	if value.SellingPrice != nil {
		result["sellingPrice"] = map[string]any{"amount": value.SellingPrice.Amount, "currency": value.SellingPrice.Currency}
	}
	return result
}

func modelOutput(snapshot *entity.CatalogSnapshot, value entity.CatalogProductModel, public bool) (map[string]any, error) {
	result := map[string]any{"id": value.ID, "name": value.Name, "categoryId": value.CategoryID, "description": value.Description, "fixedAttributes": jsonMap(value.FixedAttributes), "measurements": jsonMap(value.Measurements), "status": value.Status, "images": []map[string]any{}, "variants": []map[string]any{}, "variantDimensions": []map[string]any{}, "warrantySummary": value.WarrantySummary}
	if value.FixedPackID != nil {
		result["fixedPackId"] = *value.FixedPackID
	}
	images := make([]map[string]any, 0, len(value.Images))
	for _, image := range value.Images {
		images = append(images, map[string]any{"id": image.ID, "url": image.URL, "ordering": image.Ordering, "primary": image.PrimaryImage})
	}
	result["images"] = images
	dimensions := make([]map[string]any, 0, len(value.Dimensions))
	for _, dimension := range value.Dimensions {
		output, err := dimensionOutput(snapshot, value.ID, dimension)
		if err != nil {
			return nil, err
		}
		dimensions = append(dimensions, output)
	}
	result["variantDimensions"] = dimensions
	variants := make([]map[string]any, 0, len(value.Variants))
	for _, variant := range value.Variants {
		if !public || variant.SaleReady() {
			variants = append(variants, variantOutput(variant))
		}
	}
	result["variants"] = variants
	return result, nil
}

func (s *Service) activeMaster(snapshot *entity.CatalogSnapshot, kind, id string) error {
	master, err := findMaster(snapshot, kind, id)
	if err != nil {
		return err
	}
	if !master.Active {
		return conflict("inactive catalog master cannot be assigned")
	}
	return nil
}

func (s *Service) validateFixedAttributes(snapshot *entity.CatalogSnapshot, input map[string]any) error {
	fixed := recordMapValue(input, "fixedAttributes", "fixed_attributes")
	for key, rawValue := range fixed {
		definition, definitionErr := fixedAttributeDefinition(snapshot, key)
		if definitionErr != nil && uuidPattern.MatchString(key) {
			return definitionErr
		}
		if definitionErr == nil {
			if !definition.Active {
				return conflict("inactive attribute definition cannot be assigned")
			}
			switch definition.ValueKind {
			case "Reference":
				master, err := s.findActiveMasterByValue(snapshot, strings.ToLower(definition.ReferenceTarget), fmt.Sprint(rawValue))
				if err != nil {
					return err
				}
				fixed[key] = master.ID
			case "Enum":
				if !enumValueMatches(definition.EnumValues, fmt.Sprint(rawValue)) {
					return conflict("fixed attribute must use an active enum value")
				}
			case "Scalar":
				if err := validateScalarValue(definition, rawValue); err != nil {
					return err
				}
			}
			continue
		}
		if strings.EqualFold(key, "materialId") || strings.EqualFold(key, "finishId") || strings.EqualFold(key, "packId") {
			kind := strings.ToLower(strings.TrimSuffix(key, "Id"))
			if err := s.activeMaster(snapshot, kind, fmt.Sprint(rawValue)); err != nil {
				var apiErr *APIError
				if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound {
					continue
				}
				return err
			}
		}
	}
	return nil
}

func fixedAttributeDefinition(snapshot *entity.CatalogSnapshot, key string) (*entity.CatalogAttributeDefinition, error) {
	if uuidPattern.MatchString(key) {
		return findDefinition(snapshot, key)
	}
	for index := range snapshot.Definitions {
		definition := &snapshot.Definitions[index]
		if strings.EqualFold(definition.Key, key) || strings.EqualFold(definition.DisplayName, key) {
			return definition, nil
		}
	}
	return nil, notFound("attribute definition not found")
}

func enumValueMatches(values []entity.CatalogEnumValue, raw string) bool {
	for _, value := range values {
		if value.Active && (value.ID == raw || strings.EqualFold(value.Key, raw) || strings.EqualFold(normalizeText(value.Label), normalizeText(raw))) {
			return true
		}
	}
	return false
}

func validateScalarValue(definition *entity.CatalogAttributeDefinition, raw any) error {
	switch definition.DataType {
	case "Text":
		if strings.TrimSpace(fmt.Sprint(raw)) == "" {
			return bad("text fixed attribute cannot be empty")
		}
	case "Boolean":
		if _, ok := raw.(bool); !ok && !strings.EqualFold(strings.TrimSpace(fmt.Sprint(raw)), "true") && !strings.EqualFold(strings.TrimSpace(fmt.Sprint(raw)), "false") {
			return bad("boolean fixed attribute must be true or false")
		}
	case "Number":
		if number, valid := floatValue(raw); valid && !math.IsNaN(number) && !math.IsInf(number, 0) {
			return nil
		}
		if object := recordMap(raw); len(object) > 0 {
			if number, valid := floatValue(object["value"]); valid && !math.IsNaN(number) && !math.IsInf(number, 0) {
				unit := stringValue(object, "unit")
				if unit == "" || compatibleUnit(definition.UnitFamily, unit) {
					return nil
				}
			}
		}
		if _, family, valid := canonicalUnit(fmt.Sprint(raw)); !valid || (definition.UnitFamily != "" && !strings.EqualFold(family, definition.UnitFamily)) {
			return bad("numeric fixed attribute has an incompatible value")
		}
	default:
		return bad("unsupported Scalar data type")
	}
	return nil
}

func (s *Service) findActiveMasterByValue(snapshot *entity.CatalogSnapshot, kind, value string) (*entity.CatalogMaster, error) {
	for index := range snapshot.Masters {
		master := &snapshot.Masters[index]
		if master.Kind != kind {
			continue
		}
		if master.ID == value || strings.EqualFold(normalizeText(master.Name), normalizeText(value)) {
			if !master.Active {
				return nil, conflict("inactive catalog master cannot be assigned")
			}
			return master, nil
		}
	}
	return nil, notFound("catalog master not found")
}

func recordMapValue(input map[string]any, names ...string) map[string]any {
	value, _ := mapValue(input, names...)
	return recordMap(value)
}

func (s *Service) validateModelInput(snapshot *entity.CatalogSnapshot, input map[string]any, existing *entity.CatalogProductModel) error {
	categoryID := stringValue(input, "categoryId", "category_id")
	if categoryID != "" {
		category, err := findCategory(snapshot, categoryID)
		if err != nil {
			return err
		}
		if !category.Active && (existing == nil || existing.CategoryID != categoryID) {
			return conflict("inactive category cannot be assigned to a new ProductModel")
		}
	}
	if err := s.validateFixedAttributes(snapshot, input); err != nil {
		return err
	}
	if fixedPackID := stringValue(input, "fixedPackId", "fixed_pack_id"); fixedPackID != "" {
		if err := s.activeMaster(snapshot, "pack", fixedPackID); err != nil {
			return err
		}
	}
	if value, ok := mapValue(input, "measurements"); ok {
		if err := validateMeasurements(value); err != nil {
			return err
		}
	}
	_, err := warrantyValue(input, "warrantySummary", "warranty")
	return err
}

// Categories -----------------------------------------------------------------

func (s *Service) ListCategories(ctx context.Context) ([]map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(snapshot.Categories))
	for _, value := range snapshot.Categories {
		result = append(result, categoryOutput(value))
	}
	return result, nil
}

func (s *Service) CreateCategory(ctx context.Context, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	name, slug := stringValue(input, "name"), stringValue(input, "slug")
	if name == "" || slug == "" {
		return nil, bad("category name and slug are required")
	}
	for _, value := range snapshot.Categories {
		if value.Slug == slug {
			return nil, conflict("category slug already exists")
		}
	}
	var parentID *string
	if value := stringValue(input, "parentId", "parent_id"); value != "" {
		if _, err := findCategory(&snapshot, value); err != nil {
			return nil, err
		}
		parentID = &value
	}
	value := entity.CatalogCategory{ID: uuid.NewString(), Name: name, Slug: slug, ParentID: parentID, Position: intValue(input, 0, "position"), Active: boolValue(input, true, "active"), CreatedAt: now(), UpdatedAt: now()}
	snapshot.Categories = append(snapshot.Categories, value)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return categoryOutput(value), nil
}

func (s *Service) UpdateCategory(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	if slug := stringValue(input, "slug"); slug != "" {
		for _, candidate := range snapshot.Categories {
			if candidate.ID != id && candidate.Slug == slug {
				return nil, conflict("category slug already exists")
			}
		}
	}
	value, err := findCategory(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if text := stringValue(input, "name"); text != "" {
		value.Name = text
	}
	if text := stringValue(input, "slug"); text != "" {
		value.Slug = text
	}
	if _, ok := mapValue(input, "position"); ok {
		value.Position = intValue(input, value.Position, "position")
	}
	if _, ok := mapValue(input, "active"); ok {
		value.Active = boolValue(input, value.Active, "active")
	}
	if parentID, ok := mapValue(input, "parentId", "parent_id"); ok {
		if parentID == nil || strings.TrimSpace(fmt.Sprint(parentID)) == "" {
			value.ParentID = nil
		} else {
			parent := strings.TrimSpace(fmt.Sprint(parentID))
			if parent == id {
				return nil, bad("category cannot be its own parent")
			}
			if _, parentErr := findCategory(&snapshot, parent); parentErr != nil {
				return nil, parentErr
			}
			for ancestor := parent; ancestor != ""; {
				candidate, ancestorErr := findCategory(&snapshot, ancestor)
				if ancestorErr != nil || candidate.ParentID == nil {
					break
				}
				if *candidate.ParentID == id {
					return nil, bad("category hierarchy cannot contain a cycle")
				}
				ancestor = *candidate.ParentID
			}
			value.ParentID = &parent
		}
	}
	value.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return categoryOutput(*value), nil
}

func (s *Service) DeactivateCategory(ctx context.Context, id string) (map[string]any, error) {
	return s.UpdateCategory(ctx, id, map[string]any{"active": false})
}

func (s *Service) DeleteCategory(context.Context, string) (map[string]any, error) {
	return nil, methodNotAllowed("category deletion is not supported; deactivate it instead")
}

// Attribute definitions ------------------------------------------------------

func (s *Service) ListDefinitions(ctx context.Context) ([]map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(snapshot.Definitions))
	for _, value := range snapshot.Definitions {
		result = append(result, definitionOutput(value))
	}
	return result, nil
}

func (s *Service) CreateDefinition(ctx context.Context, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	key, valueKind := stringValue(input, "key"), stringValue(input, "valueKind", "value_kind")
	if key == "" || valueKind == "" {
		return nil, bad("attribute key and valueKind are required")
	}
	dataType, referenceTarget := stringValue(input, "dataType", "data_type"), stringValue(input, "referenceTarget", "reference_target")
	unitFamily, unit := stringValue(input, "unitFamily", "unit_family"), stringValue(input, "unit")
	if err := validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit); err != nil {
		return nil, err
	}
	for _, value := range snapshot.Definitions {
		if value.Key == key {
			return nil, conflict("attribute key already exists")
		}
	}
	value := entity.CatalogAttributeDefinition{ID: uuid.NewString(), Key: key, DisplayName: stringValue(input, "displayName", "display_name"), Description: stringValue(input, "description"), Ordering: intValue(input, 0, "ordering"), ValueKind: valueKind, DataType: dataType, ReferenceTarget: referenceTarget, UnitFamily: unitFamily, Unit: unit, Active: boolValue(input, true, "active"), EnumValues: []entity.CatalogEnumValue{}, CreatedAt: now(), UpdatedAt: now()}
	if value.DisplayName == "" {
		value.DisplayName = key
	}
	snapshot.Definitions = append(snapshot.Definitions, value)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return definitionOutput(value), nil
}

func definitionUsed(snapshot entity.CatalogSnapshot, id string) bool {
	for _, model := range snapshot.Models {
		if _, ok := model.FixedAttributes[id]; ok {
			return true
		}
		for _, dimension := range model.Dimensions {
			if dimension.DefinitionID == id {
				return true
			}
		}
	}
	return false
}

func (s *Service) UpdateDefinition(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findDefinition(&snapshot, id)
	if err != nil {
		return nil, err
	}
	valueKind, dataType, referenceTarget := value.ValueKind, value.DataType, value.ReferenceTarget
	unitFamily, unit := value.UnitFamily, value.Unit
	if raw, ok := mapValue(input, "valueKind", "value_kind"); ok {
		valueKind = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "dataType", "data_type"); ok {
		dataType = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "referenceTarget", "reference_target"); ok {
		referenceTarget = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "unitFamily", "unit_family"); ok {
		unitFamily = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "unit"); ok {
		unit = fmt.Sprint(raw)
	}
	if definitionUsed(snapshot, id) && (valueKind != value.ValueKind || dataType != value.DataType || referenceTarget != value.ReferenceTarget || unitFamily != value.UnitFamily || unit != value.Unit) {
		return nil, conflict("used attribute semantic structure cannot change")
	}
	if err := validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit); err != nil {
		return nil, err
	}
	value.ValueKind, value.DataType, value.ReferenceTarget, value.UnitFamily, value.Unit = valueKind, dataType, referenceTarget, unitFamily, unit
	if text := stringValue(input, "displayName", "display_name"); text != "" {
		value.DisplayName = text
	}
	if raw, ok := mapValue(input, "description"); ok {
		value.Description = fmt.Sprint(raw)
	}
	if _, ok := mapValue(input, "ordering"); ok {
		value.Ordering = intValue(input, value.Ordering, "ordering")
	}
	value.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return definitionOutput(*value), nil
}

func (s *Service) DeactivateDefinition(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findDefinition(&snapshot, id)
	if err != nil {
		return nil, err
	}
	value.Active, value.UpdatedAt = false, now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return definitionOutput(*value), nil
}

func (s *Service) AddEnumValue(ctx context.Context, definitionID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	definition, err := findDefinition(&snapshot, definitionID)
	if err != nil {
		return nil, err
	}
	if definition.ValueKind != "Enum" {
		return nil, bad("enum values require an Enum definition")
	}
	key, label := stringValue(input, "key"), stringValue(input, "label")
	if key == "" || label == "" {
		return nil, bad("enum key and label are required")
	}
	for _, value := range definition.EnumValues {
		if strings.EqualFold(value.Key, key) {
			return nil, conflict("enum key already exists")
		}
	}
	item := entity.CatalogEnumValue{ID: uuid.NewString(), Key: key, Label: label, Active: boolValue(input, true, "active")}
	definition.EnumValues = append(definition.EnumValues, item)
	definition.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return map[string]any{"id": item.ID, "key": item.Key, "label": item.Label, "active": item.Active}, nil
}

func (s *Service) DeactivateEnumValue(ctx context.Context, definitionID, enumID string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	definition, err := findDefinition(&snapshot, definitionID)
	if err != nil {
		return nil, err
	}
	for index := range definition.EnumValues {
		if definition.EnumValues[index].ID == enumID {
			definition.EnumValues[index].Active = false
			for modelIndex := range snapshot.Models {
				for dimensionIndex := range snapshot.Models[modelIndex].Dimensions {
					dimension := &snapshot.Models[modelIndex].Dimensions[dimensionIndex]
					if dimension.DefinitionID != definitionID {
						continue
					}
					for valueIndex := range dimension.AllowedValues {
						if dimension.AllowedValues[valueIndex].ID == enumID {
							dimension.AllowedValues[valueIndex].Active = false
						}
					}
				}
			}
			definition.UpdatedAt = now()
			if err := s.save(ctx, snapshot); err != nil {
				return nil, err
			}
			return map[string]any{"id": enumID, "key": definition.EnumValues[index].Key, "label": definition.EnumValues[index].Label, "active": false}, nil
		}
	}
	return nil, notFound("enum value not found")
}

// Masters --------------------------------------------------------------------

func validMasterKind(kind string) bool {
	return kind == "material" || kind == "finish" || kind == "pack"
}

func (s *Service) ListMasters(ctx context.Context, kind string) ([]map[string]any, error) {
	if !validMasterKind(kind) {
		return nil, bad("unsupported catalog master kind")
	}
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0)
	for _, value := range snapshot.Masters {
		if value.Kind == kind {
			result = append(result, masterOutput(value))
		}
	}
	sort.Slice(result, func(i, j int) bool { return fmt.Sprint(result[i]["name"]) < fmt.Sprint(result[j]["name"]) })
	return result, nil
}

func (s *Service) CreateMaster(ctx context.Context, kind string, input map[string]any) (map[string]any, error) {
	if !validMasterKind(kind) {
		return nil, bad("unsupported catalog master kind")
	}
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	name := stringValue(input, "name")
	if name == "" {
		return nil, bad("master name is required")
	}
	for _, value := range snapshot.Masters {
		if value.Kind == kind && value.Name == name {
			return nil, conflict("catalog master already exists")
		}
	}
	var quantity *float64
	if raw, ok := mapValue(input, "quantity"); ok && raw != nil {
		parsed, valid := floatValue(raw)
		if !valid || parsed <= 0 {
			return nil, bad("pack quantity must be positive")
		}
		quantity = &parsed
	}
	media := []string{}
	if raw, ok := mapValue(input, "swatchMedia", "swatch_media"); ok {
		for _, item := range recordSlice(raw) {
			media = append(media, fmt.Sprint(item))
		}
	}
	value := entity.CatalogMaster{ID: uuid.NewString(), Kind: kind, Name: name, Description: stringValue(input, "description"), SwatchMedia: media, SellingUnit: stringValue(input, "sellingUnit", "selling_unit"), Quantity: quantity, BaseUnit: stringValue(input, "baseUnit", "base_unit"), Active: boolValue(input, true, "active"), CreatedAt: now(), UpdatedAt: now()}
	snapshot.Masters = append(snapshot.Masters, value)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return masterOutput(value), nil
}

func recordSlice(value any) []any {
	switch values := value.(type) {
	case []any:
		return values
	case []string:
		result := make([]any, len(values))
		for index, item := range values {
			result[index] = item
		}
		return result
	}
	return nil
}

func (s *Service) UpdateMaster(ctx context.Context, kind, id string, input map[string]any) (map[string]any, error) {
	if !validMasterKind(kind) {
		return nil, bad("unsupported catalog master kind")
	}
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findMaster(&snapshot, kind, id)
	if err != nil {
		return nil, err
	}
	if text := stringValue(input, "name"); text != "" {
		value.Name = text
	}
	if raw, ok := mapValue(input, "description"); ok {
		value.Description = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "swatchMedia", "swatch_media"); ok {
		value.SwatchMedia = make([]string, 0)
		for _, item := range recordSlice(raw) {
			value.SwatchMedia = append(value.SwatchMedia, fmt.Sprint(item))
		}
	}
	if raw, ok := mapValue(input, "sellingUnit", "selling_unit"); ok {
		value.SellingUnit = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "baseUnit", "base_unit"); ok {
		value.BaseUnit = fmt.Sprint(raw)
	}
	if raw, ok := mapValue(input, "quantity"); ok {
		parsed, valid := floatValue(raw)
		if !valid || parsed <= 0 {
			return nil, bad("pack quantity must be positive")
		}
		value.Quantity = &parsed
	}
	if _, ok := mapValue(input, "active"); ok {
		value.Active = boolValue(input, value.Active, "active")
	}
	value.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return masterOutput(*value), nil
}

func (s *Service) DeactivateMaster(ctx context.Context, kind, id string) (map[string]any, error) {
	return s.UpdateMaster(ctx, kind, id, map[string]any{"active": false})
}

// ProductModel ---------------------------------------------------------------

func (s *Service) CreateModel(ctx context.Context, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	name, categoryID := stringValue(input, "name"), stringValue(input, "categoryId", "category_id")
	if name == "" || categoryID == "" {
		return nil, bad("ProductModel name and categoryId are required")
	}
	if err := s.validateModelInput(&snapshot, input, nil); err != nil {
		return nil, err
	}
	warranty, err := warrantyValue(input, "warrantySummary", "warranty")
	if err != nil {
		return nil, err
	}
	fixedPackID := optionalString(input, "fixedPackId", "fixed_pack_id")
	value := entity.CatalogProductModel{ID: uuid.NewString(), Name: name, CategoryID: categoryID, Description: stringValue(input, "description"), WarrantySummary: warranty, FixedAttributes: recordMapValue(input, "fixedAttributes", "fixed_attributes"), FixedPackID: fixedPackID, Measurements: recordMapValue(input, "measurements"), Status: entity.CatalogDraft, Images: []entity.CatalogProductImage{}, Dimensions: []entity.CatalogVariantDimension{}, Variants: []entity.CatalogVariant{}, CreatedAt: now(), UpdatedAt: now()}
	snapshot.Models = append(snapshot.Models, value)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, value, false)
}

func optionalString(input map[string]any, names ...string) *string {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return nil
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "" {
		return nil
	}
	return &text
}

func (s *Service) ListModels(ctx context.Context) ([]map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(snapshot.Models))
	for _, value := range snapshot.Models {
		output, outputErr := modelOutput(&snapshot, value, false)
		if outputErr != nil {
			return nil, outputErr
		}
		result = append(result, output)
	}
	return result, nil
}

func (s *Service) GetModel(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

func (s *Service) UpdateModel(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if value.Status == entity.CatalogDiscontinued && (stringValue(input, "status") != entity.CatalogDiscontinued || boolValue(input, false, "is_active")) {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	if err := s.validateModelInput(&snapshot, input, value); err != nil {
		return nil, err
	}
	if text := stringValue(input, "name"); text != "" {
		value.Name = text
	}
	if text := stringValue(input, "categoryId", "category_id"); text != "" {
		value.CategoryID = text
	}
	if raw, ok := mapValue(input, "description"); ok {
		value.Description = fmt.Sprint(raw)
	}
	if _, ok := mapValue(input, "fixedAttributes", "fixed_attributes"); ok {
		value.FixedAttributes = recordMapValue(input, "fixedAttributes", "fixed_attributes")
	}
	if raw, ok := mapValue(input, "fixedPackId", "fixed_pack_id"); ok {
		if raw == nil || strings.TrimSpace(fmt.Sprint(raw)) == "" {
			value.FixedPackID = nil
		} else {
			packID := strings.TrimSpace(fmt.Sprint(raw))
			value.FixedPackID = &packID
		}
	}
	if _, ok := mapValue(input, "measurements"); ok {
		value.Measurements = recordMapValue(input, "measurements")
	}
	if _, ok := mapValue(input, "warrantySummary", "warranty"); ok {
		warranty, warrantyErr := warrantyValue(map[string]any{"warrantySummary": mapValueMust(input, "warrantySummary", "warranty")}, "warrantySummary")
		if warrantyErr != nil {
			return nil, warrantyErr
		}
		value.WarrantySummary = warranty
	}
	if status := stringValue(input, "status"); status != "" {
		if status != entity.CatalogDraft && status != entity.CatalogActive && status != entity.CatalogInactive && status != entity.CatalogDiscontinued {
			return nil, bad("unsupported ProductModel status")
		}
		if value.Status == entity.CatalogDiscontinued && status != entity.CatalogDiscontinued {
			return nil, conflict("Discontinued ProductModel is terminal")
		}
		value.Status = status
	}
	value.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

func mapValueMust(input map[string]any, names ...string) any {
	value, _ := mapValue(input, names...)
	return value
}

func (s *Service) PublishModel(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if value.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel cannot be published")
	}
	if !value.HasPrimaryImage() {
		return nil, conflict("ProductModel requires exactly one primary model image")
	}
	if value.SaleReadyVariantCount() == 0 {
		return nil, conflict("ProductModel requires a sale-ready Variant")
	}
	value.Status, value.UpdatedAt = entity.CatalogActive, now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

func (s *Service) UnpublishModel(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if value.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	value.Status, value.UpdatedAt = entity.CatalogInactive, now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

func (s *Service) DiscontinueModel(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if value.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	value.Status, value.UpdatedAt = entity.CatalogDiscontinued, now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

func (s *Service) DeleteModel(context.Context, string) (map[string]any, error) {
	return nil, methodNotAllowed("ProductModel deletion is not supported; use lifecycle transitions")
}

func (s *Service) ReplaceMedia(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	value, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	raw, ok := mapValue(input, "images")
	if !ok {
		return nil, bad("images are required")
	}
	items := recordSlice(raw)
	if items == nil {
		return nil, bad("images must be an array")
	}
	images := make([]entity.CatalogProductImage, 0, len(items))
	primaryCount := 0
	for _, item := range items {
		image := recordMap(item)
		url := stringValue(image, "url")
		if url == "" {
			return nil, bad("image url is required")
		}
		primary := boolValue(image, false, "primary")
		if primary {
			primaryCount++
		}
		images = append(images, entity.CatalogProductImage{ID: uuid.NewString(), URL: url, Ordering: intValue(image, 0, "ordering"), PrimaryImage: primary, CreatedAt: now()})
	}
	if primaryCount > 1 || (value.Status == entity.CatalogActive && primaryCount != 1) {
		return nil, conflict("ProductModel requires exactly one primary model image")
	}
	value.Images = images
	value.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return modelOutput(&snapshot, *value, false)
}

// Variant dimensions ---------------------------------------------------------

func decodeDimensionValues(value any) ([]entity.CatalogDimensionValue, error) {
	items := recordSlice(value)
	if items == nil {
		return nil, bad("allowedValues must be an array")
	}
	values := make([]entity.CatalogDimensionValue, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		object := recordMap(item)
		id, label := stringValue(object, "id"), stringValue(object, "label")
		if id == "" || label == "" {
			return nil, bad("dimension value id and label are required")
		}
		if _, exists := seen[id]; exists {
			return nil, conflict("duplicate dimension value")
		}
		seen[id] = struct{}{}
		values = append(values, entity.CatalogDimensionValue{ID: id, Label: label, Active: boolValue(object, true, "active")})
	}
	return values, nil
}

func findDimension(model *entity.CatalogProductModel, id string) (*entity.CatalogVariantDimension, error) {
	for index := range model.Dimensions {
		if model.Dimensions[index].ID == id {
			return &model.Dimensions[index], nil
		}
	}
	return nil, notFound("VariantDimension not found")
}

func (s *Service) CreateDimension(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	if len(model.Variants) > 0 {
		return nil, conflict("VariantDimension structure cannot change after Variant creation")
	}
	definitionID := stringValue(input, "definitionId", "definition_id")
	definition, err := findDefinition(&snapshot, definitionID)
	if err != nil {
		return nil, err
	}
	if !definition.Active {
		return nil, conflict("inactive attribute definition cannot become a VariantDimension")
	}
	for _, dimension := range model.Dimensions {
		if dimension.DefinitionID == definitionID {
			return nil, conflict("VariantDimension already exists")
		}
	}
	if _, exists := model.FixedAttributes[definitionID]; exists {
		return nil, conflict("a fixed attribute cannot also be a VariantDimension")
	}
	values := []entity.CatalogDimensionValue{}
	if raw, ok := mapValue(input, "allowedValues", "allowed_values"); ok {
		values, err = decodeDimensionValues(raw)
		if err != nil {
			return nil, err
		}
	} else {
		for _, value := range definition.EnumValues {
			values = append(values, entity.CatalogDimensionValue{ID: value.ID, Label: value.Label, Active: value.Active})
		}
	}
	if len(values) == 0 {
		return nil, bad("VariantDimension requires at least one allowed value")
	}
	if definition.ValueKind == "Reference" {
		for _, value := range values {
			if err := s.activeMaster(&snapshot, strings.ToLower(definition.ReferenceTarget), value.ID); err != nil {
				return nil, err
			}
		}
	}
	dimension := entity.CatalogVariantDimension{ID: uuid.NewString(), DefinitionID: definitionID, AllowedValues: values, CreatedAt: now(), UpdatedAt: now()}
	model.Dimensions = append(model.Dimensions, dimension)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return dimensionOutput(&snapshot, modelID, dimension)
}

func (s *Service) UpdateDimension(ctx context.Context, modelID, dimensionID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	dimension, err := findDimension(model, dimensionID)
	if err != nil {
		return nil, err
	}
	hasVariants := len(model.Variants) > 0
	if raw, ok := mapValue(input, "definitionId", "definition_id"); ok {
		definitionID := strings.TrimSpace(fmt.Sprint(raw))
		if definitionID == "" {
			return nil, bad("definitionId is required")
		}
		if hasVariants && definitionID != dimension.DefinitionID {
			return nil, conflict("VariantDimension definition cannot be replaced after Variant creation")
		}
		definition, definitionErr := findDefinition(&snapshot, definitionID)
		if definitionErr != nil {
			return nil, definitionErr
		}
		if !definition.Active {
			return nil, conflict("inactive attribute definition cannot become a VariantDimension")
		}
		if _, exists := model.FixedAttributes[definitionID]; exists {
			return nil, conflict("a fixed attribute cannot also be a VariantDimension")
		}
		if definition.ValueKind == "Reference" {
			for _, value := range dimension.AllowedValues {
				if err := s.activeMaster(&snapshot, strings.ToLower(definition.ReferenceTarget), value.ID); err != nil {
					return nil, err
				}
			}
		}
		dimension.DefinitionID = definitionID
	}
	if raw, ok := mapValue(input, "allowedValues", "allowed_values"); ok {
		if hasVariants {
			return nil, conflict("VariantDimension values cannot be replaced after Variant creation")
		}
		values, decodeErr := decodeDimensionValues(raw)
		if decodeErr != nil {
			return nil, decodeErr
		}
		if len(values) == 0 {
			return nil, bad("VariantDimension requires at least one allowed value")
		}
		definition, definitionErr := findDefinition(&snapshot, dimension.DefinitionID)
		if definitionErr != nil {
			return nil, definitionErr
		}
		if definition.ValueKind == "Reference" {
			for _, value := range values {
				if err := s.activeMaster(&snapshot, strings.ToLower(definition.ReferenceTarget), value.ID); err != nil {
					return nil, err
				}
			}
		}
		dimension.AllowedValues = values
	}
	dimension.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return dimensionOutput(&snapshot, modelID, *dimension)
}

func (s *Service) AddDimensionValue(ctx context.Context, modelID, dimensionID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	dimension, err := findDimension(model, dimensionID)
	if err != nil {
		return nil, err
	}
	id, label := stringValue(input, "id"), stringValue(input, "label")
	if id == "" || label == "" {
		return nil, bad("dimension value id and label are required")
	}
	for index := range dimension.AllowedValues {
		if dimension.AllowedValues[index].ID == id {
			dimension.AllowedValues[index].Label = label
			dimension.AllowedValues[index].Active = boolValue(input, true, "active")
			dimension.UpdatedAt = now()
			if err := s.save(ctx, snapshot); err != nil {
				return nil, err
			}
			return dimensionOutput(&snapshot, modelID, *dimension)
		}
	}
	definition, err := findDefinition(&snapshot, dimension.DefinitionID)
	if err != nil {
		return nil, err
	}
	if definition.ValueKind == "Reference" {
		if err := s.activeMaster(&snapshot, strings.ToLower(definition.ReferenceTarget), id); err != nil {
			return nil, err
		}
	}
	dimension.AllowedValues = append(dimension.AllowedValues, entity.CatalogDimensionValue{ID: id, Label: label, Active: boolValue(input, true, "active")})
	dimension.UpdatedAt = now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return dimensionOutput(&snapshot, modelID, *dimension)
}

func (s *Service) DeactivateDimensionValue(ctx context.Context, modelID, dimensionID, valueID string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	dimension, err := findDimension(model, dimensionID)
	if err != nil {
		return nil, err
	}
	for index := range dimension.AllowedValues {
		if dimension.AllowedValues[index].ID == valueID {
			dimension.AllowedValues[index].Active = false
			dimension.UpdatedAt = now()
			if err := s.save(ctx, snapshot); err != nil {
				return nil, err
			}
			return dimensionOutput(&snapshot, modelID, *dimension)
		}
	}
	return nil, notFound("Variant value not found")
}

// Variants -------------------------------------------------------------------

func (s *Service) normalizeSelection(snapshot *entity.CatalogSnapshot, model *entity.CatalogProductModel, input map[string]string) (map[string]string, string, error) {
	if len(model.Dimensions) == 0 {
		return nil, "", bad("ProductModel has no VariantDimensions")
	}
	selected := map[string]string{}
	parts := make([]string, 0, len(model.Dimensions))
	for _, dimension := range model.Dimensions {
		definition, err := findDefinition(snapshot, dimension.DefinitionID)
		if err != nil {
			return nil, "", err
		}
		provided, found := "", false
		for key, value := range input {
			if strings.EqualFold(key, definition.DisplayName) || strings.EqualFold(key, definition.Key) {
				provided, found = strings.TrimSpace(value), true
				break
			}
		}
		if !found || provided == "" {
			return nil, "", bad("Variant must select one value for every dimension")
		}
		matched := false
		for _, value := range dimension.AllowedValues {
			if !value.Active {
				continue
			}
			matches := value.ID == provided || strings.EqualFold(normalizeText(value.Label), normalizeText(provided)) || canonicalValue(value.Label) == canonicalValue(provided)
			if definition.ValueKind == "Reference" && !matches {
				if master, masterErr := s.findActiveMasterByValue(snapshot, strings.ToLower(definition.ReferenceTarget), value.ID); masterErr == nil {
					matches = strings.EqualFold(normalizeText(master.Name), normalizeText(provided))
				}
			}
			if !matches {
				continue
			}
			matched = true
			if definition.ValueKind == "Reference" {
				selected[definition.DisplayName] = value.ID
				parts = append(parts, dimension.ID+"="+canonicalValue(value.ID))
			} else {
				selected[definition.DisplayName] = value.Label
				parts = append(parts, dimension.ID+"="+canonicalValue(value.Label))
			}
			break
		}
		if !matched {
			return nil, "", conflict("Variant selection is not an active allowed value")
		}
	}
	if len(input) != len(model.Dimensions) {
		return nil, "", bad("Variant selection contains an unknown dimension")
	}
	sort.Strings(parts)
	return selected, strings.Join(parts, "|"), nil
}

func parsePrice(input map[string]any, names ...string) (*entity.CatalogMoney, bool, error) {
	value, ok := mapValue(input, names...)
	if !ok {
		return nil, false, nil
	}
	if value == nil {
		return nil, true, bad("sellingPrice must be a positive VND amount")
	}
	object := recordMap(value)
	amount, valid := floatValue(object["amount"])
	currency := strings.ToUpper(stringValue(object, "currency"))
	if !valid || amount <= 0 || math.Trunc(amount) != amount || currency != "VND" {
		return nil, true, bad("sellingPrice must be a positive VND amount")
	}
	integer := int64(amount)
	return &entity.CatalogMoney{Amount: integer, Currency: currency}, true, nil
}

func appendHistory(history []entity.CatalogHistoryEntry, action string) []entity.CatalogHistoryEntry {
	result := append([]entity.CatalogHistoryEntry{}, history...)
	return append(result, entity.CatalogHistoryEntry{Action: action, At: now()})
}

func (s *Service) validatePackAssignment(snapshot *entity.CatalogSnapshot, model *entity.CatalogProductModel, selected map[string]string, requested *string) (*string, error) {
	packID := requested
	if packID == nil && model.FixedPackID != nil {
		packID = model.FixedPackID
	}
	for _, dimension := range model.Dimensions {
		definition, err := findDefinition(snapshot, dimension.DefinitionID)
		if err != nil {
			return nil, err
		}
		if definition.ValueKind == "Reference" && definition.ReferenceTarget == "Pack" {
			if value := selected[definition.DisplayName]; value != "" {
				if packID != nil && *packID != value {
					return nil, conflict("Variant Pack reference conflicts with ProductModel Pack")
				}
				packID = &value
			}
		}
	}
	if packID != nil {
		if err := s.activeMaster(snapshot, "pack", *packID); err != nil {
			return nil, err
		}
	}
	return packID, nil
}

func normalizeSKU(snapshot entity.CatalogSnapshot, value, currentID string) (string, error) {
	sku := strings.ToLower(strings.TrimSpace(value))
	if sku == "" {
		return "", nil
	}
	for _, model := range snapshot.Models {
		for _, variant := range model.Variants {
			if variant.ID != currentID && variant.SKU == sku {
				return "", conflict("SKU is already reserved by another Variant")
			}
		}
	}
	return sku, nil
}

func (s *Service) CreateVariant(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel cannot receive new Variants")
	}
	raw, ok := mapValue(input, "selectedOptions", "selected_options")
	if !ok {
		return nil, bad("selectedOptions is required")
	}
	selected, canonical, err := s.normalizeSelection(&snapshot, model, parseStringMap(raw))
	if err != nil {
		return nil, err
	}
	for _, variant := range model.Variants {
		if variant.CanonicalCombination == canonical {
			return nil, conflict("canonical Variant combination already exists")
		}
	}
	sku, err := normalizeSKU(snapshot, stringValue(input, "sku"), "")
	if err != nil {
		return nil, err
	}
	price, hasPrice, err := parsePrice(input, "sellingPrice", "selling_price")
	if err != nil {
		return nil, err
	}
	var requestedPack *string
	if rawPack, ok := mapValue(input, "packId", "pack_id"); ok && rawPack != nil && strings.TrimSpace(fmt.Sprint(rawPack)) != "" {
		pack := strings.TrimSpace(fmt.Sprint(rawPack))
		requestedPack = &pack
	}
	packID, err := s.validatePackAssignment(&snapshot, model, selected, requestedPack)
	if err != nil {
		return nil, err
	}
	technicalValues := recordMapValue(input, "technicalValues", "technical_values", "variantAttributes", "variant_attributes")
	variant := entity.CatalogVariant{ID: uuid.NewString(), SelectedOptions: selected, TechnicalValues: technicalValues, SKU: sku, SellingPrice: price, PackID: packID, Status: entity.CatalogVariantActive, CanonicalCombination: canonical, History: appendHistory(nil, "created"), CreatedAt: now(), UpdatedAt: now()}
	if !hasPrice {
		variant.SellingPrice = nil
	}
	model.Variants = append(model.Variants, variant)
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return variantOutput(variant), nil
}

func (s *Service) ListVariants(ctx context.Context, modelID string) ([]map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(model.Variants))
	for _, value := range model.Variants {
		result = append(result, variantOutput(value))
	}
	return result, nil
}

func (s *Service) GetVariant(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	variant, _, err := findVariant(&snapshot, id)
	if err != nil {
		return nil, err
	}
	return variantOutput(*variant), nil
}

func (s *Service) UpdateVariant(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	variant, model, err := findVariant(&snapshot, id)
	if err != nil {
		return nil, err
	}
	wasSaleReady := variant.SaleReady()
	if raw, ok := mapValue(input, "selectedOptions", "selected_options"); ok {
		_, canonical, canonicalErr := s.normalizeSelection(&snapshot, model, parseStringMap(raw))
		if canonicalErr != nil {
			return nil, canonicalErr
		}
		if canonical != variant.CanonicalCombination {
			return nil, conflict("Variant selected combination is immutable")
		}
	}
	if _, ok := mapValue(input, "sku"); ok {
		variant.SKU, err = normalizeSKU(snapshot, stringValue(input, "sku"), id)
		if err != nil {
			return nil, err
		}
	}
	if price, hasPrice, priceErr := parsePrice(input, "sellingPrice", "selling_price"); priceErr != nil {
		return nil, priceErr
	} else if hasPrice {
		variant.SellingPrice = price
	}
	if _, ok := mapValue(input, "technicalValues", "technical_values", "variantAttributes", "variant_attributes"); ok {
		variant.TechnicalValues = recordMapValue(input, "technicalValues", "technical_values", "variantAttributes", "variant_attributes")
	}
	if rawPack, ok := mapValue(input, "packId", "pack_id"); ok {
		if rawPack == nil || strings.TrimSpace(fmt.Sprint(rawPack)) == "" {
			if model.FixedPackID != nil {
				return nil, conflict("Variant must retain the ProductModel fixed Pack")
			}
			variant.PackID = nil
		} else {
			pack := strings.TrimSpace(fmt.Sprint(rawPack))
			if model.FixedPackID != nil && pack != *model.FixedPackID {
				return nil, conflict("Variant Pack reference conflicts with ProductModel Pack")
			}
			if err := s.activeMaster(&snapshot, "pack", pack); err != nil {
				return nil, err
			}
			variant.PackID = &pack
		}
	}
	if model.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel cannot be mutated")
	}
	if model.Status == entity.CatalogActive && wasSaleReady && !variant.SaleReady() {
		remainingSaleReady := 0
		for index := range model.Variants {
			if model.Variants[index].ID != variant.ID && model.Variants[index].SaleReady() {
				remainingSaleReady++
			}
		}
		if remainingSaleReady == 0 {
			return nil, conflict("Active ProductModel must keep one sale-ready Variant")
		}
	}
	variant.History, variant.UpdatedAt = appendHistory(variant.History, "updated"), now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return variantOutput(*variant), nil
}

func (s *Service) ActivateVariant(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	variant, model, err := findVariant(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if model.Status == entity.CatalogDiscontinued {
		return nil, conflict("Discontinued ProductModel cannot receive activated Variants")
	}
	variant.Status, variant.History, variant.UpdatedAt = entity.CatalogVariantActive, appendHistory(variant.History, "activated"), now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return variantOutput(*variant), nil
}

func (s *Service) InactivateVariant(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	variant, model, err := findVariant(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if model.Status == entity.CatalogActive && variant.SaleReady() && model.SaleReadyVariantCount() == 1 {
		return nil, conflict("Active ProductModel must keep one sale-ready Variant")
	}
	variant.Status, variant.History, variant.UpdatedAt = entity.CatalogVariantInactive, appendHistory(variant.History, "inactivated"), now()
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	return variantOutput(*variant), nil
}

func (s *Service) BulkSetPrice(ctx context.Context, input map[string]any) ([]map[string]any, error) {
	idsValue, ok := mapValue(input, "variantIds", "variant_ids")
	if !ok {
		return nil, bad("variantIds are required")
	}
	items := recordSlice(idsValue)
	if len(items) == 0 {
		return nil, bad("variantIds must not be empty")
	}
	price, hasPrice, err := parsePrice(input, "sellingPrice", "selling_price")
	if err != nil {
		return nil, err
	}
	if !hasPrice {
		return nil, bad("sellingPrice is required")
	}
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, strings.TrimSpace(fmt.Sprint(item)))
	}
	variants := make([]*entity.CatalogVariant, 0, len(ids))
	for _, id := range ids {
		variant, _, findErr := findVariant(&snapshot, id)
		if findErr != nil {
			return nil, conflict("bulk price contains an unknown Variant")
		}
		variants = append(variants, variant)
	}
	for _, variant := range variants {
		variant.SellingPrice, variant.History, variant.UpdatedAt = price, appendHistory(variant.History, "price_updated"), now()
	}
	if err := s.save(ctx, snapshot); err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(variants))
	for _, variant := range variants {
		result = append(result, variantOutput(*variant))
	}
	return result, nil
}

// Public catalog -------------------------------------------------------------

type PublicFilter struct {
	Page       int
	Limit      int
	CategoryID string
	MaterialID string
	FinishID   string
	MinPrice   *int64
	MaxPrice   *int64
	Search     string
	Sort       string
}

func publicVariants(model entity.CatalogProductModel) []entity.CatalogVariant {
	result := make([]entity.CatalogVariant, 0)
	for _, variant := range model.Variants {
		if variant.SaleReady() {
			result = append(result, variant)
		}
	}
	return result
}

func priceInRange(variants []entity.CatalogVariant, minPrice, maxPrice *int64) bool {
	for _, variant := range variants {
		if variant.SellingPrice == nil {
			continue
		}
		if minPrice != nil && variant.SellingPrice.Amount < *minPrice {
			continue
		}
		if maxPrice != nil && variant.SellingPrice.Amount > *maxPrice {
			continue
		}
		return true
	}
	return false
}

func mapMatches(value map[string]any, key, expected string) bool {
	if expected == "" {
		return true
	}
	for actualKey, actualValue := range value {
		if strings.EqualFold(actualKey, key) && fmt.Sprint(actualValue) == expected {
			return true
		}
	}
	return false
}

func (s *Service) referenceMatches(snapshot *entity.CatalogSnapshot, fixed map[string]any, key, expected string) (bool, error) {
	if expected == "" || mapMatches(fixed, key, expected) {
		return true, nil
	}
	target := strings.TrimSuffix(strings.ToLower(key), "id")
	for definitionID, rawValue := range fixed {
		if !uuidPattern.MatchString(definitionID) {
			continue
		}
		definition, err := findDefinition(snapshot, definitionID)
		if err != nil {
			continue
		}
		if definition.ValueKind == "Reference" && strings.EqualFold(definition.ReferenceTarget, target) && fmt.Sprint(rawValue) == expected {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) ListPublicModels(ctx context.Context, filter PublicFilter) (map[string]any, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	type candidate struct {
		model    entity.CatalogProductModel
		variants []entity.CatalogVariant
	}
	candidates := make([]candidate, 0)
	for _, model := range snapshot.Models {
		if model.Status != entity.CatalogActive || (filter.CategoryID != "" && model.CategoryID != filter.CategoryID) {
			continue
		}
		variants := publicVariants(model)
		if len(variants) == 0 || !priceInRange(variants, filter.MinPrice, filter.MaxPrice) {
			continue
		}
		materialMatches, materialErr := s.referenceMatches(&snapshot, model.FixedAttributes, "materialId", filter.MaterialID)
		if materialErr != nil {
			return nil, materialErr
		}
		finishMatches, finishErr := s.referenceMatches(&snapshot, model.FixedAttributes, "finishId", filter.FinishID)
		if finishErr != nil {
			return nil, finishErr
		}
		if !materialMatches || !finishMatches {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(model.Name+" "+model.Description), strings.ToLower(filter.Search)) {
			continue
		}
		candidates = append(candidates, candidate{model: model, variants: variants})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if filter.Sort == "newest" {
			return candidates[i].model.CreatedAt.After(candidates[j].model.CreatedAt)
		}
		price := func(item candidate) int64 {
			if len(item.variants) == 0 {
				return 0
			}
			lowest := int64(1<<63 - 1)
			for _, variant := range item.variants {
				if variant.SellingPrice != nil && variant.SellingPrice.Amount < lowest {
					lowest = variant.SellingPrice.Amount
				}
			}
			if lowest == int64(1<<63-1) {
				return 0
			}
			return lowest
		}
		if filter.Sort == "price_desc" {
			return price(candidates[i]) > price(candidates[j])
		}
		if filter.Sort == "price_asc" {
			return price(candidates[i]) < price(candidates[j])
		}
		return candidates[i].model.CreatedAt.Before(candidates[j].model.CreatedAt)
	})
	total := len(candidates)
	start := (filter.Page - 1) * filter.Limit
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}
	items := make([]map[string]any, 0, end-start)
	for _, item := range candidates[start:end] {
		output, outputErr := modelOutput(&snapshot, item.model, true)
		if outputErr != nil {
			return nil, outputErr
		}
		items = append(items, output)
	}
	return map[string]any{"items": items, "page": filter.Page, "limit": filter.Limit, "total": total}, nil
}

func (s *Service) GetPublicModel(ctx context.Context, id string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, id)
	if err != nil {
		return nil, err
	}
	if model.Status != entity.CatalogActive || len(publicVariants(*model)) == 0 {
		return nil, notFound("public ProductModel not found")
	}
	return modelOutput(&snapshot, *model, true)
}

func selectedCompatible(variant entity.CatalogVariant, selected map[string]string) bool {
	for key, value := range selected {
		found := false
		for optionKey, optionValue := range variant.SelectedOptions {
			if strings.EqualFold(optionKey, key) && canonicalValue(optionValue) == canonicalValue(value) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (s *Service) AvailableOptions(ctx context.Context, modelID string, selected map[string]string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status != entity.CatalogActive {
		return nil, notFound("public ProductModel not found")
	}
	variants := publicVariants(*model)
	compatible := make([]entity.CatalogVariant, 0)
	for _, variant := range variants {
		if selectedCompatible(variant, selected) {
			compatible = append(compatible, variant)
		}
	}
	if len(compatible) == 0 {
		return map[string]any{"options": []any{}}, nil
	}
	options := make([]map[string]any, 0, len(model.Dimensions))
	for _, dimension := range model.Dimensions {
		definition, definitionErr := findDefinition(&snapshot, dimension.DefinitionID)
		if definitionErr != nil {
			return nil, definitionErr
		}
		valueSet := map[string]entity.CatalogDimensionValue{}
		for _, variant := range compatible {
			for key, value := range variant.SelectedOptions {
				if strings.EqualFold(key, definition.DisplayName) {
					for _, allowed := range dimension.AllowedValues {
						if allowed.ID == value || canonicalValue(allowed.Label) == canonicalValue(value) {
							valueSet[allowed.ID] = allowed
						}
					}
				}
			}
		}
		values := make([]entity.CatalogDimensionValue, 0, len(valueSet))
		for _, value := range valueSet {
			values = append(values, value)
		}
		sort.Slice(values, func(i, j int) bool { return values[i].Label < values[j].Label })
		valueOutput := make([]map[string]any, 0, len(values))
		for _, value := range values {
			valueOutput = append(valueOutput, map[string]any{"id": value.ID, "label": value.Label})
		}
		if len(valueOutput) > 0 {
			options = append(options, map[string]any{"key": definition.DisplayName, "values": valueOutput})
		}
	}
	return map[string]any{"options": options}, nil
}

func (s *Service) ResolvePublicVariant(ctx context.Context, modelID string, selected map[string]string) (map[string]any, error) {
	snapshot, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	model, err := findModel(&snapshot, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status != entity.CatalogActive {
		return nil, notFound("public Variant not found")
	}
	variants := publicVariants(*model)
	_, canonical, normalizeErr := s.normalizeSelection(&snapshot, model, selected)
	if normalizeErr != nil {
		parts := make([]string, 0, len(model.Dimensions))
		for _, dimension := range model.Dimensions {
			definition, definitionErr := findDefinition(&snapshot, dimension.DefinitionID)
			if definitionErr != nil {
				return nil, definitionErr
			}
			value := selected[definition.DisplayName]
			if value == "" {
				return nil, notFound("public Variant not found")
			}
			parts = append(parts, dimension.ID+"="+canonicalValue(value))
		}
		sort.Strings(parts)
		canonical = strings.Join(parts, "|")
	}
	for _, variant := range variants {
		if variant.CanonicalCombination == canonical {
			return variantOutput(variant), nil
		}
	}
	return nil, notFound("public Variant not found")
}
