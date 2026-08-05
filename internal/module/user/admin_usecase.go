package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

type adminStore interface {
	repo.AdminRepository
	ListProducts(ctx context.Context, page entity.Pagination) ([]entity.Product, int, error)
	GetProduct(ctx context.Context, productID string) (entity.Product, error)
	UpsertProduct(ctx context.Context, product entity.Product) (entity.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListCategories(ctx context.Context) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	ListSettings(ctx context.Context) ([]entity.Setting, error)
	DeleteSetting(ctx context.Context, key string) error
	ListRefundRequests(ctx context.Context, status string) ([]entity.RefundRequest, error)
	ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	DeleteOrder(ctx context.Context, orderID string) error
	ListReviews(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	UpdateReviewStatus(ctx context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error)
	DeleteReview(ctx context.Context, reviewID int64) error
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
