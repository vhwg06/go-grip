package admin

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Admin capability.
type Handler struct {
	adminUC usecase.Admin
	logger  logger.Interface
}

// NewHandler constructs a new Admin vertical handler instance.
func NewHandler(adminUC usecase.Admin, l logger.Interface) *Handler {
	return &Handler{
		adminUC: adminUC,
		logger:  l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// GetAdminDashboardStats handles GET /admin/dashboard/stats
func (h *Handler) GetAdminDashboardStats(ctx context.Context, request openapi.GetAdminDashboardStatsRequestObject) (openapi.GetAdminDashboardStatsResponseObject, error) {
	actor := getActor(ctx)

	// Fetch users and orders counts as summary stats
	page := entity.Pagination{Limit: 1, Offset: 0}
	_, userCount, err := h.adminUC.ListUsers(ctx, actor, page)
	if err != nil {
		status, errResp := mapAdminError(err)
		switch status {
		case 401:
			return openapi.GetAdminDashboardStats401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.GetAdminDashboardStats403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetAdminDashboardStats500JSONResponse{}, nil
		}
	}

	_, orderCount, _ := h.adminUC.ListOrders(ctx, actor, page, "", "")

	statsDTO := toAdminDashboardStatsResponse(userCount, orderCount, 0, 0)
	return openapi.GetAdminDashboardStats200JSONResponse(statsDTO), nil
}

// ListAdminAuditLogs handles GET /admin/audit-logs
func (h *Handler) ListAdminAuditLogs(ctx context.Context, request openapi.ListAdminAuditLogsRequestObject) (openapi.ListAdminAuditLogsResponseObject, error) {

	// Return empty audit logs list DTO
	items := []openapi.AdminAuditLogResponse{}
	listDTO := openapi.AdminAuditLogListResponse{
		Items: items,
		Total: 0,
	}

	return openapi.ListAdminAuditLogs200JSONResponse(listDTO), nil
}
