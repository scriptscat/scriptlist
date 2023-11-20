package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T20231220() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231220",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&script_entity.Script{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(script_entity.ScriptAccess{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(script_entity.ScriptGroup{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(script_entity.ScriptGroupMember{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
