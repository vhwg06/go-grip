package profile

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toAccountProfileResponse maps entity.User to openapi.AccountProfileResponse DTO.
func toAccountProfileResponse(u entity.User) openapi.AccountProfileResponse {
	displayName := u.Username
	desktopNotif := false
	return openapi.AccountProfileResponse{
		Id:                           u.ID,
		Username:                     u.Username,
		Email:                        u.Email,
		DisplayName:                  &displayName,
		DesktopNotificationsEnabled: &desktopNotif,
		CreatedAt:                    &u.CreatedAt,
	}
}
