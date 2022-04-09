package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/safe/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/gorm"
)

type SafeRepositories struct {
	db *gorm.DB
	repository.Rate
}

func NewRepositories(redis *redis.Client) *SafeRepositories {
	return &SafeRepositories{
		Rate: NewRate(redis),
	}
}

func (r *SafeRepositories) AutoMigrate() error {
	return utils.Errs()
}
