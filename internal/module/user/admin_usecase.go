package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

type adminStore interface {
	ListUsers(ctx context.Context, page entity.Pagination) ([]entity.User, int, error)
	UpdateUserStatus(ctx context.Context, userID string, status entity.UserStatus) error
	ListOrders(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	RebuildProductAggregates(ctx context.Context) error
}

// AdminUseCase defines the application service interface for backoffice administrative tools.
type AdminUseCase interface {
	ListUsers(ctx context.Context, actor Actor, page pagination.Pagination) ([]User, int, error)
	UpdateUserStatus(ctx context.Context, actor Actor, userID string, status UserStatus) error
	ListOrders(ctx context.Context, actor Actor, page pagination.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrder(ctx context.Context, actor Actor, orderID string) (entity.Order, error)
	RepairAggregates(ctx context.Context, actor Actor) error
}

type adminUseCase struct {
	repo       adminStore
	notifier   webapi.AdminNotifier
	adminUsers map[string]struct{}
}

// NewAdminUseCase constructs a new AdminUseCase instance.
func NewAdminUseCase(adminRepo adminStore, notifier webapi.AdminNotifier, adminUsersCSV string) AdminUseCase {
	adminUsers := make(map[string]struct{})
	for raw := range strings.SplitSeq(adminUsersCSV, ",") {
		trimmed := strings.ToLower(strings.TrimSpace(raw))
		if trimmed == "" {
			continue
		}
		adminUsers[trimmed] = struct{}{}
	}
	return &adminUseCase{repo: adminRepo, notifier: notifier, adminUsers: adminUsers}
}

func (uc *adminUseCase) ListUsers(ctx context.Context, actor Actor, page pagination.Pagination) ([]User, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	users, total, err := uc.repo.ListUsers(ctx, entity.Pagination(page))
	if err != nil {
		return nil, 0, err
	}
	res := make([]User, len(users))
	for i, u := range users {
		res[i] = User{
			ID:          u.ID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Role:        RoleName(u.Role),
			Status:      UserStatus(u.Status),
			IsAdmin:     u.IsAdmin,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		}
	}
	return res, total, nil
}

func (uc *adminUseCase) UpdateUserStatus(ctx context.Context, actor Actor, userID string, status UserStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateUserStatus(ctx, userID, entity.UserStatus(status))
}

func (uc *adminUseCase) ListOrders(ctx context.Context, actor Actor, page pagination.Pagination, query, status string) ([]entity.Order, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListOrders(ctx, entity.Pagination(page), query, status)
}

func (uc *adminUseCase) GetOrder(ctx context.Context, actor Actor, orderID string) (entity.Order, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Order{}, err
	}
	return uc.repo.GetOrderByID(ctx, orderID)
}

func (uc *adminUseCase) RepairAggregates(ctx context.Context, actor Actor) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if err := uc.repo.RebuildProductAggregates(ctx); err != nil {
		return fmt.Errorf("AdminUseCase.RepairAggregates: %w", err)
	}
	return nil
}

func (uc *adminUseCase) ensureAdmin(actor Actor) error {
	if actor.IsAdmin {
		return nil
	}
	if _, ok := uc.adminUsers[strings.ToLower(strings.TrimSpace(actor.Username))]; ok {
		return nil
	}
	return ErrForbidden
}
