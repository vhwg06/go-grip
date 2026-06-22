package v1

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

type BaseAdminUseCaseStub struct{}

func (BaseAdminUseCaseStub) ListProducts(context.Context, entity.Actor, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}
func (BaseAdminUseCaseStub) GetProduct(context.Context, entity.Actor, string) (entity.Product, error) {
	return entity.Product{}, nil
}
func (BaseAdminUseCaseStub) UpsertProduct(context.Context, entity.Actor, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}
func (BaseAdminUseCaseStub) DeleteProduct(context.Context, entity.Actor, string) error {
	return nil
}
func (BaseAdminUseCaseStub) ListCategories(context.Context, entity.Actor) ([]entity.Category, error) {
	return nil, nil
}
func (BaseAdminUseCaseStub) UpsertCategory(context.Context, entity.Actor, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}
func (BaseAdminUseCaseStub) DeleteCategory(context.Context, entity.Actor, string) error {
	return nil
}
func (BaseAdminUseCaseStub) ListOrders(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}
func (BaseAdminUseCaseStub) GetOrder(context.Context, entity.Actor, string) (entity.Order, error) {
	return entity.Order{}, nil
}
func (BaseAdminUseCaseStub) UpdateOrderStatus(context.Context, entity.Actor, string, entity.OrderStatus) error {
	return nil
}
func (BaseAdminUseCaseStub) DeleteOrder(context.Context, entity.Actor, string) error {
	return nil
}
func (BaseAdminUseCaseStub) ListRefunds(context.Context, entity.Actor, string) ([]entity.RefundRequest, error) {
	return nil, nil
}
func (BaseAdminUseCaseStub) GetRefund(context.Context, entity.Actor, int64) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (BaseAdminUseCaseStub) GetOrderRefundStatus(context.Context, entity.Actor, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (BaseAdminUseCaseStub) DecideRefund(context.Context, entity.Actor, int64, bool, string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}
func (BaseAdminUseCaseStub) ListReviews(context.Context, entity.Actor, entity.Pagination, string, string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	return nil, repo.ReviewModerationStats{}, 0, nil
}
func (BaseAdminUseCaseStub) UpdateReviewStatus(context.Context, entity.Actor, int64, entity.ReviewStatus) (entity.Review, error) {
	return entity.Review{}, nil
}
func (BaseAdminUseCaseStub) BulkPublishReviews(context.Context, entity.Actor, []int64) (int, error) {
	return 0, nil
}
func (BaseAdminUseCaseStub) DeleteReview(context.Context, entity.Actor, int64) error {
	return nil
}
func (BaseAdminUseCaseStub) ListSettings(context.Context, entity.Actor) ([]entity.Setting, error) {
	return nil, nil
}
func (BaseAdminUseCaseStub) SetSetting(context.Context, entity.Actor, string, string) error {
	return nil
}
func (BaseAdminUseCaseStub) DeleteSetting(context.Context, entity.Actor, string) error {
	return nil
}
func (BaseAdminUseCaseStub) SendBroadcast(context.Context, entity.Actor, string, string) error {
	return nil
}
func (BaseAdminUseCaseStub) SendTargeted(context.Context, entity.Actor, string, string, string) error {
	return nil
}
func (BaseAdminUseCaseStub) ListAdminMessages(context.Context, entity.Actor) ([]entity.AdminMessage, error) {
	return nil, nil
}
func (BaseAdminUseCaseStub) ListUsers(context.Context, entity.Actor, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}
func (BaseAdminUseCaseStub) UpdateUserStatus(context.Context, entity.Actor, string, entity.UserStatus) error {
	return nil
}
func (BaseAdminUseCaseStub) RepairAggregates(context.Context, entity.Actor) error {
	return nil
}
