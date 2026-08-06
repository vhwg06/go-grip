package user

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the User capability.
type Handler struct {
	userUC usermodule.UserUseCase
	logger logger.Interface
}

// NewHandler constructs a new User vertical handler instance.
func NewHandler(userUC usermodule.UserUseCase, l logger.Interface) *Handler {
	return &Handler{
		userUC: userUC,
		logger: l,
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

// GetMyProfile handles GET /users/profile
func (h *Handler) GetMyProfile(ctx context.Context, request openapi.GetMyProfileRequestObject) (openapi.GetMyProfileResponseObject, error) {
	actor := getActor(ctx)
	user, err := h.userUC.GetUser(ctx, actor.UserID)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 401:
			return openapi.GetMyProfile401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetMyProfile500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.GetMyProfile200JSONResponse(userDTO), nil
}

// UpdateMyProfile handles PUT /users/profile
func (h *Handler) UpdateMyProfile(ctx context.Context, request openapi.UpdateMyProfileRequestObject) (openapi.UpdateMyProfileResponseObject, error) {
	if request.Body == nil {
		return openapi.UpdateMyProfile400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	user, err := h.userUC.UpdateProfile(ctx, actor.UserID, request.Body.DisplayName)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 400:
			return openapi.UpdateMyProfile400JSONResponse{}, nil
		case 401:
			return openapi.UpdateMyProfile401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UpdateMyProfile500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.UpdateMyProfile200JSONResponse(userDTO), nil
}

// ListUsers handles GET /users
func (h *Handler) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	limit := 10
	if request.Params.Limit != nil && *request.Params.Limit > 0 {
		limit = *request.Params.Limit
	}
	offset := 0
	if request.Params.Offset != nil && *request.Params.Offset >= 0 {
		offset = *request.Params.Offset
	}

	users, total, err := h.userUC.List(ctx, limit, offset)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 401:
			return openapi.ListUsers401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.ListUsers403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.ListUsers500JSONResponse{}, nil
		}
	}

	listDTO := toUserListResponse(users, total)
	return openapi.ListUsers200JSONResponse(listDTO), nil
}

// CreateAdminUser handles POST /users/admin
func (h *Handler) CreateAdminUser(ctx context.Context, request openapi.CreateAdminUserRequestObject) (openapi.CreateAdminUserResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateAdminUser400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	role := usermodule.RoleAdministrator
	if request.Body.Role != nil && *request.Body.Role != "" {
		role = usermodule.RoleName(*request.Body.Role)
	}

	user, err := h.userUC.CreateAdminUser(ctx, actor.UserID, request.Body.Username, request.Body.Email, request.Body.Password, role)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 400:
			return openapi.CreateAdminUser400JSONResponse{}, nil
		case 401:
			return openapi.CreateAdminUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.CreateAdminUser403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 409:
			return openapi.CreateAdminUser409JSONResponse{
				ConflictResponseJSONResponse: openapi.ConflictResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.CreateAdminUser500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.CreateAdminUser201JSONResponse(userDTO), nil
}

// GetUserByID handles GET /users/{id}
func (h *Handler) GetUserByID(ctx context.Context, request openapi.GetUserByIDRequestObject) (openapi.GetUserByIDResponseObject, error) {
	user, err := h.userUC.GetUser(ctx, request.Id)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 401:
			return openapi.GetUserByID401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.GetUserByID403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.GetUserByID404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetUserByID500JSONResponse{}, nil
		}
	}

	userDTO := toUserResponse(user)
	return openapi.GetUserByID200JSONResponse(userDTO), nil
}

// LockUser handles POST /users/{id}/lock
func (h *Handler) LockUser(ctx context.Context, request openapi.LockUserRequestObject) (openapi.LockUserResponseObject, error) {
	actor := getActor(ctx)
	err := h.userUC.Lock(ctx, actor.UserID, request.Id)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 401:
			return openapi.LockUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.LockUser403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.LockUser404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.LockUser500JSONResponse{}, nil
		}
	}

	return openapi.LockUser200Response{}, nil
}

// UnlockUser handles POST /users/{id}/unlock
func (h *Handler) UnlockUser(ctx context.Context, request openapi.UnlockUserRequestObject) (openapi.UnlockUserResponseObject, error) {
	actor := getActor(ctx)
	err := h.userUC.Unlock(ctx, actor.UserID, request.Id)
	if err != nil {
		status, errResp := mapUserError(err)
		switch status {
		case 401:
			return openapi.UnlockUser401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		case 403:
			return openapi.UnlockUser403JSONResponse{
				ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp),
			}, nil
		case 404:
			return openapi.UnlockUser404JSONResponse{
				NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UnlockUser500JSONResponse{}, nil
		}
	}

	return openapi.UnlockUser200Response{}, nil
}
