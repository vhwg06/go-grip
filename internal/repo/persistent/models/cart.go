package models

import "time"

type Cart struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	SessionID string    `gorm:"type:text;uniqueIndex"`
	Status    string    `gorm:"type:text;not null;default:'active'"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Cart) TableName() string { return "carts" }

type CartItem struct {
	ID              string    `gorm:"type:uuid;primaryKey"`
	CartID          string    `gorm:"type:uuid;index"`
	ProductID       string    `gorm:"type:uuid;index"`
	Quantity        int       `gorm:"not null"`
	UnitPrice       int64     `gorm:"not null;default:0"`
	ProductSnapshot string    `gorm:"type:jsonb;not null;default:'{}'"`
	Blocked         bool      `gorm:"not null;default:false"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (CartItem) TableName() string { return "cart_items" }
