package persistent

import (
	"fmt"

	"gorm.io/gorm"
)

func catalogDB(db *gorm.DB) (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("catalog repository: database is not configured")
	}
	return db, nil
}
