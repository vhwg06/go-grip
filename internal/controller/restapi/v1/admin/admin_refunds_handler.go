package admin

import (
	"context"
	"strconv"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// AdminListRefunds handles GET /admin/refunds
func (h *Handler) AdminListRefunds(ctx context.Context, request openapi.AdminListRefundsRequestObject) (openapi.AdminListRefundsResponseObject, error) {
	actor := getActor(ctx)

	status := ""
	if request.Params.Status != nil {
		status = string(*request.Params.Status)
	}

	refunds, err := h.adminUC.ListRefunds(ctx, actor, status)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListRefunds401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListRefunds403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListRefunds500JSONResponse{}, nil
		}
	}

	items := make([]openapi.AdminRefundResponse, 0, len(refunds))
	for _, r := range refunds {
		items = append(items, toAdminRefundResponse(r))
	}
	total := len(items)
	resp := openapi.AdminRefundListResponse{Items: &items, Total: &total}
	return openapi.AdminListRefunds200JSONResponse(resp), nil
}

// AdminApproveRefund handles POST /admin/refunds/{refundId}/approve
func (h *Handler) AdminApproveRefund(ctx context.Context, request openapi.AdminApproveRefundRequestObject) (openapi.AdminApproveRefundResponseObject, error) {
	actor := getActor(ctx)

	refundID, err := strconv.ParseInt(request.RefundId, 10, 64)
	if err != nil {
		_, errResp := mapAdminError(entity.ErrNotFound)
		return openapi.AdminApproveRefund404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
	}

	note := ""
	if request.Body != nil && request.Body.Note != nil {
		note = *request.Body.Note
	}

	refund, err := h.adminUC.ProcessRefund(ctx, actor, refundID, true, note)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminApproveRefund401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminApproveRefund403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminApproveRefund404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminApproveRefund500JSONResponse{}, nil
		}
	}

	resp := toAdminRefundResponse(refund)
	return openapi.AdminApproveRefund200JSONResponse(resp), nil
}

// AdminRejectRefund handles POST /admin/refunds/{refundId}/reject
func (h *Handler) AdminRejectRefund(ctx context.Context, request openapi.AdminRejectRefundRequestObject) (openapi.AdminRejectRefundResponseObject, error) {
	actor := getActor(ctx)

	refundID, err := strconv.ParseInt(request.RefundId, 10, 64)
	if err != nil {
		_, errResp := mapAdminError(entity.ErrNotFound)
		return openapi.AdminRejectRefund404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
	}

	note := ""
	if request.Body != nil && request.Body.Note != nil {
		note = *request.Body.Note
	}

	refund, err := h.adminUC.ProcessRefund(ctx, actor, refundID, false, note)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminRejectRefund401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminRejectRefund403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminRejectRefund404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminRejectRefund500JSONResponse{}, nil
		}
	}

	resp := toAdminRefundResponse(refund)
	return openapi.AdminRejectRefund200JSONResponse(resp), nil
}

// toAdminRefundResponse maps entity.RefundRequest to openapi.AdminRefundResponse.
func toAdminRefundResponse(r entity.RefundRequest) openapi.AdminRefundResponse {
	id := r.ID
	orderID := r.OrderID
	userID := r.UserID
	username := r.Username
	reason := r.Reason
	status := string(r.Status)
	adminUsername := r.AdminUsername
	adminNote := r.AdminNote
	productName := r.ProductName
	amount := int64(r.Amount)
	tradeNo := r.TradeNo
	orderStatus := r.OrderStatus

	return openapi.AdminRefundResponse{
		Id:            &id,
		OrderId:       &orderID,
		UserId:        &userID,
		Username:      &username,
		Reason:        &reason,
		Status:        &status,
		AdminUsername: &adminUsername,
		AdminNote:     &adminNote,
		ProductName:   &productName,
		Amount:        &amount,
		TradeNo:       &tradeNo,
		OrderStatus:   &orderStatus,
		ProcessedAt:   r.ProcessedAt,
		CreatedAt:     &r.CreatedAt,
	}
}
