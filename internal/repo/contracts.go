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
)
