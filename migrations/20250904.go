package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/feedback_entity"
	"gorm.io/gorm"
)

func T20250904() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20250904",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&feedback_entity.Feedback{},
			); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&feedback_entity.Feedback{},
			)
		},
	}
}
