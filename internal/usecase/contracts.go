// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(ctx context.Context, userID string, t entity.Translation) (entity.Translation, error)
		History(ctx context.Context, userID string) (entity.TranslationHistory, error)
	}

	// User -.
	User interface {
		Register(ctx context.Context, username, email, password string) (entity.User, error)
		Login(ctx context.Context, email, password string) (string, error)
		GetUser(ctx context.Context, userID string) (entity.User, error)
		List(ctx context.Context, limit, offset int) ([]entity.User, int, error)
		CreateAdminUser(ctx context.Context, actorID, username, email, password string, role entity.RoleName) (entity.User, error)
		UpdateProfile(ctx context.Context, userID, displayName string) (entity.User, error)
		Lock(ctx context.Context, actorID, userID string) error
		Unlock(ctx context.Context, actorID, userID string) error
	}

	// Task -.
	Task interface {
		Create(ctx context.Context, userID, title, description string) (entity.Task, error)
		Get(ctx context.Context, userID, taskID string) (entity.Task, error)
		List(ctx context.Context, userID string, status *entity.TaskStatus, limit, offset int) ([]entity.Task, int, error)
		Update(ctx context.Context, userID, taskID, title, description string) (entity.Task, error)
		Transition(ctx context.Context, userID, taskID string, newStatus entity.TaskStatus) (entity.Task, error)
		Delete(ctx context.Context, userID, taskID string) error
	}

	Catalog interface {
		CreateProduct(ctx context.Context, product entity.Product) (entity.Product, error)
		ListProducts(ctx context.Context, filter entity.ProductFilter) ([]entity.Product, int, error)
		GetProduct(ctx context.Context, id string) (entity.Product, error)
		UpdateProduct(ctx context.Context, product entity.Product) (entity.Product, error)
		DeleteProduct(ctx context.Context, id string) error
		CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
		ListCategories(ctx context.Context) ([]entity.Category, error)
		CreateTag(ctx context.Context, tag entity.Tag) (entity.Tag, error)
		ListTags(ctx context.Context) ([]entity.Tag, error)
	}

	Media interface {
		Store(ctx context.Context, media entity.MediaAsset) (entity.MediaAsset, error)
		List(ctx context.Context, page entity.Pagination) ([]entity.MediaAsset, int, error)
		Delete(ctx context.Context, id string) error
		GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)
	}

	MediaStorage interface {
		GeneratePresignedURL(ctx context.Context, fileName string, contentType string) (uploadURL string, publicURL string, fileID string, err error)
		Delete(ctx context.Context, key string) error
	}

	Homepage interface {
		StoreBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error)
		ListBlocks(ctx context.Context, activeOnly bool) ([]entity.HomepageBlock, error)
		UpdateBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error)
		DeleteBlock(ctx context.Context, id string) error
		ListSupport(ctx context.Context, enabledOnly bool) ([]entity.SupportChannel, error)
		UpdateSupport(ctx context.Context, channel entity.SupportChannel) (entity.SupportChannel, error)
	}

	Cart interface {
		Create(ctx context.Context, sessionID string) (entity.Cart, error)
		Get(ctx context.Context, sessionID string) (entity.Cart, error)
		AddItem(ctx context.Context, sessionID, productID string, quantity int) (entity.Cart, error)
		UpdateItem(ctx context.Context, sessionID, itemID string, quantity int) (entity.Cart, error)
		RemoveItem(ctx context.Context, sessionID, itemID string) (entity.Cart, error)
		SubmitOrder(ctx context.Context, order entity.OrderRequest) (entity.OrderRequest, error)
	}

	Lead interface {
		Submit(ctx context.Context, lead entity.LeadSubmission) (entity.LeadSubmission, error)
		Get(ctx context.Context, id string) (entity.LeadSubmission, error)
	}

	Content interface {
		CreateArticle(ctx context.Context, article entity.ContentArticle) (entity.ContentArticle, error)
		UpdateArticle(ctx context.Context, article entity.ContentArticle) (entity.ContentArticle, error)
		ListArticles(ctx context.Context, filter entity.ArticleFilter) ([]entity.ContentArticle, int, error)
		GetArticle(ctx context.Context, idOrSlug string) (entity.ContentArticle, error)
		DeleteArticle(ctx context.Context, id string) error
		CreatePage(ctx context.Context, page entity.StaticPage) (entity.StaticPage, error)
		UpdatePage(ctx context.Context, page entity.StaticPage) (entity.StaticPage, error)
		GetPage(ctx context.Context, slug string) (entity.StaticPage, error)
		PublishDue(ctx context.Context) (int, error)
	}

	Notification interface {
		Dispatch(ctx context.Context, notification entity.Notification) error
	}

	Importer interface {
		Import(ctx context.Context, items []entity.ImportItem) (entity.ImportResult, error)
	}

	Auth interface {
		Login(ctx context.Context, email, password string) (entity.User, string, string, error)
		Refresh(ctx context.Context, refreshToken string) (string, string, error)
		Logout(ctx context.Context, actor entity.Actor, refreshToken string) error
		Me(ctx context.Context, actor entity.Actor) (entity.User, error)
	}

	Checkout interface {
		Preview(ctx context.Context, actor entity.Actor, productID string, quantity int, usePoints bool) (AmountBreakdown, error)
		CreateOrder(ctx context.Context, actor entity.Actor, productID string, quantity int, email string, usePoints bool) (entity.Order, error)
		PaymentParams(ctx context.Context, actor entity.Actor, orderID string) (PaymentParams, error)
		PaymentNotify(ctx context.Context, payload map[string]string) error
		PaymentStatus(ctx context.Context, orderID string) (entity.Order, error)
		Cancel(ctx context.Context, actor entity.Actor, orderID string) error
	}

	Orders interface {
		List(ctx context.Context, actor entity.Actor, email string, page entity.Pagination) ([]entity.Order, int, error)
		Get(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error)
		RequestRefund(ctx context.Context, actor entity.Actor, orderID, reason string) (entity.RefundRequest, error)
	}

	Profile interface {
		Get(ctx context.Context, actor entity.Actor) (entity.User, error)
		Update(ctx context.Context, actor entity.Actor, email string, desktopNotificationsEnabled bool) (entity.User, error)
		Checkin(ctx context.Context, actor entity.Actor) (entity.DailyCheckin, error)
	}

	Wishlist interface {
		List(ctx context.Context, page entity.Pagination) ([]entity.WishlistItem, int, error)
		Create(ctx context.Context, actor entity.Actor, title, description string) (entity.WishlistItem, error)
		Update(ctx context.Context, actor entity.Actor, item entity.WishlistItem) (entity.WishlistItem, error)
		Delete(ctx context.Context, actor entity.Actor, itemID int64) error
		ToggleVote(ctx context.Context, actor entity.Actor, itemID int64) error
		CreateReview(ctx context.Context, actor entity.Actor, review entity.Review) (entity.Review, error)
		ListReviews(ctx context.Context, productID string) ([]entity.Review, error)
	}

	NotificationCenter interface {
		Inbox(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.UserNotification, int, error)
		UnreadCount(ctx context.Context, actor entity.Actor) (int, error)
		MarkRead(ctx context.Context, actor entity.Actor, notificationID int64) error
		MarkAllRead(ctx context.Context, actor entity.Actor) error
		Clear(ctx context.Context, actor entity.Actor) error
	}

	Admin interface {
		ListUsers(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.User, int, error)
		UpdateUserStatus(ctx context.Context, actor entity.Actor, userID string, status entity.UserStatus) error
		UpdateUserPoints(ctx context.Context, actor entity.Actor, userID string, points int) error
		ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error)
		GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error)
		RepairAggregates(ctx context.Context, actor entity.Actor) error
	}

	Maintenance interface {
		CancelExpiredPendingOrders(ctx context.Context) error
		CleanupExpiredCards(ctx context.Context) error
		SyncProductAggregates(ctx context.Context) error
	}
)

type AmountBreakdown struct {
	Subtotal    entity.Amount `json:"subtotal"`
	PointsToUse int           `json:"points_to_use"`
	FinalPrice  entity.Amount `json:"final_price"`
}

type PaymentParams struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
