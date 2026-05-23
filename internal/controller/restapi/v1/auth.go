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

// @Summary     Begin LinuxDO OAuth
// @Description Generates LinuxDO OAuth authorization URL
// @ID          grip_auth_oauth_linuxdo
// @Tags        auth
// @Produce     json
// @Success     200 {object} envelope
// @Failure     500 {object} envelope
// @Router      /auth/oauth/linuxdo [get]
func (r *V1) gripBeginLinuxDO(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}
	url, err := r.authUC.BeginLinuxDO(ctx.UserContext())
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"url": url}))
}

// @Summary     Begin GitHub OAuth
// @Description Generates GitHub OAuth authorization URL
// @ID          grip_auth_oauth_github
// @Tags        auth
// @Produce     json
// @Success     200 {object} envelope
// @Failure     500 {object} envelope
// @Router      /auth/oauth/github [get]
func (r *V1) gripBeginGitHub(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}
	url, err := r.authUC.BeginGitHub(ctx.UserContext())
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	return ctx.JSON(apiSuccessEnvelope(fiber.Map{"url": url}))
}

// @Summary     Complete LinuxDO OAuth callback
// @Description Exchanges callback code and issues access/refresh tokens
// @ID          grip_auth_callback_linuxdo
// @Tags        auth
// @Produce     json
// @Param       code query string true "OAuth code"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Router      /auth/callback/linuxdo [get]
func (r *V1) gripCompleteLinuxDO(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}

	user, accessToken, refreshToken, err := r.authUC.CompleteLinuxDO(ctx.UserContext(), ctx.Query("code"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	return ctx.JSON(apiSuccessEnvelope(gripTokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}))
}

// @Summary     Complete GitHub OAuth callback
// @Description Exchanges callback code and issues access/refresh tokens
// @ID          grip_auth_callback_github
// @Tags        auth
// @Produce     json
// @Param       code query string true "OAuth code"
// @Success     200 {object} envelope
// @Failure     400 {object} envelope
// @Failure     500 {object} envelope
// @Router      /auth/callback/github [get]
func (r *V1) gripCompleteGitHub(ctx *fiber.Ctx) error {
	if r.authUC == nil {
		return ctx.Status(http.StatusInternalServerError).JSON(envelope{Error: "auth_usecase_not_configured"})
	}

	user, accessToken, refreshToken, err := r.authUC.CompleteGitHub(ctx.UserContext(), ctx.Query("code"))
	if err != nil {
		status, body := mapDomainError(err)
		return ctx.Status(status).JSON(body)
	}
	return ctx.JSON(apiSuccessEnvelope(gripTokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}))
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
