package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"gorm.io/gorm"
)

// UnitOfWork coordinates one infrastructure transaction for an application
// operation. Repositories receive the transaction through the context and
// remain independent of this coordinator's concrete database handle.
type UnitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork creates a transaction coordinator for the configured
// PostgreSQL connection.
func NewUnitOfWork(pg *postgres.Postgres) *UnitOfWork {
	if pg == nil {
		return &UnitOfWork{}
	}
	return &UnitOfWork{db: pg.Gorm}
}

var _ repo.UnitOfWork = (*UnitOfWork)(nil)

// Within executes fn atomically and rolls back when any operation returns an
// error. The callback context is the only transaction handle visible to
// repository adapters.
func (u *UnitOfWork) Within(ctx context.Context, fn func(context.Context) error) error {
	if u == nil || u.db == nil {
		return fmt.Errorf("unit of work: database is not configured")
	}
	if fn == nil {
		return fmt.Errorf("unit of work: callback is nil")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(withTransactionContext(ctx, tx))
	})
}
