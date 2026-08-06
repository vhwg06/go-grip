// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"time"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	pagination "github.com/evrone/go-clean-template/internal/shared/pagination"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepo interface {
		Store(ctx context.Context, user *usermodule.User) error
		GetByID(ctx context.Context, id string) (usermodule.User, error)
		GetByEmail(ctx context.Context, email string) (usermodule.User, error)
		List(ctx context.Context, filter UserFilter) ([]usermodule.User, int, error)
		Update(ctx context.Context, user *usermodule.User) error
		SetStatus(ctx context.Context, id string, status usermodule.UserStatus) error
	}

	// RoleRepo -.
	RoleRepo interface {
		List(ctx context.Context) ([]usermodule.Role, error)
		GetByName(ctx context.Context, name usermodule.RoleName) (usermodule.Role, error)
	}
	UserFilter struct {
		Limit  uint64
		Offset uint64
	}

	CatalogRepo interface {
		StoreProduct(ctx context.Context, product *catalogmodule.Product) error
		GetProduct(ctx context.Context, id string) (catalogmodule.Product, error)
		ListProducts(ctx context.Context, filter ProductFilter) ([]catalogmodule.Product, int, error)
		UpdateProduct(ctx context.Context, product *catalogmodule.Product) error
		DeleteProduct(ctx context.Context, id string) error
		StoreCategory(ctx context.Context, category *catalogmodule.Category) error
		ListCategories(ctx context.Context) ([]catalogmodule.Category, error)
		StoreTag(ctx context.Context, tag *catalogmodule.Tag) error
		ListTags(ctx context.Context) ([]catalogmodule.Tag, error)
	}

	ProductFilter struct {
		Keyword  string
		Category string
		Brand    string
		MinPrice *int64
		MaxPrice *int64
		Sort     string
		Limit    uint64
		Offset   uint64
	}

	SEORepo interface {
		Store(ctx context.Context, meta *catalogmodule.SeoMetadata) error
		GetByOwner(ctx context.Context, ownerType, ownerID string) (catalogmodule.SeoMetadata, error)
	}

	HomepageRepo interface {
		Store(ctx context.Context, block *contentmodule.HomepageBlock) error
		List(ctx context.Context, activeOnly bool) ([]contentmodule.HomepageBlock, error)
		Update(ctx context.Context, block *contentmodule.HomepageBlock) error
		Delete(ctx context.Context, id string) error
	}

	SupportChannelRepo interface {
		List(ctx context.Context, enabledOnly bool) ([]contentmodule.SupportChannel, error)
		Update(ctx context.Context, channel *contentmodule.SupportChannel) error
	}

	CartRepo interface {
		Store(ctx context.Context, cart *cartmodule.Cart) error
		GetBySession(ctx context.Context, sessionID string) (cartmodule.Cart, error)
		AddItem(ctx context.Context, sessionID string, item *cartmodule.CartItem) error
		UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) error
		RemoveItem(ctx context.Context, sessionID, itemID string) error
		Convert(ctx context.Context, cartID string) error
	}

	OrderRequestRepo interface {
		Store(ctx context.Context, order *cartmodule.OrderRequest) error
	}
	ContentRepo interface {
		StoreArticle(ctx context.Context, article *contentmodule.ContentArticle) error
		UpdateArticle(ctx context.Context, article *contentmodule.ContentArticle) error
		ListArticles(ctx context.Context, filter contentmodule.ArticleFilter) ([]contentmodule.ContentArticle, int, error)
		GetArticle(ctx context.Context, idOrSlug string) (contentmodule.ContentArticle, error)
		DeleteArticle(ctx context.Context, id string) error
		StorePage(ctx context.Context, page *contentmodule.StaticPage) error
		UpdatePage(ctx context.Context, page *contentmodule.StaticPage) error
		GetPageBySlug(ctx context.Context, slug string) (contentmodule.StaticPage, error)
		PublishDue(ctx context.Context) (int, error)
	}
	AuthRepository interface {
		GetUserByID(ctx context.Context, userID string) (usermodule.User, error)
		GetUserByEmail(ctx context.Context, email string) (usermodule.User, error)
		GetUserByUsername(ctx context.Context, username string) (usermodule.User, error)
		UpsertUser(ctx context.Context, user usermodule.User) (usermodule.User, error)
		StoreRefreshSession(ctx context.Context, session usermodule.RefreshSession) error
		RevokeRefreshSession(ctx context.Context, tokenID string) error
		GetRefreshSession(ctx context.Context, tokenID string) (usermodule.RefreshSession, error)
	}

	CatalogRepository interface {
		ListVisibleProducts(ctx context.Context, actor usermodule.Actor, filter ProductFilter) ([]catalogmodule.Product, int, error)
		GetVisibleProduct(ctx context.Context, actor usermodule.Actor, productID string) (catalogmodule.Product, error)
		ListCategories(ctx context.Context) ([]catalogmodule.Category, error)
		ListSettings(ctx context.Context) ([]catalogmodule.Setting, error)
		GetSetting(ctx context.Context, key string) (catalogmodule.Setting, error)
	}

	// UnitOfWork executes an application operation with one transaction-bound
	// context without exposing infrastructure transaction handles.
	UnitOfWork interface {
		Within(ctx context.Context, fn func(context.Context) error) error
	}

	CheckoutRepository interface {
		CreateOrderWithReservation(ctx context.Context, actor usermodule.Actor, order ordermodule.Order) (ordermodule.Order, error)
		AttachPayment(ctx context.Context, payment ordermodule.Payment) error
		UpdateOrderStatus(ctx context.Context, orderID string, status ordermodule.OrderStatus) error
		ReleaseReservation(ctx context.Context, orderID string) error
	}

	OrderRepository interface {
		ListOrdersByOwner(ctx context.Context, userID, email string, page pagination.Pagination) ([]ordermodule.Order, int, error)
		GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error)
		CancelPendingOrder(ctx context.Context, actor usermodule.Actor, orderID string) error
		SubmitRefundRequest(ctx context.Context, refund *ordermodule.RefundRequest) error
	}

	ProfileRepository interface {
		GetProfile(ctx context.Context, userID string) (usermodule.User, error)
		UpdateProfile(ctx context.Context, user usermodule.User) (usermodule.User, error)
	}

	WishlistRepository interface {
		ListWishlistItems(ctx context.Context, page pagination.Pagination) ([]wishlistmodule.WishlistItem, int, error)
		StoreWishlistItem(ctx context.Context, item wishlistmodule.WishlistItem) (wishlistmodule.WishlistItem, error)
		UpdateWishlistItem(ctx context.Context, item wishlistmodule.WishlistItem) (wishlistmodule.WishlistItem, error)
		DeleteWishlistItem(ctx context.Context, itemID int64) error
		ToggleWishlistVote(ctx context.Context, itemID int64, userID string) (bool, error)
		StoreReview(ctx context.Context, review wishlistmodule.Review) (wishlistmodule.Review, error)
		ListReviews(ctx context.Context, productID string) ([]wishlistmodule.Review, error)
	}

	NotificationRepository interface {
		ListUserNotifications(ctx context.Context, userID string, page pagination.Pagination) ([]notificationmodule.UserNotification, int, error)
		ListBroadcastMessages(ctx context.Context, page pagination.Pagination) ([]notificationmodule.BroadcastMessage, int, error)
		MarkNotificationRead(ctx context.Context, userID string, notificationID int64) error
		MarkAllRead(ctx context.Context, userID string) error
		ClearAll(ctx context.Context, userID string) error
	}

	AdminRepository interface {
		ListUsers(ctx context.Context, page pagination.Pagination) ([]usermodule.User, int, error)
		UpdateUserStatus(ctx context.Context, userID string, status usermodule.UserStatus) error
		ListOrders(ctx context.Context, page pagination.Pagination, query, status string) ([]ordermodule.Order, int, error)
		GetOrderByID(ctx context.Context, orderID string) (ordermodule.Order, error)
		ListRefundRequests(ctx context.Context, status string) ([]ordermodule.RefundRequest, error)
		GetRefundRequest(ctx context.Context, refundID int64) (ordermodule.RefundRequest, error)
		GetOrderRefundStatus(ctx context.Context, orderID string) (ordermodule.RefundRequest, error)
		ProcessRefund(ctx context.Context, refundID int64, approve bool, adminUsername, note string) (ordermodule.RefundRequest, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status ordermodule.OrderStatus) error
		DeleteOrder(ctx context.Context, orderID string) error
		ListSettings(ctx context.Context) ([]catalogmodule.Setting, error)
		StoreSetting(ctx context.Context, setting catalogmodule.Setting) error
		DeleteSetting(ctx context.Context, key string) error
		ListReviews(ctx context.Context, page pagination.Pagination, query, status string) ([]wishlistmodule.Review, wishlistmodule.ReviewModerationStats, int, error)
		UpdateReviewStatus(ctx context.Context, reviewID int64, status wishlistmodule.ReviewStatus) (wishlistmodule.Review, error)
		BulkUpdateReviewStatus(ctx context.Context, reviewIDs []int64, status wishlistmodule.ReviewStatus) (int, error)
		DeleteReview(ctx context.Context, reviewID int64) error
		StoreAdminMessage(ctx context.Context, msg notificationmodule.AdminMessage) (notificationmodule.AdminMessage, error)
		ListAdminMessages(ctx context.Context) ([]notificationmodule.AdminMessage, error)
		RebuildProductAggregates(ctx context.Context) error
	}

	MaintenanceRepository interface {
		CancelExpiredPendingOrders(ctx context.Context, olderThan time.Time) error
		SyncProductAggregates(ctx context.Context) error
	}

	ReviewModerationStats = wishlistmodule.ReviewModerationStats
)
