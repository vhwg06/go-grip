package wishlist

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type UseCase struct {
	repo      repo.WishlistRepository
	orderRepo repo.OrderRepository
}

func New(wishlistRepo repo.WishlistRepository, orderRepo repo.OrderRepository) *UseCase {
	return &UseCase{repo: wishlistRepo, orderRepo: orderRepo}
}

var _ usecase.Wishlist = (*UseCase)(nil)

func (uc *UseCase) List(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error) {
	items, total, err := uc.repo.ListWishlistItems(ctx, page)
	if err != nil {
		return nil, 0, fmt.Errorf("WishlistUseCase.List - repo.ListWishlistItems: %w", err)
	}
	return items, total, nil
}

func (uc *UseCase) Create(ctx context.Context, actor entity.Actor, title, description string) (entity.WishlistItem, error) {
	if actor.UserID == "" {
		return entity.WishlistItem{}, entity.ErrUnauthorized
	}
	if title == "" {
		return entity.WishlistItem{}, entity.ErrInvalidInput
	}
	item := entity.WishlistItem{
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
		return entity.WishlistItem{}, fmt.Errorf("WishlistUseCase.Create - repo.StoreWishlistItem: %w", err)
	}
	return stored, nil
}

func (uc *UseCase) Update(ctx context.Context, actor entity.Actor, item entity.WishlistItem) (entity.WishlistItem, error) {
	if actor.UserID == "" {
		return entity.WishlistItem{}, entity.ErrUnauthorized
	}
	item.UpdatedAt = time.Now().UTC()
	updated, err := uc.repo.UpdateWishlistItem(ctx, item)
	if err != nil {
		return entity.WishlistItem{}, fmt.Errorf("WishlistUseCase.Update - repo.UpdateWishlistItem: %w", err)
	}
	return updated, nil
}

func (uc *UseCase) Delete(ctx context.Context, actor entity.Actor, itemID int64) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if err := uc.repo.DeleteWishlistItem(ctx, itemID); err != nil {
		return fmt.Errorf("WishlistUseCase.Delete - repo.DeleteWishlistItem: %w", err)
	}
	return nil
}

func (uc *UseCase) ToggleVote(ctx context.Context, actor entity.Actor, itemID int64) error {
	if actor.UserID == "" {
		return entity.ErrUnauthorized
	}
	if _, err := uc.repo.ToggleWishlistVote(ctx, itemID, actor.UserID); err != nil {
		return fmt.Errorf("WishlistUseCase.ToggleVote - repo.ToggleWishlistVote: %w", err)
	}
	return nil
}

func (uc *UseCase) CreateReview(ctx context.Context, actor entity.Actor, review entity.Review) (entity.Review, error) {
	if actor.UserID == "" {
		return entity.Review{}, entity.ErrUnauthorized
	}
	if review.ProductID == "" || review.Rating < 1 || review.Rating > 5 {
		return entity.Review{}, entity.ErrInvalidInput
	}
	if review.OrderID == "" {
		review.OrderID = fmt.Sprintf("no_order_%d_%s", time.Now().UnixNano(), actor.UserID)
	} else if uc.orderRepo != nil {
		order, err := uc.orderRepo.GetOrderByID(ctx, review.OrderID)
		if err != nil {
			return entity.Review{}, fmt.Errorf("WishlistUseCase.CreateReview - orderRepo.GetOrderByID: %w", err)
		}
		if order.Status != entity.OrderStatusDelivered {
			return entity.Review{}, entity.ErrForbidden
		}
	}

	review.UserID = actor.UserID
	review.Username = actor.Username
	review.CreatedAt = time.Now().UTC()
	review.UpdatedAt = time.Now().UTC()

	stored, err := uc.repo.StoreReview(ctx, review)
	if err != nil {
		return entity.Review{}, fmt.Errorf("WishlistUseCase.CreateReview - repo.StoreReview: %w", err)
	}
	return stored, nil
}

func (uc *UseCase) ListReviews(ctx context.Context, productID string) ([]entity.Review, error) {
	if productID == "" {
		return nil, entity.ErrInvalidInput
	}
	return uc.repo.ListReviews(ctx, productID)
}
