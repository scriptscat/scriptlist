package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T20240804() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240804",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&script_entity.ScriptInvite{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
