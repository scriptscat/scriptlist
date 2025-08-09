package script_repo

import (
	"context"
	"strconv"
	"time"

	"github.com/cago-frame/cago/database/cache"
	cache2 "github.com/cago-frame/cago/database/cache/cache"
	"github.com/cago-frame/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source=script_statistics.go -destination=mock/script_statistics.go
type ScriptStatisticsRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptStatistics, error)
	Create(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error
	Update(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error
	Delete(ctx context.Context, id int64) error

	FindByScriptID(ctx context.Context, scriptId int64) (*entity.ScriptStatistics, error)
	// IncrDownload 增加下载量,不会去重
	IncrDownload(ctx context.Context, scriptId, num int64) error
	IncrUpdate(ctx context.Context, scriptId int64, num int64) error
	// IncrScore 分数统计,当用户分数变更时可以使用之前的分数和之后的分数进行计算,num为0
	IncrScore(ctx context.Context, scriptId, score int64, num int) error
}

var defaultScriptStatistics ScriptStatisticsRepo

func ScriptStatistics() ScriptStatisticsRepo {
	return defaultScriptStatistics
}

func RegisterScriptStatistics(i ScriptStatisticsRepo) {
	defaultScriptStatistics = i
}

type scriptStatisticsRepo struct {
}

func NewScriptStatistics() ScriptStatisticsRepo {
	return &scriptStatisticsRepo{}
}

func (u *scriptStatisticsRepo) Find(ctx context.Context, id int64) (*entity.ScriptStatistics, error) {
	ret := &entity.ScriptStatistics{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptStatisticsRepo) Create(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error {
	return db.Ctx(ctx).Create(scriptStatistics).Error
}

func (u *scriptStatisticsRepo) Update(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error {
	return db.Ctx(ctx).Updates(scriptStatistics).Error
}

func (u *scriptStatisticsRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.ScriptStatistics{ID: id}).Error
}

func (u *scriptStatisticsRepo) key(id int64) string {
	return "script:statistics:" + strconv.FormatInt(id, 10)
}

func (u *scriptStatisticsRepo) FindByScriptID(ctx context.Context, scriptId int64) (*entity.ScriptStatistics, error) {
	ret := &entity.ScriptStatistics{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId), func() (interface{}, error) {
		if err := db.Ctx(ctx).First(ret, "script_id=?", scriptId).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Minute)).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *scriptStatisticsRepo) IncrDownload(ctx context.Context, scriptId, num int64) error {
	if db.Ctx(ctx).Model(&entity.ScriptStatistics{}).Where("script_id=?", scriptId).
		Update("download", gorm.Expr("`download`+?", num)).RowsAffected == 0 {
		return db.Ctx(ctx).Save(&entity.ScriptStatistics{
			ScriptID: scriptId,
			Download: 1,
		}).Error
	}
	return nil
}

func (u *scriptStatisticsRepo) IncrUpdate(ctx context.Context, scriptId int64, num int64) error {
	if db.Ctx(ctx).Model(&entity.ScriptStatistics{}).Where("script_id=?", scriptId).
		Update("update", gorm.Expr("`update`+?", num)).RowsAffected == 0 {
		return db.Ctx(ctx).Save(&entity.ScriptStatistics{
			ScriptID: scriptId,
			Update:   1,
		}).Error
	}
	return nil
}

func (u *scriptStatisticsRepo) IncrScore(ctx context.Context, scriptId, score int64, num int) error {
	if db.Ctx(ctx).Model(&entity.ScriptStatistics{}).Where("script_id=?", scriptId).
		Updates(map[string]interface{}{
			"score":       gorm.Expr("`score`+?", score),
			"score_count": gorm.Expr("`score_count`+?", num),
		}).RowsAffected == 0 {
		return db.Ctx(ctx).Save(&entity.ScriptStatistics{
			ScriptID:   scriptId,
			Score:      int64(score),
			ScoreCount: int64(num),
		}).Error
	}
	return nil
}
