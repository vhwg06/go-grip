package models

import (
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
)

func UserToEntity(m User) entity.User {
	return entity.User{
		ID:                          m.ID,
		Username:                    m.Username,
		Email:                       m.Email,
		PasswordHash:                m.PasswordHash,
		RoleID:                      m.RoleID,
		Role:                        entity.RoleName(m.Role),
		Status:                      entity.UserStatus(m.Status),
		Provider:                    m.Provider,
		ProviderID:                  m.ProviderID,
		Points:                      m.Points,
		TrustLevel:                  m.TrustLevel,
		IsAdmin:                     m.IsAdmin,
		DesktopNotificationsEnabled: m.DesktopNotificationsEnabled,
		LastLoginAt:                 nullableTime(m.LastLoginAt),
		LastCheckinAt:               nullableTime(m.LastCheckinAt),
		ConsecutiveDays:             m.ConsecutiveDays,
		CreatedAt:                   m.CreatedAt,
		UpdatedAt:                   m.UpdatedAt,
	}
}

func EntityToUser(e entity.User) User {
	return User{
		ID:                          e.ID,
		Username:                    e.Username,
		Email:                       e.Email,
		PasswordHash:                e.PasswordHash,
		RoleID:                      e.RoleID,
		Role:                        string(e.Role),
		Status:                      string(e.Status),
		Provider:                    e.Provider,
		ProviderID:                  e.ProviderID,
		Points:                      e.Points,
		TrustLevel:                  e.TrustLevel,
		IsAdmin:                     e.IsAdmin,
		DesktopNotificationsEnabled: e.DesktopNotificationsEnabled,
		LastLoginAt:                 denullTime(e.LastLoginAt),
		LastCheckinAt:               denullTime(e.LastCheckinAt),
		ConsecutiveDays:             e.ConsecutiveDays,
		CreatedAt:                   e.CreatedAt,
		UpdatedAt:                   e.UpdatedAt,
	}
}

func ProductToEntity(m Product) entity.Product {
	var images []string
	if m.Image != "" {
		images = []string{m.Image}
	} else {
		images = []string{}
	}

	return entity.Product{
		ID:              m.ID,
		Title:           m.Name,
		Description:     m.Description,
		Price:           m.Price,
		CategoryID:      m.Category,
		ImageURL:        m.Image,
		Images:          images,
		IsHot:           m.IsHot,
		IsActive:        m.IsActive,
		IsShared:        m.IsShared,
		SortOrder:       m.SortOrder,
		PurchaseLimit:   m.PurchaseLimit,
		PurchaseWarning: m.PurchaseWarning,
		VisibilityLevel: m.VisibilityLevel,
		StockCount:      m.StockCount,
		LockedCount:     m.LockedCount,
		SoldCount:       m.SoldCount,
		Rating:          m.Rating,
		ReviewCount:     m.ReviewCount,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func EntityToProduct(e entity.Product) Product {
	return Product{
		ID:              e.ID,
		Name:            e.Title,
		Description:     e.Description,
		Price:           e.Price,
		Category:        e.CategoryID,
		Image:           e.ImageURL,
		IsHot:           e.IsHot,
		IsActive:        e.IsActive,
		IsShared:        e.IsShared,
		SortOrder:       e.SortOrder,
		PurchaseLimit:   e.PurchaseLimit,
		PurchaseWarning: e.PurchaseWarning,
		VisibilityLevel: e.VisibilityLevel,
		StockCount:      e.StockCount,
		LockedCount:     e.LockedCount,
		SoldCount:       e.SoldCount,
		Rating:          e.Rating,
		ReviewCount:     e.ReviewCount,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func CategoryToEntity(m Category) entity.Category {
	var parentID *string
	if m.ParentID != "" {
		parent := m.ParentID
		parentID = &parent
	}

	return entity.Category{
		ID:       m.ID,
		Name:     m.Name,
		ParentID: parentID,
		Position: m.SortOrder,
		IsActive: m.IsActive,
	}
}

func EntityToCategory(e entity.Category) Category {
	var parentID string
	if e.ParentID != nil {
		parentID = *e.ParentID
	}

	return Category{
		ID:        e.ID,
		Name:      e.Name,
		ParentID:  parentID,
		SortOrder: e.Position,
		IsActive:  e.IsActive,
	}
}

func OrderToEntity(m Order) entity.Order {
	return entity.Order{
		ID:               m.OrderID,
		ProductID:        m.ProductID,
		ProductName:      m.ProductName,
		Amount:           entity.Amount(m.Amount),
		Email:            m.Email,
		Status:           entity.OrderStatus(m.Status),
		TradeNo:          m.TradeNo,
		CardKey:          m.CardKey,
		UserID:           m.UserID,
		Username:         m.Username,
		Payee:            m.Payee,
		PointsUsed:       m.PointsUsed,
		Quantity:         m.Quantity,
		CurrentPaymentID: m.CurrentPaymentID,
		StatusText:       m.StatusText,
		StatusColor:      m.StatusColor,
		PaidAt:           nullableTime(m.PaidAt),
		DeliveredAt:      nullableTime(m.DeliveredAt),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func EntityToOrder(e entity.Order) Order {
	return Order{
		OrderID:          e.ID,
		ProductID:        e.ProductID,
		ProductName:      e.ProductName,
		Amount:           int64(e.Amount),
		Email:            e.Email,
		Status:           string(e.Status),
		TradeNo:          e.TradeNo,
		CardKey:          e.CardKey,
		UserID:           e.UserID,
		Username:         e.Username,
		Payee:            e.Payee,
		PointsUsed:       e.PointsUsed,
		Quantity:         e.Quantity,
		CurrentPaymentID: e.CurrentPaymentID,
		StatusText:       e.StatusText,
		StatusColor:      e.StatusColor,
		PaidAt:           denullTime(e.PaidAt),
		DeliveredAt:      denullTime(e.DeliveredAt),
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func SettingToEntity(m Setting) entity.Setting {
	return entity.Setting{
		Key:       m.Key,
		Value:     m.Value,
		UpdatedAt: m.UpdatedAt,
	}
}

func EntityToSetting(e entity.Setting) Setting {
	return Setting{
		Key:       e.Key,
		Value:     e.Value,
		UpdatedAt: e.UpdatedAt,
	}
}

func nullableTime(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}

	copied := v

	return &copied
}

func denullTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}

	return *v
}

func MediaAssetToEntity(m MediaAsset) entity.MediaAsset {
	return entity.MediaAsset{
		ID:        m.ID,
		FileName:  m.FileName,
		MimeType:  m.MimeType,
		SizeBytes: m.SizeBytes,
		URL:       m.URL,
		AltText:   stringPointerValue(m.AltText),
		OwnerType: stringPointerValue(m.OwnerType),
		OwnerID:   stringPointerValue(m.OwnerID),
		CreatedAt: m.CreatedAt,
	}
}

func EntityToMediaAsset(e entity.MediaAsset) MediaAsset {
	return MediaAsset{
		ID:        e.ID,
		FileName:  e.FileName,
		MimeType:  e.MimeType,
		SizeBytes: e.SizeBytes,
		URL:       e.URL,
		AltText:   valueStringPointer(e.AltText),
		OwnerType: valueStringPointer(e.OwnerType),
		OwnerID:   valueStringPointer(e.OwnerID),
		CreatedAt: e.CreatedAt,
	}
}

func stringPointerValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func valueStringPointer(v string) *string {
	if v == "" {
		return nil
	}
	copied := v
	return &copied
}
