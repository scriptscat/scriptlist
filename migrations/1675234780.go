package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"gorm.io/gorm"
)

func T1675234780() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1675234780",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&script_entity.Code{})
		}, Rollback: func(db *gorm.DB) error {
			return nil
		},
	}
}
