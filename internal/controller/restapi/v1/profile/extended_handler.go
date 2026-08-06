package profile

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// GetProfile handles GET /profile (alias for account profile).
func (h *Handler) GetProfile(ctx context.Context, _ openapi.GetProfileRequestObject) (openapi.GetProfileResponseObject, error) {
	actor := getActor(ctx)
	user, err := h.profileUC.Get(ctx, actor)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.GetProfile401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.GetProfile500JSONResponse{}, nil
		}
	}
	dto := toAccountProfileResponse(user)
	return openapi.GetProfile200JSONResponse(dto), nil
}

// GetUserProfile handles GET /users/{userId}/profile (admin view of any user profile).
func (h *Handler) GetUserProfile(ctx context.Context, _ openapi.GetUserProfileRequestObject) (openapi.GetUserProfileResponseObject, error) {
	actor := getActor(ctx)
	user, err := h.profileUC.Get(ctx, actor)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.GetUserProfile401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.GetUserProfile500JSONResponse{}, nil
		}
	}
	dto := toAccountProfileResponse(user)
	return openapi.GetUserProfile200JSONResponse(dto), nil
}

// UpdateProfile handles PUT /profile
func (h *Handler) UpdateProfile(ctx context.Context, request openapi.UpdateProfileRequestObject) (openapi.UpdateProfileResponseObject, error) {
	actor := getActor(ctx)

	email := ""
	displayName := ""
	desktopNotif := false

	if request.Body != nil {
		if request.Body.Email != nil {
			email = *request.Body.Email
		}
		if request.Body.DisplayName != nil {
			displayName = *request.Body.DisplayName
		}
		if request.Body.DesktopNotificationsEnabled != nil {
			desktopNotif = *request.Body.DesktopNotificationsEnabled
		}
	}

	user, err := h.profileUC.Update(ctx, actor, email, displayName, desktopNotif)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.UpdateProfile401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.UpdateProfile500JSONResponse{}, nil
		}
	}

	dto := toAccountProfileResponse(user)
	return openapi.UpdateProfile200JSONResponse(dto), nil
}

// UpdateProfileEmail handles PUT /profile/email
func (h *Handler) UpdateProfileEmail(ctx context.Context, request openapi.UpdateProfileEmailRequestObject) (openapi.UpdateProfileEmailResponseObject, error) {
	actor := getActor(ctx)

	email := ""
	if request.Body != nil {
		email = request.Body.Email
	}

	_, err := h.profileUC.Update(ctx, actor, email, "", false)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 400:
			return openapi.UpdateProfileEmail400JSONResponse{BadRequestResponseJSONResponse: openapi.BadRequestResponseJSONResponse(errResp)}, nil
		case 401:
			return openapi.UpdateProfileEmail401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.UpdateProfileEmail500JSONResponse{}, nil
		}
	}

	return openapi.UpdateProfileEmail200Response{}, nil
}

// GetProfileSecurity handles GET /profile/security
func (h *Handler) GetProfileSecurity(ctx context.Context, _ openapi.GetProfileSecurityRequestObject) (openapi.GetProfileSecurityResponseObject, error) {
	actor := getActor(ctx)
	user, err := h.profileUC.Get(ctx, actor)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.GetProfileSecurity401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.GetProfileSecurity500JSONResponse{}, nil
		}
	}

	email := user.Email
	resp := openapi.ProfileSecurityResponse{Email: &email}
	return openapi.GetProfileSecurity200JSONResponse(resp), nil
}

// GetProfileSessions handles GET /profile/sessions
func (h *Handler) GetProfileSessions(ctx context.Context, _ openapi.GetProfileSessionsRequestObject) (openapi.GetProfileSessionsResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		_, errResp := mapProfileError(usermodule.ErrUnauthorized)
		return openapi.GetProfileSessions401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
	}

	// Sessions listing is a read-only operation; return empty list stub until infrastructure is wired.
	items := []openapi.ProfileSessionResponse{}
	resp := openapi.ProfileSessionListResponse{Items: &items}
	return openapi.GetProfileSessions200JSONResponse(resp), nil
}

// UpdateProfileNotifications handles PUT /profile/notifications
func (h *Handler) UpdateProfileNotifications(ctx context.Context, request openapi.UpdateProfileNotificationsRequestObject) (openapi.UpdateProfileNotificationsResponseObject, error) {
	actor := getActor(ctx)

	desktopEnabled := false
	if request.Body != nil {
		if v, ok := (*request.Body)["desktop_notifications_enabled"]; ok {
			if b, ok := v.(bool); ok {
				desktopEnabled = b
			}
		}
	}

	_, err := h.profileUC.Update(ctx, actor, "", "", desktopEnabled)
	if err != nil {
		status, errResp := mapProfileError(err)
		switch status {
		case 401:
			return openapi.UpdateProfileNotifications401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		default:
			return openapi.UpdateProfileNotifications500JSONResponse{}, nil
		}
	}

	return openapi.UpdateProfileNotifications200Response{}, nil
}
