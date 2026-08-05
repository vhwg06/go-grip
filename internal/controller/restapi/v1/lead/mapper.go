package lead

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
)

// toLeadResponse maps leadmodule.LeadSubmission to openapi.LeadResponse DTO.
func toLeadResponse(l leadmodule.LeadSubmission) openapi.LeadResponse {
	phone := l.CustomerPhone
	email := l.CustomerEmail
	msg := l.Message

	return openapi.LeadResponse{
		Id:        l.ID,
		Name:      l.CustomerName,
		Email:     email,
		Phone:     &phone,
		Message:   &msg,
		CreatedAt: &l.CreatedAt,
	}
}

// toLeadListResponse maps []leadmodule.LeadSubmission to openapi.LeadListResponse DTO.
func toLeadListResponse(leads []leadmodule.LeadSubmission, total int) openapi.LeadListResponse {
	items := make([]openapi.LeadResponse, len(leads))
	for i, l := range leads {
		items[i] = toLeadResponse(l)
	}

	return openapi.LeadListResponse{
		Items: items,
		Total: total,
	}
}
