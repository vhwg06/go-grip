package notification

import (
	"context"
	"strconv"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Notification capability.
type Handler struct {
	notificationUC notificationmodule.NotificationCenterUseCase
	logger         logger.Interface
}

// NewHandler constructs a new Notification vertical handler instance.
func NewHandler(notificationUC notificationmodule.NotificationCenterUseCase, l logger.Interface) *Handler {
	return &Handler{
		notificationUC: notificationUC,
		logger:         l,
	}
}

func getActor(ctx context.Context) usermodule.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(usermodule.Actor); ok {
			return a
		}
	}
	return usermodule.Actor{}
}

// ListNotifications handles GET /notifications
func (h *Handler) ListNotifications(ctx context.Context, request openapi.ListNotificationsRequestObject) (openapi.ListNotificationsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.ListNotifications401JSONResponse{}, nil
	}

	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	pag := pagination.Pagination{Limit: limit, Offset: offset}

	items, total, err := h.notificationUC.Inbox(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, pag)
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

	unread, _ := h.notificationUC.UnreadCount(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin})
	listDTO := toNotificationListResponse(items, total, unread)
	return openapi.ListNotifications200JSONResponse(listDTO), nil
}

// MarkAllNotificationsRead handles POST /notifications/read-all
func (h *Handler) MarkAllNotificationsRead(ctx context.Context, request openapi.MarkAllNotificationsReadRequestObject) (openapi.MarkAllNotificationsReadResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.MarkAllNotificationsRead401JSONResponse{}, nil
	}

	err := h.notificationUC.MarkAllRead(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin})
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

// GetUnreadNotificationCount handles GET /notifications/unread-count
func (h *Handler) GetUnreadNotificationCount(ctx context.Context, request openapi.GetUnreadNotificationCountRequestObject) (openapi.GetUnreadNotificationCountResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetUnreadNotificationCount401JSONResponse{}, nil
	}

	count, err := h.notificationUC.UnreadCount(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin})
	if err != nil {
		return openapi.GetUnreadNotificationCount500JSONResponse{}, nil
	}

	return openapi.GetUnreadNotificationCount200JSONResponse(openapi.UnreadNotificationCountResponse{
		Count: count,
	}), nil
}

// MarkNotificationRead handles POST /notifications/{id}/read
func (h *Handler) MarkNotificationRead(ctx context.Context, request openapi.MarkNotificationReadRequestObject) (openapi.MarkNotificationReadResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.MarkNotificationRead401JSONResponse{}, nil
	}

	notifID, err := strconv.ParseInt(request.Id, 10, 64)
	if err != nil {
		return openapi.MarkNotificationRead400JSONResponse{}, nil
	}

	err = h.notificationUC.MarkRead(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin}, notifID)
	if err != nil {
		status, errResp := mapNotificationError(err)
		switch status {
		case 400:
			return openapi.MarkNotificationRead400JSONResponse{}, nil
		case 401:
			return openapi.MarkNotificationRead401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.MarkNotificationRead404JSONResponse{}, nil
		default:
			return openapi.MarkNotificationRead500JSONResponse{}, nil
		}
	}

	return openapi.MarkNotificationRead200Response{}, nil
}


// ClearNotificationInbox handles DELETE /notifications
func (h *Handler) ClearNotificationInbox(ctx context.Context, _ openapi.ClearNotificationInboxRequestObject) (openapi.ClearNotificationInboxResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.ClearNotificationInbox401Response{}, nil
	}

	err := h.notificationUC.Clear(ctx, notificationmodule.Actor{UserID: actor.UserID, Username: actor.Username, IsAdmin: actor.IsAdmin})
	if err != nil {
		return openapi.ClearNotificationInbox500Response{}, nil
	}

	return openapi.ClearNotificationInbox204Response{}, nil
}

// QueueAdminNotificationTest handles POST /admin/notifications/test
func (h *Handler) QueueAdminNotificationTest(ctx context.Context, _ openapi.QueueAdminNotificationTestRequestObject) (openapi.QueueAdminNotificationTestResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.QueueAdminNotificationTest401JSONResponse{}, nil
	}
	if !actor.IsAdmin {
		return openapi.QueueAdminNotificationTest403JSONResponse{}, nil
	}

	return openapi.QueueAdminNotificationTest200JSONResponse{}, nil
}
