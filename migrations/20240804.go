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
			var migratorInstance gorm.Migrator = tx.Migrator()
			if migratorInstance.HasIndex(&script_entity.ScriptInvite{}, "invite_status") {
				err := migratorInstance.DropIndex(&script_entity.ScriptInvite{}, "invite_status")
				if err != nil {
					return err
				}
			}
			if migratorInstance.HasIndex(&script_entity.ScriptInvite{}, "expiretime") {
				err := migratorInstance.DropIndex(&script_entity.ScriptInvite{}, "expiretime")
				if err != nil {
					return err
				}
			}
			return tx.AutoMigrate(
				&script_entity.ScriptInvite{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
