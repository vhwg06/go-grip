package notification

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toNotificationResponse maps entity.UserNotification to openapi.NotificationResponse DTO.
func toNotificationResponse(n entity.UserNotification) openapi.NotificationResponse {
	idInt := int(n.ID)
	body := n.ContentKey
	return openapi.NotificationResponse{
		Id:        idInt,
		Title:     n.TitleKey,
		Body:      &body,
		IsRead:    &n.IsRead,
		CreatedAt: &n.CreatedAt,
	}
}

// toNotificationListResponse maps []entity.UserNotification to openapi.NotificationListResponse DTO.
func toNotificationListResponse(notifications []entity.UserNotification, total, unreadCount int) openapi.NotificationListResponse {
	items := make([]openapi.NotificationResponse, len(notifications))
	for i, n := range notifications {
		items[i] = toNotificationResponse(n)
	}

	return openapi.NotificationListResponse{
		Items:       items,
		Total:       total,
		UnreadCount: unreadCount,
	}
}
