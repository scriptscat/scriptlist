package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func T1670770912() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1670770912",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate()
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable()
		},
	}
}
