package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
	"gorm.io/gorm"
)

func T20230309() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230309",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&script_entity.ScriptDomain{}); err != nil {
				return err
			}
			return tx.AutoMigrate(&statistics_entity.StatisticsInfo{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&statistics_entity.StatisticsInfo{})
		},
	}
}
