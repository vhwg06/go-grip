package wishlist

import (
	"context"
	"fmt"
	"time"

	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

// WishlistUseCase defines business operations for Wishlist & Reviews.
type WishlistUseCase interface {
	List(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error)
	Create(ctx context.Context, actor Actor, title, description string) (WishlistItem, error)
	Update(ctx context.Context, actor Actor, item WishlistItem) (WishlistItem, error)
	Delete(ctx context.Context, actor Actor, itemID int64) error
	ToggleVote(ctx context.Context, actor Actor, itemID int64) error
	CreateReview(ctx context.Context, actor Actor, review Review) (Review, error)
	ListReviews(ctx context.Context, productID string) ([]Review, error)
}

type wishlistUseCase struct {
	repo      WishlistRepo
	orderRepo OrderReader
}

// NewWishlistUseCase constructs a new WishlistUseCase instance.
func NewWishlistUseCase(wishlistRepo WishlistRepo, orderRepo OrderReader) WishlistUseCase {
	return &wishlistUseCase{repo: wishlistRepo, orderRepo: orderRepo}
}

func (uc *wishlistUseCase) List(ctx context.Context, page pagination.Pagination) ([]WishlistItem, int, error) {
	items, total, err := uc.repo.ListWishlistItems(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("WishlistUseCase.List - repo.ListWishlistItems: %w", err)
	}
	return items, total, nil
}

func (uc *wishlistUseCase) Create(ctx context.Context, actor Actor, title, description string) (WishlistItem, error) {
	if actor.UserID == "" {
		return WishlistItem{}, ErrUnauthorized
	}
	if title == "" {
		return WishlistItem{}, ErrInvalidInput
	}
	item := WishlistItem{
		Title:       title,
		Description: description,
		UserID:      actor.UserID,
		Username:    actor.Username,
		VoteCount:   0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	stored, err := uc.repo.StoreWishlistItem(ctx, item)
	if err != nil {
		return WishlistItem{}, fmt.Errorf("WishlistUseCase.Create - repo.StoreWishlistItem: %w", err)
	}
	return stored, nil
}

func (uc *wishlistUseCase) Update(ctx context.Context, actor Actor, item WishlistItem) (WishlistItem, error) {
	if actor.UserID == "" {
		return WishlistItem{}, ErrUnauthorized
	}
	item.UpdatedAt = time.Now().UTC()
	updated, err := uc.repo.UpdateWishlistItem(ctx, item)
	if err != nil {
		return WishlistItem{}, fmt.Errorf("WishlistUseCase.Update - repo.UpdateWishlistItem: %w", err)
	}
	return updated, nil
}

func (uc *wishlistUseCase) Delete(ctx context.Context, actor Actor, itemID int64) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if err := uc.repo.DeleteWishlistItem(ctx, itemID); err != nil {
		return fmt.Errorf("WishlistUseCase.Delete - repo.DeleteWishlistItem: %w", err)
	}
	return nil
}

func (uc *wishlistUseCase) ToggleVote(ctx context.Context, actor Actor, itemID int64) error {
	if actor.UserID == "" {
		return ErrUnauthorized
	}
	if _, err := uc.repo.ToggleWishlistVote(ctx, itemID, actor.UserID); err != nil {
		return fmt.Errorf("WishlistUseCase.ToggleVote - repo.ToggleWishlistVote: %w", err)
	}
	return nil
}

func (uc *wishlistUseCase) CreateReview(ctx context.Context, actor Actor, review Review) (Review, error) {
	if actor.UserID == "" {
		return Review{}, ErrUnauthorized
	}
	if review.ProductID == "" || review.Rating < 1 || review.Rating > 5 {
		return Review{}, ErrInvalidInput
	}
	if review.OrderID == "" {
		review.OrderID = fmt.Sprintf("no_order_%d_%s", time.Now().UnixNano(), actor.UserID)
	} else if uc.orderRepo != nil {
		order, err := uc.orderRepo.GetOrderByID(ctx, review.OrderID)
		if err != nil {
			return Review{}, fmt.Errorf("WishlistUseCase.CreateReview - orderRepo.GetOrderByID: %w", err)
		}
		if order.Status != ordermodule.OrderStatusDelivered {
			return Review{}, ErrForbidden
		}
	}

	review.UserID = actor.UserID
	review.Username = actor.Username
	review.CreatedAt = time.Now().UTC()
	review.UpdatedAt = time.Now().UTC()

	stored, err := uc.repo.StoreReview(ctx, review)
	if err != nil {
		return Review{}, fmt.Errorf("WishlistUseCase.CreateReview - repo.StoreReview: %w", err)
	}
	return stored, nil
}

func (uc *wishlistUseCase) ListReviews(ctx context.Context, productID string) ([]Review, error) {
	if productID == "" {
		return nil, ErrInvalidInput
	}
	return uc.repo.ListReviews(ctx, productID)
}
