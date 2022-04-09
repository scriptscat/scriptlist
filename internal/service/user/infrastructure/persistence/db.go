package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
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

func (r *UserRepositories) AutoMigrate() error {
	return utils.Errs(
		r.db.AutoMigrate(entity.UserConfig{}),
	)
}
