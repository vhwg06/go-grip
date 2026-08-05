package admin

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// AdminListUsers handles GET /admin/users (backoffice listing with stats)
func (h *Handler) AdminListUsers(ctx context.Context, request openapi.AdminListUsersRequestObject) (openapi.AdminListUsersResponseObject, error) {
	actor := getActor(ctx)

	page := 1
	pageSize := 20
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}
	offset := (page - 1) * pageSize
	pag := pagination.New(pageSize, offset)

	users, total, err := h.adminUC.ListUsers(ctx, actor, pag)
	if err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminListUsers401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminListUsers403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminListUsers500JSONResponse{}, nil
		}
	}

	items := make([]openapi.AdminUserResponse, 0, len(users))
	for _, u := range users {
		uid := u.ID
		username := u.Username
		displayName := u.DisplayName
		email := u.Email
		role := string(u.Role)
		status := string(u.Status)
		isAdmin := u.IsAdmin

		items = append(items, openapi.AdminUserResponse{
			Id:          &uid,
			Username:    &username,
			DisplayName: &displayName,
			Email:       &email,
			Role:        &role,
			Status:      &status,
			IsAdmin:     &isAdmin,
			CreatedAt:   &u.CreatedAt,
		})
	}

	totalInt := total
	pageInt := page
	sizeInt := pageSize
	resp := openapi.AdminUserListResponse{
		Items:    &items,
		Total:    &totalInt,
		Page:     &pageInt,
		PageSize: &sizeInt,
	}
	return openapi.AdminListUsers200JSONResponse(resp), nil
}

// AdminBlockUser handles PATCH /admin/users/{userId}/block
func (h *Handler) AdminBlockUser(ctx context.Context, request openapi.AdminBlockUserRequestObject) (openapi.AdminBlockUserResponseObject, error) {
	actor := getActor(ctx)

	if request.Body == nil {
		return openapi.AdminBlockUser500JSONResponse{}, nil
	}

	var targetStatus string
	if request.Body.Blocked {
		targetStatus = "locked"
	} else {
		targetStatus = "active"
	}

	if err := h.adminUC.UpdateUserStatus(ctx, actor, request.UserId, userStatusFromString(targetStatus)); err != nil {
		statusCode, errResp := mapAdminError(err)
		switch statusCode {
		case 401:
			return openapi.AdminBlockUser401JSONResponse{UnauthorizedResponseJSONResponse: openapi.UnauthorizedResponseJSONResponse(errResp)}, nil
		case 403:
			return openapi.AdminBlockUser403JSONResponse{ForbiddenResponseJSONResponse: openapi.ForbiddenResponseJSONResponse(errResp)}, nil
		case 404:
			return openapi.AdminBlockUser404JSONResponse{NotFoundResponseJSONResponse: openapi.NotFoundResponseJSONResponse(errResp)}, nil
		default:
			return openapi.AdminBlockUser500JSONResponse{}, nil
		}
	}
	return openapi.AdminBlockUser200Response{}, nil
}
