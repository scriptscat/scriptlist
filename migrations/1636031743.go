package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/domain/statistics/entity"
	"gorm.io/gorm"
)

func T1636031743() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1636031743",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AddColumn(&entity.StatisticsDownload{}, "statistics_token"); err != nil {
				return err
			}
			if err := tx.Migrator().AddColumn(&entity.StatisticsUpdate{}, "statistics_token"); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&entity.StatisticsDownload{}, "statistics_token"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&entity.StatisticsUpdate{}, "statistics_token"); err != nil {
				return err
			}
			return nil
		},
	}
}
