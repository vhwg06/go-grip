package user

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// toUserResponse maps usermodule.User to openapi.UserResponse DTO.
func toUserResponse(u usermodule.User) openapi.UserResponse {
	role := string(u.Role)
	return openapi.UserResponse{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      &role,
		IsAdmin:   &u.IsAdmin,
		IsBlocked: &u.IsBlocked,
		CreatedAt: &u.CreatedAt,
	}
}

// toUserListResponse maps a slice of usermodule.User and total count to openapi.UserListResponse.
func toUserListResponse(users []usermodule.User, total int) openapi.UserListResponse {
	items := make([]openapi.UserResponse, len(users))
	for i, u := range users {
		items[i] = toUserResponse(u)
	}
	return openapi.UserListResponse{
		Items: items,
		Total: total,
	}
}
