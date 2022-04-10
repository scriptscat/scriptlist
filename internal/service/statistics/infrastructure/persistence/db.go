package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/statistics/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/statistics/entity"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/gorm"
)

type StatisRepositories struct {
	db *gorm.DB
	repository.Statistics
}

func NewRepositories(db *gorm.DB, redis *redis.Client) *StatisRepositories {
	return &StatisRepositories{
		db:         db,
		Statistics: NewStatistics(db, redis),
	}
}

func (r *StatisRepositories) AutoMigrate() error {
	return utils.ErrFunc(
		func() error {
			return r.db.AutoMigrate(&entity.StatisticsDownload{})
		},
		func() error {
			return r.db.AutoMigrate(&entity.StatisticsUpdate{})
		},
	)
}
