package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T20250107() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "T20250107",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&script_entity.ScriptScoreReply{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&script_entity.ScriptScoreReply{},
			)
		},
	}
}
