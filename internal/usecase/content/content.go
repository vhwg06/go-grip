package content

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repo.ContentRepo
}

func New(r repo.ContentRepo) *UseCase { return &UseCase{repo: r} }

func (uc *UseCase) CreateArticle(ctx context.Context, article entity.ContentArticle) (entity.ContentArticle, error) {
	now := time.Now().UTC()
	article.ID = uuid.New().String()
	if article.Status == "" {
		article.Status = entity.ContentStatusDraft
	}
	article.CreatedAt = now
	article.UpdatedAt = now
	if err := uc.repo.StoreArticle(ctx, &article); err != nil {
		return entity.ContentArticle{}, err
	}
	return article, nil
}

func (uc *UseCase) UpdateArticle(ctx context.Context, article entity.ContentArticle) (entity.ContentArticle, error) {
	article.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdateArticle(ctx, &article); err != nil {
		return entity.ContentArticle{}, err
	}
	return article, nil
}

func (uc *UseCase) ListArticles(ctx context.Context, filter entity.ArticleFilter) ([]entity.ContentArticle, int, error) {
	return uc.repo.ListArticles(ctx, filter)
}

func (uc *UseCase) GetArticle(ctx context.Context, idOrSlug string) (entity.ContentArticle, error) {
	return uc.repo.GetArticle(ctx, idOrSlug)
}

func (uc *UseCase) DeleteArticle(ctx context.Context, id string) error {
	return uc.repo.DeleteArticle(ctx, id)
}

func (uc *UseCase) CreatePage(ctx context.Context, page entity.StaticPage) (entity.StaticPage, error) {
	page.ID = uuid.New().String()
	page.UpdatedAt = time.Now().UTC()
	if err := uc.repo.StorePage(ctx, &page); err != nil {
		return entity.StaticPage{}, err
	}
	return page, nil
}

func (uc *UseCase) UpdatePage(ctx context.Context, page entity.StaticPage) (entity.StaticPage, error) {
	page.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdatePage(ctx, &page); err != nil {
		return entity.StaticPage{}, err
	}
	return page, nil
}

func (uc *UseCase) GetPage(ctx context.Context, slug string) (entity.StaticPage, error) {
	return uc.repo.GetPageBySlug(ctx, slug)
}

func (uc *UseCase) PublishDue(ctx context.Context) (int, error) {
	return uc.repo.PublishDue(ctx)
}
