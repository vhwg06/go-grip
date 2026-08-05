package profile

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// toAccountProfileResponse maps usermodule.User to openapi.AccountProfileResponse DTO.
func toAccountProfileResponse(u usermodule.User) openapi.AccountProfileResponse {
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
