package persistent

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestAdminRepo(t *testing.T) *AdminRepo {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Product{}, &models.ProductDetail{}))

	return NewAdminRepo(&postgres.Postgres{Gorm: db})
}

func TestAdminRepoUpsertProductPersistsInactiveState(t *testing.T) {
	t.Parallel()

	repo := newTestAdminRepo(t)
	ctx := context.Background()
	now := time.Now().UTC()

	created, err := repo.UpsertProduct(ctx, entity.Product{
		ID:         "p-inactive",
		Title:      "Initial",
		SKU:        "sku-inactive",
		Price:      100,
		CategoryID: "test-category",
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	require.NoError(t, err)
	require.True(t, created.IsActive)

	updated, err := repo.UpsertProduct(ctx, entity.Product{
		ID:         created.ID,
		Title:      created.Title,
		SKU:        created.SKU,
		Price:      created.Price,
		CategoryID: created.CategoryID,
		IsActive:   false,
		CreatedAt:  created.CreatedAt,
		UpdatedAt:  time.Now().UTC(),
	})
	require.NoError(t, err)
	require.False(t, updated.IsActive)

	stored, err := repo.GetProduct(ctx, created.ID)
	require.NoError(t, err)
	require.False(t, stored.IsActive)
}

func TestAdminRepoUpsertProductPersistsIntroArticleLink(t *testing.T) {
	t.Parallel()

	repo := newTestAdminRepo(t)
	ctx := context.Background()
	now := time.Now().UTC()

	created, err := repo.UpsertProduct(ctx, entity.Product{
		ID:             "p-intro",
		Title:          "Editorial Product",
		SKU:            "sku-intro",
		Price:          100,
		CategoryID:     "test-category",
		IntroArticleID: "article-1",
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	require.NoError(t, err)
	require.Equal(t, "article-1", created.IntroArticleID)

	stored, err := repo.GetProduct(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, "article-1", stored.IntroArticleID)

	cleared, err := repo.UpsertProduct(ctx, entity.Product{
		ID:             created.ID,
		Title:          created.Title,
		SKU:            created.SKU,
		Price:          created.Price,
		CategoryID:     created.CategoryID,
		IntroArticleID: "",
		IsActive:       true,
		CreatedAt:      created.CreatedAt,
		UpdatedAt:      time.Now().UTC(),
	})
	require.NoError(t, err)
	require.Empty(t, cleared.IntroArticleID)

	stored, err = repo.GetProduct(ctx, created.ID)
	require.NoError(t, err)
	require.Empty(t, stored.IntroArticleID)
}
