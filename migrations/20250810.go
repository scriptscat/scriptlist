package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T20250810() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20250810",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&script_entity.ScriptCategoryList{},
			); err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&script_entity.ScriptCategoryList{},
			)
		},
	}
}
