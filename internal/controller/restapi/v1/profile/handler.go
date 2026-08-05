package profile

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Handler implements strict OpenAPI handlers for the Profile capability.
type Handler struct {
	profileUC usecase.Profile
	logger    logger.Interface
}

// NewHandler constructs a new Profile vertical handler instance.
func NewHandler(profileUC usecase.Profile, l logger.Interface) *Handler {
	return &Handler{
		profileUC: profileUC,
		logger:    l,
	}
}

func getActor(ctx context.Context) entity.Actor {
	if val := ctx.Value("actor"); val != nil {
		if a, ok := val.(entity.Actor); ok {
			return a
		}
	}
	return entity.Actor{}
}

// GetAccountProfile handles GET /account/profile
func (h *Handler) GetAccountProfile(ctx context.Context, request openapi.GetAccountProfileRequestObject) (openapi.GetAccountProfileResponseObject, error) {
	actor := getActor(ctx)
	user, err := h.profileUC.Get(ctx, actor)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.GetAccountProfile401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.GetAccountProfile500JSONResponse{}, nil
		}
	}

	profileDTO := toAccountProfileResponse(user)
	return openapi.GetAccountProfile200JSONResponse(profileDTO), nil
}

// UpdateAccountProfile handles PUT /account/profile
func (h *Handler) UpdateAccountProfile(ctx context.Context, request openapi.UpdateAccountProfileRequestObject) (openapi.UpdateAccountProfileResponseObject, error) {
	if request.Body == nil {
		return openapi.UpdateAccountProfile400JSONResponse{}, nil
	}

	actor := getActor(ctx)
	email := ""
	if request.Body.Email != nil {
		email = *request.Body.Email
	}
	displayName := ""
	if request.Body.DisplayName != nil {
		displayName = *request.Body.DisplayName
	}
	desktopNotif := false
	if request.Body.DesktopNotificationsEnabled != nil {
		desktopNotif = *request.Body.DesktopNotificationsEnabled
	}

	user, err := h.profileUC.Update(ctx, actor, email, displayName, desktopNotif)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 400:
			return openapi.UpdateAccountProfile400JSONResponse{}, nil
		case 401:
			return openapi.UpdateAccountProfile401JSONResponse{
				UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp),
			}, nil
		default:
			return openapi.UpdateAccountProfile500JSONResponse{}, nil
		}
	}

	profileDTO := toAccountProfileResponse(user)
	return openapi.UpdateAccountProfile200JSONResponse(profileDTO), nil
}
