package migrations

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"gorm.io/gorm"
)

func T1654137068() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654137068",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&user_entity.User{},
				&script_entity.Script{},
				&script_entity.Code{},
				&script_entity.ScriptCategory{},
				&script_entity.ScriptCategoryList{},
				&script_entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&script_entity.ScriptDomain{},
				&script_entity.LibDefinition{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&script_entity.Script{},
				&script_entity.Code{},
				&script_entity.ScriptCategory{},
				&script_entity.ScriptCategoryList{},
				&script_entity.ScriptScore{},
				&entity.ScriptStatistics{},
				&entity.ScriptDateStatistics{},
				&script_entity.ScriptDomain{},
				&script_entity.LibDefinition{},
			)
		},
	}
}
