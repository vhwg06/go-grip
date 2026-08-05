package content

import (
	"context"

	"github.com/google/uuid"
)

// HomepageUseCase defines management operations for HomepageBlocks and SupportChannels.
type HomepageUseCase interface {
	StoreBlock(ctx context.Context, block HomepageBlock) (HomepageBlock, error)
	ListBlocks(ctx context.Context, activeOnly bool) ([]HomepageBlock, error)
	UpdateBlock(ctx context.Context, block HomepageBlock) (HomepageBlock, error)
	DeleteBlock(ctx context.Context, id string) error
	ListSupport(ctx context.Context, enabledOnly bool) ([]SupportChannel, error)
	UpdateSupport(ctx context.Context, channel SupportChannel) (SupportChannel, error)
}

type homepageUseCase struct {
	homepage HomepageRepo
	support  SupportChannelRepo
}

// NewHomepageUseCase constructs a new HomepageUseCase instance.
func NewHomepageUseCase(homepage HomepageRepo, support SupportChannelRepo) HomepageUseCase {
	return &homepageUseCase{homepage: homepage, support: support}
}

func (uc *homepageUseCase) StoreBlock(ctx context.Context, block HomepageBlock) (HomepageBlock, error) {
	block.ID = uuid.New().String()
	if err := uc.homepage.Store(ctx, &block); err != nil {
		return HomepageBlock{}, err
	}
	return block, nil
}

func (uc *homepageUseCase) ListBlocks(ctx context.Context, activeOnly bool) ([]HomepageBlock, error) {
	return uc.homepage.List(ctx, activeOnly)
}

func (uc *homepageUseCase) UpdateBlock(ctx context.Context, block HomepageBlock) (HomepageBlock, error) {
	if err := uc.homepage.Update(ctx, &block); err != nil {
		return HomepageBlock{}, err
	}
	return block, nil
}

func (uc *homepageUseCase) DeleteBlock(ctx context.Context, id string) error {
	return uc.homepage.Delete(ctx, id)
}

func (uc *homepageUseCase) ListSupport(ctx context.Context, enabledOnly bool) ([]SupportChannel, error) {
	return uc.support.List(ctx, enabledOnly)
}

func (uc *homepageUseCase) UpdateSupport(ctx context.Context, channel SupportChannel) (SupportChannel, error) {
	if channel.ID == "" {
		channel.ID = uuid.New().String()
	}
	if err := uc.support.Update(ctx, &channel); err != nil {
		return SupportChannel{}, err
	}
	return channel, nil
}
