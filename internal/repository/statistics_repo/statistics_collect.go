package statistics_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/clickhouse"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsCollectRepo interface {
	FindPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*statistics_entity.StatisticsCollect, int64, error)

	Create(ctx context.Context, statistic *statistics_entity.StatisticsCollect) error
	CheckLimit(ctx context.Context, scriptId int64) (bool, error)
	GetLimit(ctx context.Context, scriptId int64) (int64, error)
	RealtimeChart(ctx context.Context, scriptId int64, now time.Time) ([]*Realtime, error)
	Pv(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (int64, error)
	Uv(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (int64, error)
	UseTimeAvg(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (float64, error)
	// OperationHostList 操作域列表
	OperationHostList(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time, page httputils.PageRequest) ([]*api.PieChart, int64, error)
}

var defaultStatisticsCollect StatisticsCollectRepo

func StatisticsCollect() StatisticsCollectRepo {
	return defaultStatisticsCollect
}

func RegisterStatisticsCollect(i StatisticsCollectRepo) {
	defaultStatisticsCollect = i
}

type statisticsCollectRepo struct {
}

func NewStatisticsCollect() StatisticsCollectRepo {
	return &statisticsCollectRepo{}
}

func (u *statisticsCollectRepo) Create(ctx context.Context, statistic *statistics_entity.StatisticsCollect) error {
	return clickhouse.Ctx(ctx).Create(statistic).Error
}

func (u *statisticsCollectRepo) CheckLimit(ctx context.Context, scriptId int64) (bool, error) {
	key := fmt.Sprintf("statistics:limit:%s:%d", time.Now().Format("2006-01"), scriptId)
	n, err := redis.Ctx(ctx).Incr(key).Result()
	if err != nil {
		return false, err
	}
	return n < 1000000, nil
}

func (u *statisticsCollectRepo) GetLimit(ctx context.Context, scriptId int64) (int64, error) {
	key := fmt.Sprintf("statistics:limit:%s:%d", time.Now().Format("2006-01"), scriptId)
	num, err := redis.Ctx(ctx).Get(key).Int64()
	if err != nil {
		if redis.Nil(err) {
			return 0, nil
		}
		return 0, err
	}
	return num, nil
}

type Realtime struct {
	Time int
	Num  int64
}

func (u *statisticsCollectRepo) RealtimeChart(ctx context.Context, scriptId int64, now time.Time) ([]*Realtime, error) {
	result := make([]*Realtime, 0)
	if err := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).Select(
		"FROM_UNIXTIME(visit_time, '%M') as time, count(*) as num",
	).Group("FROM_UNIXTIME(visit_time, '%M')").
		Where("script_id=? and visit_time >= ?", scriptId, now.Add(-time.Minute*15).Unix()).
		Order("time desc").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u *statisticsCollectRepo) Pv(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (int64, error) {
	var num int64
	err := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix()).
		Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (u *statisticsCollectRepo) Uv(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (int64, error) {
	var num int64
	err := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).
		Select("count(distinct visitor_id) as num").
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix()).
		Scan(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (u *statisticsCollectRepo) UseTimeAvg(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) (float64, error) {
	var t float64
	err := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).
		Select("avg(duration) as t").
		Where("script_id=? and visit_time >= ? and visit_time <= ? and duration!=0", scriptId, startTime.Unix(), endTime.Unix()).
		Scan(&t).Error
	if err != nil {
		return 0, err
	}
	return t, nil
}

func (u *statisticsCollectRepo) OperationHostList(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time, page httputils.PageRequest) ([]*api.PieChart, int64, error) {
	var total int64
	result := make([]*api.PieChart, 0)
	query := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).Select(
		"operation_host as key, count(*) as value",
	).Group("operation_host").
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix())
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(page.GetLimit()).Offset(page.GetOffset()).
		Order("value desc").Scan(&result).Error; err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (u *statisticsCollectRepo) FindPage(ctx context.Context, scriptId int64, page httputils.PageRequest) ([]*statistics_entity.StatisticsCollect, int64, error) {
	var list []*statistics_entity.StatisticsCollect
	var count int64
	query := clickhouse.Ctx(ctx).Model(&statistics_entity.StatisticsCollect{}).
		Where("script_id=?", scriptId)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("visit_time desc").
		Offset(page.GetOffset()).Limit(page.GetLimit()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
