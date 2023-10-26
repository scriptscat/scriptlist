package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T20231026() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231026",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&script_entity.Script{})
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	}
}
