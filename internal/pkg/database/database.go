package database

import (
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"gorm.io/driver/mysql"
	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDatabase(cfg config.MySQL, debug bool) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.Dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
	})
	if err == nil {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
	}
	if err != nil {
		return nil, err
	}
	if debug {
		db = db.Debug()
	}
	return db, nil
}
