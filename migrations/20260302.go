package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/audit_entity"
	"gorm.io/gorm"
)

func T20260302() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20260302",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&audit_entity.AuditLog{},
			); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&audit_entity.AuditLog{},
			)
		},
	}
}
