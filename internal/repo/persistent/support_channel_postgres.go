package persistent

import (
	"context"
	"sync"

	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type SupportChannelRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]contentmodule.SupportChannel
}

func NewSupportChannelRepo(pg *postgres.Postgres) *SupportChannelRepo {
	return &SupportChannelRepo{Postgres: pg, items: map[string]contentmodule.SupportChannel{}}
}

var _ contentmodule.SupportChannelRepo = (*SupportChannelRepo)(nil)

func (r *SupportChannelRepo) List(ctx context.Context, enabledOnly bool) ([]contentmodule.SupportChannel, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]contentmodule.SupportChannel, 0, len(r.items))
	for _, item := range r.items {
		if enabledOnly && !item.IsEnabled {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *SupportChannelRepo) Update(ctx context.Context, channel *contentmodule.SupportChannel) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[channel.ID] = *channel
	return nil
}
