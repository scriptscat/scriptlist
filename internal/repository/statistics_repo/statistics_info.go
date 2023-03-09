package statistics_repo

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsInfoRepo interface {
	Find(ctx context.Context, id int64) (*statistics_entity.StatisticsInfo, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*statistics_entity.StatisticsInfo, int64, error)
	Create(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error
	Update(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error
	Delete(ctx context.Context, id int64) error
}

var defaultStatisticsInfo StatisticsInfoRepo

func StatisticsInfo() StatisticsInfoRepo {
	return defaultStatisticsInfo
}

func RegisterStatisticsInfo(i StatisticsInfoRepo) {
	defaultStatisticsInfo = i
}

type statisticsInfoRepo struct {
}

func NewStatisticsInfo() StatisticsInfoRepo {
	return &statisticsInfoRepo{}
}

func (u *statisticsInfoRepo) Find(ctx context.Context, id int64) (*statistics_entity.StatisticsInfo, error) {
	ret := &statistics_entity.StatisticsInfo{}
	if err := db.Ctx(ctx).Where("id=? and status=?", id, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *statisticsInfoRepo) Create(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error {
	return db.Ctx(ctx).Create(statisticsInfo).Error
}

func (u *statisticsInfoRepo) Update(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error {
	return db.Ctx(ctx).Updates(statisticsInfo).Error
}

func (u *statisticsInfoRepo) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Model(&statistics_entity.StatisticsInfo{}).
		Where("id=?", id).Update("status", consts.DELETE).Error
}

func (u *statisticsInfoRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*statistics_entity.StatisticsInfo, int64, error) {
	var list []*statistics_entity.StatisticsInfo
	var count int64
	find := db.Ctx(ctx).Model(&statistics_entity.StatisticsInfo{}).Where("status=?", consts.ACTIVE)
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := find.Order("createtime desc").Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
