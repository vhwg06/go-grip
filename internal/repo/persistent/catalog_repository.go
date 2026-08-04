package persistent

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func catalogDB(db *gorm.DB) (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("catalog repository: database is not configured")
	}
	return db, nil
}

type transactionContextKey struct{}

func withTransactionContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

func catalogDBForContext(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	if ctx != nil {
		if tx, ok := ctx.Value(transactionContextKey{}).(*gorm.DB); ok && tx != nil {
			return tx, nil
		}
	}
	return catalogDB(db)
}
