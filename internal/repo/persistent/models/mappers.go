package models

import (
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

func OrderToModule(m Order) ordermodule.Order {
	return ordermodule.Order{
		ID:               m.OrderID,
		ProductID:        m.ProductID,
		ProductName:      m.ProductName,
		Amount:           ordermodule.Amount(m.Amount),
		Email:            m.Email,
		Status:           ordermodule.OrderStatus(m.Status),
		TradeNo:          m.TradeNo,
		UserID:           m.UserID,
		Username:         m.Username,
		Payee:            m.Payee,
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

func ModuleToOrder(e ordermodule.Order) Order {
	return Order{
		OrderID:          e.ID,
		ProductID:        e.ProductID,
		ProductName:      e.ProductName,
		Amount:           int64(e.Amount),
		Email:            e.Email,
		Status:           string(e.Status),
		TradeNo:          e.TradeNo,
		UserID:           e.UserID,
		Username:         e.Username,
		Payee:            e.Payee,
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

func ProductToModule(m Product) catalogmodule.Product {
	var images []string
	var comparePrice *int64
	var introArticleID string
	if m.Image != "" {
		images = []string{m.Image}
	} else {
		images = []string{}
	}
	if m.CompareAtPrice > 0 {
		value := m.CompareAtPrice
		comparePrice = &value
	}
	if m.IntroArticleID != nil {
		introArticleID = *m.IntroArticleID
	}

	return catalogmodule.Product{
		ID:              m.ID,
		Title:           m.Name,
		SKU:             m.SKU,
		Description:     m.Description,
		Price:           m.Price,
		ComparePrice:    comparePrice,
		CategoryID:      m.Category,
		ImageURL:        m.Image,
		Images:          images,
		IsHot:           m.IsHot,
		IsActive:        m.IsActive,
		SortOrder:       m.SortOrder,
		PurchaseLimit:   m.PurchaseLimit,
		PurchaseWarning: m.PurchaseWarning,
		IntroArticleID:  introArticleID,
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

func ModuleToProduct(e catalogmodule.Product) Product {
	var compareAtPrice int64
	if e.ComparePrice != nil {
		compareAtPrice = *e.ComparePrice
	}
	var introArticleID *string
	if e.IntroArticleID != "" {
		introArticleID = &e.IntroArticleID
	}
	return Product{
		ID:              e.ID,
		Name:            e.Title,
		SKU:             e.SKU,
		Description:     e.Description,
		Price:           e.Price,
		CompareAtPrice:  compareAtPrice,
		Category:        e.CategoryID,
		Image:           e.ImageURL,
		IsHot:           e.IsHot,
		IsActive:        e.IsActive,
		SortOrder:       e.SortOrder,
		PurchaseLimit:   e.PurchaseLimit,
		PurchaseWarning: e.PurchaseWarning,
		IntroArticleID:  introArticleID,
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

func CategoryToModule(m Category) catalogmodule.Category {
	return catalogmodule.Category{
		ID:       m.ID,
		Name:     m.Name,
		ParentID: m.ParentID,
		Position: m.SortOrder,
		IsActive: m.IsActive,
	}
}

func SettingToModule(m Setting) catalogmodule.Setting {
	return catalogmodule.Setting{
		Key:       m.Key,
		Value:     m.Value,
		UpdatedAt: m.UpdatedAt,
	}
}

func UserToEntity(m User) entity.User {
	return entity.User{
		ID:                          m.ID,
		Username:                    m.Username,
		DisplayName:                 m.DisplayName,
		Email:                       m.Email,
		PasswordHash:                m.PasswordHash,
		RoleID:                      m.RoleID,
		Role:                        entity.RoleName(m.Role),
		Status:                      entity.UserStatus(m.Status),
		Provider:                    m.Provider,
		ProviderID:                  m.ProviderID,
		TrustLevel:                  m.TrustLevel,
		IsAdmin:                     m.IsAdmin,
		DesktopNotificationsEnabled: m.DesktopNotificationsEnabled,
		LastLoginAt:                 nullableTime(m.LastLoginAt),
		IsBlocked:                   m.Status == string(entity.UserStatusLocked),
		CreatedAt:                   m.CreatedAt,
		UpdatedAt:                   m.UpdatedAt,
	}
}

func UserToModule(m User) usermodule.User {
	return usermodule.User{
		ID:                          m.ID,
		Username:                    m.Username,
		DisplayName:                 m.DisplayName,
		Email:                       m.Email,
		PasswordHash:                m.PasswordHash,
		RoleID:                      m.RoleID,
		Role:                        usermodule.RoleName(m.Role),
		Status:                      usermodule.UserStatus(m.Status),
		Provider:                    m.Provider,
		ProviderID:                  m.ProviderID,
		TrustLevel:                  m.TrustLevel,
		IsAdmin:                     m.IsAdmin,
		DesktopNotificationsEnabled: m.DesktopNotificationsEnabled,
		LastLoginAt:                 nullableTime(m.LastLoginAt),
		IsBlocked:                   m.Status == string(usermodule.UserStatusLocked),
		CreatedAt:                   m.CreatedAt,
		UpdatedAt:                   m.UpdatedAt,
	}
}

func EntityToUser(e entity.User) User {
	return User{
		ID:                          e.ID,
		Username:                    e.Username,
		DisplayName:                 e.DisplayName,
		Email:                       e.Email,
		PasswordHash:                e.PasswordHash,
		RoleID:                      e.RoleID,
		Role:                        string(e.Role),
		Status:                      string(e.Status),
		Provider:                    e.Provider,
		ProviderID:                  e.ProviderID,
		TrustLevel:                  e.TrustLevel,
		IsAdmin:                     e.IsAdmin,
		DesktopNotificationsEnabled: e.DesktopNotificationsEnabled,
		LastLoginAt:                 denullTime(e.LastLoginAt),
		CreatedAt:                   e.CreatedAt,
		UpdatedAt:                   e.UpdatedAt,
	}
}

func ModuleToUser(e usermodule.User) User {
	return User{
		ID:                          e.ID,
		Username:                    e.Username,
		DisplayName:                 e.DisplayName,
		Email:                       e.Email,
		PasswordHash:                e.PasswordHash,
		RoleID:                      e.RoleID,
		Role:                        string(e.Role),
		Status:                      string(e.Status),
		Provider:                    e.Provider,
		ProviderID:                  e.ProviderID,
		TrustLevel:                  e.TrustLevel,
		IsAdmin:                     e.IsAdmin,
		DesktopNotificationsEnabled: e.DesktopNotificationsEnabled,
		LastLoginAt:                 denullTime(e.LastLoginAt),
		CreatedAt:                   e.CreatedAt,
		UpdatedAt:                   e.UpdatedAt,
	}
}

func ProductToEntity(m Product) entity.Product {
	var images []string
	var comparePrice *int64
	var introArticleID string
	if m.Image != "" {
		images = []string{m.Image}
	} else {
		images = []string{}
	}
	if m.CompareAtPrice > 0 {
		value := m.CompareAtPrice
		comparePrice = &value
	}
	if m.IntroArticleID != nil {
		introArticleID = *m.IntroArticleID
	}

	return entity.Product{
		ID:              m.ID,
		Title:           m.Name,
		SKU:             m.SKU,
		Description:     m.Description,
		Price:           m.Price,
		ComparePrice:    comparePrice,
		CategoryID:      m.Category,
		ImageURL:        m.Image,
		Images:          images,
		IsHot:           m.IsHot,
		IsActive:        m.IsActive,
		SortOrder:       m.SortOrder,
		PurchaseLimit:   m.PurchaseLimit,
		PurchaseWarning: m.PurchaseWarning,
		IntroArticleID:  introArticleID,
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
	var compareAtPrice int64
	var introArticleID *string
	if e.ComparePrice != nil {
		compareAtPrice = *e.ComparePrice
	}
	if e.IntroArticleID != "" {
		introArticleID = &e.IntroArticleID
	}

	return Product{
		ID:              e.ID,
		Name:            e.Title,
		SKU:             e.SKU,
		Description:     e.Description,
		Price:           e.Price,
		CompareAtPrice:  compareAtPrice,
		Category:        e.CategoryID,
		Image:           e.ImageURL,
		IsHot:           e.IsHot,
		IsActive:        e.IsActive,
		SortOrder:       e.SortOrder,
		PurchaseLimit:   e.PurchaseLimit,
		PurchaseWarning: e.PurchaseWarning,
		IntroArticleID:  introArticleID,
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

func DetailToEntity(m ProductDetail) entity.ProductSpecItem {
	return entity.ProductSpecItem{
		Key:   m.Key,
		Value: m.Value,
	}
}

func EntityToDetail(productID string, e entity.ProductSpecItem) ProductDetail {
	return ProductDetail{
		ProductID: productID,
		Key:       e.Key,
		Value:     e.Value,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func CategoryToEntity(m Category) entity.Category {
	return entity.Category{
		ID:       m.ID,
		Name:     m.Name,
		ParentID: m.ParentID,
		Position: m.SortOrder,
		IsActive: m.IsActive,
	}
}

func EntityToCategory(e entity.Category) Category {
	return Category{
		ID:        e.ID,
		Name:      e.Name,
		ParentID:  e.ParentID,
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
		UserID:           m.UserID,
		Username:         m.Username,
		Payee:            m.Payee,
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
		UserID:           e.UserID,
		Username:         e.Username,
		Payee:            e.Payee,
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

func MediaAssetToModule(m MediaAsset) mediamodule.MediaAsset {
	return mediamodule.MediaAsset{
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

func ModuleToMediaAsset(e mediamodule.MediaAsset) MediaAsset {
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

func AdminMessageToEntity(m AdminMessage) entity.AdminMessage {
	return entity.AdminMessage{
		ID:          m.ID,
		TargetType:  m.TargetType,
		TargetValue: m.TargetValue,
		Title:       m.Title,
		Body:        m.Body,
		Sender:      m.Sender,
		CreatedAt:   m.CreatedAt,
	}
}

func EntityToAdminMessage(e entity.AdminMessage) AdminMessage {
	return AdminMessage{
		ID:          e.ID,
		TargetType:  e.TargetType,
		TargetValue: e.TargetValue,
		Title:       e.Title,
		Body:        e.Body,
		Sender:      e.Sender,
		CreatedAt:   e.CreatedAt,
	}
}
