package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/internal/repo/webapi"
	"github.com/evrone/go-clean-template/internal/usecase"
)

type adminStore interface {
	repo.AdminRepository
	ListProducts(ctx context.Context, page entity.Pagination) ([]entity.Product, int, error)
	GetProduct(ctx context.Context, productID string) (entity.Product, error)
	UpsertProduct(ctx context.Context, product entity.Product) (entity.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	ListCategories(ctx context.Context) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	ImportCards(ctx context.Context, productID string, keys []string) (int, error)
}

type UseCase struct {
	repo       adminStore
	notifier   webapi.AdminNotifier
	adminUsers map[string]struct{}
}

func New(adminRepo adminStore, notifier webapi.AdminNotifier, adminUsersCSV string) *UseCase {
	adminUsers := make(map[string]struct{})
	for raw := range strings.SplitSeq(adminUsersCSV, ",") {
		trimmed := strings.ToLower(strings.TrimSpace(raw))
		if trimmed == "" {
			continue
		}
		adminUsers[trimmed] = struct{}{}
	}
	return &UseCase{repo: adminRepo, notifier: notifier, adminUsers: adminUsers}
}

var _ usecase.Admin = (*UseCase)(nil)

func (uc *UseCase) ListUsers(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.User, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListUsers(ctx, page)
}

func (uc *UseCase) UpdateUserStatus(ctx context.Context, actor entity.Actor, userID string, status entity.UserStatus) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateUserStatus(ctx, userID, status)
}

func (uc *UseCase) UpdateUserPoints(ctx context.Context, actor entity.Actor, userID string, points int) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.UpdateUserPoints(ctx, userID, points)
}

func (uc *UseCase) ListOrders(ctx context.Context, actor entity.Actor, page entity.Pagination, query, status string) ([]entity.Order, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListOrders(ctx, page, query, status)
}

func (uc *UseCase) GetOrder(ctx context.Context, actor entity.Actor, orderID string) (entity.Order, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Order{}, err
	}
	return uc.repo.GetOrderByID(ctx, orderID)
}

func (uc *UseCase) RepairAggregates(ctx context.Context, actor entity.Actor) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.RebuildProductAggregates(ctx)
}

func (uc *UseCase) ListProducts(ctx context.Context, actor entity.Actor, page entity.Pagination) ([]entity.Product, int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListProducts(ctx, page)
}

func (uc *UseCase) UpsertProduct(ctx context.Context, actor entity.Actor, product entity.Product) (entity.Product, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Product{}, err
	}
	return uc.repo.UpsertProduct(ctx, product)
}

func (uc *UseCase) GetProduct(ctx context.Context, actor entity.Actor, productID string) (entity.Product, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Product{}, err
	}
	return uc.repo.GetProduct(ctx, productID)
}

func (uc *UseCase) DeleteProduct(ctx context.Context, actor entity.Actor, productID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteProduct(ctx, productID)
}

func (uc *UseCase) ListCategories(ctx context.Context, actor entity.Actor) ([]entity.Category, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return nil, err
	}
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) UpsertCategory(ctx context.Context, actor entity.Actor, category entity.Category) (entity.Category, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return entity.Category{}, err
	}
	return uc.repo.UpsertCategory(ctx, category)
}

func (uc *UseCase) DeleteCategory(ctx context.Context, actor entity.Actor, categoryID string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	return uc.repo.DeleteCategory(ctx, categoryID)
}

func (uc *UseCase) ImportCards(ctx context.Context, actor entity.Actor, productID string, keys []string) (int, error) {
	if err := uc.ensureAdmin(actor); err != nil {
		return 0, err
	}
	return uc.repo.ImportCards(ctx, productID, keys)
}

func (uc *UseCase) SendBroadcast(ctx context.Context, actor entity.Actor, title, body string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if uc.notifier == nil {
		return nil
	}
	return uc.notifier.SendBroadcast(ctx, title, body)
}

func (uc *UseCase) SendTargeted(ctx context.Context, actor entity.Actor, userID, title, body string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if uc.notifier == nil {
		return nil
	}
	return uc.notifier.SendTargeted(ctx, userID, title, body)
}

func (uc *UseCase) ensureAdmin(actor entity.Actor) error {
	if actor.IsAdmin {
		return nil
	}
	if _, ok := uc.adminUsers[strings.ToLower(strings.TrimSpace(actor.Username))]; ok {
		return nil
	}
	return entity.ErrForbidden
}

func (uc *UseCase) setSetting(ctx context.Context, actor entity.Actor, key, value string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if key == "" {
		return entity.ErrInvalidInput
	}
	return uc.repo.StoreSetting(ctx, entity.Setting{Key: key, Value: value})
}

func (uc *UseCase) deleteSetting(ctx context.Context, actor entity.Actor, key string) error {
	if err := uc.ensureAdmin(actor); err != nil {
		return err
	}
	if key == "" {
		return entity.ErrInvalidInput
	}
	return uc.repo.StoreSetting(ctx, entity.Setting{Key: key, Value: ""})
}

func (uc *UseCase) requireInput(v string) error {
	if strings.TrimSpace(v) == "" {
		return entity.ErrInvalidInput
	}
	return nil
}

func (uc *UseCase) ensureNoErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
