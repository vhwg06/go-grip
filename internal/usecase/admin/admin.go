package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
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

type UseCase struct {
	repo       adminStore
	notifier   webapi.AdminNotifier
	adminUsers map[string]struct{}
}

func New(adminRepo adminStore, notifier webapi.AdminNotifier, adminUsersCSV string) *UseCase {
	adminUsers := make(map[string]struct{})
	for raw := range strings.SplitSeq(adminUsersCSV, ",") {
		trimmed := strings.ToLower(strings.TrimSpace(raw))
		if trimmed == "" {
			continue
		}
		adminUsers[trimmed] = struct{}{}
	}
	return &UseCase{repo: adminRepo, notifier: notifier, adminUsers: adminUsers}
}

var _ usecase.Admin = (*UseCase)(nil)

func (uc *UseCase) ListUsers(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.User, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListUsers(ctx, page)
}

func (uc *UseCase) UpdateUserStatus(ctx context.Context, actor entity.Actor, userID string, status entity.UserStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateUserStatus(ctx, userID, status)
}

func (uc *UseCase) UpdateUserPoints(ctx context.Context, actor entity.Actor, userID string, points int) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateUserPoints(ctx, userID, points)
}

func (uc *UseCase) ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListOrders(ctx, page, query, status)
}

func (uc *UseCase) GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Order{}, err
	}
	return uc.repo.GetOrderByID(ctx, orderID)
}

func (uc *UseCase) UpdateOrderStatus(ctx context.Context, actor entity.Actor, orderID string, status entity.OrderStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	switch status {
	case entity.OrderStatusPaid, entity.OrderStatusDelivered, entity.OrderStatusCancelled:
		order, err := uc.repo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		if status == entity.OrderStatusDelivered && order.Status == entity.OrderStatusPending {
			return entity.ErrInvalidInput
		}
		return uc.repo.UpdateOrderStatus(ctx, orderID, status)
	default:
		return entity.ErrInvalidInput
	}
}

func (uc *UseCase) DeleteOrder(ctx context.Context, actor entity.Actor, orderID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteOrder(ctx, orderID)
}

func (uc *UseCase) ListRefunds(ctx context.Context, actor entity.Actor, status string) ([]entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListRefundRequests(ctx, status)
}

func (uc *UseCase) GetRefund(ctx context.Context, actor entity.Actor, refundID int64) (entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.RefundRequest{}, err
	}
	return uc.repo.GetRefundRequest(ctx, refundID)
}

func (uc *UseCase) GetOrderRefundStatus(ctx context.Context, actor entity.Actor, orderID string) (entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.RefundRequest{}, err
	}
	return uc.repo.GetOrderRefundStatus(ctx, orderID)
}

func (uc *UseCase) ListCards(ctx context.Context, actor entity.Actor) ([]entity.Card, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListCards(ctx)
}

func (uc *UseCase) DecideRefund(ctx context.Context, actor entity.Actor, refundID int64, approve bool, note string) (entity.RefundRequest, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.RefundRequest{}, err
	}
	if refundID <= 0 {
		return entity.RefundRequest{}, entity.ErrInvalidInput
	}
	return uc.repo.ProcessRefund(ctx, refundID, approve, actor.Username, note)
}

func (uc *UseCase) RepairAggregates(ctx context.Context, actor entity.Actor) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.RebuildProductAggregates(ctx)
}

func (uc *UseCase) ListProducts(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.Product, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListProducts(ctx, page)
}

func (uc *UseCase) UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Product{}, err
	}
	return uc.repo.UpsertProduct(ctx, product)
}

func (uc *UseCase) GetProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Product{}, err
	}
	return uc.repo.GetProduct(ctx, productID)
}

func (uc *UseCase) DeleteProduct(ctx context.Context, actor entity.Actor, productID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteProduct(ctx, productID)
}

func (uc *UseCase) ListCategories(ctx context.Context, actor entity.Actor) ([]entity.Category, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) UpsertCategory(ctx context.Context, actor entity.Actor, category entity.Category) (entity.Category, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Category{}, err
	}
	return uc.repo.UpsertCategory(ctx, category)
}

func (uc *UseCase) DeleteCategory(ctx context.Context, actor entity.Actor, categoryID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteCategory(ctx, categoryID)
}


func (uc *UseCase) ListSettings(ctx context.Context, actor entity.Actor) ([]entity.Setting, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListSettings(ctx)
}

func (uc *UseCase) SetSetting(ctx context.Context, actor entity.Actor, key, value string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if strings.TrimSpace(key) == "" {
		return entity.ErrInvalidInput
	}
	return uc.repo.StoreSetting(ctx, entity.Setting{Key: key, Value: value})
}

func (uc *UseCase) DeleteSetting(ctx context.Context, actor entity.Actor, key string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if strings.TrimSpace(key) == "" {
		return entity.ErrInvalidInput
	}
	return uc.repo.DeleteSetting(ctx, key)
}

func (uc *UseCase) ListReviews(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, repo.ReviewModerationStats{}, 0, err
	}
	return uc.repo.ListReviews(ctx, page, query, status)
}

func (uc *UseCase) UpdateReviewStatus(ctx context.Context, actor entity.Actor, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Review{}, err
	}
	if reviewID <= 0 {
		return entity.Review{}, entity.ErrInvalidInput
	}
	switch status {
	case entity.ReviewStatusApproved, entity.ReviewStatusHidden, entity.ReviewStatusFeatured:
		return uc.repo.UpdateReviewStatus(ctx, reviewID, status)
	default:
		return entity.Review{}, entity.ErrInvalidInput
	}
}

func (uc *UseCase) BulkPublishReviews(ctx context.Context, actor entity.Actor, reviewIDs []int64) (int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return 0, err
	}
	if len(reviewIDs) == 0 {
		return 0, entity.ErrInvalidInput
	}
	return uc.repo.BulkUpdateReviewStatus(ctx, reviewIDs, entity.ReviewStatusApproved)
}

func (uc *UseCase) DeleteReview(ctx context.Context, actor entity.Actor, reviewID int64) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if reviewID <= 0 {
		return entity.ErrInvalidInput
	}
	return uc.repo.DeleteReview(ctx, reviewID)
}

func (uc *UseCase) SendBroadcast(ctx context.Context, actor entity.Actor, title, body string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	_, _ = uc.repo.StoreAdminMessage(ctx, entity.AdminMessage{
		TargetType:  "broadcast",
		TargetValue: "",
		Title:       title,
		Body:        body,
		Sender:      actor.Username,
	})
	if uc.notifier == nil {
		return nil
	}
	return uc.notifier.SendBroadcast(ctx, title, body)
}

func (uc *UseCase) SendTargeted(ctx context.Context, actor entity.Actor, userID, title, body string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	_, _ = uc.repo.StoreAdminMessage(ctx, entity.AdminMessage{
		TargetType:  "targeted",
		TargetValue: userID,
		Title:       title,
		Body:        body,
		Sender:      actor.Username,
	})
	if uc.notifier == nil {
		return nil
	}
	return uc.notifier.SendTargeted(ctx, userID, title, body)
}

func (uc *UseCase) ListAdminMessages(ctx context.Context, actor entity.Actor) ([]entity.AdminMessage, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListAdminMessages(ctx)
}


func (uc *UseCase) ensureAdmin(actor entity.Actor) error {
	if actor.IsAdmin {
		return nil
	}
	if _, ok := uc.adminUsers[strings.ToLower(strings.TrimSpace(actor.Username))]; ok {
		return nil
	}
	return entity.ErrForbidden
}

func (uc *UseCase) setSetting(ctx context.Context, actor entity.Actor, key, value string) error {
	return uc.SetSetting(ctx, actor, key, value)
}

func (uc *UseCase) deleteSetting(ctx context.Context, actor entity.Actor, key string) error {
	return uc.DeleteSetting(ctx, actor, key)
}

func (uc *UseCase) requireInput(v string) error {
	if strings.TrimSpace(v) == "" {
		return entity.ErrInvalidInput
	}
	return nil
}

func (uc *UseCase) ensureNoErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
