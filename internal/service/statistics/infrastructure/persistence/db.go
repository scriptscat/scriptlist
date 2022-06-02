package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/service/statistics/domain/repository"
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
