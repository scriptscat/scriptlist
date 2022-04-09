package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/resource/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/gorm"
)

type Repositories struct {
	db *gorm.DB
	repository.Resource
}

func NewRepositories(redis *redis.Client) *Repositories {
	return &Repositories{
		Resource: NewResource(redis),
	}
}

func (r *Repositories) AutoMigrate() error {
	return utils.Errs()
}
