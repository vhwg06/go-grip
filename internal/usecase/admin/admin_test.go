package admin

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/stretchr/testify/require"
)

type adminStoreStub struct {
	listSettingsFunc  func(ctx context.Context) ([]entity.Setting, error)
	storeSettingFunc  func(ctx context.Context, setting entity.Setting) error
	deleteSettingFunc func(ctx context.Context, key string) error
	listRefundsFunc   func(ctx context.Context, status string) ([]entity.RefundRequest, error)
	processRefundFunc func(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error)
	updateOrderFunc   func(ctx context.Context, orderID string, status entity.OrderStatus) error
	deleteOrderFunc   func(ctx context.Context, orderID string) error
	listReviewsFunc   func(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error)
	updateReviewFunc  func(ctx context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error)
	bulkReviewsFunc   func(ctx context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error)
	deleteReviewFunc  func(ctx context.Context, reviewID int64) error
	getOrderFunc      func(ctx context.Context, orderID string) (entity.Order, error)
}

func (s *adminStoreStub) ListUsers(context.Context, entity.Pagination) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (s *adminStoreStub) UpdateUserStatus(context.Context, string, entity.UserStatus) error {
	return nil
}

func (s *adminStoreStub) UpdateUserPoints(context.Context, string, int) error {
	return nil
}

func (s *adminStoreStub) ListOrders(context.Context, entity.Pagination, string, string) ([]entity.Order, int, error) {
	return nil, 0, nil
}

func (s *adminStoreStub) GetOrderByID(ctx context.Context, id string) (entity.Order, error) {
	if s.getOrderFunc != nil {
		return s.getOrderFunc(ctx, id)
	}
	return entity.Order{}, nil
}

func (s *adminStoreStub) ListRefundRequests(ctx context.Context, status string) ([]entity.RefundRequest, error) {
	if s.listRefundsFunc != nil {
		return s.listRefundsFunc(ctx, status)
	}
	return nil, nil
}

func (s *adminStoreStub) ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error) {
	if s.processRefundFunc != nil {
		return s.processRefundFunc(ctx, refundID, approve, adminUsername, note)
	}
	return entity.RefundRequest{}, nil
}

func (s *adminStoreStub) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	if s.updateOrderFunc != nil {
		return s.updateOrderFunc(ctx, orderID, status)
	}
	return nil
}

func (s *adminStoreStub) DeleteOrder(ctx context.Context, orderID string) error {
	if s.deleteOrderFunc != nil {
		return s.deleteOrderFunc(ctx, orderID)
	}
	return nil
}

func (s *adminStoreStub) ListReviews(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	if s.listReviewsFunc != nil {
		return s.listReviewsFunc(ctx, page, query, status)
	}
	return nil, repo.ReviewModerationStats{}, 0, nil
}

func (s *adminStoreStub) UpdateReviewStatus(ctx context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
	if s.updateReviewFunc != nil {
		return s.updateReviewFunc(ctx, reviewID, status)
	}
	return entity.Review{}, nil
}

func (s *adminStoreStub) BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error) {
	if s.bulkReviewsFunc != nil {
		return s.bulkReviewsFunc(ctx, reviewIDs, status)
	}
	return len(reviewIDs), nil
}

func (s *adminStoreStub) DeleteReview(ctx context.Context, reviewID int64) error {
	if s.deleteReviewFunc != nil {
		return s.deleteReviewFunc(ctx, reviewID)
	}
	return nil
}

func (s *adminStoreStub) ListSettings(ctx context.Context) ([]entity.Setting, error) {
	if s.listSettingsFunc != nil {
		return s.listSettingsFunc(ctx)
	}
	return nil, nil
}

func (s *adminStoreStub) StoreSetting(ctx context.Context, setting entity.Setting) error {
	if s.storeSettingFunc != nil {
		return s.storeSettingFunc(ctx, setting)
	}
	return nil
}

func (s *adminStoreStub) DeleteSetting(ctx context.Context, key string) error {
	if s.deleteSettingFunc != nil {
		return s.deleteSettingFunc(ctx, key)
	}
	return nil
}

func (s *adminStoreStub) GetRefundRequest(ctx context.Context, refundID int64) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminStoreStub) GetOrderRefundStatus(ctx context.Context, orderID string) (entity.RefundRequest, error) {
	return entity.RefundRequest{}, nil
}

func (s *adminStoreStub) RebuildProductAggregates(context.Context) error {
	return nil
}

func (s *adminStoreStub) ListProducts(context.Context, entity.Pagination) ([]entity.Product, int, error) {
	return nil, 0, nil
}

func (s *adminStoreStub) GetProduct(context.Context, string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminStoreStub) UpsertProduct(context.Context, entity.Product) (entity.Product, error) {
	return entity.Product{}, nil
}

func (s *adminStoreStub) DeleteProduct(context.Context, string) error {
	return nil
}

func (s *adminStoreStub) ListCategories(context.Context) ([]entity.Category, error) {
	return nil, nil
}

func (s *adminStoreStub) UpsertCategory(context.Context, entity.Category) (entity.Category, error) {
	return entity.Category{}, nil
}

func (s *adminStoreStub) DeleteCategory(context.Context, string) error {
	return nil
}

func (s *adminStoreStub) StoreAdminMessage(ctx context.Context, msg entity.AdminMessage) (entity.AdminMessage, error) {
	return entity.AdminMessage{}, nil
}

func (s *adminStoreStub) ListAdminMessages(ctx context.Context) ([]entity.AdminMessage, error) {
	return nil, nil
}

func TestUseCase_SettingsRequireAdminAndPersist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	adminActor := entity.Actor{IsAdmin: true, Username: "admin"}
	userActor := entity.Actor{Username: "user"}

	t.Run("list settings rejects non-admin", func(t *testing.T) {
		uc := New(&adminStoreStub{}, nil, "")
		_, err := uc.ListSettings(ctx, userActor)
		require.ErrorIs(t, err, entity.ErrForbidden)
	})

	t.Run("list settings returns repo values", func(t *testing.T) {
		uc := New(&adminStoreStub{
			listSettingsFunc: func(context.Context) ([]entity.Setting, error) {
				return []entity.Setting{{Key: "shopName", Value: "Grip Store"}}, nil
			},
		}, nil, "")

		settings, err := uc.ListSettings(ctx, adminActor)
		require.NoError(t, err)
		require.Equal(t, []entity.Setting{{Key: "shopName", Value: "Grip Store"}}, settings)
	})

	t.Run("set setting validates key and stores value", func(t *testing.T) {
		var stored entity.Setting
		uc := New(&adminStoreStub{
			storeSettingFunc: func(_ context.Context, setting entity.Setting) error {
				stored = setting
				return nil
			},
		}, nil, "")

		require.ErrorIs(t, uc.SetSetting(ctx, adminActor, "", "x"), entity.ErrInvalidInput)
		require.NoError(t, uc.SetSetting(ctx, adminActor, "shopName", "Updated Store"))
		require.Equal(t, entity.Setting{Key: "shopName", Value: "Updated Store"}, stored)
	})

	t.Run("delete setting validates key and deletes value", func(t *testing.T) {
		var deletedKey string
		uc := New(&adminStoreStub{
			deleteSettingFunc: func(_ context.Context, key string) error {
				deletedKey = key
				return nil
			},
		}, nil, "")

		require.ErrorIs(t, uc.DeleteSetting(ctx, adminActor, " "), entity.ErrInvalidInput)
		require.NoError(t, uc.DeleteSetting(ctx, adminActor, "announcement"))
		require.Equal(t, "announcement", deletedKey)
	})
}

func TestUseCase_OrderAndRefundAdminActions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	adminActor := entity.Actor{IsAdmin: true, Username: "admin"}
	userActor := entity.Actor{Username: "user"}

	t.Run("update order status validates allowed transitions", func(t *testing.T) {
		var gotOrderID string
		var gotStatus entity.OrderStatus
		uc := New(&adminStoreStub{
			updateOrderFunc: func(_ context.Context, orderID string, status entity.OrderStatus) error {
				gotOrderID = orderID
				gotStatus = status
				return nil
			},
			getOrderFunc: func(_ context.Context, orderID string) (entity.Order, error) {
				return entity.Order{ID: orderID, Status: entity.OrderStatusPending}, nil
			},
		}, nil, "")

		require.ErrorIs(t, uc.UpdateOrderStatus(ctx, userActor, "o1", entity.OrderStatusPaid), entity.ErrForbidden)
		require.ErrorIs(t, uc.UpdateOrderStatus(ctx, adminActor, "o1", entity.OrderStatusRefunded), entity.ErrInvalidInput)
		require.ErrorIs(t, uc.UpdateOrderStatus(ctx, adminActor, "o1", entity.OrderStatusDelivered), entity.ErrInvalidTransition)
		require.NoError(t, uc.UpdateOrderStatus(ctx, adminActor, "o1", entity.OrderStatusCancelled))
		require.Equal(t, "o1", gotOrderID)
		require.Equal(t, entity.OrderStatusCancelled, gotStatus)
	})

	t.Run("delete order requires admin", func(t *testing.T) {
		var deleted string
		uc := New(&adminStoreStub{
			deleteOrderFunc: func(_ context.Context, orderID string) error {
				deleted = orderID
				return nil
			},
		}, nil, "")

		require.ErrorIs(t, uc.DeleteOrder(ctx, userActor, "o1"), entity.ErrForbidden)
		require.NoError(t, uc.DeleteOrder(ctx, adminActor, "o1"))
		require.Equal(t, "o1", deleted)
	})

	t.Run("list refunds passes status through", func(t *testing.T) {
		uc := New(&adminStoreStub{
			listRefundsFunc: func(_ context.Context, status string) ([]entity.RefundRequest, error) {
				require.Equal(t, "pending", status)
				return []entity.RefundRequest{{ID: 7, Status: entity.RefundStatusPending}}, nil
			},
		}, nil, "")

		refunds, err := uc.ListRefunds(ctx, adminActor, "pending")
		require.NoError(t, err)
		require.Len(t, refunds, 1)
	})

	t.Run("decide refund validates inputs and forwards actor metadata", func(t *testing.T) {
		var gotID int64
		var gotApprove bool
		var gotAdmin string
		var gotNote string
		uc := New(&adminStoreStub{
			processRefundFunc: func(_ context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error) {
				gotID = refundID
				gotApprove = approve
				gotAdmin = adminUsername
				gotNote = note
				return entity.RefundRequest{ID: refundID, Status: entity.RefundStatusApproved, AdminUsername: adminUsername, AdminNote: note}, nil
			},
		}, nil, "")

		_, err := uc.DecideRefund(ctx, userActor, 1, true, "")
		require.ErrorIs(t, err, entity.ErrForbidden)
		_, err = uc.DecideRefund(ctx, adminActor, 0, true, "")
		require.ErrorIs(t, err, entity.ErrInvalidInput)

		refund, err := uc.DecideRefund(ctx, adminActor, 9, true, "approved")
		require.NoError(t, err)
		require.Equal(t, int64(9), gotID)
		require.True(t, gotApprove)
		require.Equal(t, "admin", gotAdmin)
		require.Equal(t, "approved", gotNote)
		require.Equal(t, entity.RefundStatusApproved, refund.Status)
	})
}

func TestUseCase_ReviewModerationAdminActions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	adminActor := entity.Actor{IsAdmin: true}
	userActor := entity.Actor{}

	t.Run("list reviews requires admin", func(t *testing.T) {
		uc := New(&adminStoreStub{
			listReviewsFunc: func(_ context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
				require.Equal(t, "needle", query)
				require.Equal(t, "PENDING", status)
				return []entity.Review{{ID: 1, Status: entity.ReviewStatusPending}}, repo.ReviewModerationStats{Pending: 1}, 1, nil
			},
		}, nil, "")

		_, _, _, err := uc.ListReviews(ctx, userActor, entity.Pagination{}, "needle", "PENDING")
		require.ErrorIs(t, err, entity.ErrForbidden)

		items, stats, total, err := uc.ListReviews(ctx, adminActor, entity.Pagination{}, "needle", "PENDING")
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, 1, stats.Pending)
		require.Equal(t, 1, total)
	})

	t.Run("update review status validates inputs", func(t *testing.T) {
		var gotID int64
		var gotStatus entity.ReviewStatus
		uc := New(&adminStoreStub{
			updateReviewFunc: func(_ context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
				gotID = reviewID
				gotStatus = status
				return entity.Review{ID: reviewID, Status: status}, nil
			},
		}, nil, "")

		_, err := uc.UpdateReviewStatus(ctx, userActor, 1, entity.ReviewStatusApproved)
		require.ErrorIs(t, err, entity.ErrForbidden)
		_, err = uc.UpdateReviewStatus(ctx, adminActor, 0, entity.ReviewStatusApproved)
		require.ErrorIs(t, err, entity.ErrInvalidInput)
		_, err = uc.UpdateReviewStatus(ctx, adminActor, 1, entity.ReviewStatus("BOGUS"))
		require.ErrorIs(t, err, entity.ErrInvalidInput)

		review, err := uc.UpdateReviewStatus(ctx, adminActor, 11, entity.ReviewStatusFeatured)
		require.NoError(t, err)
		require.Equal(t, int64(11), gotID)
		require.Equal(t, entity.ReviewStatusFeatured, gotStatus)
		require.Equal(t, entity.ReviewStatusFeatured, review.Status)
	})

	t.Run("bulk publish and delete validate inputs", func(t *testing.T) {
		var bulkIDs []int64
		var deletedID int64
		uc := New(&adminStoreStub{
			bulkReviewsFunc: func(_ context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error) {
				require.Equal(t, entity.ReviewStatusApproved, status)
				bulkIDs = reviewIDs
				return len(reviewIDs), nil
			},
			deleteReviewFunc: func(_ context.Context, reviewID int64) error {
				deletedID = reviewID
				return nil
			},
		}, nil, "")

		_, err := uc.BulkPublishReviews(ctx, adminActor, nil)
		require.ErrorIs(t, err, entity.ErrInvalidInput)

		count, err := uc.BulkPublishReviews(ctx, adminActor, []int64{1, 2})
		require.NoError(t, err)
		require.Equal(t, 2, count)
		require.Equal(t, []int64{1, 2}, bulkIDs)

		require.ErrorIs(t, uc.DeleteReview(ctx, adminActor, 0), entity.ErrInvalidInput)
		require.NoError(t, uc.DeleteReview(ctx, adminActor, 9))
		require.Equal(t, int64(9), deletedID)
	})
}
