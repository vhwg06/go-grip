package persistent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func buildDetailRows(product entity.Product) []models.ProductDetail {
	details := make([]entity.ProductSpecItem, 0, len(product.Specs)+2)
	if product.SKU != "" {
		details = append(details, entity.ProductSpecItem{Key: "sku", Value: product.SKU})
	}
	if product.Brand != "" {
		details = append(details, entity.ProductSpecItem{Key: "brand", Value: product.Brand})
	}
	for _, spec := range product.Specs {
		key := strings.TrimSpace(spec.Key)
		value := strings.TrimSpace(spec.Value)
		if key == "" {
			continue
		}
		if key == "sku" || key == "brand" {
			continue
		}
		details = append(details, entity.ProductSpecItem{Key: key, Value: value})
	}

	rows := make([]models.ProductDetail, 0, len(details))
	for i, detail := range details {
		row := models.EntityToDetail(product.ID, detail)
		row.SortOrder = i
		rows = append(rows, row)
	}

	return rows
}

type AdminRepo struct {
	*postgres.Postgres
}

func NewAdminRepo(pg *postgres.Postgres) *AdminRepo {
	return &AdminRepo{Postgres: pg}
}

var _ repo.AdminRepository = (*AdminRepo)(nil)

func (r *AdminRepo) ListUsers(ctx context.Context, page entity.Pagination) ([]entity.User, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.User{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListUsers(count): %w", err)
	}

	normalized := page.Normalize()
	var rows []models.User
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListUsers(find): %w", err)
	}

	users := make([]entity.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, models.UserToEntity(row))
	}
	return users, int(total), nil
}

func (r *AdminRepo) UpdateUserStatus(ctx context.Context, userID string, status entity.UserStatus) error {
	if err := r.Gorm.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{"status": string(status), "updated_at": time.Now().UTC()}).Error; err != nil {
		return fmt.Errorf("AdminRepo.UpdateUserStatus: %w", err)
	}
	return nil
}

func (r *AdminRepo) UpdateUserPoints(ctx context.Context, userID string, points int) error {
	if err := r.Gorm.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{"points": points, "updated_at": time.Now().UTC()}).Error; err != nil {
		return fmt.Errorf("AdminRepo.UpdateUserPoints: %w", err)
	}
	return nil
}

func (r *AdminRepo) ListOrders(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Order, int, error) {
	db := r.Gorm.WithContext(ctx).Model(&models.Order{})

	if trimmed := strings.TrimSpace(query); trimmed != "" {
		like := "%" + strings.ToLower(trimmed) + "%"
		db = db.Where(
			"LOWER(order_id) LIKE ? OR LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(product_name) LIKE ?",
			like, like, like, like,
		)
	}
	if trimmedStatus := strings.TrimSpace(strings.ToLower(status)); trimmedStatus != "" && trimmedStatus != "all" {
		db = db.Where("LOWER(status) = ?", trimmedStatus)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListOrders(count): %w", err)
	}

	normalized := page.Normalize()
	var rows []models.Order
	if err := db.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListOrders(find): %w", err)
	}

	orders := make([]entity.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, models.OrderToEntity(row))
	}
	return orders, int(total), nil
}

func (r *AdminRepo) ListReviews(ctx context.Context, page entity.Pagination, query, status string) ([]entity.Review, repo.ReviewModerationStats, int, error) {
	db := r.Gorm.WithContext(ctx).Table("reviews r").
		Select("r.*, COALESCE(p.title, '') AS product_name").
		Joins("LEFT JOIN products p ON p.id::text = r.product_id")

	if trimmed := strings.TrimSpace(query); trimmed != "" {
		like := "%" + strings.ToLower(trimmed) + "%"
		db = db.Where("LOWER(COALESCE(p.title, '')) LIKE ? OR LOWER(r.username) LIKE ? OR LOWER(r.comment) LIKE ?", like, like, like)
	}
	if trimmedStatus := strings.TrimSpace(strings.ToUpper(status)); trimmedStatus != "" {
		db = db.Where("UPPER(r.status) = ?", trimmedStatus)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, repo.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(count): %w", err)
	}

	type reviewRow struct {
		models.Review
		ProductName string `gorm:"column:product_name"`
	}
	var rows []reviewRow
	normalized := page.Normalize()
	if err := db.Order("r.created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, repo.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(find): %w", err)
	}

	var statsRows []struct {
		Status string
		Count  int
	}
	if err := r.Gorm.WithContext(ctx).Table("reviews").
		Select("UPPER(status) AS status, COUNT(*) AS count").
		Group("UPPER(status)").
		Scan(&statsRows).Error; err != nil {
		return nil, repo.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(stats): %w", err)
	}

	stats := repo.ReviewModerationStats{}
	for _, row := range statsRows {
		switch row.Status {
		case string(entity.ReviewStatusPending):
			stats.Pending = row.Count
		case string(entity.ReviewStatusFeatured):
			stats.Featured = row.Count
		case string(entity.ReviewStatusHidden):
			stats.Hidden = row.Count
		}
	}

	items := make([]entity.Review, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.Review{
			ID:                 row.ID,
			ProductID:          row.ProductID,
			ProductName:        row.ProductName,
			OrderID:            row.OrderID,
			UserID:             row.UserID,
			Username:           row.Username,
			Rating:             row.Rating,
			Comment:            row.Comment,
			Status:             entity.ReviewStatus(strings.ToUpper(row.Status)),
			Attachments:        []string{},
			IsVerifiedPurchase: row.OrderID != "" && !strings.HasPrefix(row.OrderID, "no_order_"),
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
		})
	}

	return items, stats, int(total), nil
}

func (r *AdminRepo) GetOrderByID(ctx context.Context, orderID string) (entity.Order, error) {
	var row models.Order
	if err := r.Gorm.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Order{}, entity.ErrOrderNotFound
		}
		return entity.Order{}, fmt.Errorf("AdminRepo.GetOrderByID: %w", err)
	}
	return models.OrderToEntity(row), nil
}

func (r *AdminRepo) ListRefundRequests(ctx context.Context, status string) ([]entity.RefundRequest, error) {
	db := r.Gorm.WithContext(ctx).Model(&models.RefundRequest{})
	if trimmed := strings.TrimSpace(strings.ToLower(status)); trimmed != "" && trimmed != "all" {
		db = db.Where("LOWER(status) = ?", trimmed)
	}

	var rows []models.RefundRequest
	if err := db.Order("created_at DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListRefundRequests: %w", err)
	}

	items := make([]entity.RefundRequest, 0, len(rows))
	for _, row := range rows {
		item := entity.RefundRequest{
			ID:            row.ID,
			OrderID:       row.OrderID,
			UserID:        row.UserID,
			Username:      row.Username,
			Reason:        row.Reason,
			Status:        entity.RefundStatus(row.Status),
			AdminUsername: row.AdminUsername,
			AdminNote:     row.AdminNote,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
		if !row.ProcessedAt.IsZero() {
			processedAt := row.ProcessedAt
			item.ProcessedAt = &processedAt
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *AdminRepo) ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (entity.RefundRequest, error) {
	var result entity.RefundRequest

	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var refund models.RefundRequest
		if err := forUpdate(tx).Where("id = ?", refundID).First(&refund).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.ErrNotFound
			}
			return fmt.Errorf("AdminRepo.ProcessRefund(find refund): %w", err)
		}
		if strings.ToLower(refund.Status) != string(entity.RefundStatusPending) {
			return entity.ErrOrderStateConflict
		}

		var order models.Order
		if err := forUpdate(tx).Where("order_id = ?", refund.OrderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.ErrOrderNotFound
			}
			return fmt.Errorf("AdminRepo.ProcessRefund(find order): %w", err)
		}

		now := time.Now().UTC()
		nextRefundStatus := entity.RefundStatusRejected
		nextOrderStatus := entity.OrderStatusDelivered
		if approve {
			nextRefundStatus = entity.RefundStatusApproved
			nextOrderStatus = entity.OrderStatusRefunded
		}

		if err := tx.Model(&models.RefundRequest{}).
			Where("id = ?", refundID).
			Updates(map[string]any{
				"status":         string(nextRefundStatus),
				"admin_username": adminUsername,
				"admin_note":     note,
				"processed_at":   now,
				"updated_at":     now,
			}).Error; err != nil {
			return fmt.Errorf("AdminRepo.ProcessRefund(update refund): %w", err)
		}

		orderUpdates := map[string]any{
			"status":     string(nextOrderStatus),
			"updated_at": now,
		}
		if approve {
			orderUpdates["card_key"] = ""
		}
		if err := tx.Model(&models.Order{}).
			Where("order_id = ?", order.OrderID).
			Updates(orderUpdates).Error; err != nil {
			return fmt.Errorf("AdminRepo.ProcessRefund(update order): %w", err)
		}

		if approve {
			if order.UserID != "" && order.PointsUsed > 0 {
				if err := tx.Model(&models.User{}).
					Where("id = ?", order.UserID).
					Update("points", gorm.Expr("points + ?", order.PointsUsed)).Error; err != nil {
					return fmt.Errorf("AdminRepo.ProcessRefund(refund points): %w", err)
				}
			}

			cardIDs := splitCardIDs(order.CardIDs)
			if len(cardIDs) > 0 {
				if err := tx.Model(&models.Card{}).
					Where("id IN ?", cardIDs).
					Updates(map[string]any{
						"is_used":           false,
						"reserved_order_id": "",
						"reserved_at":       time.Time{},
						"used_at":           time.Time{},
					}).Error; err != nil {
					return fmt.Errorf("AdminRepo.ProcessRefund(reclaim cards): %w", err)
				}
			}

			if err := tx.Exec(`
				UPDATE products p
				SET
					stock_count = COALESCE(src.stock_count, 0),
					locked_count = COALESCE(src.locked_count, 0),
					sold_count = COALESCE(src.sold_count, 0),
					updated_at = NOW()
				FROM (
					SELECT
						product_id,
						COUNT(*) FILTER (WHERE is_used = FALSE AND (reserved_order_id = '' OR reserved_order_id IS NULL)) AS stock_count,
						COUNT(*) FILTER (WHERE is_used = FALSE AND reserved_order_id <> '') AS locked_count,
						COUNT(*) FILTER (WHERE is_used = TRUE) AS sold_count
					FROM cards
					WHERE product_id = ?
					GROUP BY product_id
				) src
				WHERE p.id::text = src.product_id
			`, order.ProductID).Error; err != nil {
				return fmt.Errorf("AdminRepo.ProcessRefund(rebuild aggregate): %w", err)
			}
		}

		result = entity.RefundRequest{
			ID:            refund.ID,
			OrderID:       refund.OrderID,
			UserID:        refund.UserID,
			Username:      refund.Username,
			Reason:        refund.Reason,
			Status:        nextRefundStatus,
			AdminUsername: adminUsername,
			AdminNote:     note,
			ProcessedAt:   &now,
			CreatedAt:     refund.CreatedAt,
			UpdatedAt:     now,
		}
		return nil
	})
	if err != nil {
		return entity.RefundRequest{}, err
	}

	return result, nil
}

func (r *AdminRepo) UpdateReviewStatus(ctx context.Context, reviewID int64, status entity.ReviewStatus) (entity.Review, error) {
	var row models.Review
	if err := r.Gorm.WithContext(ctx).Where("id = ?", reviewID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Review{}, entity.ErrNotFound
		}
		return entity.Review{}, fmt.Errorf("AdminRepo.UpdateReviewStatus(find): %w", err)
	}

	row.Status = string(status)
	row.UpdatedAt = time.Now().UTC()
	if err := r.Gorm.WithContext(ctx).Save(&row).Error; err != nil {
		return entity.Review{}, fmt.Errorf("AdminRepo.UpdateReviewStatus(save): %w", err)
	}

	return entity.Review{
		ID:        row.ID,
		ProductID: row.ProductID,
		OrderID:   row.OrderID,
		UserID:    row.UserID,
		Username:  row.Username,
		Rating:    row.Rating,
		Comment:   row.Comment,
		Status:    entity.ReviewStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *AdminRepo) BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status entity.ReviewStatus) (int, error) {
	result := r.Gorm.WithContext(ctx).Model(&models.Review{}).
		Where("id IN ?", reviewIDs).
		Updates(map[string]any{"status": string(status), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return 0, fmt.Errorf("AdminRepo.BulkUpdateReviewStatus: %w", result.Error)
	}
	return int(result.RowsAffected), nil
}

func (r *AdminRepo) DeleteReview(ctx context.Context, reviewID int64) error {
	result := r.Gorm.WithContext(ctx).Delete(&models.Review{}, reviewID)
	if result.Error != nil {
		return fmt.Errorf("AdminRepo.DeleteReview: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (r *AdminRepo) UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var order models.Order
		if err := forUpdate(tx).Where("order_id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.ErrOrderNotFound
			}
			return fmt.Errorf("AdminRepo.UpdateOrderStatus(find): %w", err)
		}

		current := entity.OrderStatus(order.Status)
		now := time.Now().UTC()
		updates := map[string]any{
			"status":     string(status),
			"updated_at": now,
		}

		switch status {
		case entity.OrderStatusPaid:
			if current != entity.OrderStatusPending {
				return entity.ErrOrderStateConflict
			}
			updates["paid_at"] = now
		case entity.OrderStatusDelivered:
			if current != entity.OrderStatusPending && current != entity.OrderStatusPaid {
				return entity.ErrOrderStateConflict
			}
			if order.PaidAt.IsZero() {
				updates["paid_at"] = now
			}
			updates["delivered_at"] = now
		case entity.OrderStatusCancelled:
			if current != entity.OrderStatusPending && current != entity.OrderStatusPaid {
				return entity.ErrOrderStateConflict
			}
		default:
			return entity.ErrInvalidInput
		}

		if err := tx.Model(&models.Order{}).Where("order_id = ?", orderID).Updates(updates).Error; err != nil {
			return fmt.Errorf("AdminRepo.UpdateOrderStatus(update): %w", err)
		}

		if status == entity.OrderStatusCancelled {
			if order.UserID != "" && order.PointsUsed > 0 {
				if err := tx.Model(&models.User{}).
					Where("id = ?", order.UserID).
					Update("points", gorm.Expr("points + ?", order.PointsUsed)).Error; err != nil {
					return fmt.Errorf("AdminRepo.UpdateOrderStatus(refund points): %w", err)
				}
			}

			if err := tx.Model(&models.Card{}).
				Where("reserved_order_id = ?", orderID).
				Updates(map[string]any{
					"reserved_order_id": "",
					"reserved_at":       time.Time{},
				}).Error; err != nil {
				return fmt.Errorf("AdminRepo.UpdateOrderStatus(release cards): %w", err)
			}
		}

		return nil
	})
}

func (r *AdminRepo) DeleteOrder(ctx context.Context, orderID string) error {
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", orderID).Delete(&models.RefundRequest{}).Error; err != nil {
			return fmt.Errorf("AdminRepo.DeleteOrder(delete refunds): %w", err)
		}
		if err := tx.Model(&models.Card{}).
			Where("reserved_order_id = ?", orderID).
			Updates(map[string]any{
				"reserved_order_id": "",
				"reserved_at":       time.Time{},
			}).Error; err != nil {
			return fmt.Errorf("AdminRepo.DeleteOrder(release cards): %w", err)
		}
		if err := tx.Where("order_id = ?", orderID).Delete(&models.Order{}).Error; err != nil {
			return fmt.Errorf("AdminRepo.DeleteOrder(delete order): %w", err)
		}
		return nil
	})
}

func (r *AdminRepo) StoreSetting(ctx context.Context, setting entity.Setting) error {
	model := models.EntityToSetting(setting)
	if model.UpdatedAt.IsZero() {
		model.UpdatedAt = time.Now().UTC()
	}
	if err := r.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
		}).
		Create(&model).Error; err != nil {
		return fmt.Errorf("AdminRepo.StoreSetting: %w", err)
	}
	return nil
}

func splitCardIDs(raw string) []int64 {
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func (r *AdminRepo) ListSettings(ctx context.Context) ([]entity.Setting, error) {
	var rows []models.Setting
	if err := r.Gorm.WithContext(ctx).Order("\"key\" ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListSettings: %w", err)
	}

	settings := make([]entity.Setting, 0, len(rows))
	for _, row := range rows {
		settings = append(settings, models.SettingToEntity(row))
	}
	return settings, nil
}

func (r *AdminRepo) DeleteSetting(ctx context.Context, key string) error {
	if err := r.Gorm.WithContext(ctx).Where("key = ?", key).Delete(&models.Setting{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteSetting: %w", err)
	}
	return nil
}

func (r *AdminRepo) RebuildProductAggregates(ctx context.Context) error {
	sql := `
		UPDATE products p
		SET
			stock_count = COALESCE(src.stock_count, 0),
			locked_count = COALESCE(src.locked_count, 0),
			sold_count = COALESCE(src.sold_count, 0),
			updated_at = NOW()
		FROM (
			SELECT
				product_id,
				COUNT(*) FILTER (WHERE is_used = FALSE AND (reserved_order_id = '' OR reserved_order_id IS NULL)) AS stock_count,
				COUNT(*) FILTER (WHERE is_used = FALSE AND reserved_order_id <> '') AS locked_count,
				COUNT(*) FILTER (WHERE is_used = TRUE) AS sold_count
			FROM cards
			GROUP BY product_id
		) src
		WHERE p.id::text = src.product_id
	`
	if err := r.Gorm.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("AdminRepo.RebuildProductAggregates: %w", err)
	}
	return nil
}

func (r *AdminRepo) ListProducts(ctx context.Context, page entity.Pagination) ([]entity.Product, int, error) {
	query := r.Gorm.WithContext(ctx).Model(&models.Product{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListProducts(count): %w", err)
	}
	normalized := page.Normalize()
	var rows []models.Product
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListProducts(find): %w", err)
	}
	items := make([]entity.Product, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.ProductToEntity(row))
	}
	return items, int(total), nil
}

func (r *AdminRepo) GetProduct(ctx context.Context, productID string) (entity.Product, error) {
	var row models.Product
	if err := r.Gorm.WithContext(ctx).
		Where("id = ?", productID).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, entity.ErrNotFound
		}
		return entity.Product{}, fmt.Errorf("AdminRepo.GetProduct: %w", err)
	}

	result := models.ProductToEntity(row)

	var detailRows []models.ProductDetail
	if err := r.Gorm.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&detailRows).Error; err != nil {
		return entity.Product{}, fmt.Errorf("AdminRepo.GetProduct(specs): %w", err)
	}
	result.Specs = make([]entity.ProductSpecItem, 0, len(detailRows))
	for _, detail := range detailRows {
		switch detail.Key {
		case "sku":
			result.SKU = detail.Value
		case "brand":
			result.Brand = detail.Value
		default:
			result.Specs = append(result.Specs, models.DetailToEntity(detail))
		}
	}

	return result, nil
}

func (r *AdminRepo) UpsertProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	if _, err := uuid.Parse(strings.TrimSpace(product.ID)); err != nil {
		product.ID = uuid.NewString()
	}

	if strings.TrimSpace(product.SKU) == "" {
		product.SKU = product.ID
	}

	var existing models.Product
	err := r.Gorm.WithContext(ctx).Where("id = ?", product.ID).First(&existing).Error
	if err != nil {
		product.IsActive = true
	} else {
		product.IsActive = existing.IsActive
	}

	model := models.EntityToProduct(product)
	now := time.Now().UTC()
	if model.CreatedAt.IsZero() {
		model.CreatedAt = now
	}
	model.UpdatedAt = now

	if err := r.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"title", "sku", "description", "price", "compare_price", "category", "image", "is_hot", "is_active",
				"is_shared", "sort_order", "purchase_limit", "purchase_warning", "visibility_level",
				"stock_count", "locked_count", "sold_count", "rating", "review_count", "updated_at",
			}),
		}).
		Create(&model).Error; err != nil {
		return entity.Product{}, fmt.Errorf("AdminRepo.UpsertProduct: %w", err)
	}

	// Persist specs
	if err := r.Gorm.WithContext(ctx).Where("product_id = ?", product.ID).Delete(&models.ProductDetail{}).Error; err != nil {
		return entity.Product{}, fmt.Errorf("AdminRepo.UpsertProduct(delete specs): %w", err)
	}

	detailRows := buildDetailRows(product)
	if len(detailRows) > 0 {
		if err := r.Gorm.WithContext(ctx).Create(&detailRows).Error; err != nil {
			return entity.Product{}, fmt.Errorf("AdminRepo.UpsertProduct(create specs): %w", err)
		}
	}

	result := models.ProductToEntity(model)
	result.SKU = product.SKU
	result.Brand = product.Brand
	result.Specs = make([]entity.ProductSpecItem, 0, len(detailRows))
	for _, detail := range detailRows {
		if detail.Key == "sku" || detail.Key == "brand" {
			continue
		}
		result.Specs = append(result.Specs, models.DetailToEntity(detail))
	}

	return result, nil
}

func (r *AdminRepo) DeleteProduct(ctx context.Context, productID string) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", productID).Delete(&models.Product{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteProduct: %w", err)
	}
	return nil
}

func (r *AdminRepo) ListCategories(ctx context.Context) ([]entity.Category, error) {
	var rows []models.Category
	if err := r.Gorm.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListCategories: %w", err)
	}
	items := make([]entity.Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.CategoryToEntity(row))
	}
	return items, nil
}

func (r *AdminRepo) UpsertCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	model := models.EntityToCategory(category)
	now := time.Now().UTC()
	if model.CreatedAt.IsZero() {
		model.CreatedAt = now
	}
	model.UpdatedAt = now
	if err := r.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "parent_id", "sort_order", "is_active", "updated_at"}),
		}).
		Create(&model).Error; err != nil {
		return entity.Category{}, fmt.Errorf("AdminRepo.UpsertCategory: %w", err)
	}
	return models.CategoryToEntity(model), nil
}

func (r *AdminRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", categoryID).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteCategory: %w", err)
	}
	return nil
}

func (r *AdminRepo) ImportCards(ctx context.Context, productID string, keys []string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	now := time.Now().UTC()
	modelsToCreate := make([]models.Card, 0, len(keys))
	for _, key := range keys {
		modelsToCreate = append(modelsToCreate, models.Card{
			ProductID: productID,
			CardKey:   key,
			IsUsed:    false,
			CreatedAt: now,
		})
	}
	if err := r.Gorm.WithContext(ctx).Create(&modelsToCreate).Error; err != nil {
		return 0, fmt.Errorf("AdminRepo.ImportCards: %w", err)
	}
	return len(modelsToCreate), nil
}

func (r *AdminRepo) ListCards(ctx context.Context, productID string) ([]entity.Card, error) {
	db := r.Gorm.WithContext(ctx).Model(&models.Card{})
	if trimmed := strings.TrimSpace(productID); trimmed != "" {
		db = db.Where("product_id = ?", trimmed)
	}

	var rows []models.Card
	if err := db.Order("created_at DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListCards: %w", err)
	}

	items := make([]entity.Card, 0, len(rows))
	for _, row := range rows {
		item := entity.Card{
			ID:              row.ID,
			ProductID:       row.ProductID,
			CardKey:         row.CardKey,
			IsUsed:          row.IsUsed,
			ReservedOrderID: row.ReservedOrderID,
			CreatedAt:       row.CreatedAt,
		}
		if !row.ReservedAt.IsZero() {
			reservedAt := row.ReservedAt
			item.ReservedAt = &reservedAt
		}
		if !row.ExpiresAt.IsZero() {
			expiresAt := row.ExpiresAt
			item.ExpiresAt = &expiresAt
		}
		if !row.UsedAt.IsZero() {
			usedAt := row.UsedAt
			item.UsedAt = &usedAt
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *AdminRepo) CreateCard(ctx context.Context, productID, cardKey string) (entity.Card, error) {
	now := time.Now().UTC()
	model := models.Card{
		ProductID: productID,
		CardKey:   cardKey,
		IsUsed:    false,
		CreatedAt: now,
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.Card{}, fmt.Errorf("AdminRepo.CreateCard: %w", err)
	}
	return entity.Card{
		ID:        model.ID,
		ProductID: productID,
		CardKey:   cardKey,
		IsUsed:    false,
		CreatedAt: now,
	}, nil
}

func (r *AdminRepo) DeleteCard(ctx context.Context, cardID int64) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", cardID).Delete(&models.Card{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteCard: %w", err)
	}
	return nil
}
