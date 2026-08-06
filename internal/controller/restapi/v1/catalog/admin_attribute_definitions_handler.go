package catalog

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// AdminListAttributeDefinitions handles GET /admin/catalog/attribute-definitions
func (h *Handler) AdminListAttributeDefinitions(ctx context.Context, _ openapi.AdminListAttributeDefinitionsRequestObject) (openapi.AdminListAttributeDefinitionsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminListAttributeDefinitions401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminListAttributeDefinitions403JSONResponse{}, nil
	}

	_, err := h.catalogBase.ListDefinitions(ctx)
	if err != nil {
		return openapi.AdminListAttributeDefinitions500JSONResponse{}, nil
	}

	return openapi.AdminListAttributeDefinitions200Response{}, nil
}

// AdminCreateAttributeDefinition handles POST /admin/catalog/attribute-definitions
func (h *Handler) AdminCreateAttributeDefinition(ctx context.Context, request openapi.AdminCreateAttributeDefinitionRequestObject) (openapi.AdminCreateAttributeDefinitionResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminCreateAttributeDefinition401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminCreateAttributeDefinition403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.CreateDefinition(ctx, input)
	if err != nil {
		return openapi.AdminCreateAttributeDefinition400JSONResponse{}, nil
	}

	return openapi.AdminCreateAttributeDefinition201Response{}, nil
}

// AdminUpdateAttributeDefinition handles PATCH /admin/catalog/attribute-definitions/{definitionId}
func (h *Handler) AdminUpdateAttributeDefinition(ctx context.Context, request openapi.AdminUpdateAttributeDefinitionRequestObject) (openapi.AdminUpdateAttributeDefinitionResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminUpdateAttributeDefinition401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminUpdateAttributeDefinition403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.UpdateDefinition(ctx, request.DefinitionId, input)
	if err != nil {
		return openapi.AdminUpdateAttributeDefinition400JSONResponse{}, nil
	}

	return openapi.AdminUpdateAttributeDefinition200Response{}, nil
}

// AdminDeactivateAttributeDefinition handles POST /admin/catalog/attribute-definitions/{definitionId}/deactivate
func (h *Handler) AdminDeactivateAttributeDefinition(ctx context.Context, request openapi.AdminDeactivateAttributeDefinitionRequestObject) (openapi.AdminDeactivateAttributeDefinitionResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeactivateAttributeDefinition401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeactivateAttributeDefinition403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeactivateDefinition(ctx, request.DefinitionId)
	if err != nil {
		return openapi.AdminDeactivateAttributeDefinition404JSONResponse{}, nil
	}

	return openapi.AdminDeactivateAttributeDefinition200Response{}, nil
}

// AdminAddAttributeEnumValue handles POST /admin/catalog/attribute-definitions/{definitionId}/enum-values
func (h *Handler) AdminAddAttributeEnumValue(ctx context.Context, request openapi.AdminAddAttributeEnumValueRequestObject) (openapi.AdminAddAttributeEnumValueResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminAddAttributeEnumValue401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminAddAttributeEnumValue403JSONResponse{}, nil
	}

	input := map[string]any{}
	if request.Body != nil {
		input = map[string]any(*request.Body)
	}

	_, err := h.catalogBase.AddEnumValue(ctx, request.DefinitionId, input)
	if err != nil {
		return openapi.AdminAddAttributeEnumValue400JSONResponse{}, nil
	}

	return openapi.AdminAddAttributeEnumValue201Response{}, nil
}

// AdminDeactivateAttributeEnumValue handles POST /admin/catalog/attribute-definitions/{definitionId}/enum-values/{enumValueId}/deactivate
func (h *Handler) AdminDeactivateAttributeEnumValue(ctx context.Context, request openapi.AdminDeactivateAttributeEnumValueRequestObject) (openapi.AdminDeactivateAttributeEnumValueResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.AdminDeactivateAttributeEnumValue401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.AdminDeactivateAttributeEnumValue403JSONResponse{}, nil
	}

	_, err := h.catalogBase.DeactivateEnumValue(ctx, request.DefinitionId, request.EnumValueId)
	if err != nil {
		return openapi.AdminDeactivateAttributeEnumValue404JSONResponse{}, nil
	}

	return openapi.AdminDeactivateAttributeEnumValue200Response{}, nil
}
