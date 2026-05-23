package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

type gripTokenPairResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         entity.User `json:"user,omitempty"`
}

type gripRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// @Summary     Refresh access token
// @Description Rotates refresh token and returns a new token pair
// @ID          grip_auth_refresh
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body gripRefreshRequest true "Refresh token payload"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Router      /auth/refresh [post]
func (r *V1) gripRefresh(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}

	var body gripRefreshRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	accessToken, refreshToken, err := r.authUC.Refresh(ctx.UserContext(), body.RefreshToken)
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(gripTokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

// @Summary     Logout
// @Description Revokes the provided refresh token for the current user
// @ID          grip_auth_logout
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body gripRefreshRequest true "Refresh token payload"
// @Success     204 {string} string ""
// @Failure     400 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /auth/logout [post]
func (r *V1) gripLogout(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}

	var body gripRefreshRequest
	if err := ctx.BodyParser(&body); err != nil {
		status, payload := mapDomainError(entity.ErrInvalidInput)
		return ctx.Status(status).JSON(payload)
	}

	if err := r.authUC.Logout(ctx.UserContext(), r.gripActor(ctx), body.RefreshToken); err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Current user profile
// @Description Returns current authenticated user profile and admin flag
// @ID          grip_auth_me
// @Tags        auth
// @Produce     json
// @Success     200 {object} envelope
// @Failure     401 {object} envelope
// @Failure     500 {object} envelope
// @Security    BearerAuth
// @Router      /auth/me [get]
func (r *V1) gripMe(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}

	user, err := r.authUC.Me(ctx.UserContext(), r.gripActor(ctx))
	if err != nil {
		status, payload := mapDomainError(err)
		return ctx.Status(status).JSON(payload)
	}

	return ctx.JSON(apiSuccessEnvelope(user))
}
