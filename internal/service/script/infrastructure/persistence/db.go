package persistence

import (
	"github.com/go-redis/redis/v8"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
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

func NewRepositories(db *gorm.DB, redis *redis.Client, cache cache.Cache) *Repositories {
	return &Repositories{
		db:          db,
		Script:      NewScript(db, redis),
		Code:        NewCode(db, cache),
		Score:       NewScore(db, redis),
		Category:    NewCategory(db),
		ScriptWatch: NewScriptWatch(db, redis),
		Statistics:  NewStatistics(db),
	}
}

func (r *Repositories) AutoMigrate() error {
	return r.db.AutoMigrate(
		&entity.Script{},
		&entity.ScriptCode{},
		&entity.ScriptCategory{},
		&entity.ScriptCategoryList{},
		&entity.ScriptScore{},
		&entity.ScriptStatistics{},
		&entity.ScriptDateStatistics{},
		&entity.ScriptDomain{},
		&entity.LibDefinition{},
	)
}
