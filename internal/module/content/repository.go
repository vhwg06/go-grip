package content

import (
	"context"
)

// ContentRepo defines persistence operations for Articles and StaticPages.
type ContentRepo interface {
	StoreArticle(ctx context.Context, article *ContentArticle) error
	UpdateArticle(ctx context.Context, article *ContentArticle) error
	ListArticles(ctx context.Context, filter ArticleFilter) ([]ContentArticle, int, error)
	GetArticle(ctx context.Context, idOrSlug string) (ContentArticle, error)
	DeleteArticle(ctx context.Context, id string) error
	StorePage(ctx context.Context, page *StaticPage) error
	UpdatePage(ctx context.Context, page *StaticPage) error
	GetPageBySlug(ctx context.Context, slug string) (StaticPage, error)
	PublishDue(ctx context.Context) (int, error)
}

// HomepageRepo defines persistence operations for HomepageBlocks.
type HomepageRepo interface {
	Store(ctx context.Context, block *HomepageBlock) error
	List(ctx context.Context, activeOnly bool) ([]HomepageBlock, error)
	Update(ctx context.Context, block *HomepageBlock) error
	Delete(ctx context.Context, id string) error
}

// SupportChannelRepo defines persistence operations for SupportChannels.
type SupportChannelRepo interface {
	List(ctx context.Context, enabledOnly bool) ([]SupportChannel, error)
	Update(ctx context.Context, channel *SupportChannel) error
}
