package persistent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/internal/shared/pagination"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func buildDetailRows(product catalogmodule.Product) []models.ProductDetail {
	details := make([]catalogmodule.ProductSpecItem, 0, len(product.Specs)+2)
	if product.SKU != "" {
		details = append(details, catalogmodule.ProductSpecItem{Key: "sku", Value: product.SKU})
	}
	if product.Brand != "" {
		details = append(details, catalogmodule.ProductSpecItem{Key: "brand", Value: product.Brand})
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
		details = append(details, catalogmodule.ProductSpecItem{Key: key, Value: value})
	}

	rows := make([]models.ProductDetail, 0, len(details))
	for i, detail := range details {
		row := models.ModuleToDetail(product.ID, catalogmodule.ProductSpecItem{Key: detail.Key, Value: detail.Value})
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

func (r *AdminRepo) ListUsers(ctx context.Context, page pagination.Pagination) ([]usermodule.User, int, error) {
	var q string
	if val := ctx.Value("query"); val != nil {
		if s, ok := val.(string); ok {
			q = s
		}
	}
	var role string
	if val := ctx.Value("role"); val != nil {
		if s, ok := val.(string); ok {
			role = s
		}
	}

	query := r.Gorm.WithContext(ctx).Model(&models.User{})

	// Exclude admins if role is customer, OR if role is not user/admin and q does not contain "admin"
	if role == "customer" || (role != "user" && role != "admin" && (q == "" || !strings.Contains(strings.ToLower(q), "admin"))) {
		query = query.Where("is_admin = ?", false).Where("role != ? OR role IS NULL", "Administrator")
	}

	if trimmed := strings.TrimSpace(q); trimmed != "" {
		like := "%" + strings.ToLower(trimmed) + "%"
		query = query.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(display_name) LIKE ? OR id = ?", like, like, like, trimmed)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListUsers(count): %w", err)
	}

	normalized := page.Normalize()
	var rows []models.User
	if err := query.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("AdminRepo.ListUsers(find): %w", err)
	}

	users := make([]usermodule.User, 0, len(rows))
	for _, row := range rows {
		u := models.UserToModule(row)
		u.CustomerID = &u.ID
		if role == "customer" || (role != "user" && q != "" && !strings.Contains(strings.ToLower(q), "admin")) {
			var oCount int64
			var refCount int64
			var revCount int64

			r.Gorm.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", u.ID).Count(&oCount)
			r.Gorm.WithContext(ctx).Model(&models.RefundRequest{}).Where("user_id = ?", u.ID).Count(&refCount)
			r.Gorm.WithContext(ctx).Model(&models.Review{}).Where("user_id = ?", u.ID).Count(&revCount)

			oVal := int(oCount)
			refVal := int(refCount)
			revVal := int(revCount)

			u.OrderCount = &oVal
			u.RefundCount = &refVal
			u.ReviewCount = &revVal
		}
		users = append(users, u)
	}
	return users, int(total), nil
}

func (r *AdminRepo) UpdateUserStatus(ctx context.Context, userID string, status usermodule.UserStatus) error {
	if err := r.Gorm.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{"status": string(status), "updated_at": time.Now().UTC()}).Error; err != nil {
		return fmt.Errorf("AdminRepo.UpdateUserStatus: %w", err)
	}
	return nil
}

func (r *AdminRepo) ListOrders(ctx context.Context, page pagination.Pagination, query, status string) ([]ordermodule.Order, int, error) {
	db := r.Gorm.WithContext(ctx).Model(&models.Order{})

	if trimmed := strings.TrimSpace(query); trimmed != "" {
		like := "%" + strings.ToLower(trimmed) + "%"
		db = db.Where(
			"LOWER(order_id) LIKE ? OR LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(product_name) LIKE ? OR user_id = ?",
			like, like, like, like, trimmed,
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

	orders := make([]ordermodule.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, models.OrderToModule(row))
	}
	return orders, int(total), nil
}

func (r *AdminRepo) ListReviews(ctx context.Context, page pagination.Pagination, query, status string) ([]wishlistmodule.Review, wishlistmodule.ReviewModerationStats, int, error) {
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
		return nil, wishlistmodule.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(count): %w", err)
	}

	type reviewRow struct {
		models.Review
		ProductName string `gorm:"column:product_name"`
	}
	var rows []reviewRow
	normalized := page.Normalize()
	if err := db.Order("r.created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, wishlistmodule.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(find): %w", err)
	}

	var statsRows []struct {
		Status string
		Count  int
	}
	if err := r.Gorm.WithContext(ctx).Table("reviews").
		Select("UPPER(status) AS status, COUNT(*) AS count").
		Group("UPPER(status)").
		Scan(&statsRows).Error; err != nil {
		return nil, wishlistmodule.ReviewModerationStats{}, 0, fmt.Errorf("AdminRepo.ListReviews(stats): %w", err)
	}

	stats := wishlistmodule.ReviewModerationStats{}
	for _, row := range statsRows {
		switch row.Status {
		case string(wishlistmodule.ReviewStatusPending):
			stats.Pending = row.Count
		case string(wishlistmodule.ReviewStatusFeatured):
			stats.Featured = row.Count
		case string(wishlistmodule.ReviewStatusHidden):
			stats.Hidden = row.Count
		}
	}

	items := make([]wishlistmodule.Review, 0, len(rows))
	for _, row := range rows {
		items = append(items, wishlistmodule.Review{
			ID:                 row.ID,
			ProductID:          row.ProductID,
			ProductName:        row.ProductName,
			OrderID:            row.OrderID,
			UserID:             row.UserID,
			Username:           row.Username,
			Rating:             row.Rating,
			Comment:            row.Comment,
			Status:             wishlistmodule.ReviewStatus(strings.ToUpper(row.Status)),
			Attachments:        []string{},
			IsVerifiedPurchase: row.OrderID != "" && !strings.HasPrefix(row.OrderID, "no_order_"),
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
		})
	}

	return items, stats, int(total), nil
}

func (r *AdminRepo) GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error) {
	var row models.Order
	if err := r.Gorm.WithContext(ctx).Where("order_id = ?", orderID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ordermodule.Order{}, ordermodule.ErrNotFound
		}
		return ordermodule.Order{}, fmt.Errorf("AdminRepo.GetOrderByID: %w", err)
	}
	return models.OrderToModule(row), nil
}

func (r *AdminRepo) ListRefundRequests(ctx context.Context, status string) ([]ordermodule.RefundRequest, error) {
	db := r.Gorm.WithContext(ctx).Table("refund_requests")
	if trimmed := strings.TrimSpace(strings.ToLower(status)); trimmed != "" && trimmed != "all" {
		db = db.Where("LOWER(refund_requests.status) = ?", trimmed)
	}

	type refundRow struct {
		models.RefundRequest
		ProductName string `gorm:"column:product_name"`
		Amount      int64  `gorm:"column:amount"`
		TradeNo     string `gorm:"column:trade_no"`
		OrderStatus string `gorm:"column:order_status"`
	}

	var rows []refundRow
	if err := db.Select("refund_requests.*, orders.product_name, orders.amount, orders.trade_no, orders.status AS order_status").
		Joins("LEFT JOIN orders ON orders.order_id = refund_requests.order_id").
		Order("refund_requests.created_at DESC, refund_requests.id DESC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListRefundRequests: %w", err)
	}

	items := make([]ordermodule.RefundRequest, 0, len(rows))
	for _, row := range rows {
		item := ordermodule.RefundRequest{
			ID:            row.ID,
			OrderID:       row.OrderID,
			UserID:        row.UserID,
			Username:      row.Username,
			Reason:        row.Reason,
			Status:        ordermodule.RefundStatus(row.Status),
			AdminUsername: row.AdminUsername,
			AdminNote:     row.AdminNote,
			ProductName:   row.ProductName,
			Amount:        ordermodule.Amount(row.Amount),
			TradeNo:       row.TradeNo,
			OrderStatus:   row.OrderStatus,
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

func (r *AdminRepo) GetRefundRequest(ctx context.Context, refundID int64) (ordermodule.RefundRequest, error) {
	type refundRow struct {
		models.RefundRequest
		ProductName string `gorm:"column:product_name"`
		Amount      int64  `gorm:"column:amount"`
		TradeNo     string `gorm:"column:trade_no"`
		OrderStatus string `gorm:"column:order_status"`
	}

	var row refundRow
	if err := r.Gorm.WithContext(ctx).Table("refund_requests").
		Select("refund_requests.*, orders.product_name, orders.amount, orders.trade_no, orders.status AS order_status").
		Joins("LEFT JOIN orders ON orders.order_id = refund_requests.order_id").
		Where("refund_requests.id = ?", refundID).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ordermodule.RefundRequest{}, usermodule.ErrNotFound
		}
		return ordermodule.RefundRequest{}, fmt.Errorf("AdminRepo.GetRefundRequest: %w", err)
	}

	item := ordermodule.RefundRequest{
		ID:            row.ID,
		OrderID:       row.OrderID,
		UserID:        row.UserID,
		Username:      row.Username,
		Reason:        row.Reason,
		Status:        ordermodule.RefundStatus(row.Status),
		AdminUsername: row.AdminUsername,
		AdminNote:     row.AdminNote,
		ProductName:   row.ProductName,
		Amount:        ordermodule.Amount(row.Amount),
		TradeNo:       row.TradeNo,
		OrderStatus:   row.OrderStatus,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
	if !row.ProcessedAt.IsZero() {
		processedAt := row.ProcessedAt
		item.ProcessedAt = &processedAt
	}
	return item, nil
}

func (r *AdminRepo) GetOrderRefundStatus(ctx context.Context, orderID string) (ordermodule.RefundRequest, error) {
	var row models.RefundRequest
	if err := r.Gorm.WithContext(ctx).Model(&models.RefundRequest{}).
		Where("order_id = ? AND status = ?", orderID, "pending").
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ordermodule.RefundRequest{}, usermodule.ErrNotFound
		}
		return ordermodule.RefundRequest{}, fmt.Errorf("AdminRepo.GetOrderRefundStatus: %w", err)
	}

	return ordermodule.RefundRequest{
		ID:      row.ID,
		OrderID: row.OrderID,
		Status:  ordermodule.RefundStatus(row.Status),
	}, nil
}

func (r *AdminRepo) ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (ordermodule.RefundRequest, error) {
	var result ordermodule.RefundRequest

	err := withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var refund models.RefundRequest
		if err := forUpdate(tx).Where("id = ?", refundID).First(&refund).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return usermodule.ErrNotFound
			}
			return fmt.Errorf("AdminRepo.ProcessRefund(find refund): %w", err)
		}
		if strings.ToLower(refund.Status) != string(ordermodule.RefundStatusPending) {
			return ordermodule.ErrInvalidInput
		}

		var order models.Order
		if err := forUpdate(tx).Where("order_id = ?", refund.OrderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ordermodule.ErrNotFound
			}
			return fmt.Errorf("AdminRepo.ProcessRefund(find order): %w", err)
		}

		now := time.Now().UTC()
		nextRefundStatus := ordermodule.RefundStatusRejected
		nextOrderStatus := ordermodule.OrderStatusDelivered
		if approve {
			nextRefundStatus = ordermodule.RefundStatusApproved
			nextOrderStatus = ordermodule.OrderStatusRefunded
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
		if err := tx.Model(&models.Order{}).
			Where("order_id = ?", order.OrderID).
			Updates(orderUpdates).Error; err != nil {
			return fmt.Errorf("AdminRepo.ProcessRefund(update order): %w", err)
		}

		if approve {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", order.ProductID).
				Updates(map[string]any{
					"stock_count": gorm.Expr("stock_count + ?", order.Quantity),
					"sold_count":  gorm.Expr("sold_count - ?", order.Quantity),
				}).Error; err != nil {
				return fmt.Errorf("AdminRepo.ProcessRefund(restore product stock): %w", err)
			}

		}

		result = ordermodule.RefundRequest{
			ID:            refund.ID,
			OrderID:       refund.OrderID,
			UserID:        refund.UserID,
			Username:      refund.Username,
			Reason:        refund.Reason,
			Status:        nextRefundStatus,
			AdminUsername: adminUsername,
			AdminNote:     note,
			ProductName:   order.ProductName,
			Amount:        ordermodule.Amount(order.Amount),
			TradeNo:       order.TradeNo,
			OrderStatus:   string(nextOrderStatus),
			ProcessedAt:   &now,
			CreatedAt:     refund.CreatedAt,
			UpdatedAt:     now,
		}
		return nil
	})
	if err != nil {
		return ordermodule.RefundRequest{}, err
	}

	return result, nil
}

func (r *AdminRepo) UpdateReviewStatus(ctx context.Context, reviewID int64, status wishlistmodule.ReviewStatus) (wishlistmodule.Review, error) {
	var row models.Review
	if err := r.Gorm.WithContext(ctx).Where("id = ?", reviewID).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wishlistmodule.Review{}, usermodule.ErrNotFound
		}
		return wishlistmodule.Review{}, fmt.Errorf("AdminRepo.UpdateReviewStatus(find): %w", err)
	}

	row.Status = string(status)
	row.UpdatedAt = time.Now().UTC()
	if err := r.Gorm.WithContext(ctx).Save(&row).Error; err != nil {
		return wishlistmodule.Review{}, fmt.Errorf("AdminRepo.UpdateReviewStatus(save): %w", err)
	}

	return wishlistmodule.Review{
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
	}, nil
}

func (r *AdminRepo) BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status wishlistmodule.ReviewStatus) (int, error) {
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
		return usermodule.ErrNotFound
	}
	return nil
}

func (r *AdminRepo) UpdateOrderStatus(ctx context.Context, orderID string, status ordermodule.OrderStatus) error {
	return withTransaction(ctx, r.Gorm, func(tx *gorm.DB) error {
		var order models.Order
		if err := forUpdate(tx).Where("order_id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ordermodule.ErrNotFound
			}
			return fmt.Errorf("AdminRepo.UpdateOrderStatus(find): %w", err)
		}

		current := ordermodule.OrderStatus(order.Status)
		now := time.Now().UTC()
		updates := map[string]any{
			"status":     string(status),
			"updated_at": now,
		}

		switch status {
		case ordermodule.OrderStatusPaid:
			if current != ordermodule.OrderStatusPending {
				return ordermodule.ErrInvalidInput
			}
			updates["paid_at"] = now
		case ordermodule.OrderStatusDelivered:
			if current != ordermodule.OrderStatusPaid {
				return ordermodule.ErrInvalidInput
			}
			updates["delivered_at"] = now
		case ordermodule.OrderStatusCancelled:
			if current != ordermodule.OrderStatusPending && current != ordermodule.OrderStatusPaid {
				return ordermodule.ErrInvalidInput
			}
		default:
			return usermodule.ErrInvalidInput
		}

		if err := tx.Model(&models.Order{}).Where("order_id = ?", orderID).Updates(updates).Error; err != nil {
			return fmt.Errorf("AdminRepo.UpdateOrderStatus(update): %w", err)
		}

		if status == ordermodule.OrderStatusCancelled {
			if current == ordermodule.OrderStatusPending {
				if err := tx.Model(&models.Product{}).
					Where("id = ?", order.ProductID).
					UpdateColumn("locked_count", gorm.Expr("locked_count - ?", order.Quantity)).Error; err != nil {
					return fmt.Errorf("AdminRepo.UpdateOrderStatus(release locked count): %w", err)
				}
			} else if current == ordermodule.OrderStatusPaid {
				if err := tx.Model(&models.Product{}).
					Where("id = ?", order.ProductID).
					Updates(map[string]any{
						"stock_count": gorm.Expr("stock_count + ?", order.Quantity),
						"sold_count":  gorm.Expr("sold_count - ?", order.Quantity),
					}).Error; err != nil {
					return fmt.Errorf("AdminRepo.UpdateOrderStatus(restore stock from paid): %w", err)
				}
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
		if err := tx.Where("order_id = ?", orderID).Delete(&models.Order{}).Error; err != nil {
			return fmt.Errorf("AdminRepo.DeleteOrder(delete order): %w", err)
		}
		return nil
	})
}

func (r *AdminRepo) StoreSetting(ctx context.Context, setting catalogmodule.Setting) error {
	model := models.ModuleToSetting(setting)
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

func (r *AdminRepo) ListSettings(ctx context.Context) ([]catalogmodule.Setting, error) {
	var rows []models.Setting
	if err := r.Gorm.WithContext(ctx).Order("\"key\" ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListSettings: %w", err)
	}

	settings := make([]catalogmodule.Setting, 0, len(rows))
	for _, row := range rows {
		settings = append(settings, models.SettingToModule(row))
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
	return nil
}

func (r *AdminRepo) ListProducts(ctx context.Context, page pagination.Pagination) ([]catalogmodule.Product, int, error) {
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
	items := make([]catalogmodule.Product, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.ProductToModule(row))
	}
	return items, int(total), nil
}

func (r *AdminRepo) GetProduct(ctx context.Context, productID string) (catalogmodule.Product, error) {
	var row models.Product
	if err := r.Gorm.WithContext(ctx).
		Where("id = ?", productID).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return catalogmodule.Product{}, catalogmodule.ErrNotFound
		}
		return catalogmodule.Product{}, fmt.Errorf("AdminRepo.GetProduct: %w", err)
	}

	result := models.ProductToModule(row)

	var detailRows []models.ProductDetail
	if err := r.Gorm.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&detailRows).Error; err != nil {
		return catalogmodule.Product{}, fmt.Errorf("AdminRepo.GetProduct(specs): %w", err)
	}
	result.Specs = make([]catalogmodule.ProductSpecItem, 0, len(detailRows))
	for _, detail := range detailRows {
		switch detail.Key {
		case "sku":
			result.SKU = detail.Value
		case "brand":
			result.Brand = detail.Value
		default:
			d := models.DetailToModule(detail)
			result.Specs = append(result.Specs, catalogmodule.ProductSpecItem{Key: d.Key, Value: d.Value})
		}
	}

	return result, nil
}

func (r *AdminRepo) UpsertProduct(ctx context.Context, product catalogmodule.Product) (catalogmodule.Product, error) {
	if _, err := uuid.Parse(strings.TrimSpace(product.ID)); err != nil {
		product.ID = uuid.NewString()
	}

	if strings.TrimSpace(product.SKU) == "" {
		product.SKU = product.ID
	}

	var existing models.Product
	err := r.Gorm.WithContext(ctx).Where("id = ?", product.ID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return catalogmodule.Product{}, fmt.Errorf("AdminRepo.UpsertProduct(find existing): %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		product.IsActive = true
	}

	model := models.ModuleToProduct(product)
	now := time.Now().UTC()
	if model.CreatedAt.IsZero() {
		model.CreatedAt = now
	}
	model.UpdatedAt = now

	if err := r.Gorm.WithContext(ctx).Save(&model).Error; err != nil {
		return catalogmodule.Product{}, fmt.Errorf("AdminRepo.UpsertProduct: %w", err)
	}

	// Persist specs
	if err := r.Gorm.WithContext(ctx).Where("product_id = ?", product.ID).Delete(&models.ProductDetail{}).Error; err != nil {
		return catalogmodule.Product{}, fmt.Errorf("AdminRepo.UpsertProduct(delete specs): %w", err)
	}

	detailRows := buildDetailRows(product)
	if len(detailRows) > 0 {
		if err := r.Gorm.WithContext(ctx).Create(&detailRows).Error; err != nil {
			return catalogmodule.Product{}, fmt.Errorf("AdminRepo.UpsertProduct(create specs): %w", err)
		}
	}

	result := models.ProductToModule(model)
	result.SKU = product.SKU
	result.Brand = product.Brand
	result.Specs = make([]catalogmodule.ProductSpecItem, 0, len(detailRows))
	for _, detail := range detailRows {
		if detail.Key == "sku" || detail.Key == "brand" {
			continue
		}
		d := models.DetailToModule(detail)
		result.Specs = append(result.Specs, catalogmodule.ProductSpecItem{Key: d.Key, Value: d.Value})
	}

	return result, nil
}

func (r *AdminRepo) DeleteProduct(ctx context.Context, productID string) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", productID).Delete(&models.Product{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteProduct: %w", err)
	}
	return nil
}

func (r *AdminRepo) ListCategories(ctx context.Context) ([]catalogmodule.Category, error) {
	var rows []models.Category
	if err := r.Gorm.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListCategories: %w", err)
	}
	items := make([]catalogmodule.Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.CategoryToModule(row))
	}
	return items, nil
}

func (r *AdminRepo) UpsertCategory(ctx context.Context, category catalogmodule.Category) (catalogmodule.Category, error) {
	if category.ID == "" {
		category.ID = uuid.NewString()
	}
	model := models.ModuleToCategory(category)
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
		return catalogmodule.Category{}, fmt.Errorf("AdminRepo.UpsertCategory: %w", err)
	}
	return models.CategoryToModule(model), nil
}

func (r *AdminRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	if err := r.Gorm.WithContext(ctx).Where("id = ?", categoryID).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("AdminRepo.DeleteCategory: %w", err)
	}
	return nil
}

func (r *AdminRepo) StoreAdminMessage(ctx context.Context, msg notificationmodule.AdminMessage) (notificationmodule.AdminMessage, error) {
	model := models.ModuleToAdminMessage(msg)
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now().UTC()
	}
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return notificationmodule.AdminMessage{}, fmt.Errorf("AdminRepo.StoreAdminMessage: %w", err)
	}
	return models.AdminMessageToModule(model), nil
}

func (r *AdminRepo) ListAdminMessages(ctx context.Context) ([]notificationmodule.AdminMessage, error) {
	var rows []models.AdminMessage
	if err := r.Gorm.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("AdminRepo.ListAdminMessages: %w", err)
	}

	msgs := make([]notificationmodule.AdminMessage, 0, len(rows))
	for _, row := range rows {
		msgs = append(msgs, models.AdminMessageToModule(row))
	}
	return msgs, nil
}
