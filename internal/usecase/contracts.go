// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	pagination "github.com/evrone/go-clean-template/internal/shared/pagination"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// User -.
	User interface {
		Register(ctx context.Context, username, email, password string) (usermodule.User, error)
		Login(ctx context.Context, email, password string) (string, error)
		GetUser(ctx context.Context, userID string) (usermodule.User, error)
		List(ctx context.Context, limit, offset int) ([]usermodule.User, int, error)
		CreateAdminUser(ctx context.Context, actorID, username, email, password string, role usermodule.RoleName) (usermodule.User, error)
		UpdateProfile(ctx context.Context, userID, displayName string) (usermodule.User, error)
		Lock(ctx context.Context, actorID, userID string) error
		Unlock(ctx context.Context, actorID, userID string) error
	}
	Catalog interface {
		CreateProduct(ctx context.Context, product catalogmodule.Product) (catalogmodule.Product, error)
		ListProducts(ctx context.Context, filter catalogmodule.ProductFilter) ([]catalogmodule.Product, int, error)
		GetProduct(ctx context.Context, id string) (catalogmodule.Product, error)
		UpdateProduct(ctx context.Context, product catalogmodule.Product) (catalogmodule.Product, error)
		DeleteProduct(ctx context.Context, id string) error
		CreateCategory(ctx context.Context, category catalogmodule.Category) (catalogmodule.Category, error)
		ListCategories(ctx context.Context) ([]catalogmodule.Category, error)
		CreateTag(ctx context.Context, tag catalogmodule.Tag) (catalogmodule.Tag, error)
		ListTags(ctx context.Context) ([]catalogmodule.Tag, error)
	}

	Homepage interface {
		StoreBlock(ctx context.Context, block contentmodule.HomepageBlock) (contentmodule.HomepageBlock, error)
		ListBlocks(ctx context.Context, activeOnly bool) ([]contentmodule.HomepageBlock, error)
		UpdateBlock(ctx context.Context, block contentmodule.HomepageBlock) (contentmodule.HomepageBlock, error)
		DeleteBlock(ctx context.Context, id string) error
		ListSupport(ctx context.Context, enabledOnly bool) ([]contentmodule.SupportChannel, error)
		UpdateSupport(ctx context.Context, channel contentmodule.SupportChannel) (contentmodule.SupportChannel, error)
	}

	Cart interface {
		Create(ctx context.Context, sessionID string) (cartmodule.Cart, error)
		Get(ctx context.Context, sessionID string) (cartmodule.Cart, error)
		AddItem(ctx context.Context, sessionID, productID string, quantity int) (cartmodule.Cart, error)
		UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) (cartmodule.Cart, error)
		RemoveItem(ctx context.Context, sessionID, itemID string) (cartmodule.Cart, error)
		SubmitOrder(ctx context.Context, order cartmodule.OrderRequest) (cartmodule.OrderRequest, error)
	}

	Content interface {
		CreateArticle(ctx context.Context, article contentmodule.ContentArticle) (contentmodule.ContentArticle, error)
		UpdateArticle(ctx context.Context, article contentmodule.ContentArticle) (contentmodule.ContentArticle, error)
		ListArticles(ctx context.Context, filter contentmodule.ArticleFilter) ([]contentmodule.ContentArticle, int, error)
		GetArticle(ctx context.Context, idOrSlug string) (contentmodule.ContentArticle, error)
		DeleteArticle(ctx context.Context, id string) error
		CreatePage(ctx context.Context, page contentmodule.StaticPage) (contentmodule.StaticPage, error)
		UpdatePage(ctx context.Context, page contentmodule.StaticPage) (contentmodule.StaticPage, error)
		GetPage(ctx context.Context, slug string) (contentmodule.StaticPage, error)
		PublishDue(ctx context.Context) (int, error)
	}

	Notification interface {
		Dispatch(ctx context.Context, notification notificationmodule.Notification) error
	}

	Auth interface {
		Login(ctx context.Context, email, password string) (usermodule.User, string, string, error)
		Refresh(ctx context.Context, refreshToken string) (string, string, error)
		Logout(ctx context.Context, actor usermodule.Actor, refreshToken string) error
		Me(ctx context.Context, actor usermodule.Actor) (usermodule.User, error)
	}

	Checkout interface {
		Preview(ctx context.Context, actor usermodule.Actor, productID string, quantity int) (AmountBreakdown, error)
		CreateOrder(ctx context.Context, actor usermodule.Actor, productID string, quantity int, email string) (ordermodule.Order, error)
		PaymentParams(ctx context.Context, actor usermodule.Actor, orderID string) (PaymentParams, error)
		PaymentNotify(ctx context.Context, payload map[string]string) error
		PaymentStatus(ctx context.Context, orderID string) (ordermodule.Order, error)
		Cancel(ctx context.Context, actor usermodule.Actor, orderID string) error
	}

	Orders interface {
		List(ctx context.Context, actor usermodule.Actor, email string, page pagination.Pagination) ([]ordermodule.Order, int, error)
		Get(ctx context.Context, actor usermodule.Actor, orderID string) (ordermodule.Order, error)
		RequestRefund(ctx context.Context, actor usermodule.Actor, orderID, reason string) (ordermodule.RefundRequest, error)
	}

	Profile interface {
		Get(ctx context.Context, actor usermodule.Actor) (usermodule.User, error)
		Update(ctx context.Context, actor usermodule.Actor, email string, displayName string, desktopNotificationsEnabled bool) (usermodule.User, error)
	}

	Wishlist interface {
		List(ctx context.Context, page pagination.Pagination) ([]wishlistmodule.WishlistItem, int, error)
		Create(ctx context.Context, actor usermodule.Actor, title, description string) (wishlistmodule.WishlistItem, error)
		Update(ctx context.Context, actor usermodule.Actor, item wishlistmodule.WishlistItem) (wishlistmodule.WishlistItem, error)
		Delete(ctx context.Context, actor usermodule.Actor, itemID int64) error
		ToggleVote(ctx context.Context, actor usermodule.Actor, itemID int64) error
		CreateReview(ctx context.Context, actor usermodule.Actor, review wishlistmodule.Review) (wishlistmodule.Review, error)
		ListReviews(ctx context.Context, productID string) ([]wishlistmodule.Review, error)
	}

	NotificationCenter interface {
		Inbox(ctx context.Context, actor usermodule.Actor, page pagination.Pagination) ([]notificationmodule.UserNotification, int, error)
		UnreadCount(ctx context.Context, actor usermodule.Actor) (int, error)
		MarkRead(ctx context.Context, actor usermodule.Actor, notificationID int64) error
		MarkAllRead(ctx context.Context, actor usermodule.Actor) error
		Clear(ctx context.Context, actor usermodule.Actor) error
	}

	Admin interface {
		ListUsers(ctx context.Context, actor usermodule.Actor, page pagination.Pagination) ([]usermodule.User, int, error)
		UpdateUserStatus(ctx context.Context, actor usermodule.Actor, userID string, status usermodule.UserStatus) error
		ListOrders(ctx context.Context, actor usermodule.Actor, page pagination.Pagination, query, status string) ([]ordermodule.Order, int, error)
		GetOrder(ctx context.Context, actor usermodule.Actor, orderID string) (ordermodule.Order, error)
		RepairAggregates(ctx context.Context, actor usermodule.Actor) error
	}

	Maintenance interface {
		CancelExpiredPendingOrders(ctx context.Context) error
		SyncProductAggregates(ctx context.Context) error
	}
)

type AmountBreakdown struct {
	Subtotal   ordermodule.Amount `json:"subtotal"`
	FinalPrice ordermodule.Amount `json:"final_price"`
}

type PaymentParams struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
