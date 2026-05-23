package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

type WishlistRepo struct {
	*postgres.Postgres
}

func NewWishlistRepo(pg *postgres.Postgres) *WishlistRepo {
	return &WishlistRepo{Postgres: pg}
}

var _ repo.WishlistRepository = (*WishlistRepo)(nil)

func (r *WishlistRepo) ListWishlistItems(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error) {
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

	items := make([]entity.WishlistItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.WishlistItem{
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

func (r *WishlistRepo) StoreWishlistItem(ctx context.Context, item entity.WishlistItem) (entity.WishlistItem, error) {
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
		return entity.WishlistItem{}, fmt.Errorf("WishlistRepo.StoreWishlistItem: %w", err)
	}

	item.ID = model.ID
	item.CreatedAt = model.CreatedAt
	item.UpdatedAt = model.UpdatedAt
	return item, nil
}

func (r *WishlistRepo) UpdateWishlistItem(ctx context.Context, item entity.WishlistItem) (entity.WishlistItem, error) {
	if err := r.Gorm.WithContext(ctx).
		Model(&models.WishlistItem{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"title":       item.Title,
			"description": item.Description,
			"updated_at":  time.Now().UTC(),
		}).Error; err != nil {
		return entity.WishlistItem{}, fmt.Errorf("WishlistRepo.UpdateWishlistItem: %w", err)
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

func (r *WishlistRepo) StoreReview(ctx context.Context, review entity.Review) (entity.Review, error) {
	model := models.Review{
		ProductID: review.ProductID,
		OrderID:   review.OrderID,
		UserID:    review.UserID,
		Username:  review.Username,
		Rating:    review.Rating,
		Comment:   review.Comment,
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
		return entity.Review{}, err
	}

	review.ID = model.ID
	review.CreatedAt = model.CreatedAt
	review.UpdatedAt = model.UpdatedAt
	return review, nil
}
