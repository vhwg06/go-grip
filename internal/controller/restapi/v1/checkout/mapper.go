package checkout

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
)

// toCheckoutPreviewResponse maps usecase.AmountBreakdown to openapi.CheckoutPreviewResponse DTO.
func toCheckoutPreviewResponse(b usecase.AmountBreakdown) openapi.CheckoutPreviewResponse {
	subtotalInt := int(b.Subtotal)
	totalInt := int(b.FinalPrice)
	zero := 0

	return openapi.CheckoutPreviewResponse{
		Subtotal:    subtotalInt,
		ShippingFee: &zero,
		Tax:         &zero,
		Total:       totalInt,
	}
}

// toCheckoutOrderResponse maps entity.Order to openapi.CheckoutOrderResponse DTO.
func toCheckoutOrderResponse(o entity.Order) openapi.CheckoutOrderResponse {
	totalInt := int(o.Amount)
	statusStr := string(o.Status)
	email := o.Email

	return openapi.CheckoutOrderResponse{
		Id:          o.ID,
		OrderSn:     o.ID,
		Status:      statusStr,
		TotalAmount: totalInt,
		Email:       &email,
		CreatedAt:   &o.CreatedAt,
	}
}

// toPaymentParamsResponse maps usecase.PaymentParams to openapi.PaymentParamsResponse DTO.
func toPaymentParamsResponse(p usecase.PaymentParams) openapi.PaymentParamsResponse {
	var paramsMap *map[string]string
	if p.Fields != nil {
		fieldsMap := make(map[string]string, len(p.Fields))
		for k, v := range p.Fields {
			fieldsMap[k] = v
		}
		paramsMap = &fieldsMap
	}

	return openapi.PaymentParamsResponse{
		PayUrl: p.URL,
		Params: paramsMap,
	}
}
