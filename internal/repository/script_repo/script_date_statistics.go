package script_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cago-frame/cago/database/cache"
	cache2 "github.com/cago-frame/cago/database/cache/cache"
	"github.com/cago-frame/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"gorm.io/gorm"
)

type ScriptDateStatisticsRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptDateStatistics, error)
	Create(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error
	Update(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error
	Delete(ctx context.Context, id int64) error

	FindByScriptID(ctx context.Context, scriptId int64, t time.Time) (*entity.ScriptDateStatistics, error)
	IncrDownload(ctx context.Context, scriptId int64, t time.Time, num int64) error
	IncrUpdate(ctx context.Context, scriptId int64, t time.Time, num int64) error
}

var defaultScriptDateStatistics ScriptDateStatisticsRepo

func ScriptDateStatistics() ScriptDateStatisticsRepo {
	return defaultScriptDateStatistics
}

func RegisterScriptDateStatistics(i ScriptDateStatisticsRepo) {
	defaultScriptDateStatistics = i
}

type scriptDateStatisticsRepo struct {
}

func NewScriptDateStatistics() ScriptDateStatisticsRepo {
	return &scriptDateStatisticsRepo{}
}

func (u *scriptDateStatisticsRepo) Find(ctx context.Context, id int64) (*entity.ScriptDateStatistics, error) {
	ret := &entity.ScriptDateStatistics{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *scriptDateStatisticsRepo) Create(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error {
	return db.Ctx(ctx).Create(scriptDateStatistics).Error
}

func (u *scriptDateStatisticsRepo) Update(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error {
	return db.Ctx(ctx).Updates(scriptDateStatistics).Error
}

func (u *scriptDateStatisticsRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.ScriptDateStatistics{ID: id}).Error
}

func (u *scriptDateStatisticsRepo) key(id int64) string {
	return "script:date:statistics:" + fmt.Sprintf("%d", id)
}

func (u *scriptDateStatisticsRepo) FindByScriptID(ctx context.Context, scriptId int64, t time.Time) (*entity.ScriptDateStatistics, error) {
	ret := &entity.ScriptDateStatistics{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(scriptId)+":"+t.Format("2006-01-02"), func() (interface{}, error) {
		if err := db.Ctx(ctx).Where("script_id=? and date=?", scriptId, t.Format("2006-01-02")).
			First(ret).Error; err != nil {
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

func (u *scriptDateStatisticsRepo) IncrDownload(ctx context.Context, scriptId int64, t time.Time, num int64) error {
	date := t.Format("2006-01-02")
	if db.Ctx(ctx).Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", scriptId, date).
		Update("download", gorm.Expr("`download`+?", num)).RowsAffected == 0 {
		if err := db.Ctx(ctx).Save(&entity.ScriptDateStatistics{
			ScriptID: scriptId,
			Date:     date,
			Download: 1,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (u *scriptDateStatisticsRepo) IncrUpdate(ctx context.Context, scriptId int64, t time.Time, num int64) error {
	date := t.Format("2006-01-02")
	if db.Ctx(ctx).Model(&entity.ScriptDateStatistics{}).Where("script_id=? and date=?", scriptId, date).
		Update("update", gorm.Expr("`update`+?", num)).RowsAffected == 0 {
		if err := db.Ctx(ctx).Save(&entity.ScriptDateStatistics{
			ScriptID: scriptId,
			Date:     date,
			Update:   1,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
