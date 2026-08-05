package notification

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Notification capability.
type Handler struct {
	notificationUC usecase.NotificationCenter
	logger         logger.Interface
}

// NewHandler constructs a new Notification vertical handler instance.
func NewHandler(notificationUC usecase.NotificationCenter, l logger.Interface) *Handler {
	return &Handler{
		notificationUC: notificationUC,
		logger:         l,
	}
}

// ListNotifications handles GET /notifications
func (h *Handler) ListNotifications(ctx context.Context, request openapi.ListNotificationsRequestObject) (openapi.ListNotificationsResponseObject, error) {
	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	actor := entity.Actor{UserID: "usr-1"}
	pag := entity.Pagination{Limit: limit, Offset: offset}

	items, total, err := h.notificationUC.Inbox(ctx, actor, pag)
	if err != nil {
		status, errResp := mapNotificationError(err)
		switch status {
		case 401:
			return openapi.ListNotifications401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ListNotifications500JSONResponse{}, nil
		}
	}

	unread, _ := h.notificationUC.UnreadCount(ctx, actor)
	listDTO := toNotificationListResponse(items, total, unread)
	return openapi.ListNotifications200JSONResponse(listDTO), nil
}

// MarkAllNotificationsRead handles POST /notifications/read-all
func (h *Handler) MarkAllNotificationsRead(ctx context.Context, request openapi.MarkAllNotificationsReadRequestObject) (openapi.MarkAllNotificationsReadResponseObject, error) {
	actor := entity.Actor{UserID: "usr-1"}
	err := h.notificationUC.MarkAllRead(ctx, actor)
	if err != nil {
		status, errResp := mapNotificationError(err)
		switch status {
		case 401:
			return openapi.MarkAllNotificationsRead401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.MarkAllNotificationsRead500JSONResponse{}, nil
		}
	}

	return openapi.MarkAllNotificationsRead200Response{}, nil
}
