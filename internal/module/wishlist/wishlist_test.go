package wishlist

import (
	"context"
	"testing"

	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/stretchr/testify/require"
)

type wishlistRepoStub struct {
	listFunc       func(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error)
	storeFunc      func(ctx context.Context, item WishlistItem) (WishlistItem, error)
	updateFunc     func(ctx context.Context, item WishlistItem) (WishlistItem, error)
	deleteFunc     func(ctx context.Context, itemID int64) error
	toggleVoteFunc func(ctx context.Context, itemID int64, userID string) (bool, error)
	storeReview    func(ctx context.Context, review Review) (Review, error)
	listReviews    func(ctx context.Context, productID string) ([]Review, error)
	getReview      func(ctx context.Context, reviewID int64) (Review, error)
	deleteReview   func(ctx context.Context, reviewID int64) error
}

func (s *wishlistRepoStub) ListWishlistItems(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error) {
	if s.listFunc != nil {
		return s.listFunc(ctx, page)
	}
	return nil, 0, nil
}

func (s *wishlistRepoStub) StoreWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error) {
	if s.storeFunc != nil {
		return s.storeFunc(ctx, item)
	}
	return item, nil
}

func (s *wishlistRepoStub) UpdateWishlistItem(ctx context.Context, item WishlistItem) (WishlistItem, error) {
	if s.updateFunc != nil {
		return s.updateFunc(ctx, item)
	}
	return item, nil
}

func (s *wishlistRepoStub) DeleteWishlistItem(ctx context.Context, itemID int64) error {
	if s.deleteFunc != nil {
		return s.deleteFunc(ctx, itemID)
	}
	return nil
}

func (s *wishlistRepoStub) ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error) {
	if s.toggleVoteFunc != nil {
		return s.toggleVoteFunc(ctx, itemID, userID)
	}
	return true, nil
}

func (s *wishlistRepoStub) StoreReview(ctx context.Context, review Review) (Review, error) {
	if s.storeReview != nil {
		return s.storeReview(ctx, review)
	}
	return review, nil
}

func (s *wishlistRepoStub) ListReviews(ctx context.Context, productID string) ([]Review, error) {
	if s.listReviews != nil {
		return s.listReviews(ctx, productID)
	}
	return nil, nil
}

func (s *wishlistRepoStub) GetReview(ctx context.Context, reviewID int64) (Review, error) {
	if s.getReview != nil {
		return s.getReview(ctx, reviewID)
	}
	return Review{ID: reviewID, UserID: "u1"}, nil
}

func (s *wishlistRepoStub) DeleteReview(ctx context.Context, reviewID int64) error {
	if s.deleteReview != nil {
		return s.deleteReview(ctx, reviewID)
	}
	return nil
}

type orderReaderStub struct {
	getFunc func(ctx context.Context, orderID string) (ordermodule.Order, error)
}

func (s *orderReaderStub) GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error) {
	if s.getFunc != nil {
		return s.getFunc(ctx, orderID)
	}
	return ordermodule.Order{}, nil
}

func TestWishlistUseCase_Lifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	actor := Actor{UserID: "u1", Username: "alice"}

	t.Run("create requires auth and title", func(t *testing.T) {
		uc := NewWishlistUseCase(&wishlistRepoStub{}, nil)

		_, err := uc.Create(ctx, Actor{}, "wish", "desc")
		require.ErrorIs(t, err, ErrUnauthorized)

		_, err = uc.Create(ctx, actor, "", "desc")
		require.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("create persists actor-owned item", func(t *testing.T) {
		var stored WishlistItem
		uc := NewWishlistUseCase(&wishlistRepoStub{
			storeFunc: func(_ context.Context, item WishlistItem) (WishlistItem, error) {
				stored = item
				item.ID = 9
				return item, nil
			},
		}, nil)

		item, err := uc.Create(ctx, actor, "Need key", "Weekend")
		require.NoError(t, err)
		require.Equal(t, "u1", stored.UserID)
		require.Equal(t, "alice", stored.Username)
		require.Equal(t, "Need key", stored.Title)
		require.Equal(t, int64(9), item.ID)
	})

	t.Run("toggle vote requires auth and forwards user id", func(t *testing.T) {
		var gotItemID int64
		var gotUserID string
		uc := NewWishlistUseCase(&wishlistRepoStub{
			toggleVoteFunc: func(_ context.Context, itemID int64, userID string) (bool, error) {
				gotItemID = itemID
				gotUserID = userID
				return true, nil
			},
		}, nil)

		err := uc.ToggleVote(ctx, Actor{}, 1)
		require.ErrorIs(t, err, ErrUnauthorized)

		require.NoError(t, uc.ToggleVote(ctx, actor, 7))
		require.Equal(t, int64(7), gotItemID)
		require.Equal(t, "u1", gotUserID)
	})

	t.Run("create review validates rating and delivered order", func(t *testing.T) {
		uc := NewWishlistUseCase(&wishlistRepoStub{}, &orderReaderStub{
			getFunc: func(_ context.Context, orderID string) (ordermodule.Order, error) {
				require.Equal(t, "o1", orderID)
				return ordermodule.Order{ID: orderID, Status: ordermodule.OrderStatusPending}, nil
			},
		})

		_, err := uc.CreateReview(ctx, Actor{}, Review{ProductID: "p1", Rating: 5})
		require.ErrorIs(t, err, ErrUnauthorized)

		_, err = uc.CreateReview(ctx, actor, Review{ProductID: "", Rating: 5})
		require.ErrorIs(t, err, ErrInvalidInput)

		_, err = uc.CreateReview(ctx, actor, Review{ProductID: "p1", Rating: 0})
		require.ErrorIs(t, err, ErrInvalidInput)

		_, err = uc.CreateReview(ctx, actor, Review{ProductID: "p1", OrderID: "o1", Rating: 5})
		require.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("create review stores actor review for delivered order", func(t *testing.T) {
		var stored Review
		uc := NewWishlistUseCase(&wishlistRepoStub{
			storeReview: func(_ context.Context, review Review) (Review, error) {
				stored = review
				review.ID = 3
				return review, nil
			},
		}, &orderReaderStub{
			getFunc: func(_ context.Context, orderID string) (ordermodule.Order, error) {
				return ordermodule.Order{ID: orderID, Status: ordermodule.OrderStatusDelivered}, nil
			},
		})

		review, err := uc.CreateReview(ctx, actor, Review{ProductID: "p1", OrderID: "o2", Rating: 5, Comment: "Great"})
		require.NoError(t, err)
		require.Equal(t, "u1", stored.UserID)
		require.Equal(t, "alice", stored.Username)
		require.Equal(t, "o2", stored.OrderID)
		require.Equal(t, int64(3), review.ID)
	})

	t.Run("list reviews validates product id", func(t *testing.T) {
		uc := NewWishlistUseCase(&wishlistRepoStub{
			listReviews: func(_ context.Context, productID string) ([]Review, error) {
				require.Equal(t, "p1", productID)
				return []Review{{ID: 1, ProductID: productID}}, nil
			},
		}, nil)

		_, err := uc.ListReviews(ctx, "")
		require.ErrorIs(t, err, ErrInvalidInput)

		reviews, err := uc.ListReviews(ctx, "p1")
		require.NoError(t, err)
		require.Len(t, reviews, 1)
	})

	t.Run("delete review enforces ownership and allows admins", func(t *testing.T) {
		deletedID := int64(0)
		uc := NewWishlistUseCase(&wishlistRepoStub{
			getReview: func(_ context.Context, reviewID int64) (Review, error) {
				return Review{ID: reviewID, UserID: "u1"}, nil
			},
			deleteReview: func(_ context.Context, reviewID int64) error {
				deletedID = reviewID
				return nil
			},
		}, nil)

		err := uc.DeleteReview(ctx, Actor{}, 9)
		require.ErrorIs(t, err, ErrUnauthorized)
		err = uc.DeleteReview(ctx, Actor{UserID: "u2"}, 9)
		require.ErrorIs(t, err, ErrForbidden)
		require.NoError(t, uc.DeleteReview(ctx, Actor{UserID: "admin", IsAdmin: true}, 9))
		require.Equal(t, int64(9), deletedID)
	})
}
