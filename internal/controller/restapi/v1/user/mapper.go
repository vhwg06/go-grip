package user

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toUserResponse maps domain entity.User to openapi.UserResponse DTO.
func toUserResponse(u entity.User) openapi.UserResponse {
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

// toUserListResponse maps a slice of domain entity.User and total count to openapi.UserListResponse.
func toUserListResponse(users []entity.User, total int) openapi.UserListResponse {
	items := make([]openapi.UserResponse, len(users))
	for i, u := range users {
		items[i] = toUserResponse(u)
	}
	return openapi.UserListResponse{
		Items: items,
		Total: total,
	}
}
