package migrations

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script"
	"gorm.io/gorm"
)

func T1654137068() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654137068",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&script.Script{},
				&script.Code{},
				&script.ScriptCategory{},
				&script.ScriptCategoryList{},
				&entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&script.ScriptDomain{},
				&script.LibDefinition{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&script.Script{},
				&script.Code{},
				&script.ScriptCategory{},
				&script.ScriptCategoryList{},
				&entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&script.ScriptDomain{},
				&script.LibDefinition{},
			)
		},
	}
}
