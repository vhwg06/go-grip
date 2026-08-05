package auth

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Auth capability.
type Handler struct {
	authUC usecase.Auth
	userUC usecase.User
	logger logger.Interface
}

// NewHandler constructs a new Auth vertical handler instance.
func NewHandler(authUC usecase.Auth, userUC usecase.User, l logger.Interface) *Handler {
	return &Handler{
		authUC: authUC,
		userUC: userUC,
		logger: l,
	}
}

// RegisterUser handles POST /auth/register
func (h *Handler) RegisterUser(ctx context.Context, request openapi.RegisterUserRequestObject) (openapi.RegisterUserResponseObject, error) {
	if request.Body == nil {
		return openapi.RegisterUser400JSONResponse{}, nil
	}

	user, err := h.userUC.Register(ctx, request.Body.Username, request.Body.Email, request.Body.Password)
	if err != nil {
		status, errResp := mapAuthError(err)
		switch status {
		case 400:
			return openapi.RegisterUser400JSONResponse{}, nil
		case 409:
			return openapi.RegisterUser409JSONResponse{
				ConflictResponseJSONResponse: openapi.ConflictResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.RegisterUser500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.RegisterUser201JSONResponse(userDTO), nil
}

// LoginUser handles POST /auth/login
func (h *Handler) LoginUser(ctx context.Context, request openapi.LoginUserRequestObject) (openapi.LoginUserResponseObject, error) {
	if request.Body == nil {
		return openapi.LoginUser400JSONResponse{}, nil
	}

	user, accessToken, refreshToken, err := h.authUC.Login(ctx, request.Body.Email, request.Body.Password)
	if err != nil {
		status, errResp := mapAuthError(err)
		switch status {
		case 400:
			return openapi.LoginUser400JSONResponse{}, nil
		case 401:
			return openapi.LoginUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.LoginUser500JSONResponse{}, nil
		}
	}

	tokenPairDTO := toTokenPairResponse(accessToken, refreshToken, user)
	return openapi.LoginUser200JSONResponse(tokenPairDTO), nil
}

// RefreshToken handles POST /auth/refresh
func (h *Handler) RefreshToken(ctx context.Context, request openapi.RefreshTokenRequestObject) (openapi.RefreshTokenResponseObject, error) {
	if request.Body == nil {
		return openapi.RefreshToken400JSONResponse{}, nil
	}

	var token string
	if request.Body.RefreshToken != nil {
		token = *request.Body.RefreshToken
	}
	if token == "" && request.Body.LegacyRefreshToken != nil {
		token = *request.Body.LegacyRefreshToken
	}

	if token == "" {
		return openapi.RefreshToken400JSONResponse{}, nil
	}

	accessToken, refreshToken, err := h.authUC.Refresh(ctx, token)
	if err != nil {
		status, errResp := mapAuthError(err)
		switch status {
		case 400:
			return openapi.RefreshToken400JSONResponse{}, nil
		case 401:
			return openapi.RefreshToken401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.RefreshToken500JSONResponse{}, nil
		}
	}

	tokenPairDTO := toTokenPairResponse(accessToken, refreshToken, entity.User{})
	return openapi.RefreshToken200JSONResponse(tokenPairDTO), nil
}

// LogoutUser handles POST /auth/logout
func (h *Handler) LogoutUser(ctx context.Context, request openapi.LogoutUserRequestObject) (openapi.LogoutUserResponseObject, error) {
	if request.Body == nil {
		return openapi.LogoutUser400JSONResponse{}, nil
	}

	var token string
	if request.Body.RefreshToken != nil {
		token = *request.Body.RefreshToken
	}
	if token == "" && request.Body.LegacyRefreshToken != nil {
		token = *request.Body.LegacyRefreshToken
	}

	var actor entity.Actor
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			actor = a
		}
	}

	err := h.authUC.Logout(ctx, actor, token)
	if err != nil {
		status, errResp := mapAuthError(err)
		switch status {
		case 400:
			return openapi.LogoutUser400JSONResponse{}, nil
		case 401:
			return openapi.LogoutUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.LogoutUser500JSONResponse{}, nil
		}
	}

	return openapi.LogoutUser204Response{}, nil
}

// GetCurrentUser handles GET /auth/me
func (h *Handler) GetCurrentUser(ctx context.Context, request openapi.GetCurrentUserRequestObject) (openapi.GetCurrentUserResponseObject, error) {
	var actor entity.Actor
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			actor = a
		}
	}

	user, err := h.authUC.Me(ctx, actor)
	if err != nil {
		status, errResp := mapAuthError(err)
		switch status {
		case 401:
			return openapi.GetCurrentUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetCurrentUser500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.GetCurrentUser200JSONResponse(userDTO), nil
}
