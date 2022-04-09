package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return run(db)
}

func run(db *gorm.DB, fs ...func() *gormigrate.Migration) error {
	var ms []*gormigrate.Migration
	for _, f := range fs {
		ms = append(ms, f())
	}
	m := gormigrate.New(db, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              200,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}, ms)
	if err := m.Migrate(); err != nil {
		return err
	}
	return nil
}
