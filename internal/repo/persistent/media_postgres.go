package persistent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo/persistent/models"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

type MediaRepo struct {
	*postgres.Postgres
	mu    sync.RWMutex
	items map[string]entity.MediaAsset
}

func NewMediaRepo(pg *postgres.Postgres) *MediaRepo {
	return &MediaRepo{Postgres: pg, items: make(map[string]entity.MediaAsset)}
}

func (r *MediaRepo) Store(ctx context.Context, media *entity.MediaAsset) error {
	if r.Postgres == nil || r.Gorm == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.items == nil {
			r.items = make(map[string]entity.MediaAsset)
		}
		r.items[media.ID] = *media
		return nil
	}

	model := models.EntityToMediaAsset(*media)
	if err := r.Gorm.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("MediaRepo.Store: %w", err)
	}
	return nil
}

func (r *MediaRepo) Get(ctx context.Context, id string) (entity.MediaAsset, error) {
	if r.Postgres == nil || r.Gorm == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		item, exists := r.items[id]
		if !exists {
			return entity.MediaAsset{}, fmt.Errorf("MediaRepo.Get: not found")
		}
		return item, nil
	}

	var row models.MediaAsset
	if err := r.Gorm.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return entity.MediaAsset{}, fmt.Errorf("MediaRepo.Get: %w", err)
	}
	return models.MediaAssetToEntity(row), nil
}

func (r *MediaRepo) List(ctx context.Context, page entity.Pagination, q string) ([]entity.MediaAsset, int, error) {
	if r.Postgres == nil || r.Gorm == nil {
		_ = page.Normalize()
		r.mu.RLock()
		defer r.mu.RUnlock()
		items := make([]entity.MediaAsset, 0, len(r.items))
		for _, item := range r.items {
			if q == "" || strings.Contains(strings.ToLower(item.FileName), strings.ToLower(q)) {
				items = append(items, item)
			}
		}
		return items, len(items), nil
	}

	db := r.Gorm.WithContext(ctx).Model(&models.MediaAsset{})
	if q != "" {
		db = db.Where("LOWER(file_name) LIKE ?", "%"+strings.ToLower(q)+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("MediaRepo.List(count): %w", err)
	}

	normalized := page.Normalize()
	var rows []models.MediaAsset
	if err := db.Order("created_at DESC").Limit(normalized.Limit).Offset(normalized.Offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("MediaRepo.List(find): %w", err)
	}

	items := make([]entity.MediaAsset, 0, len(rows))
	for _, row := range rows {
		items = append(items, models.MediaAssetToEntity(row))
	}
	return items, int(total), nil
}

func (r *MediaRepo) Delete(ctx context.Context, id string) error {
	if r.Postgres == nil || r.Gorm == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.items, id)
		return nil
	}

	if err := r.Gorm.WithContext(ctx).Where("id = ?", id).Delete(&models.MediaAsset{}).Error; err != nil {
		return fmt.Errorf("MediaRepo.Delete: %w", err)
	}
	return nil
}
