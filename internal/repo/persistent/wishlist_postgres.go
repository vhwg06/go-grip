package persistent

import (
	"context"
	"fmt"
	"time"

	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type WishlistRepo struct {
	*postgres.Postgres
}

func NewWishlistRepo(pg *postgres.Postgres) *WishlistRepo {
	return &WishlistRepo{Postgres: pg}
}

var _ wishlistmodule.WishlistRepo = (*WishlistRepo)(nil)

func (r *WishlistRepo) ListWishlistItems(ctx context.Context, page pagination.Pagination) ([]wishlistmodule.WishlistItem, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.WishlistItem{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("WishlistRepo.ListWishlistItems(count): %w", err)
	}
	normalized := page.Normalize()
	var rows []models.WishlistItem
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("WishlistRepo.ListWishlistItems(find): %w", err)
	}

	items := make([]wishlistmodule.WishlistItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, wishlistmodule.WishlistItem{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			UserID:      row.UserID,
			Username:    row.Username,
			VoteCount:   row.VoteCount,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}
	return items, int(total), nil
}

func (r *WishlistRepo) StoreWishlistItem(ctx context.Context, item wishlistmodule.WishlistItem) (wishlistmodule.WishlistItem, error) {
	model := models.WishlistItem{
		Title:       item.Title,
		Description: item.Description,
		UserID:      item.UserID,
		Username:    item.Username,
		VoteCount:   item.VoteCount,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return wishlistmodule.WishlistItem{}, fmt.Errorf("WishlistRepo.StoreWishlistItem: %w", err)
	}

	item.ID = model.ID
	item.CreatedAt = model.CreatedAt
	item.UpdatedAt = model.UpdatedAt
	return item, nil
}

func (r *WishlistRepo) UpdateWishlistItem(ctx context.Context, item wishlistmodule.WishlistItem) (wishlistmodule.WishlistItem, error) {
	if err := r.Gorm.WithContext(ctx).
		Model(&models.WishlistItem{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"title":       item.Title,
			"description": item.Description,
			"updated_at":  time.Now().UTC(),
		}).Error; err != nil {
		return wishlistmodule.WishlistItem{}, fmt.Errorf("WishlistRepo.UpdateWishlistItem: %w", err)
	}
	return item, nil
}

func (r *WishlistRepo) DeleteWishlistItem(ctx context.Context, itemID int64) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", itemID).Delete(&models.WishlistItem{}).Error; err != nil {
		return fmt.Errorf("WishlistRepo.DeleteWishlistItem: %w", err)
	}
	return nil
}

func (r *WishlistRepo) ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error) {
	var added bool
	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var vote models.WishlistVote
		if err := tx.Where("item_id = ? AND user_id = ?", itemID, userID).First(&vote).Error; err == nil {
			if err := tx.Delete(&vote).Error; err != nil {
				return fmt.Errorf("WishlistRepo.ToggleWishlistVote(delete): %w", err)
			}
			if err := tx.Model(&models.WishlistItem{}).
				Where("id = ? AND vote_count > 0", itemID).
				Update("vote_count", gorm.Expr("vote_count - 1")).Error; err != nil {
				return fmt.Errorf("WishlistRepo.ToggleWishlistVote(decrement): %w", err)
			}
			added = false
			return nil
		}

		newVote := models.WishlistVote{ItemID: itemID, UserID: userID, CreatedAt: time.Now().UTC()}
		if err := tx.Create(&newVote).Error; err != nil {
			return fmt.Errorf("WishlistRepo.ToggleWishlistVote(insert): %w", err)
		}
		if err := tx.Model(&models.WishlistItem{}).
			Where("id = ?", itemID).
			Update("vote_count", gorm.Expr("vote_count + 1")).Error; err != nil {
			return fmt.Errorf("WishlistRepo.ToggleWishlistVote(increment): %w", err)
		}
		added = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return added, nil
}

func (r *WishlistRepo) StoreReview(ctx context.Context, review wishlistmodule.Review) (wishlistmodule.Review, error) {
	model := models.Review{
		ProductID: review.ProductID,
		OrderID:   review.OrderID,
		UserID:    review.UserID,
		Username:  review.Username,
		Rating:    review.Rating,
		Comment:   review.Comment,
		Status:    string(wishlistmodule.ReviewStatusPending),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("WishlistRepo.StoreReview(create): %w", err)
		}

		var agg struct {
			Count int64
			Avg   float64
		}
		if err := tx.Model(&models.Review{}).
			Where("product_id = ?", review.ProductID).
			Select("COUNT(*) as count, COALESCE(AVG(rating), 0) as avg").
			Scan(&agg).Error; err != nil {
			return fmt.Errorf("WishlistRepo.StoreReview(aggregate): %w", err)
		}

		if err := tx.Model(&models.Product{}).
			Where("id = ?", review.ProductID).
			Updates(map[string]any{
				"rating":       agg.Avg,
				"review_count": int(agg.Count),
				"updated_at":   time.Now().UTC(),
			}).Error; err != nil {
			return fmt.Errorf("WishlistRepo.StoreReview(update product): %w", err)
		}
		return nil
	})
	if err != nil {
		return wishlistmodule.Review{}, err
	}

	review.ID = model.ID
	review.CreatedAt = model.CreatedAt
	review.UpdatedAt = model.UpdatedAt
	review.Status = wishlistmodule.ReviewStatus(model.Status)
	return review, nil
}

func (r *WishlistRepo) ListReviews(ctx context.Context, productID string) ([]wishlistmodule.Review, error) {
	var rows []models.Review
	if err := r.Gorm.WithContext(ctx).
		Where("product_id = ?", productID).
		Where("status IN ?", []string{string(wishlistmodule.ReviewStatusApproved), string(wishlistmodule.ReviewStatusFeatured)}).
		Order("CASE WHEN status = 'FEATURED' THEN 0 ELSE 1 END, created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("WishlistRepo.ListReviews: %w", err)
	}

	reviews := make([]wishlistmodule.Review, 0, len(rows))
	for _, row := range rows {
		reviews = append(reviews, wishlistmodule.Review{
			ID:        row.ID,
			ProductID: row.ProductID,
			OrderID:   row.OrderID,
			UserID:    row.UserID,
			Username:  row.Username,
			Rating:    row.Rating,
			Comment:   row.Comment,
			Status:    wishlistmodule.ReviewStatus(row.Status),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		})
	}
	return reviews, nil
}
