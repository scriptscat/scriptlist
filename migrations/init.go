package migrations

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// 添加分布式锁
	if ok, err := redis.Default().
		SetNX(context.Background(), "migrations", "lock", time.Minute).Result(); err != nil {
		return err
	} else if !ok {
		return nil
	}
	defer redis.Default().Del(context.Background(), "migrations")
	return run(db,
		T1654137068,
		T1654137843,
		T1654138003,
		T1654138087,
	)
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
