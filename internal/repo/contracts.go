// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(ctx context.Context, userID string, t entity.Translation) error
		GetHistory(ctx context.Context, userID string) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(ctx context.Context, t entity.Translation) (entity.Translation, error)
	}

	// UserRepo -.
	UserRepo interface {
		Store(ctx context.Context, user *entity.User) error
		GetByID(ctx context.Context, id string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		List(ctx context.Context, filter UserFilter) ([]entity.User, int, error)
		Update(ctx context.Context, user *entity.User) error
		SetStatus(ctx context.Context, id string, status entity.UserStatus) error
	}

	// RoleRepo -.
	RoleRepo interface {
		List(ctx context.Context) ([]entity.Role, error)
		GetByName(ctx context.Context, name entity.RoleName) (entity.Role, error)
	}

	// TaskRepo -.
	TaskRepo interface {
		Store(ctx context.Context, task *entity.Task) error
		GetByID(ctx context.Context, userID, taskID string) (entity.Task, error)
		List(ctx context.Context, userID string, filter TaskFilter) ([]entity.Task, int, error)
		Update(ctx context.Context, task *entity.Task) error
		Delete(ctx context.Context, userID, taskID string) error
	}

	// TaskFilter -.
	TaskFilter struct {
		Status *entity.TaskStatus
		Limit  uint64
		Offset uint64
	}

	UserFilter struct {
		Limit  uint64
		Offset uint64
	}

	CatalogRepo interface {
		StoreProduct(ctx context.Context, product *entity.Product) error
		GetProduct(ctx context.Context, id string) (entity.Product, error)
		ListProducts(ctx context.Context, filter ProductFilter) ([]entity.Product, int, error)
		UpdateProduct(ctx context.Context, product *entity.Product) error
		DeleteProduct(ctx context.Context, id string) error
		StoreCategory(ctx context.Context, category *entity.Category) error
		ListCategories(ctx context.Context) ([]entity.Category, error)
		StoreTag(ctx context.Context, tag *entity.Tag) error
		ListTags(ctx context.Context) ([]entity.Tag, error)
	}

	ProductFilter struct {
		Keyword  string
		Brand    string
		MinPrice *int64
		MaxPrice *int64
		Sort     string
		Limit    uint64
		Offset   uint64
	}

	MediaRepo interface {
		Store(ctx context.Context, media *entity.MediaAsset) error
		List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error)
		Delete(ctx context.Context, id string) error
	}

	SEORepo interface {
		Store(ctx context.Context, meta *entity.SeoMetadata) error
		GetByOwner(ctx context.Context, ownerType, ownerID string) (entity.SeoMetadata, error)
	}

	HomepageRepo interface {
		Store(ctx context.Context, block *entity.HomepageBlock) error
		List(ctx context.Context, activeOnly bool) ([]entity.HomepageBlock, error)
		Update(ctx context.Context, block *entity.HomepageBlock) error
		Delete(ctx context.Context, id string) error
	}

	SupportChannelRepo interface {
		List(ctx context.Context, enabledOnly bool) ([]entity.SupportChannel, error)
		Update(ctx context.Context, channel *entity.SupportChannel) error
	}

	CartRepo interface {
		Store(ctx context.Context, cart *entity.Cart) error
		GetBySession(ctx context.Context, sessionID string) (entity.Cart, error)
		AddItem(ctx context.Context, sessionID string, item *entity.CartItem) error
		UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) error
		RemoveItem(ctx context.Context, sessionID, itemID string) error
		Convert(ctx context.Context, cartID string) error
	}

	OrderRequestRepo interface {
		Store(ctx context.Context, order *entity.OrderRequest) error
	}

	LeadRepo interface {
		Store(ctx context.Context, lead *entity.LeadSubmission) error
		Get(ctx context.Context, id string) (entity.LeadSubmission, error)
	}

	ContentRepo interface {
		StoreArticle(ctx context.Context, article *entity.ContentArticle) error
		UpdateArticle(ctx context.Context, article *entity.ContentArticle) error
		ListArticles(ctx context.Context, publicOnly bool, page entity.Pagination) ([]entity.ContentArticle, int, error)
		GetArticle(ctx context.Context, idOrSlug string) (entity.ContentArticle, error)
		StorePage(ctx context.Context, page *entity.StaticPage) error
		UpdatePage(ctx context.Context, page *entity.StaticPage) error
		GetPageBySlug(ctx context.Context, slug string) (entity.StaticPage, error)
		PublishDue(ctx context.Context) (int, error)
	}

	ImportRepo interface {
		StoreImportedProduct(ctx context.Context, product *entity.Product) error
		StoreImportedPost(ctx context.Context, article *entity.ContentArticle) error
	}

	AuthRepository interface {
		GetUserByID(ctx context.Context, userID string) (entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (entity.User, error)
		UpsertUser(ctx context.Context, user entity.User) (entity.User, error)
		StoreRefreshSession(ctx context.Context, session entity.RefreshSession) error
		RevokeRefreshSession(ctx context.Context, tokenID string) error
		GetRefreshSession(ctx context.Context, tokenID string) (entity.RefreshSession, error)
	}

	CatalogRepository interface {
		ListVisibleProducts(ctx context.Context, actor entity.Actor, filter ProductFilter) ([]entity.Product, int, error)
		GetVisibleProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error)
		ListCategories(ctx context.Context) ([]entity.Category, error)
		ListSettings(ctx context.Context) ([]entity.Setting, error)
		GetSetting(ctx context.Context, key string) (entity.Setting, error)
	}

	CheckoutRepository interface {
		CreateOrderWithReservation(ctx context.Context, actor entity.Actor, order entity.Order) (entity.Order, error)
		AttachPayment(ctx context.Context, payment entity.Payment) error
		UpdateOrderStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
		ReserveCards(ctx context.Context, orderID, productID string, quantity int, isShared bool) ([]entity.Card, error)
		DeductPoints(ctx context.Context, userID string, points int) error
		ReleaseReservation(ctx context.Context, orderID string) error
	}

	OrderRepository interface {
		ListOrdersByOwner(ctx context.Context, userID, email string, page entity.Pagination) ([]entity.Order, int, error)
		GetOrderByID(ctx context.Context, orderID string) (entity.Order, error)
		CancelPendingOrder(ctx context.Context, actor entity.Actor, orderID string) error
		SubmitRefundRequest(ctx context.Context, refund entity.RefundRequest) error
	}

	ProfileRepository interface {
		GetProfile(ctx context.Context, userID string) (entity.User, error)
		UpdateProfile(ctx context.Context, user entity.User) (entity.User, error)
		RecordDailyCheckin(ctx context.Context, checkin entity.DailyCheckin) error
	}

	WishlistRepository interface {
		ListWishlistItems(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error)
		StoreWishlistItem(ctx context.Context, item entity.WishlistItem) (entity.WishlistItem, error)
		UpdateWishlistItem(ctx context.Context, item entity.WishlistItem) (entity.WishlistItem, error)
		DeleteWishlistItem(ctx context.Context, itemID int64) error
		ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error)
		StoreReview(ctx context.Context, review entity.Review) (entity.Review, error)
	}

	NotificationRepository interface {
		ListUserNotifications(ctx context.Context, userID string, page entity.Pagination) ([]entity.UserNotification, int, error)
		ListBroadcastMessages(ctx context.Context, page entity.Pagination) ([]entity.BroadcastMessage, int, error)
		MarkNotificationRead(ctx context.Context, userID string, notificationID int64) error
		MarkAllRead(ctx context.Context, userID string) error
		ClearAll(ctx context.Context, userID string) error
	}

	AdminRepository interface {
		ListUsers(ctx context.Context, page entity.Pagination) ([]entity.User, int, error)
		UpdateUserStatus(ctx context.Context, userID string, status entity.UserStatus) error
		UpdateUserPoints(ctx context.Context, userID string, points int) error
		StoreSetting(ctx context.Context, setting entity.Setting) error
		RebuildProductAggregates(ctx context.Context) error
	}
)
