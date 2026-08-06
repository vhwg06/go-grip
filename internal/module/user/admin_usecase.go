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
	ListUsers(ctx context.Context, page entity.Pagination) ([]entity.User, int, error)
	UpdateUserStatus(ctx context.Context, userID string, status entity.UserStatus) error
	ListOrders(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	DeleteOrder(ctx context.Context, orderID string) error
	ListRefundRequests(ctx context.Context, status string) ([]entity.RefundRequest, error)
	GetRefundRequest(ctx context.Context, refundID int64) (entity.RefundRequest, error)
	GetOrderRefundStatus(ctx context.Context, orderID string) (entity.RefundRequest, error)
	ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error)
	ListReviews(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	UpdateReviewStatus(ctx context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error)
	DeleteReview(ctx context.Context, reviewID int64) error
	ListSettings(ctx context.Context) ([]entity.Setting, error)
	StoreSetting(ctx context.Context, setting entity.Setting) error
	DeleteSetting(ctx context.Context, key string) error
	StoreAdminMessage(ctx context.Context, msg entity.AdminMessage) (entity.AdminMessage, error)
	ListAdminMessages(ctx context.Context) ([]entity.AdminMessage, error)
	ListProducts(ctx context.Context, page entity.Pagination) ([]entity.Product, int, error)
	ListCategories(ctx context.Context) ([]entity.Category, error)
	RebuildProductAggregates(ctx context.Context) error
}

// AdminUseCase defines the application service interface for backoffice administrative tools.
type AdminUseCase interface {
	ListUsers(ctx context.Context, actor Actor, page pagination.Pagination) ([]User, int, error)
	UpdateUserStatus(ctx context.Context, actor Actor, userID string, status UserStatus) error
	ListOrders(ctx context.Context, actor Actor, page pagination.Pagination, query, status string) ([]entity.Order, int, error)
	GetOrder(ctx context.Context, actor Actor, orderID string) (entity.Order, error)
	UpdateOrderStatus(ctx context.Context, actor Actor, orderID string, status entity.OrderStatus) error
	DeleteOrder(ctx context.Context, actor Actor, orderID string) error
	ListRefunds(ctx context.Context, actor Actor, status string) ([]entity.RefundRequest, error)
	GetRefund(ctx context.Context, actor Actor, refundID int64) (entity.RefundRequest, error)
	ProcessRefund(ctx context.Context, actor Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error)
	ListReviews(ctx context.Context, actor Actor, page pagination.Pagination, query, status string) ([]entity.Review, int, error)
	UpdateReviewStatus(ctx context.Context, actor Actor, reviewID int64, status entity.ReviewStatus) error
	BulkPublishReviews(ctx context.Context, actor Actor, reviewIDs []int64) (int, error)
	DeleteAdminReview(ctx context.Context, actor Actor, reviewID int64) error
	ListSettings(ctx context.Context, actor Actor) ([]entity.Setting, error)
	UpsertSetting(ctx context.Context, actor Actor, key, value string) error
	DeleteSetting(ctx context.Context, actor Actor, key string) error
	ListMessages(ctx context.Context, actor Actor) ([]entity.AdminMessage, error)
	BroadcastMessage(ctx context.Context, actor Actor, title, body, imageURL string) (entity.AdminMessage, error)
	ListProducts(ctx context.Context, actor Actor, page pagination.Pagination) ([]entity.Product, int, error)
	ListAdminCategories(ctx context.Context, actor Actor) ([]entity.Category, error)
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

// UpdateOrderStatus transitions an order to a new status, enforcing admin access and
// state-machine rules (paid → delivered, pending → cancelled, etc.) in the repository layer.
func (uc *adminUseCase) UpdateOrderStatus(ctx context.Context, actor Actor, orderID string, status entity.OrderStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateOrderStatus(ctx, orderID, status)
}

// DeleteOrder removes an order and its associated refund requests atomically.
func (uc *adminUseCase) DeleteOrder(ctx context.Context, actor Actor, orderID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteOrder(ctx, orderID)
}

// ListRefunds returns all refund requests optionally filtered by status.
func (uc *adminUseCase) ListRefunds(ctx context.Context, actor Actor, status string) ([]entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListRefundRequests(ctx, status)
}

// GetRefund returns a single refund request by ID.
func (uc *adminUseCase) GetRefund(ctx context.Context, actor Actor, refundID int64) (entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.RefundRequest{}, err
	}
	return uc.repo.GetRefundRequest(ctx, refundID)
}

// ProcessRefund atomically approves or rejects a pending refund and adjusts the
// associated order status and product stock in one transaction.
func (uc *adminUseCase) ProcessRefund(ctx context.Context, actor Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.RefundRequest{}, err
	}
	return uc.repo.ProcessRefund(ctx, refundID, approve, actor.Username, note)
}

// ListReviews returns paginated reviews for moderation, optionally filtered by query and status.
func (uc *adminUseCase) ListReviews(ctx context.Context, actor Actor, page pagination.Pagination, query, status string) ([]entity.Review, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	items, _, total, err := uc.repo.ListReviews(ctx, entity.Pagination(page), query, status)
	return items, total, err
}

// UpdateReviewStatus changes review visibility status (pending, featured, hidden).
func (uc *adminUseCase) UpdateReviewStatus(ctx context.Context, actor Actor, reviewID int64, status entity.ReviewStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	_, err := uc.repo.UpdateReviewStatus(ctx, reviewID, status)
	return err
}

// BulkPublishReviews sets all selected review IDs to Featured status and returns the count updated.
func (uc *adminUseCase) BulkPublishReviews(ctx context.Context, actor Actor, reviewIDs []int64) (int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return 0, err
	}
	return uc.repo.BulkUpdateReviewStatus(ctx, reviewIDs, entity.ReviewStatusFeatured)
}

// DeleteAdminReview permanently deletes a review by ID.
func (uc *adminUseCase) DeleteAdminReview(ctx context.Context, actor Actor, reviewID int64) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteReview(ctx, reviewID)
}

// ListSettings returns all key-value store settings.
func (uc *adminUseCase) ListSettings(ctx context.Context, actor Actor) ([]entity.Setting, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListSettings(ctx)
}

// UpsertSetting creates or updates a setting by key.
func (uc *adminUseCase) UpsertSetting(ctx context.Context, actor Actor, key, value string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.StoreSetting(ctx, entity.Setting{Key: key, Value: value})
}

// DeleteSetting removes a setting by key.
func (uc *adminUseCase) DeleteSetting(ctx context.Context, actor Actor, key string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteSetting(ctx, key)
}

// ListMessages returns all admin broadcast messages ordered by creation time.
func (uc *adminUseCase) ListMessages(ctx context.Context, actor Actor) ([]entity.AdminMessage, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListAdminMessages(ctx)
}

// BroadcastMessage persists a broadcast notification and dispatches it via the configured notifier.
func (uc *adminUseCase) BroadcastMessage(ctx context.Context, actor Actor, title, body, imageURL string) (entity.AdminMessage, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.AdminMessage{}, err
	}
	msg := entity.AdminMessage{
		TargetType:  "all",
		TargetValue: "",
		Title:       title,
		Body:        body,
		Sender:      actor.Username,
	}
	stored, err := uc.repo.StoreAdminMessage(ctx, msg)
	if err != nil {
		return entity.AdminMessage{}, fmt.Errorf("AdminUseCase.BroadcastMessage: %w", err)
	}
	if uc.notifier != nil {
		_ = uc.notifier.SendBroadcast(ctx, stored.Title, stored.Body)
	}
	return stored, nil
}

// ListProducts returns paginated legacy products for the admin backoffice.
func (uc *adminUseCase) ListProducts(ctx context.Context, actor Actor, page pagination.Pagination) ([]entity.Product, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListProducts(ctx, entity.Pagination(page))
}

// ListAdminCategories returns all categories for admin management.
func (uc *adminUseCase) ListAdminCategories(ctx context.Context, actor Actor) ([]entity.Category, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListCategories(ctx)
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

