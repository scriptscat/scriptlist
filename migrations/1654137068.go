package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"gorm.io/gorm"
)

func T1654137068() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654137068",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&entity.Script{},
				&entity.ScriptCode{},
				&entity.ScriptCategory{},
				&entity.ScriptCategoryList{},
				&entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&entity.ScriptDomain{},
				&entity.LibDefinition{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&entity.Script{},
				&entity.ScriptCode{},
				&entity.ScriptCategory{},
				&entity.ScriptCategoryList{},
				&entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&entity.ScriptDomain{},
				&entity.LibDefinition{},
			)
		},
	}
}
