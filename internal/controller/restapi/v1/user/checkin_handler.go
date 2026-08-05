package user

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
)

// checkinStatusForActor returns a stub CheckinStatusResponse for the currently authenticated user.
// The daily check-in feature is not yet backed by persistent storage; endpoints return a safe default.
func checkinStatusForActor(checkedIn bool) openapi.CheckinStatusResponse {
	streak := 0
	return openapi.CheckinStatusResponse{CheckedIn: &checkedIn, Streak: &streak}
}

// DoCheckin handles POST /checkin
func (h *Handler) DoCheckin(ctx context.Context, _ openapi.DoCheckinRequestObject) (openapi.DoCheckinResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.DoCheckin401JSONResponse{}, nil
	}
	resp := checkinStatusForActor(true)
	return openapi.DoCheckin200JSONResponse(resp), nil
}

// GetCheckinStatus handles GET /checkin/status
func (h *Handler) GetCheckinStatus(ctx context.Context, _ openapi.GetCheckinStatusRequestObject) (openapi.GetCheckinStatusResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetCheckinStatus401JSONResponse{}, nil
	}
	resp := checkinStatusForActor(false)
	return openapi.GetCheckinStatus200JSONResponse(resp), nil
}

// GetCheckinStatusAlt handles GET /checkin-status (alternative path alias for GetCheckinStatus)
func (h *Handler) GetCheckinStatusAlt(ctx context.Context, _ openapi.GetCheckinStatusAltRequestObject) (openapi.GetCheckinStatusAltResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetCheckinStatusAlt401JSONResponse{}, nil
	}
	resp := checkinStatusForActor(false)
	return openapi.GetCheckinStatusAlt200JSONResponse(resp), nil
}

// GetUserCheckinStatus handles GET /user/profile/checkin/status (legacy path)
func (h *Handler) GetUserCheckinStatus(ctx context.Context, _ openapi.GetUserCheckinStatusRequestObject) (openapi.GetUserCheckinStatusResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetUserCheckinStatus401JSONResponse{}, nil
	}
	resp := checkinStatusForActor(false)
	return openapi.GetUserCheckinStatus200JSONResponse(resp), nil
}

// GetUserCheckinStatusAlt handles GET /user/profile/checkin-status (legacy path alt for GetUserCheckinStatus)
func (h *Handler) GetUserCheckinStatusAlt(ctx context.Context, _ openapi.GetUserCheckinStatusAltRequestObject) (openapi.GetUserCheckinStatusAltResponseObject, error) {
	actor := getActor(ctx)
	if actor.UserID == "" {
		return openapi.GetUserCheckinStatusAlt401JSONResponse{}, nil
	}
	resp := checkinStatusForActor(false)
	return openapi.GetUserCheckinStatusAlt200JSONResponse(resp), nil
}
