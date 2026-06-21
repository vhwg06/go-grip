package models

import "time"

type User struct {
	ID                          string `gorm:"type:text;primaryKey"`
	Provider                    string `gorm:"type:text;index"`
	ProviderID                  string `gorm:"type:text;index"`
	Username                    string `gorm:"type:text;index"`
	DisplayName                 string `gorm:"type:text"`
	Email                       string `gorm:"type:text;index"`
	PasswordHash                string `gorm:"type:text"`
	RoleID                      string `gorm:"type:text"`
	Role                        string `gorm:"type:text"`
	Status                      string `gorm:"type:text"`
	Points                      int    `gorm:"not null;default:0"`
	TrustLevel                  int    `gorm:"not null;default:0"`
	IsAdmin                     bool   `gorm:"not null;default:false"`
	DesktopNotificationsEnabled bool   `gorm:"not null;default:false"`
	LastLoginAt                 time.Time
	LastCheckinAt               time.Time
	ConsecutiveDays             int       `gorm:"not null;default:0"`
	CreatedAt                   time.Time `gorm:"not null"`
	UpdatedAt                   time.Time `gorm:"not null"`
}

func (User) TableName() string { return "login_users" }

type RefreshSession struct {
	ID        string    `gorm:"type:text;primaryKey"`
	UserID    string    `gorm:"type:text;index"`
	TokenID   string    `gorm:"type:text;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt time.Time
	CreatedAt time.Time `gorm:"not null"`
}

type Product struct {
	ID              string    `gorm:"type:text;primaryKey"`
	Name            string    `gorm:"column:title;type:text;index"`
	SKU             string    `gorm:"column:sku;type:text;uniqueIndex"`
	Description     string    `gorm:"type:text"`
	Price           int64     `gorm:"not null"`
	CompareAtPrice  int64     `gorm:"column:compare_price"`
	Category        string    `gorm:"type:text;index"`
	Image           string    `gorm:"type:text"`
	IsHot           bool      `gorm:"not null;default:false"`
	IsActive        bool      `gorm:"not null;default:true"`
	SortOrder       int       `gorm:"not null;default:0"`
	PurchaseLimit   int       `gorm:"not null;default:0"`
	PurchaseWarning string    `gorm:"type:text"`
	VisibilityLevel int       `gorm:"not null;default:-1"`
	StockCount      int       `gorm:"not null;default:0"`
	LockedCount     int       `gorm:"not null;default:0"`
	SoldCount       int       `gorm:"not null;default:0"`
	Rating          float64   `gorm:"not null;default:0"`
	ReviewCount     int       `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (Product) TableName() string { return "products" }

type ProductDetail struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	ProductID string    `gorm:"type:text;index;not null"`
	Key       string    `gorm:"type:text;not null"`
	Value     string    `gorm:"type:text;not null"`
	SortOrder int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (ProductDetail) TableName() string { return "product_details" }

type Category struct {
	ID        string    `gorm:"type:text;primaryKey"`
	Name      string    `gorm:"type:text"`
	ParentID  *string   `gorm:"type:text"`
	SortOrder int       `gorm:"not null;default:0"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Category) TableName() string { return "categories" }

type Order struct {
	OrderID           string `gorm:"type:text;primaryKey"`
	ProductID         string `gorm:"type:text;index"`
	ProductName       string `gorm:"type:text"`
	Amount            int64  `gorm:"not null"`
	Email             string `gorm:"type:text;index"`
	Status            string `gorm:"type:text;index"`
	TradeNo           string `gorm:"type:text"`
	PaidAt            time.Time
	DeliveredAt       time.Time
	UserID            string    `gorm:"type:text;index"`
	Username          string    `gorm:"type:text"`
	Payee             string    `gorm:"type:text"`
	PointsUsed        int       `gorm:"not null;default:0"`
	Quantity          int       `gorm:"not null;default:1"`
	CurrentPaymentID  string    `gorm:"type:text"`
	StatusText        string    `gorm:"type:text"`
	StatusColor       string    `gorm:"type:text"`
	PaymentProviderID string    `gorm:"type:text;index"`
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}

func (Order) TableName() string { return "orders" }

type Payment struct {
	ID                     string `gorm:"type:text;primaryKey"`
	OrderID                string `gorm:"type:text;index"`
	Provider               string `gorm:"type:text"`
	ProviderPaymentID      string `gorm:"type:text;index"`
	Amount                 int64  `gorm:"not null"`
	Status                 string `gorm:"type:text;index"`
	RequestPayloadSummary  string `gorm:"type:text"`
	CallbackPayloadSummary string `gorm:"type:text"`
	IsSignatureValid       bool   `gorm:"not null;default:false"`
	ProcessedAt            time.Time
	CreatedAt              time.Time `gorm:"not null"`
}

func (Payment) TableName() string { return "payments" }

type RefundRequest struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	OrderID       string `gorm:"type:text;index"`
	UserID        string `gorm:"type:text;index"`
	Username      string `gorm:"type:text"`
	Reason        string `gorm:"type:text"`
	Status        string `gorm:"type:text;index"`
	AdminUsername string `gorm:"type:text"`
	AdminNote     string `gorm:"type:text"`
	ProcessedAt   time.Time
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (RefundRequest) TableName() string { return "refund_requests" }

type DailyCheckin struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	UserID      string    `gorm:"type:text;index"`
	CheckinDate time.Time `gorm:"not null"`
	Reward      int       `gorm:"not null;default:0"`
	StreakAfter int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (DailyCheckin) TableName() string { return "daily_checkins_v2" }

type Review struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	ProductID string    `gorm:"type:text;index"`
	OrderID   string    `gorm:"type:text;index"`
	UserID    string    `gorm:"type:text;index"`
	Username  string    `gorm:"type:text"`
	Rating    int       `gorm:"not null"`
	Comment   string    `gorm:"type:text"`
	Status    string    `gorm:"type:text;not null;default:'PENDING';index"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Review) TableName() string { return "reviews" }

type WishlistItem struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	UserID      string    `gorm:"type:text;index"`
	Username    string    `gorm:"type:text"`
	VoteCount   int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (WishlistItem) TableName() string { return "wishlist_items" }

type WishlistVote struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	ItemID    int64     `gorm:"index"`
	UserID    string    `gorm:"type:text;index"`
	CreatedAt time.Time `gorm:"not null"`
}

func (WishlistVote) TableName() string { return "wishlist_votes" }

type UserNotification struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     string    `gorm:"type:text;index"`
	Type       string    `gorm:"type:text"`
	TitleKey   string    `gorm:"type:text"`
	ContentKey string    `gorm:"type:text"`
	Data       string    `gorm:"type:text"`
	IsRead     bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (UserNotification) TableName() string { return "user_notifications" }

type BroadcastMessage struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	TitleKey   string    `gorm:"type:text"`
	ContentKey string    `gorm:"type:text"`
	Data       string    `gorm:"type:text"`
	Sender     string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (BroadcastMessage) TableName() string { return "broadcast_messages" }

type BroadcastRead struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	MessageID int64     `gorm:"index"`
	UserID    string    `gorm:"type:text;index"`
	CreatedAt time.Time `gorm:"not null"`
}

func (BroadcastRead) TableName() string { return "broadcast_reads" }

type Setting struct {
	Key       string    `gorm:"type:text;primaryKey"`
	Value     string    `gorm:"type:text"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Setting) TableName() string { return "settings" }

type MediaAsset struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	FileName  string    `gorm:"column:file_name;not null"`
	MimeType  string    `gorm:"column:mime_type;not null"`
	SizeBytes int64     `gorm:"column:size_bytes;not null"`
	URL       string    `gorm:"column:url;not null"`
	AltText   *string   `gorm:"column:alt_text"`
	OwnerType *string   `gorm:"column:owner_type"`
	OwnerID   *string   `gorm:"column:owner_id;type:uuid"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (MediaAsset) TableName() string { return "media_assets" }

type Card struct {
	ID              int64      `gorm:"primaryKey;autoIncrement"`
	ProductID       string     `gorm:"type:text;not null"`
	CardKey         string     `gorm:"type:text;not null"`
	IsUsed          bool       `gorm:"not null;default:false"`
	ReservedOrderID string     `gorm:"type:text;not null;default:''"`
	ReservedAt      *time.Time `gorm:"type:timestamptz"`
	ExpiresAt       *time.Time `gorm:"type:timestamptz"`
	UsedAt          *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (Card) TableName() string { return "cards" }

type AdminMessage struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	TargetType  string    `gorm:"not null"`
	TargetValue string    `gorm:"column:target_value"`
	Title       string    `gorm:"not null"`
	Body        string    `gorm:"not null"`
	Sender      string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
}

func (AdminMessage) TableName() string { return "admin_messages" }

