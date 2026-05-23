package persistent

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func withTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if db == nil {
		return fmt.Errorf("withTransaction: db is nil")
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("withTransaction: begin transaction: %w", tx.Error)
	}

	if err := fn(tx); err != nil {
		rollbackErr := tx.Rollback().Error
		if rollbackErr != nil {
			return fmt.Errorf("withTransaction: rollback failed (%v) after error: %w", rollbackErr, err)
		}

		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("withTransaction: commit transaction: %w", err)
	}

	return nil
}

func forUpdate(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return nil
	}

	return tx.Clauses(clause.Locking{Strength: "UPDATE"})
}
