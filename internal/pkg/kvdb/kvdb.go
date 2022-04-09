package kvdb

import (
	"context"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
)

type KvDb interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Client() interface{}
	DbType() string
}

func NewKvDb(cfg config.Redis) (*goRedis.Client, error) {
	ret := goRedis.NewClient(&goRedis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := ret.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return ret, nil
}
