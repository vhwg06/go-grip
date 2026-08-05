package auth

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// toUserResponse maps a domain usermodule.User to an openapi.UserResponse DTO.
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

// toTokenPairResponse maps tokens and user entity to an openapi.TokenPairResponse DTO.
func toTokenPairResponse(accessToken, refreshToken string, u usermodule.User) openapi.TokenPairResponse {
	userDTO := toUserResponse(u)
	return openapi.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &userDTO,
	}
}
