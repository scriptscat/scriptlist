package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/user"
	"gorm.io/gorm"
)

func T1654138003() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654138003",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(&user.UserConfig{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&user.UserConfig{})
		},
	}
}
