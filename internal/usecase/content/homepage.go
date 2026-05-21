package content

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

type HomepageUseCase struct {
	homepage repo.HomepageRepo
	support  repo.SupportChannelRepo
}

func NewHomepage(homepage repo.HomepageRepo, support repo.SupportChannelRepo) *HomepageUseCase {
	return &HomepageUseCase{homepage: homepage, support: support}
}

func (uc *HomepageUseCase) StoreBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error) {
	block.ID = uuid.New().String()
	if err := uc.homepage.Store(ctx, &block); err != nil {
		return entity.HomepageBlock{}, err
	}
	return block, nil
}

func (uc *HomepageUseCase) ListBlocks(ctx context.Context, activeOnly bool) ([]entity.HomepageBlock, error) {
	return uc.homepage.List(ctx, activeOnly)
}

func (uc *HomepageUseCase) UpdateBlock(ctx context.Context, block entity.HomepageBlock) (entity.HomepageBlock, error) {
	if err := uc.homepage.Update(ctx, &block); err != nil {
		return entity.HomepageBlock{}, err
	}
	return block, nil
}

func (uc *HomepageUseCase) DeleteBlock(ctx context.Context, id string) error {
	return uc.homepage.Delete(ctx, id)
}

func (uc *HomepageUseCase) ListSupport(ctx context.Context, enabledOnly bool) ([]entity.SupportChannel, error) {
	return uc.support.List(ctx, enabledOnly)
}

func (uc *HomepageUseCase) UpdateSupport(ctx context.Context, channel entity.SupportChannel) (entity.SupportChannel, error) {
	if channel.ID == "" {
		channel.ID = uuid.New().String()
	}
	if err := uc.support.Update(ctx, &channel); err != nil {
		return entity.SupportChannel{}, err
	}
	return channel, nil
}
