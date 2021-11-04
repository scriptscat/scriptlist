package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptweb/internal/domain/script/entity"
	entity3 "github.com/scriptscat/scriptweb/internal/domain/user/entity"
	"gorm.io/gorm"
)

func T1636014908() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1636014908",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AddColumn(&entity.ScriptCategoryList{}, "sort"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&entity3.UserConfig{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&entity.ScriptCategoryList{}, "sort"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&entity3.UserConfig{}); err != nil {
				return err
			}
			return nil
		},
	}
}
