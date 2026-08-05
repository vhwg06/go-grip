package content

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ContentUseCase defines business operations for Articles and StaticPages.
type ContentUseCase interface {
	CreateArticle(ctx context.Context, article ContentArticle) (ContentArticle, error)
	UpdateArticle(ctx context.Context, article ContentArticle) (ContentArticle, error)
	ListArticles(ctx context.Context, filter ArticleFilter) ([]ContentArticle, int, error)
	GetArticle(ctx context.Context, idOrSlug string) (ContentArticle, error)
	DeleteArticle(ctx context.Context, id string) error
	CreatePage(ctx context.Context, page StaticPage) (StaticPage, error)
	UpdatePage(ctx context.Context, page StaticPage) (StaticPage, error)
	GetPage(ctx context.Context, slug string) (StaticPage, error)
	PublishDue(ctx context.Context) (int, error)
	CatchUpScheduled(ctx context.Context) (int, error)
}

type contentUseCase struct {
	repo ContentRepo
}

// NewContentUseCase constructs a new ContentUseCase instance.
func NewContentUseCase(r ContentRepo) ContentUseCase {
	return &contentUseCase{repo: r}
}

func (uc *contentUseCase) CreateArticle(ctx context.Context, article ContentArticle) (ContentArticle, error) {
	now := time.Now().UTC()
	article.ID = uuid.New().String()
	if strings.TrimSpace(article.Slug) == "" {
		article.Slug = articleSlug(article.Title, article.ID)
	}
	if article.Status == "" {
		article.Status = ContentStatusDraft
	}
	article.CreatedAt = now
	article.UpdatedAt = now
	if err := uc.repo.StoreArticle(ctx, &article); err != nil {
		return ContentArticle{}, err
	}
	return article, nil
}

var articleSlugSeparator = regexp.MustCompile(`[^a-z0-9]+`)

func articleSlug(title, id string) string {
	base := strings.ToLower(strings.TrimSpace(title))
	base = articleSlugSeparator.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "article"
	}
	suffix := strings.ReplaceAll(id, "-", "")
	if len(suffix) > 12 {
		suffix = suffix[:12]
	}
	return base + "-" + suffix
}

func (uc *contentUseCase) UpdateArticle(ctx context.Context, article ContentArticle) (ContentArticle, error) {
	article.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdateArticle(ctx, &article); err != nil {
		return ContentArticle{}, err
	}
	return article, nil
}

func (uc *contentUseCase) ListArticles(ctx context.Context, filter ArticleFilter) ([]ContentArticle, int, error) {
	return uc.repo.ListArticles(ctx, filter)
}

func (uc *contentUseCase) GetArticle(ctx context.Context, idOrSlug string) (ContentArticle, error) {
	return uc.repo.GetArticle(ctx, idOrSlug)
}

func (uc *contentUseCase) DeleteArticle(ctx context.Context, id string) error {
	return uc.repo.DeleteArticle(ctx, id)
}

func (uc *contentUseCase) CreatePage(ctx context.Context, page StaticPage) (StaticPage, error) {
	page.ID = uuid.New().String()
	page.UpdatedAt = time.Now().UTC()
	if err := uc.repo.StorePage(ctx, &page); err != nil {
		return StaticPage{}, err
	}
	return page, nil
}

func (uc *contentUseCase) UpdatePage(ctx context.Context, page StaticPage) (StaticPage, error) {
	page.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdatePage(ctx, &page); err != nil {
		return StaticPage{}, err
	}
	return page, nil
}

func (uc *contentUseCase) GetPage(ctx context.Context, slug string) (StaticPage, error) {
	return uc.repo.GetPageBySlug(ctx, slug)
}

func (uc *contentUseCase) PublishDue(ctx context.Context) (int, error) {
	return uc.repo.PublishDue(ctx)
}

func (uc *contentUseCase) CatchUpScheduled(ctx context.Context) (int, error) {
	return uc.PublishDue(ctx)
}
