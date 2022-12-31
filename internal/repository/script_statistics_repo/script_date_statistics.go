package script_statistics_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/scriptscat/scriptlist/internal/model/entity"
)

type ScriptDateStatisticsRepo interface {
	Find(ctx context.Context, id int64) (*entity.ScriptDateStatistics, error)
	Create(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error
	Update(ctx context.Context, scriptDateStatistics *entity.ScriptDateStatistics) error
	Delete(ctx context.Context, id int64) error

	FindByScriptID(ctx context.Context, scriptId int64, date string) (*entity.ScriptDateStatistics, error)
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

func (u *scriptDateStatisticsRepo) FindByScriptID(ctx context.Context, scriptId int64, date string) (*entity.ScriptDateStatistics, error) {
	ret := &entity.ScriptDateStatistics{}
	if err := db.Ctx(ctx).Where("script_id=? and date=?", scriptId, date).
		First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
