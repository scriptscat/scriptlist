package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"gorm.io/gorm"
)

func T1654137843() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654137843",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&entity.StatisticsDownload{}, &entity.StatisticsUpdate{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&entity.StatisticsDownload{}, &entity.StatisticsUpdate{})
		},
	}
}
