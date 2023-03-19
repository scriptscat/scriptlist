package migrations

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/go-gormigrate/gormigrate/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RunMigrations 数据库迁移操作
func RunMigrations(db *gorm.DB) error {
	// 添加分布式锁
	if ok, err := redis.Default().
		SetNX(context.Background(), "migrations", "lock", time.Minute).Result(); err != nil {
		logger.Ctx(context.Background()).Error("数据库迁移失败", zap.Error(err))
		return err
	} else if !ok {
		logger.Ctx(context.Background()).Info("数据库迁移已经在执行")
		return nil
	}
	logger.Ctx(context.Background()).Info("开始执行数据库迁移")
	defer redis.Default().Del(context.Background(), "migrations")
	return run(db,
		T1654137068,
		T1654137843,
		T1654138003,
		T1654138087,
		T1670770912,
		T1674744651,
		T1674893758,
		T1675234780,
		T20230210,
		T20230228,
		T20230309,
		T20230319,
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
