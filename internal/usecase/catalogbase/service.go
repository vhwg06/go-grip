// Package catalogbase implements the Catalog Base vertical slice.
//
// Catalog Base deliberately lives beside the legacy Product use case.  The
// legacy surface is still used by the existing storefront, while this service
// owns the ProductModel/Variant contract described by the Catalog Base spec.
package catalogbase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIError is a domain error with the HTTP status required by the Catalog Base
// contract.  Controllers should use ErrorStatus when translating failures.
type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string { return e.Message }

func bad(message string) error {
	return &APIError{Status: 400, Code: "invalid_catalog_command", Message: message}
}

func conflict(message string) error {
	return &APIError{Status: 409, Code: "catalog_conflict", Message: message}
}

func notFound(message string) error {
	return &APIError{Status: 404, Code: "catalog_not_found", Message: message}
}

func methodNotAllowed(message string) error {
	return &APIError{Status: 405, Code: "catalog_method_not_allowed", Message: message}
}

// ErrorStatus converts a service error into an HTTP status and a stable error
// body.  Unknown persistence errors are intentionally reported as 500.
func ErrorStatus(err error) (int, map[string]any) {
	if err == nil {
		return 200, nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status, map[string]any{"code": apiErr.Code, "message": apiErr.Message}
	}
	return 500, map[string]any{"code": "catalog_internal_error", "message": "catalog operation failed"}
}

// Service is the transactional Catalog Base application service.
type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service { return &Service{db: db} }

type categoryRow struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Slug      string    `gorm:"not null"`
	ParentID  *string   `gorm:"type:uuid"`
	Position  int       `gorm:"not null;default:0"`
	Active    bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (categoryRow) TableName() string { return "catalog_base_categories" }

type attributeDefinitionRow struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	Key             string    `gorm:"not null"`
	DisplayName     string    `gorm:"column:display_name;not null;default:''"`
	Description     string    `gorm:"not null;default:''"`
	Ordering        int       `gorm:"not null;default:0"`
	ValueKind       string    `gorm:"column:value_kind;not null"`
	DataType        string    `gorm:"column:data_type;not null;default:''"`
	ReferenceTarget string    `gorm:"column:reference_target;not null;default:''"`
	UnitFamily      string    `gorm:"column:unit_family;not null;default:''"`
	Unit            string    `gorm:"not null;default:''"`
	Active          bool      `gorm:"not null;default:true"`
	EnumValues      string    `gorm:"column:enum_values;type:jsonb;not null"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (attributeDefinitionRow) TableName() string { return "catalog_base_attribute_definitions" }

type masterRow struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Kind        string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null;default:''"`
	SwatchMedia string    `gorm:"column:swatch_media;type:jsonb;not null"`
	SellingUnit string    `gorm:"column:selling_unit;not null;default:''"`
	Quantity    *float64  `gorm:"type:numeric"`
	BaseUnit    string    `gorm:"column:base_unit;not null;default:''"`
	Active      bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (masterRow) TableName() string { return "catalog_base_masters" }

type productModelRow struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	Name            string    `gorm:"not null"`
	CategoryID      string    `gorm:"column:category_id;type:uuid;not null"`
	Description     string    `gorm:"not null;default:''"`
	WarrantySummary string    `gorm:"column:warranty_summary;type:jsonb"`
	FixedAttributes string    `gorm:"column:fixed_attributes;type:jsonb;not null"`
	FixedPackID     *string   `gorm:"column:fixed_pack_id;type:uuid"`
	Measurements    string    `gorm:"type:jsonb;not null"`
	Status          string    `gorm:"not null;default:'Draft'"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (productModelRow) TableName() string { return "catalog_base_product_models" }

type productImageRow struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	ModelID      string    `gorm:"column:model_id;type:uuid;not null"`
	URL          string    `gorm:"not null"`
	Ordering     int       `gorm:"not null;default:0"`
	PrimaryImage bool      `gorm:"column:primary_image;not null;default:false"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (productImageRow) TableName() string { return "catalog_base_product_images" }

type variantDimensionRow struct {
	ID            string    `gorm:"type:uuid;primaryKey"`
	ModelID       string    `gorm:"column:model_id;type:uuid;not null"`
	DefinitionID  string    `gorm:"column:definition_id;type:uuid;not null"`
	AllowedValues string    `gorm:"column:allowed_values;type:jsonb;not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (variantDimensionRow) TableName() string { return "catalog_base_variant_dimensions" }

type variantRow struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	ModelID              string    `gorm:"column:model_id;type:uuid;not null"`
	SelectedOptions      string    `gorm:"column:selected_options;type:jsonb;not null"`
	TechnicalValues      string    `gorm:"column:technical_values;type:jsonb;not null"`
	SKU                  string    `gorm:"not null;default:''"`
	SellingAmount        *int64    `gorm:"column:selling_amount"`
	SellingCurrency      string    `gorm:"column:selling_currency;not null;default:''"`
	PackID               *string   `gorm:"column:pack_id;type:uuid"`
	Status               string    `gorm:"not null;default:'Active'"`
	CanonicalCombination string    `gorm:"column:canonical_combination;not null"`
	History              string    `gorm:"type:jsonb;not null"`
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

func (variantRow) TableName() string { return "catalog_base_variants" }

type enumValue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type dimensionValue struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
}

type historyEntry struct {
	Action string    `json:"action"`
	At     time.Time `json:"at"`
}

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F-]{20,}$`)

func now() time.Time { return time.Now().UTC() }

func jsonString(value any, fallback string) string {
	if value == nil {
		return fallback
	}
	b, err := json.Marshal(value)
	if err != nil {
		return fallback
	}
	return string(b)
}

func jsonMap(raw string) map[string]any {
	result := map[string]any{}
	if raw == "" {
		return result
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil || result == nil {
		return map[string]any{}
	}
	return result
}

func jsonSlice(raw string) []any {
	result := []any{}
	if raw == "" {
		return result
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil || result == nil {
		return []any{}
	}
	return result
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

func boolValue(input map[string]any, defaultValue bool, names ...string) bool {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return defaultValue
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(strings.TrimSpace(typed), "true")
	default:
		return defaultValue
	}
}

func intValue(input map[string]any, defaultValue int, names ...string) int {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return defaultValue
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			return defaultValue
		}
		return int(parsed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return defaultValue
		}
		return parsed
	default:
		return defaultValue
	}
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
	case json.Number:
		parsed, err := typed.Float64()
		return parsed, err == nil
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed, err == nil
	default:
		return 0, false
	}
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
	type unitInfo struct {
		factor float64
		family string
	}
	info, ok := map[string]unitInfo{
		"mm": {factor: 1, family: "length"},
		"cm": {factor: 10, family: "length"},
		"m":  {factor: 1000, family: "length"},
		"g":  {factor: 1, family: "mass"},
		"kg": {factor: 1000, family: "mass"},
	}[parts[1]]
	if !ok {
		return 0, "", false
	}
	return number * info.factor, info.family, true
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

func validateUnit(family, unit string) bool {
	if family == "" && unit == "" {
		return true
	}
	return family != "" && unit != "" && compatibleUnit(family, unit)
}

func validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit string) error {
	valueKind = strings.TrimSpace(valueKind)
	dataType = strings.TrimSpace(dataType)
	referenceTarget = strings.TrimSpace(referenceTarget)
	switch valueKind {
	case "Scalar":
		if dataType == "" || (dataType != "Text" && dataType != "Number" && dataType != "Boolean") {
			return bad("Scalar requires a supported dataType")
		}
		if dataType == "Number" && !validateUnit(unitFamily, unit) {
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

func warrantyValue(input map[string]any, names ...string) (string, error) {
	value, ok := mapValue(input, names...)
	if !ok || value == nil {
		return "", nil
	}
	object, ok := value.(map[string]any)
	if !ok {
		return "", bad("warrantySummary must be an object")
	}
	term := stringValue(object, "term")
	if term == "" {
		return "", bad("warrantySummary.term is required")
	}
	note := stringValue(object, "note")
	return jsonString(map[string]any{"term": term, "note": note}, "{}"), nil
}

func validateMeasurements(value any) error {
	object, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	for name, raw := range object {
		measurement, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		unit := stringValue(measurement, "unit")
		if strings.Contains(strings.ToLower(name), "length") && unit != "" && !compatibleUnit("length", unit) {
			return bad("length measurement has an incompatible unit")
		}
	}
	return nil
}

func parseStringMap(value any) map[string]string {
	result := map[string]string{}
	object, ok := value.(map[string]any)
	if !ok {
		if typed, ok := value.(map[string]string); ok {
			return typed
		}
		return result
	}
	for key, raw := range object {
		if raw != nil {
			result[key] = fmt.Sprint(raw)
		}
	}
	return result
}

func (s *Service) findModel(ctx context.Context, id string) (productModelRow, error) {
	var row productModelRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productModelRow{}, notFound("ProductModel not found")
		}
		return productModelRow{}, err
	}
	return row, nil
}

func (s *Service) findVariant(ctx context.Context, id string) (variantRow, error) {
	var row variantRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return variantRow{}, notFound("Variant not found")
		}
		return variantRow{}, err
	}
	return row, nil
}

func (s *Service) findDefinition(ctx context.Context, id string) (attributeDefinitionRow, error) {
	var row attributeDefinitionRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return attributeDefinitionRow{}, notFound("attribute definition not found")
		}
		return attributeDefinitionRow{}, err
	}
	return row, nil
}

func (s *Service) findMaster(ctx context.Context, kind, id string) (masterRow, error) {
	var row masterRow
	if err := s.db.WithContext(ctx).Where("id = ? AND kind = ?", id, kind).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterRow{}, notFound("catalog master not found")
		}
		return masterRow{}, err
	}
	return row, nil
}

func (s *Service) activeMaster(ctx context.Context, kind, id string) error {
	row, err := s.findMaster(ctx, kind, id)
	if err != nil {
		return err
	}
	if !row.Active {
		return conflict("inactive catalog master cannot be assigned")
	}
	return nil
}

func (s *Service) categoryOutput(row categoryRow) map[string]any {
	result := map[string]any{
		"id": row.ID, "name": row.Name, "slug": row.Slug, "position": row.Position, "active": row.Active,
	}
	if row.ParentID != nil {
		result["parentId"] = *row.ParentID
	} else {
		result["parentId"] = nil
	}
	return result
}

func (s *Service) definitionOutput(row attributeDefinitionRow) map[string]any {
	result := map[string]any{
		"id": row.ID, "key": row.Key, "displayName": row.DisplayName, "description": row.Description,
		"ordering": row.Ordering, "valueKind": row.ValueKind, "active": row.Active,
	}
	if row.DataType != "" {
		result["dataType"] = row.DataType
	}
	if row.ReferenceTarget != "" {
		result["referenceTarget"] = row.ReferenceTarget
	}
	if row.UnitFamily != "" {
		result["unitFamily"] = row.UnitFamily
	}
	if row.Unit != "" {
		result["unit"] = row.Unit
	}
	values := []enumValue{}
	if err := json.Unmarshal([]byte(row.EnumValues), &values); err != nil {
		values = []enumValue{}
	}
	if values == nil {
		values = []enumValue{}
	}
	result["enumValues"] = values
	return result
}

func (s *Service) masterOutput(row masterRow) map[string]any {
	result := map[string]any{
		"id": row.ID, "kind": row.Kind, "name": row.Name, "description": row.Description, "active": row.Active,
		"swatchMedia": jsonSlice(row.SwatchMedia), "sellingUnit": row.SellingUnit, "baseUnit": row.BaseUnit,
	}
	if row.Quantity != nil {
		result["quantity"] = *row.Quantity
	}
	return result
}

func (s *Service) dimensionOutput(ctx context.Context, row variantDimensionRow) (map[string]any, error) {
	definition, err := s.findDefinition(ctx, row.DefinitionID)
	if err != nil {
		return nil, err
	}
	values := []dimensionValue{}
	if err := json.Unmarshal([]byte(row.AllowedValues), &values); err != nil || values == nil {
		values = []dimensionValue{}
	}
	return map[string]any{
		"id": row.ID, "definitionId": row.DefinitionID, "key": definition.DisplayName, "allowedValues": values,
	}, nil
}

func (s *Service) variantOutput(row variantRow) map[string]any {
	result := map[string]any{
		"id": row.ID, "modelId": row.ModelID, "selectedOptions": jsonMap(row.SelectedOptions),
		"technicalValues": jsonMap(row.TechnicalValues), "sku": row.SKU, "status": row.Status,
		"saleReady": saleReady(row), "canonicalCombination": row.CanonicalCombination,
		"history": jsonSlice(row.History), "packId": nil,
	}
	if row.PackID != nil {
		result["packId"] = *row.PackID
	}
	if row.SellingAmount != nil {
		result["sellingPrice"] = map[string]any{"amount": *row.SellingAmount, "currency": row.SellingCurrency}
	} else {
		result["sellingPrice"] = nil
	}
	return result
}

func saleReady(row variantRow) bool {
	return row.Status == "Active" && strings.TrimSpace(row.SKU) != "" && row.SellingAmount != nil && *row.SellingAmount > 0 && strings.EqualFold(row.SellingCurrency, "VND")
}

func (s *Service) modelOutput(ctx context.Context, row productModelRow, public bool) (map[string]any, error) {
	result := map[string]any{
		"id": row.ID, "name": row.Name, "categoryId": row.CategoryID, "description": row.Description,
		"fixedAttributes": jsonMap(row.FixedAttributes), "measurements": jsonMap(row.Measurements),
		"status": row.Status, "images": []map[string]any{}, "variants": []map[string]any{}, "variantDimensions": []map[string]any{},
	}
	if row.FixedPackID != nil {
		result["fixedPackId"] = *row.FixedPackID
	}
	if row.WarrantySummary != "" {
		result["warrantySummary"] = jsonMap(row.WarrantySummary)
	} else {
		result["warrantySummary"] = nil
	}

	var images []productImageRow
	if err := s.db.WithContext(ctx).Where("model_id = ?", row.ID).Order("ordering ASC, id ASC").Find(&images).Error; err != nil {
		return nil, err
	}
	imageOutput := make([]map[string]any, 0, len(images))
	for _, image := range images {
		imageOutput = append(imageOutput, map[string]any{"id": image.ID, "url": image.URL, "ordering": image.Ordering, "primary": image.PrimaryImage})
	}
	result["images"] = imageOutput

	var dimensions []variantDimensionRow
	if err := s.db.WithContext(ctx).Where("model_id = ?", row.ID).Order("created_at ASC, id ASC").Find(&dimensions).Error; err != nil {
		return nil, err
	}
	dimensionOutput := make([]map[string]any, 0, len(dimensions))
	for _, dimension := range dimensions {
		value, err := s.dimensionOutput(ctx, dimension)
		if err != nil {
			return nil, err
		}
		dimensionOutput = append(dimensionOutput, value)
	}
	result["variantDimensions"] = dimensionOutput

	query := s.db.WithContext(ctx).Where("model_id = ?", row.ID).Order("created_at ASC, id ASC")
	var variants []variantRow
	if err := query.Find(&variants).Error; err != nil {
		return nil, err
	}
	variantOutput := make([]map[string]any, 0, len(variants))
	for _, variant := range variants {
		if public && !saleReady(variant) {
			continue
		}
		variantOutput = append(variantOutput, s.variantOutput(variant))
	}
	result["variants"] = variantOutput
	return result, nil
}

// Categories -----------------------------------------------------------------

func (s *Service) ListCategories(ctx context.Context) ([]map[string]any, error) {
	var rows []categoryRow
	if err := s.db.WithContext(ctx).Order("position ASC, name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		result = append(result, s.categoryOutput(row))
	}
	return result, nil
}

func (s *Service) CreateCategory(ctx context.Context, input map[string]any) (map[string]any, error) {
	name, slug := stringValue(input, "name"), stringValue(input, "slug")
	if name == "" || slug == "" {
		return nil, bad("category name and slug are required")
	}
	parentID := optionalString(input, "parentId", "parent_id")
	if parentID != nil {
		var parent categoryRow
		if err := s.db.WithContext(ctx).Where("id = ?", *parentID).First(&parent).Error; err != nil {
			return nil, notFound("parent category not found")
		}
	}
	row := categoryRow{ID: uuid.NewString(), Name: name, Slug: slug, ParentID: parentID, Position: intValue(input, 0, "position"), Active: boolValue(input, true, "active"), CreatedAt: now(), UpdatedAt: now()}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, conflict("category slug already exists")
		}
		return nil, err
	}
	return s.categoryOutput(row), nil
}

func (s *Service) UpdateCategory(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	var row categoryRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("category not found")
		}
		return nil, err
	}
	if value := stringValue(input, "name"); value != "" {
		row.Name = value
	}
	if value := stringValue(input, "slug"); value != "" {
		row.Slug = value
	}
	if _, ok := mapValue(input, "position"); ok {
		row.Position = intValue(input, row.Position, "position")
	}
	if _, ok := mapValue(input, "active"); ok {
		row.Active = boolValue(input, row.Active, "active")
	}
	if parentID, ok := mapValue(input, "parentId", "parent_id"); ok {
		if parentID == nil || strings.TrimSpace(fmt.Sprint(parentID)) == "" {
			row.ParentID = nil
		} else {
			value := strings.TrimSpace(fmt.Sprint(parentID))
			row.ParentID = &value
		}
	}
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.categoryOutput(row), nil
}

func (s *Service) DeactivateCategory(ctx context.Context, id string) (map[string]any, error) {
	return s.UpdateCategory(ctx, id, map[string]any{"active": false})
}

func (s *Service) DeleteCategory(context.Context, string) (map[string]any, error) {
	return nil, methodNotAllowed("category deletion is not supported; deactivate it instead")
}

// Attribute definitions ------------------------------------------------------

func (s *Service) ListDefinitions(ctx context.Context) ([]map[string]any, error) {
	var rows []attributeDefinitionRow
	if err := s.db.WithContext(ctx).Order("ordering ASC, display_name ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		result = append(result, s.definitionOutput(row))
	}
	return result, nil
}

func (s *Service) CreateDefinition(ctx context.Context, input map[string]any) (map[string]any, error) {
	key := stringValue(input, "key")
	valueKind := stringValue(input, "valueKind", "value_kind")
	if key == "" || valueKind == "" {
		return nil, bad("attribute key and valueKind are required")
	}
	dataType := stringValue(input, "dataType", "data_type")
	referenceTarget := stringValue(input, "referenceTarget", "reference_target")
	unitFamily, unit := stringValue(input, "unitFamily", "unit_family"), stringValue(input, "unit")
	if err := validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit); err != nil {
		return nil, err
	}
	row := attributeDefinitionRow{
		ID: uuid.NewString(), Key: key, DisplayName: stringValue(input, "displayName", "display_name"), Description: stringValue(input, "description"),
		Ordering: intValue(input, 0, "ordering"), ValueKind: valueKind, DataType: dataType, ReferenceTarget: referenceTarget,
		UnitFamily: unitFamily, Unit: unit, Active: boolValue(input, true, "active"), EnumValues: "[]", CreatedAt: now(), UpdatedAt: now(),
	}
	if row.DisplayName == "" {
		row.DisplayName = key
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, conflict("attribute key already exists")
		}
		return nil, err
	}
	return s.definitionOutput(row), nil
}

func (s *Service) definitionUsed(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&variantDimensionRow{}).Where("definition_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	var models []productModelRow
	if err := s.db.WithContext(ctx).Find(&models).Error; err != nil {
		return false, err
	}
	for _, model := range models {
		if _, ok := jsonMap(model.FixedAttributes)[id]; ok {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) UpdateDefinition(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	row, err := s.findDefinition(ctx, id)
	if err != nil {
		return nil, err
	}
	used, err := s.definitionUsed(ctx, id)
	if err != nil {
		return nil, err
	}
	valueKind, dataType, referenceTarget := row.ValueKind, row.DataType, row.ReferenceTarget
	unitFamily, unit := row.UnitFamily, row.Unit
	if value, ok := mapValue(input, "valueKind", "value_kind"); ok {
		valueKind = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := mapValue(input, "dataType", "data_type"); ok {
		dataType = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := mapValue(input, "referenceTarget", "reference_target"); ok {
		referenceTarget = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := mapValue(input, "unitFamily", "unit_family"); ok {
		unitFamily = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := mapValue(input, "unit"); ok {
		unit = strings.TrimSpace(fmt.Sprint(value))
	}
	if used && (valueKind != row.ValueKind || dataType != row.DataType || referenceTarget != row.ReferenceTarget || unitFamily != row.UnitFamily || unit != row.Unit) {
		return nil, conflict("used attribute semantic structure cannot change")
	}
	if err := validateDefinition(valueKind, dataType, referenceTarget, unitFamily, unit); err != nil {
		return nil, err
	}
	row.ValueKind, row.DataType, row.ReferenceTarget, row.UnitFamily, row.Unit = valueKind, dataType, referenceTarget, unitFamily, unit
	if value := stringValue(input, "displayName", "display_name"); value != "" {
		row.DisplayName = value
	}
	if value, ok := mapValue(input, "description"); ok {
		row.Description = fmt.Sprint(value)
	}
	if _, ok := mapValue(input, "ordering"); ok {
		row.Ordering = intValue(input, row.Ordering, "ordering")
	}
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.definitionOutput(row), nil
}

func (s *Service) DeactivateDefinition(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findDefinition(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Active, row.UpdatedAt = false, now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.definitionOutput(row), nil
}

func (s *Service) AddEnumValue(ctx context.Context, definitionID string, input map[string]any) (map[string]any, error) {
	row, err := s.findDefinition(ctx, definitionID)
	if err != nil {
		return nil, err
	}
	if row.ValueKind != "Enum" {
		return nil, bad("enum values require an Enum definition")
	}
	key, label := stringValue(input, "key"), stringValue(input, "label")
	if key == "" || label == "" {
		return nil, bad("enum key and label are required")
	}
	values := []enumValue{}
	if err := json.Unmarshal([]byte(row.EnumValues), &values); err != nil {
		values = []enumValue{}
	}
	for _, value := range values {
		if strings.EqualFold(value.Key, key) {
			return nil, conflict("enum key already exists")
		}
	}
	item := enumValue{ID: uuid.NewString(), Key: key, Label: label, Active: boolValue(input, true, "active")}
	values = append(values, item)
	row.EnumValues, row.UpdatedAt = jsonString(values, "[]"), now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return map[string]any{"id": item.ID, "key": item.Key, "label": item.Label, "active": item.Active}, nil
}

func (s *Service) DeactivateEnumValue(ctx context.Context, definitionID, enumID string) (map[string]any, error) {
	row, err := s.findDefinition(ctx, definitionID)
	if err != nil {
		return nil, err
	}
	values := []enumValue{}
	if err := json.Unmarshal([]byte(row.EnumValues), &values); err != nil {
		values = []enumValue{}
	}
	for index := range values {
		if values[index].ID == enumID {
			values[index].Active = false
			row.EnumValues, row.UpdatedAt = jsonString(values, "[]"), now()
			if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
				return nil, err
			}
			return map[string]any{"id": values[index].ID, "key": values[index].Key, "label": values[index].Label, "active": false}, nil
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
	var rows []masterRow
	if err := s.db.WithContext(ctx).Where("kind = ?", kind).Order("name ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		result = append(result, s.masterOutput(row))
	}
	return result, nil
}

func (s *Service) CreateMaster(ctx context.Context, kind string, input map[string]any) (map[string]any, error) {
	if !validMasterKind(kind) {
		return nil, bad("unsupported catalog master kind")
	}
	name := stringValue(input, "name")
	if name == "" {
		return nil, bad("master name is required")
	}
	var quantity *float64
	if value, ok := mapValue(input, "quantity"); ok && value != nil {
		parsed, valid := floatValue(value)
		if !valid || parsed <= 0 {
			return nil, bad("pack quantity must be positive")
		}
		quantity = &parsed
	}
	row := masterRow{ID: uuid.NewString(), Kind: kind, Name: name, Description: stringValue(input, "description"), SwatchMedia: jsonString(input["swatchMedia"], "[]"), SellingUnit: stringValue(input, "sellingUnit"), Quantity: quantity, BaseUnit: stringValue(input, "baseUnit"), Active: boolValue(input, true, "active"), CreatedAt: now(), UpdatedAt: now()}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, conflict("catalog master already exists")
		}
		return nil, err
	}
	return s.masterOutput(row), nil
}

func (s *Service) UpdateMaster(ctx context.Context, kind, id string, input map[string]any) (map[string]any, error) {
	if !validMasterKind(kind) {
		return nil, bad("unsupported catalog master kind")
	}
	row, err := s.findMaster(ctx, kind, id)
	if err != nil {
		return nil, err
	}
	if value := stringValue(input, "name"); value != "" {
		row.Name = value
	}
	if value, ok := mapValue(input, "description"); ok {
		row.Description = fmt.Sprint(value)
	}
	if value, ok := mapValue(input, "swatchMedia"); ok {
		row.SwatchMedia = jsonString(value, "[]")
	}
	if value, ok := mapValue(input, "sellingUnit"); ok {
		row.SellingUnit = fmt.Sprint(value)
	}
	if value, ok := mapValue(input, "baseUnit"); ok {
		row.BaseUnit = fmt.Sprint(value)
	}
	if value, ok := mapValue(input, "quantity"); ok {
		parsed, valid := floatValue(value)
		if !valid || parsed <= 0 {
			return nil, bad("pack quantity must be positive")
		}
		row.Quantity = &parsed
	}
	if _, ok := mapValue(input, "active"); ok {
		row.Active = boolValue(input, row.Active, "active")
	}
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.masterOutput(row), nil
}

func (s *Service) DeactivateMaster(ctx context.Context, kind, id string) (map[string]any, error) {
	return s.UpdateMaster(ctx, kind, id, map[string]any{"active": false})
}

// ProductModel ---------------------------------------------------------------

func (s *Service) validateFixedAttributes(ctx context.Context, input map[string]any) error {
	value, ok := mapValue(input, "fixedAttributes", "fixed_attributes")
	if !ok || value == nil {
		return nil
	}
	fixed := parseStringMap(value)
	for key, rawValue := range fixed {
		if uuidPattern.MatchString(key) {
			definition, err := s.findDefinition(ctx, key)
			if err != nil {
				return err
			}
			if !definition.Active {
				return conflict("inactive attribute definition cannot be assigned")
			}
			if definition.ValueKind == "Reference" {
				if err := s.activeMaster(ctx, strings.ToLower(definition.ReferenceTarget), rawValue); err != nil {
					return err
				}
			}
			continue
		}
		// Public discovery fixtures use explicit materialId/finishId keys.  If
		// those point at a known master, enforce its lifecycle too.
		if strings.EqualFold(key, "materialId") || strings.EqualFold(key, "finishId") || strings.EqualFold(key, "packId") {
			kind := strings.ToLower(strings.TrimSuffix(key, "Id"))
			if err := s.activeMaster(ctx, kind, rawValue); err != nil {
				var apiErr *APIError
				if errors.As(err, &apiErr) && apiErr.Status == 404 {
					continue
				}
				return err
			}
		}
	}
	return nil
}

func (s *Service) validateModelInput(ctx context.Context, input map[string]any, existing *productModelRow) error {
	categoryID := stringValue(input, "categoryId", "category_id")
	if categoryID != "" {
		var category categoryRow
		if err := s.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error; err != nil {
			return notFound("category not found")
		}
		if !category.Active && (existing == nil || existing.CategoryID != categoryID) {
			return conflict("inactive category cannot be assigned to a new ProductModel")
		}
	}
	if err := s.validateFixedAttributes(ctx, input); err != nil {
		return err
	}
	if fixedPackID := stringValue(input, "fixedPackId", "fixed_pack_id"); fixedPackID != "" {
		if err := s.activeMaster(ctx, "pack", fixedPackID); err != nil {
			return err
		}
	}
	if value, ok := mapValue(input, "measurements"); ok {
		if err := validateMeasurements(value); err != nil {
			return err
		}
	}
	if _, err := warrantyValue(input, "warrantySummary", "warranty"); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateModel(ctx context.Context, input map[string]any) (map[string]any, error) {
	name, categoryID := stringValue(input, "name"), stringValue(input, "categoryId", "category_id")
	if name == "" || categoryID == "" {
		return nil, bad("ProductModel name and categoryId are required")
	}
	if err := s.validateModelInput(ctx, input, nil); err != nil {
		return nil, err
	}
	warranty, err := warrantyValue(input, "warrantySummary", "warranty")
	if err != nil {
		return nil, err
	}
	fixedPackID := optionalString(input, "fixedPackId", "fixed_pack_id")
	fixedAttributes := map[string]any{}
	if value, ok := mapValue(input, "fixedAttributes", "fixed_attributes"); ok {
		fixedAttributes = recordMap(value)
	}
	measurements := map[string]any{}
	if value, ok := mapValue(input, "measurements"); ok {
		measurements = recordMap(value)
	}
	row := productModelRow{ID: uuid.NewString(), Name: name, CategoryID: categoryID, Description: stringValue(input, "description"), WarrantySummary: warranty, FixedAttributes: jsonString(fixedAttributes, "{}"), FixedPackID: fixedPackID, Measurements: jsonString(measurements, "{}"), Status: "Draft", CreatedAt: now(), UpdatedAt: now()}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) ListModels(ctx context.Context) ([]map[string]any, error) {
	var rows []productModelRow
	if err := s.db.WithContext(ctx).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		value, err := s.modelOutput(ctx, row, false)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func (s *Service) GetModel(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) UpdateModel(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "Discontinued" && (stringValue(input, "status") != "Discontinued" || boolValue(input, false, "is_active")) {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	if err := s.validateModelInput(ctx, input, &row); err != nil {
		return nil, err
	}
	if value := stringValue(input, "name"); value != "" {
		row.Name = value
	}
	if value := stringValue(input, "categoryId", "category_id"); value != "" {
		row.CategoryID = value
	}
	if value, ok := mapValue(input, "description"); ok {
		row.Description = fmt.Sprint(value)
	}
	if value, ok := mapValue(input, "fixedAttributes", "fixed_attributes"); ok {
		row.FixedAttributes = jsonString(value, "{}")
	}
	if value, ok := mapValue(input, "fixedPackId", "fixed_pack_id"); ok {
		if value == nil || strings.TrimSpace(fmt.Sprint(value)) == "" {
			row.FixedPackID = nil
		} else {
			packID := strings.TrimSpace(fmt.Sprint(value))
			row.FixedPackID = &packID
		}
	}
	if value, ok := mapValue(input, "measurements"); ok {
		row.Measurements = jsonString(value, "{}")
	}
	if value, ok := mapValue(input, "warrantySummary", "warranty"); ok {
		warranty, warrantyErr := warrantyValue(map[string]any{"warrantySummary": value}, "warrantySummary")
		if warrantyErr != nil {
			return nil, warrantyErr
		}
		row.WarrantySummary = warranty
	}
	if status := stringValue(input, "status"); status != "" {
		if status != "Draft" && status != "Active" && status != "Inactive" && status != "Discontinued" {
			return nil, bad("unsupported ProductModel status")
		}
		if status == "Discontinued" && row.Status == "Discontinued" {
			return nil, conflict("Discontinued ProductModel is terminal")
		}
		row.Status = status
	}
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) PublishModel(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel cannot be published")
	}
	var primaryCount int64
	if err := s.db.WithContext(ctx).Model(&productImageRow{}).Where("model_id = ? AND primary_image = ?", id, true).Count(&primaryCount).Error; err != nil {
		return nil, err
	}
	if primaryCount != 1 {
		return nil, conflict("ProductModel requires exactly one primary model image")
	}
	var variants []variantRow
	if err := s.db.WithContext(ctx).Where("model_id = ?", id).Find(&variants).Error; err != nil {
		return nil, err
	}
	saleReadyCount := 0
	for _, variant := range variants {
		if saleReady(variant) {
			saleReadyCount++
		}
	}
	if saleReadyCount == 0 {
		return nil, conflict("ProductModel requires a sale-ready Variant")
	}
	row.Status, row.UpdatedAt = "Active", now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) UnpublishModel(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	row.Status, row.UpdatedAt = "Inactive", now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) DiscontinueModel(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel is terminal")
	}
	row.Status, row.UpdatedAt = "Discontinued", now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

func (s *Service) DeleteModel(context.Context, string) (map[string]any, error) {
	return nil, methodNotAllowed("ProductModel deletion is not supported; use lifecycle transitions")
}

func (s *Service) ReplaceMedia(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	row, err := s.findModel(ctx, modelID)
	if err != nil {
		return nil, err
	}
	value, ok := mapValue(input, "images")
	if !ok {
		return nil, bad("images are required")
	}
	items, ok := value.([]any)
	if !ok {
		return nil, bad("images must be an array")
	}
	images := make([]productImageRow, 0, len(items))
	primaryCount := 0
	for _, raw := range items {
		image, ok := raw.(map[string]any)
		if !ok {
			return nil, bad("invalid ProductModel image")
		}
		url := stringValue(image, "url")
		if url == "" {
			return nil, bad("image url is required")
		}
		primary := boolValue(image, false, "primary")
		if primary {
			primaryCount++
		}
		images = append(images, productImageRow{ID: uuid.NewString(), ModelID: modelID, URL: url, Ordering: intValue(image, 0, "ordering"), PrimaryImage: primary, CreatedAt: now()})
	}
	if primaryCount > 1 {
		return nil, conflict("ProductModel may have only one primary image")
	}
	if row.Status == "Active" && primaryCount != 1 {
		return nil, conflict("Active ProductModel requires a primary image")
	}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("model_id = ?", modelID).Delete(&productImageRow{}).Error; err != nil {
			return err
		}
		for _, image := range images {
			if err := tx.Create(&image).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return s.modelOutput(ctx, row, false)
}

// Variant dimensions ---------------------------------------------------------

func (s *Service) loadDimensions(ctx context.Context, modelID string) ([]variantDimensionRow, error) {
	var rows []variantDimensionRow
	if err := s.db.WithContext(ctx).Where("model_id = ?", modelID).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Service) modelHasVariants(ctx context.Context, modelID string) (bool, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&variantRow{}).Where("model_id = ?", modelID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func decodeDimensionValues(raw any) ([]dimensionValue, error) {
	items, ok := raw.([]any)
	if !ok {
		return nil, bad("allowedValues must be an array")
	}
	values := make([]dimensionValue, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		object, ok := item.(map[string]any)
		if !ok {
			return nil, bad("invalid dimension value")
		}
		id, label := stringValue(object, "id"), stringValue(object, "label")
		if id == "" || label == "" {
			return nil, bad("dimension value id and label are required")
		}
		if _, exists := seen[id]; exists {
			return nil, conflict("duplicate dimension value")
		}
		seen[id] = struct{}{}
		values = append(values, dimensionValue{ID: id, Label: label, Active: boolValue(object, true, "active")})
	}
	return values, nil
}

func (s *Service) CreateDimension(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	if _, err := s.findModel(ctx, modelID); err != nil {
		return nil, err
	}
	definitionID := stringValue(input, "definitionId", "definition_id")
	definition, err := s.findDefinition(ctx, definitionID)
	if err != nil {
		return nil, err
	}
	if !definition.Active {
		return nil, conflict("inactive attribute definition cannot become a VariantDimension")
	}
	if used, err := s.modelHasVariants(ctx, modelID); err != nil {
		return nil, err
	} else if used {
		return nil, conflict("VariantDimension structure cannot change after Variant creation")
	}
	var existing int64
	if err := s.db.WithContext(ctx).Model(&variantDimensionRow{}).Where("model_id = ? AND definition_id = ?", modelID, definitionID).Count(&existing).Error; err != nil {
		return nil, err
	}
	if existing > 0 {
		return nil, conflict("VariantDimension already exists")
	}
	model, err := s.findModel(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if _, ok := jsonMap(model.FixedAttributes)[definitionID]; ok {
		return nil, conflict("a fixed attribute cannot also be a VariantDimension")
	}
	var values []dimensionValue
	if raw, ok := mapValue(input, "allowedValues", "allowed_values"); ok {
		values, err = decodeDimensionValues(raw)
		if err != nil {
			return nil, err
		}
	} else {
		enumValues := []enumValue{}
		if err := json.Unmarshal([]byte(definition.EnumValues), &enumValues); err != nil {
			enumValues = []enumValue{}
		}
		for _, value := range enumValues {
			values = append(values, dimensionValue{ID: value.ID, Label: value.Label, Active: value.Active})
		}
	}
	if len(values) == 0 {
		return nil, bad("VariantDimension requires at least one allowed value")
	}
	if definition.ValueKind == "Reference" {
		for _, value := range values {
			if err := s.activeMaster(ctx, strings.ToLower(definition.ReferenceTarget), value.ID); err != nil {
				return nil, err
			}
		}
	}
	row := variantDimensionRow{ID: uuid.NewString(), ModelID: modelID, DefinitionID: definitionID, AllowedValues: jsonString(values, "[]"), CreatedAt: now(), UpdatedAt: now()}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return s.dimensionOutput(ctx, row)
}

func (s *Service) UpdateDimension(ctx context.Context, modelID, dimensionID string, input map[string]any) (map[string]any, error) {
	var row variantDimensionRow
	if err := s.db.WithContext(ctx).Where("id = ? AND model_id = ?", dimensionID, modelID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("VariantDimension not found")
		}
		return nil, err
	}
	hasVariants, err := s.modelHasVariants(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if value, ok := mapValue(input, "definitionId", "definition_id"); ok {
		definitionID := strings.TrimSpace(fmt.Sprint(value))
		if definitionID != row.DefinitionID && hasVariants {
			return nil, conflict("VariantDimension definition cannot be replaced after Variant creation")
		}
		if _, err := s.findDefinition(ctx, definitionID); err != nil {
			return nil, err
		}
		row.DefinitionID = definitionID
	}
	if value, ok := mapValue(input, "allowedValues", "allowed_values"); ok {
		if hasVariants {
			return nil, conflict("VariantDimension values cannot be replaced after Variant creation")
		}
		values, decodeErr := decodeDimensionValues(value)
		if decodeErr != nil {
			return nil, decodeErr
		}
		row.AllowedValues = jsonString(values, "[]")
	}
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.dimensionOutput(ctx, row)
}

func (s *Service) AddDimensionValue(ctx context.Context, modelID, dimensionID string, input map[string]any) (map[string]any, error) {
	var row variantDimensionRow
	if err := s.db.WithContext(ctx).Where("id = ? AND model_id = ?", dimensionID, modelID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("VariantDimension not found")
		}
		return nil, err
	}
	id, label := stringValue(input, "id"), stringValue(input, "label")
	if id == "" || label == "" {
		return nil, bad("dimension value id and label are required")
	}
	values := []dimensionValue{}
	if err := json.Unmarshal([]byte(row.AllowedValues), &values); err != nil {
		values = []dimensionValue{}
	}
	for index := range values {
		if values[index].ID == id {
			values[index].Label, values[index].Active = label, boolValue(input, true, "active")
			row.AllowedValues, row.UpdatedAt = jsonString(values, "[]"), now()
			if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
				return nil, err
			}
			return s.dimensionOutput(ctx, row)
		}
	}
	values = append(values, dimensionValue{ID: id, Label: label, Active: boolValue(input, true, "active")})
	row.AllowedValues, row.UpdatedAt = jsonString(values, "[]"), now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.dimensionOutput(ctx, row)
}

func (s *Service) DeactivateDimensionValue(ctx context.Context, modelID, dimensionID, valueID string) (map[string]any, error) {
	var row variantDimensionRow
	if err := s.db.WithContext(ctx).Where("id = ? AND model_id = ?", dimensionID, modelID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("VariantDimension not found")
		}
		return nil, err
	}
	values := []dimensionValue{}
	if err := json.Unmarshal([]byte(row.AllowedValues), &values); err != nil {
		values = []dimensionValue{}
	}
	for index := range values {
		if values[index].ID == valueID {
			values[index].Active = false
			row.AllowedValues, row.UpdatedAt = jsonString(values, "[]"), now()
			if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
				return nil, err
			}
			return s.dimensionOutput(ctx, row)
		}
	}
	return nil, notFound("Variant value not found")
}

// Variants -------------------------------------------------------------------

func (s *Service) normalizeSelection(ctx context.Context, modelID string, input map[string]string) (map[string]string, string, error) {
	dimensions, err := s.loadDimensions(ctx, modelID)
	if err != nil {
		return nil, "", err
	}
	if len(dimensions) == 0 {
		return nil, "", bad("ProductModel has no VariantDimensions")
	}
	selected := map[string]string{}
	parts := make([]string, 0, len(dimensions))
	for _, dimension := range dimensions {
		definition, definitionErr := s.findDefinition(ctx, dimension.DefinitionID)
		if definitionErr != nil {
			return nil, "", definitionErr
		}
		provided, found := "", false
		for key, value := range input {
			if key == definition.DisplayName || key == definition.Key || strings.EqualFold(key, definition.DisplayName) || strings.EqualFold(key, definition.Key) {
				provided, found = strings.TrimSpace(value), true
				break
			}
		}
		if !found || provided == "" {
			return nil, "", bad("Variant must select one value for every dimension")
		}
		values := []dimensionValue{}
		if err := json.Unmarshal([]byte(dimension.AllowedValues), &values); err != nil {
			values = []dimensionValue{}
		}
		matched := false
		for _, value := range values {
			if !value.Active {
				continue
			}
			if value.ID == provided || strings.EqualFold(normalizeText(value.Label), normalizeText(provided)) || canonicalValue(value.Label) == canonicalValue(provided) {
				matched = true
				if definition.ValueKind == "Reference" {
					selected[definition.DisplayName] = value.ID
				} else {
					selected[definition.DisplayName] = value.Label
				}
				parts = append(parts, dimension.ID+"="+canonicalValue(value.Label))
				break
			}
		}
		if !matched {
			return nil, "", conflict("Variant selection is not an active allowed value")
		}
	}
	if len(input) != len(dimensions) {
		return nil, "", bad("Variant selection contains an unknown dimension")
	}
	sort.Strings(parts)
	return selected, strings.Join(parts, "|"), nil
}

func (s *Service) validatePackAssignment(ctx context.Context, model productModelRow, selected map[string]string, dimensions []variantDimensionRow, requested *string) (*string, error) {
	packID := requested
	if packID == nil && model.FixedPackID != nil {
		packID = model.FixedPackID
	}
	for _, dimension := range dimensions {
		definition, err := s.findDefinition(ctx, dimension.DefinitionID)
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
		if err := s.activeMaster(ctx, "pack", *packID); err != nil {
			return nil, err
		}
	}
	return packID, nil
}

func (s *Service) normalizeSKU(ctx context.Context, value string, currentID string) (string, error) {
	sku := strings.ToLower(strings.TrimSpace(value))
	if sku == "" {
		return "", nil
	}
	var row variantRow
	query := s.db.WithContext(ctx).Where("sku = ?", sku)
	if currentID != "" {
		query = query.Where("id <> ?", currentID)
	}
	if err := query.First(&row).Error; err == nil {
		return "", conflict("SKU is already reserved by another Variant")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	return sku, nil
}

func parsePrice(input map[string]any, names ...string) (*int64, string, bool, error) {
	value, ok := mapValue(input, names...)
	if !ok {
		return nil, "", false, nil
	}
	if value == nil {
		return nil, "", true, bad("sellingPrice must be a positive VND amount")
	}
	object, ok := value.(map[string]any)
	if !ok {
		return nil, "", true, bad("sellingPrice must be an object")
	}
	amountValue, amountOK := mapValue(object, "amount")
	currency := strings.ToUpper(stringValue(object, "currency"))
	amount, valid := floatValue(amountValue)
	if !amountOK || !valid || amount <= 0 || math.Trunc(amount) != amount || currency != "VND" {
		return nil, "", true, bad("sellingPrice must be a positive VND amount")
	}
	integer := int64(amount)
	return &integer, currency, true, nil
}

func appendHistory(raw, action string) string {
	entries := []historyEntry{}
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		entries = []historyEntry{}
	}
	entries = append(entries, historyEntry{Action: action, At: now()})
	return jsonString(entries, "[]")
}

func (s *Service) ListVariants(ctx context.Context, modelID string) ([]map[string]any, error) {
	if _, err := s.findModel(ctx, modelID); err != nil {
		return nil, err
	}
	var rows []variantRow
	if err := s.db.WithContext(ctx).Where("model_id = ?", modelID).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		result = append(result, s.variantOutput(row))
	}
	return result, nil
}

func (s *Service) CreateVariant(ctx context.Context, modelID string, input map[string]any) (map[string]any, error) {
	model, err := s.findModel(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel cannot receive new Variants")
	}
	var selectedInput map[string]string
	if raw, ok := mapValue(input, "selectedOptions", "selected_options"); ok {
		selectedInput = parseStringMap(raw)
	} else {
		return nil, bad("selectedOptions is required")
	}
	dimensions, err := s.loadDimensions(ctx, modelID)
	if err != nil {
		return nil, err
	}
	selected, canonical, err := s.normalizeSelection(ctx, modelID, selectedInput)
	if err != nil {
		return nil, err
	}
	var existing variantRow
	if err := s.db.WithContext(ctx).Where("model_id = ? AND canonical_combination = ?", modelID, canonical).First(&existing).Error; err == nil {
		return nil, conflict("canonical Variant combination already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	sku, err := s.normalizeSKU(ctx, stringValue(input, "sku"), "")
	if err != nil {
		return nil, err
	}
	amount, currency, hasPrice, err := parsePrice(input, "sellingPrice", "selling_price")
	if err != nil {
		return nil, err
	}
	var requestedPack *string
	if value, ok := mapValue(input, "packId", "pack_id"); ok && value != nil && strings.TrimSpace(fmt.Sprint(value)) != "" {
		pack := strings.TrimSpace(fmt.Sprint(value))
		requestedPack = &pack
	}
	packID, err := s.validatePackAssignment(ctx, model, selected, dimensions, requestedPack)
	if err != nil {
		return nil, err
	}
	technicalValues := map[string]any{}
	if value, ok := mapValue(input, "technicalValues", "technical_values", "variantAttributes", "variant_attributes"); ok {
		technicalValues = recordMap(value)
	}
	row := variantRow{ID: uuid.NewString(), ModelID: modelID, SelectedOptions: jsonString(selected, "{}"), TechnicalValues: jsonString(technicalValues, "{}"), SKU: sku, SellingAmount: amount, SellingCurrency: currency, PackID: packID, Status: "Active", CanonicalCombination: canonical, History: appendHistory("[]", "created"), CreatedAt: now(), UpdatedAt: now()}
	if !hasPrice {
		row.SellingAmount, row.SellingCurrency = nil, ""
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, conflict("canonical Variant combination or SKU already exists")
		}
		return nil, err
	}
	return s.variantOutput(row), nil
}

func recordMap(value any) map[string]any {
	if typed, ok := value.(map[string]any); ok && typed != nil {
		return typed
	}
	return map[string]any{}
}

func (s *Service) GetVariant(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findVariant(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.variantOutput(row), nil
}

func (s *Service) UpdateVariant(ctx context.Context, id string, input map[string]any) (map[string]any, error) {
	row, err := s.findVariant(ctx, id)
	if err != nil {
		return nil, err
	}
	wasSaleReady := saleReady(row)
	model, err := s.findModel(ctx, row.ModelID)
	if err != nil {
		return nil, err
	}
	if value, ok := mapValue(input, "selectedOptions", "selected_options"); ok {
		selectedInput := parseStringMap(value)
		_, canonical, canonicalErr := s.normalizeSelection(ctx, row.ModelID, selectedInput)
		if canonicalErr != nil {
			return nil, canonicalErr
		}
		if canonical != row.CanonicalCombination {
			return nil, conflict("Variant selected combination is immutable")
		}
	}
	if value, ok := mapValue(input, "sku"); ok {
		sku, skuErr := s.normalizeSKU(ctx, fmt.Sprint(value), id)
		if skuErr != nil {
			return nil, skuErr
		}
		row.SKU = sku
	}
	if amount, currency, hasPrice, priceErr := parsePrice(input, "sellingPrice", "selling_price"); priceErr != nil {
		return nil, priceErr
	} else if hasPrice {
		row.SellingAmount, row.SellingCurrency = amount, currency
	}
	if value, ok := mapValue(input, "technicalValues", "technical_values", "variantAttributes", "variant_attributes"); ok {
		row.TechnicalValues = jsonString(recordMap(value), "{}")
	}
	if value, ok := mapValue(input, "packId", "pack_id"); ok {
		if value == nil || strings.TrimSpace(fmt.Sprint(value)) == "" {
			if model.FixedPackID != nil {
				return nil, conflict("Variant must retain the ProductModel fixed Pack")
			}
			row.PackID = nil
		} else {
			pack := strings.TrimSpace(fmt.Sprint(value))
			if model.FixedPackID != nil && pack != *model.FixedPackID {
				return nil, conflict("Variant Pack reference conflicts with ProductModel Pack")
			}
			if packErr := s.activeMaster(ctx, "pack", pack); packErr != nil {
				return nil, packErr
			}
			row.PackID = &pack
		}
	}
	if model.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel cannot be mutated")
	}
	if model.Status == "Active" && saleReady(row) == false {
		var siblings []variantRow
		if err := s.db.WithContext(ctx).Where("model_id = ? AND id <> ?", row.ModelID, row.ID).Find(&siblings).Error; err != nil {
			return nil, err
		}
		if wasSaleReady {
			remaining := false
			for _, sibling := range siblings {
				if saleReady(sibling) {
					remaining = true
					break
				}
			}
			if !remaining {
				return nil, conflict("Active ProductModel must keep one sale-ready Variant")
			}
		}
	}
	row.History = appendHistory(row.History, "updated")
	row.UpdatedAt = now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.variantOutput(row), nil
}

func (s *Service) ActivateVariant(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findVariant(ctx, id)
	if err != nil {
		return nil, err
	}
	model, err := s.findModel(ctx, row.ModelID)
	if err != nil {
		return nil, err
	}
	if model.Status == "Discontinued" {
		return nil, conflict("Discontinued ProductModel cannot receive activated Variants")
	}
	row.Status, row.History, row.UpdatedAt = "Active", appendHistory(row.History, "activated"), now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.variantOutput(row), nil
}

func (s *Service) InactivateVariant(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findVariant(ctx, id)
	if err != nil {
		return nil, err
	}
	model, err := s.findModel(ctx, row.ModelID)
	if err != nil {
		return nil, err
	}
	if model.Status == "Active" && saleReady(row) {
		var variants []variantRow
		if err := s.db.WithContext(ctx).Where("model_id = ? AND id <> ?", row.ModelID, row.ID).Find(&variants).Error; err != nil {
			return nil, err
		}
		remaining := false
		for _, candidate := range variants {
			if saleReady(candidate) {
				remaining = true
				break
			}
		}
		if !remaining {
			return nil, conflict("Active ProductModel must keep one sale-ready Variant")
		}
	}
	row.Status, row.History, row.UpdatedAt = "Inactive", appendHistory(row.History, "inactivated"), now()
	if err := s.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, err
	}
	return s.variantOutput(row), nil
}

func (s *Service) BulkSetPrice(ctx context.Context, input map[string]any) ([]map[string]any, error) {
	value, ok := mapValue(input, "variantIds", "variant_ids")
	if !ok {
		return nil, bad("variantIds are required")
	}
	items, ok := value.([]any)
	if !ok || len(items) == 0 {
		return nil, bad("variantIds must not be empty")
	}
	amount, currency, hasPrice, err := parsePrice(input, "sellingPrice", "selling_price")
	if err != nil {
		return nil, err
	}
	if !hasPrice {
		return nil, bad("sellingPrice is required")
	}
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, strings.TrimSpace(fmt.Sprint(item)))
	}
	result := []variantRow{}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var row variantRow
			if err := tx.Where("id = ?", id).First(&row).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return conflict("bulk price contains an unknown Variant")
				}
				return err
			}
			row.SellingAmount, row.SellingCurrency, row.History, row.UpdatedAt = amount, currency, appendHistory(row.History, "price_updated"), now()
			if err := tx.Save(&row).Error; err != nil {
				return err
			}
			result = append(result, row)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	output := make([]map[string]any, 0, len(result))
	for _, row := range result {
		output = append(output, s.variantOutput(row))
	}
	return output, nil
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

func publicVariants(ctx context.Context, db *gorm.DB, modelID string) ([]variantRow, error) {
	var rows []variantRow
	if err := db.WithContext(ctx).Where("model_id = ?", modelID).Order("created_at ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]variantRow, 0, len(rows))
	for _, row := range rows {
		if saleReady(row) {
			result = append(result, row)
		}
	}
	return result, nil
}

func priceInRange(variants []variantRow, minPrice, maxPrice *int64) bool {
	for _, variant := range variants {
		if variant.SellingAmount == nil {
			continue
		}
		if minPrice != nil && *variant.SellingAmount < *minPrice {
			continue
		}
		if maxPrice != nil && *variant.SellingAmount > *maxPrice {
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
	var rows []productModelRow
	if err := s.db.WithContext(ctx).Where("status = ?", "Active").Find(&rows).Error; err != nil {
		return nil, err
	}
	type candidate struct {
		model    productModelRow
		variants []variantRow
	}
	candidates := make([]candidate, 0, len(rows))
	for _, row := range rows {
		variants, err := publicVariants(ctx, s.db, row.ID)
		if err != nil {
			return nil, err
		}
		if len(variants) == 0 || (filter.CategoryID != "" && row.CategoryID != filter.CategoryID) {
			continue
		}
		fixed := jsonMap(row.FixedAttributes)
		if !mapMatches(fixed, "materialId", filter.MaterialID) || !mapMatches(fixed, "finishId", filter.FinishID) || !priceInRange(variants, filter.MinPrice, filter.MaxPrice) {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(row.Name+" "+row.Description), strings.ToLower(filter.Search)) {
			continue
		}
		candidates = append(candidates, candidate{model: row, variants: variants})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if filter.Sort == "newest" {
			return candidates[i].model.CreatedAt.After(candidates[j].model.CreatedAt)
		}
		price := func(item candidate) int64 {
			if len(item.variants) == 0 || item.variants[0].SellingAmount == nil {
				return 0
			}
			return *item.variants[0].SellingAmount
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
		output, err := s.modelOutput(ctx, item.model, true)
		if err != nil {
			return nil, err
		}
		items = append(items, output)
	}
	return map[string]any{"items": items, "page": filter.Page, "limit": filter.Limit, "total": total}, nil
}

func (s *Service) GetPublicModel(ctx context.Context, id string) (map[string]any, error) {
	row, err := s.findModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if row.Status != "Active" {
		return nil, notFound("public ProductModel not found")
	}
	variants, err := publicVariants(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if len(variants) == 0 {
		return nil, notFound("public ProductModel not found")
	}
	return s.modelOutput(ctx, row, true)
}

func selectedCompatible(row variantRow, selected map[string]string) bool {
	options := jsonMap(row.SelectedOptions)
	for key, value := range selected {
		found := false
		for optionKey, optionValue := range options {
			if strings.EqualFold(optionKey, key) && canonicalValue(fmt.Sprint(optionValue)) == canonicalValue(value) {
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
	model, err := s.findModel(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status != "Active" {
		return nil, notFound("public ProductModel not found")
	}
	variants, err := publicVariants(ctx, s.db, modelID)
	if err != nil {
		return nil, err
	}
	compatible := make([]variantRow, 0, len(variants))
	for _, variant := range variants {
		if selectedCompatible(variant, selected) {
			compatible = append(compatible, variant)
		}
	}
	if len(compatible) == 0 {
		return map[string]any{"options": []any{}}, nil
	}
	dimensions, err := s.loadDimensions(ctx, modelID)
	if err != nil {
		return nil, err
	}
	options := make([]map[string]any, 0, len(dimensions))
	for _, dimension := range dimensions {
		definition, definitionErr := s.findDefinition(ctx, dimension.DefinitionID)
		if definitionErr != nil {
			return nil, definitionErr
		}
		valueSet := map[string]dimensionValue{}
		for _, variant := range compatible {
			selectedOptions := jsonMap(variant.SelectedOptions)
			for key, value := range selectedOptions {
				if !strings.EqualFold(key, definition.DisplayName) {
					continue
				}
				allowed := []dimensionValue{}
				if err := json.Unmarshal([]byte(dimension.AllowedValues), &allowed); err != nil {
					allowed = []dimensionValue{}
				}
				for _, item := range allowed {
					if item.ID == fmt.Sprint(value) || canonicalValue(item.Label) == canonicalValue(fmt.Sprint(value)) {
						valueSet[item.ID] = item
					}
				}
			}
		}
		values := make([]dimensionValue, 0, len(valueSet))
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
	model, err := s.findModel(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if model.Status != "Active" {
		return nil, notFound("public Variant not found")
	}
	variants, err := publicVariants(ctx, s.db, modelID)
	if err != nil {
		return nil, err
	}
	_, canonical, err := s.normalizeSelection(ctx, modelID, selected)
	if err != nil {
		// A representation such as 20 cm may not be in the allowed-value set,
		// so compare its canonical representation directly with public rows.
		dimensions, dimensionErr := s.loadDimensions(ctx, modelID)
		if dimensionErr != nil {
			return nil, dimensionErr
		}
		parts := make([]string, 0, len(dimensions))
		for _, dimension := range dimensions {
			definition, definitionErr := s.findDefinition(ctx, dimension.DefinitionID)
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
			return s.variantOutput(variant), nil
		}
	}
	return nil, notFound("public Variant not found")
}
