package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/repository"
	"gorm.io/gorm"
)

type UserRepositories struct {
	db *gorm.DB
	repository.User
	repository.Follow
}

func NewRepositories(db *gorm.DB, redis *redis.Client, cache cache.Cache) *UserRepositories {
	return &UserRepositories{
		db:     db,
		User:   NewUser(db, redis, cache),
		Follow: NewFollow(db),
	}
}
