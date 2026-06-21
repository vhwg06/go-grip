package persistent

import (
	"context"
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

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
