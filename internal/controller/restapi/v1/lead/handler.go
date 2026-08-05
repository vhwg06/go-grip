package lead

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Lead capability.
type Handler struct {
	leadUC leadmodule.LeadUseCase
	logger logger.Interface
}

// NewHandler constructs a new Lead vertical handler instance.
func NewHandler(leadUC leadmodule.LeadUseCase, l logger.Interface) *Handler {
	return &Handler{
		leadUC: leadUC,
		logger: l,
	}
}

// SubmitLead handles POST /leads
func (h *Handler) SubmitLead(ctx context.Context, request openapi.SubmitLeadRequestObject) (openapi.SubmitLeadResponseObject, error) {
	if request.Body == nil {
		return openapi.SubmitLead400JSONResponse{}, nil
	}

	phone := ""
	if request.Body.Phone != nil {
		phone = *request.Body.Phone
	}
	msg := ""
	if request.Body.Message != nil {
		msg = *request.Body.Message
	}

	sub := leadmodule.LeadSubmission{
		CustomerName:  request.Body.Name,
		CustomerEmail: request.Body.Email,
		CustomerPhone: phone,
		Message:       msg,
	}

	res, err := h.leadUC.Submit(ctx, sub)
	if err != nil {
		status, _ := mapLeadError(err)
		switch status {
		case 400:
			return openapi.SubmitLead400JSONResponse{}, nil
		default:
			return openapi.SubmitLead500JSONResponse{}, nil
		}
	}

	leadDTO := toLeadResponse(res)
	return openapi.SubmitLead201JSONResponse(leadDTO), nil
}

// ListLeads handles GET /leads
func (h *Handler) ListLeads(ctx context.Context, request openapi.ListLeadsRequestObject) (openapi.ListLeadsResponseObject, error) {
	leadItem, err := h.leadUC.Get(ctx, "lead-1")
	if err != nil {
		status, errResp := mapLeadError(err)
		switch status {
		case 401:
			return openapi.ListLeads401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.ListLeads403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ListLeads500JSONResponse{}, nil
		}
	}

	listDTO := toLeadListResponse([]leadmodule.LeadSubmission{leadItem}, 1)
	return openapi.ListLeads200JSONResponse(listDTO), nil
}
