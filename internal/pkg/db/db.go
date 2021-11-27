package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	goRedis "github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

var Redis *goRedis.Client

var Cache cache.Cache

func Init() error {
	var err error
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: config.AppConfig.Mysql.Dsn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.AppConfig.Mysql.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	db, err := Db.DB()
	if err != nil {
		return err
	}
	if config.AppConfig.Mode == "debug" {
		Db = Db.Debug()
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Second * 15)
	db.SetConnMaxIdleTime(time.Second * 30)
	Redis = goRedis.NewClient(&goRedis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		return err
	}
	redisCache := goRedis.NewClient(&goRedis.Options{
		Addr:     config.AppConfig.Cache.Addr,
		Password: config.AppConfig.Cache.Password,
		DB:       config.AppConfig.Cache.DB,
	})
	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		return err
	}
	Cache = cache.NewRedisCache(redisCache)

	return nil
}

func MockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	Db, _ = gorm.Open(mysql.Open(config.AppConfig.Mysql.Dsn), &gorm.Config{})
	return db, mock
}
