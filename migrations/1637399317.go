package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/domain/issue/entity"
	"gorm.io/gorm"
)

func T1637399317() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1637399317",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&entity.ScriptIssue{}); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&entity.ScriptIssueComment{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&entity.ScriptIssueComment{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&entity.ScriptIssueComment{}); err != nil {
				return err
			}
			return nil
		},
	}
}
