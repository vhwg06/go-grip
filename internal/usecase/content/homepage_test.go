package content

import (
	"context"
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/stretchr/testify/require"
)

func TestHomepageUseCase(t *testing.T) {
	t.Parallel()
	uc := NewHomepage(persistent.NewHomepageRepo(nil), persistent.NewSupportChannelRepo(nil))
	block, err := uc.StoreBlock(context.Background(), entity.HomepageBlock{BlockType: "banner", IsActive: true})
	require.NoError(t, err)
	require.NotEmpty(t, block.ID)
	items, err := uc.ListBlocks(context.Background(), true)
	require.NoError(t, err)
	require.Len(t, items, 1)
}
