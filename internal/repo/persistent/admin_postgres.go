package persistent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
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
		WHERE p.id = src.product_id
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
