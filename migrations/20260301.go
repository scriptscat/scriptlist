package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"gorm.io/gorm"
)

func T20260301() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20260301",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&report_entity.ScriptReport{}); err != nil {
				return err
			}
			return tx.AutoMigrate(&report_entity.ScriptReportComment{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&report_entity.ScriptReportComment{}); err != nil {
				return err
			}
			return tx.Migrator().DropTable(&report_entity.ScriptReport{})
		},
	}
}
