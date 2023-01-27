package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"gorm.io/gorm"
)

func T1654138087() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1654138087",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&issue_entity.ScriptIssue{},
				&issue_entity.ScriptIssueComment{},
			)
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&issue_entity.ScriptIssue{},
				&issue_entity.ScriptIssueComment{},
			)
		},
	}
}
