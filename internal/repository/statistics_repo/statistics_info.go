package statistics_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/cache"
	cache2 "github.com/codfrm/cago/database/cache/cache"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsInfoRepo interface {
	Find(ctx context.Context, id int64) (*statistics_entity.StatisticsInfo, error)
	FindByScriptId(ctx context.Context, scriptId int64) (*statistics_entity.StatisticsInfo, error)
	FindByStatisticsKey(ctx context.Context, statisticsKey string) (*statistics_entity.StatisticsInfo, error)
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

func (u *statisticsInfoRepo) key(statisticsKey string) string {
	return fmt.Sprintf("statistics:info:statisticsKey:%v", statisticsKey)
}

func (u *statisticsInfoRepo) FindByScriptId(ctx context.Context, scriptId int64) (*statistics_entity.StatisticsInfo, error) {
	ret := &statistics_entity.StatisticsInfo{}
	if err := db.Ctx(ctx).Where("script_id=? and status=?", scriptId, consts.ACTIVE).First(ret).Error; err != nil {
		if db.RecordNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *statisticsInfoRepo) FindByStatisticsKey(ctx context.Context, statisticsKey string) (*statistics_entity.StatisticsInfo, error) {
	ret := &statistics_entity.StatisticsInfo{}
	if err := cache.Ctx(ctx).GetOrSet(u.key(statisticsKey), func() (interface{}, error) {
		if err := db.Ctx(ctx).Where("statistics_key=? and status=?", statisticsKey, consts.ACTIVE).First(ret).Error; err != nil {
			if db.RecordNotFound(err) {
				return nil, nil
			}
			return nil, err
		}
		return ret, nil
	}, cache2.Expiration(time.Hour)).Scan(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (u *statisticsInfoRepo) Create(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error {
	if err := cache.Ctx(ctx).Del(u.key(statisticsInfo.StatisticsKey)); err != nil {
		return err
	}
	return db.Ctx(ctx).Create(statisticsInfo).Error
}

func (u *statisticsInfoRepo) Update(ctx context.Context, statisticsInfo *statistics_entity.StatisticsInfo) error {
	if err := db.Ctx(ctx).Updates(statisticsInfo).Error; err != nil {
		return err
	}
	return cache.Ctx(ctx).Del(u.key(statisticsInfo.StatisticsKey))
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
