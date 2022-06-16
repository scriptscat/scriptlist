package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/pkg/gofound"
	"gorm.io/gorm"
)

type Repositories struct {
	db          *gorm.DB
	Script      repository.Script
	Code        repository.ScriptCode
	Score       repository.Score
	Category    repository.Category
	ScriptWatch repository.ScriptWatch
	Statistics  repository.Statistics
}

func NewRepositories(db *gorm.DB, redis *redis.Client, cache cache.Cache, goFound *gofound.GOFound) *Repositories {
	return &Repositories{
		db:          db,
		Script:      NewScript(db, redis, goFound),
		Code:        NewCode(db, cache),
		Score:       NewScore(db, redis),
		Category:    NewCategory(db, cache),
		ScriptWatch: NewScriptWatch(db, redis),
		Statistics:  NewStatistics(db),
	}
}
