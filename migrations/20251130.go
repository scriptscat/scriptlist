package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"gorm.io/gorm"
)

func T20251130() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20251130",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&notification_entity.Notification{},
			); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&notification_entity.Notification{},
			)
		},
	}
}
