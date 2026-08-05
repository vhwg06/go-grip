package auth

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
)

// toUserResponse maps a domain entity.User to an openapi.UserResponse DTO.
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

// toTokenPairResponse maps tokens and user entity to an openapi.TokenPairResponse DTO.
func toTokenPairResponse(accessToken, refreshToken string, user entity.User) openapi.TokenPairResponse {
	userDTO := toUserResponse(user)
	return openapi.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &userDTO,
	}
}
