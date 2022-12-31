package script_statistics_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity"
)

type ScriptStatisticsRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptStatistics, error)
	Create(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error
	Update(ctx context.Context, scriptStatistics *entity.ScriptStatistics) error
	Delete(ctx context.Context, id int64) error

	FindByScriptID(ctx context.Context, scriptId int64) (*entity.ScriptStatistics, error)
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

func (u *scriptStatisticsRepo) FindByScriptID(ctx context.Context, scriptId int64) (*entity.ScriptStatistics, error) {
	ret := &entity.ScriptStatistics{}
	if err := db.Ctx(ctx).First(ret, "script_id=?", scriptId).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
