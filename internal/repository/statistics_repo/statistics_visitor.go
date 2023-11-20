package statistics_repo

import (
	"context"
	"github.com/codfrm/cago/database/db"
	"time"

	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsVisitorRepo interface {
	Create(ctx context.Context, statistic []*statistics_entity.StatisticsVisitor) error
	FirstUserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error)
	UserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error)
	IpNumber(ctx context.Context, id int64, start time.Time, end time.Time) (int64, error)
	// OriginList 安装来源列表
	OriginList(ctx context.Context, scriptId int64, startTime, endTime time.Time, page httputils.PageRequest) ([]*api.PieChart, int64, error)
	VersionPie(ctx context.Context, scriptId int64, startTime, endTime time.Time) ([]*api.PieChart, error)
	DriverPie(ctx context.Context, scriptId int64, startTime, endTime time.Time) ([]*api.PieChart, error)
	BrowserPie(ctx context.Context, scriptId int64, startTime, endTime time.Time) ([]*api.PieChart, error)
}

var defaultStatisticsVisitorRepo StatisticsVisitorRepo

func StatisticsVisitor() StatisticsVisitorRepo {
	return defaultStatisticsVisitorRepo
}

func RegisterStatisticsVisitorRepo(i StatisticsVisitorRepo) {
	defaultStatisticsVisitorRepo = i
}

type statisticsVisitorRepo struct {
}

func NewStatisticVistior() StatisticsVisitorRepo {
	return &statisticsVisitorRepo{}
}

func (u *statisticsVisitorRepo) Create(ctx context.Context, statistic []*statistics_entity.StatisticsVisitor) error {
	return db.CtxWith(ctx, "clickhouse").CreateInBatches(statistic, len(statistic)+1).Error
}

func (u *statisticsVisitorRepo) FirstUserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	err := db.CtxWith(ctx, "clickhouse").
		Model(&statistics_entity.StatisticsVisitor{}).
		Where("script_id = ? and first_visit_time >= ? and first_visit_time <= ?", scriptId, startTime, endTime).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *statisticsVisitorRepo) UserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	err := db.CtxWith(ctx, "clickhouse").
		Model(&statistics_entity.StatisticsVisitor{}).
		Where("script_id = ? and visit_time >= ? and visit_time <= ?", scriptId, startTime, endTime).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *statisticsVisitorRepo) OriginList(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time, page httputils.PageRequest) ([]*api.PieChart, int64, error) {
	var total int64
	result := make([]*api.PieChart, 0)
	query := db.CtxWith(ctx, "clickhouse").Model(&statistics_entity.StatisticsVisitor{}).Select(
		"install_page as key, count(*) as value",
	).Group("install_page").
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

func (u *statisticsVisitorRepo) VersionPie(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) ([]*api.PieChart, error) {
	result := make([]*api.PieChart, 0)
	if err := db.CtxWith(ctx, "clickhouse").Model(&statistics_entity.StatisticsVisitor{}).Select(
		"version as key, count(*) as value",
	).Group("version").
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix()).
		Limit(10).
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u *statisticsVisitorRepo) IpNumber(ctx context.Context, id int64, start time.Time, end time.Time) (int64, error) {
	var count int64
	err := db.CtxWith(ctx, "clickhouse").
		Model(&statistics_entity.StatisticsVisitor{}).
		Where("script_id = ? and visit_time >= ? and visit_time <= ?", id, start.Unix(), end.Unix()).
		Group("ip").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *statisticsVisitorRepo) DriverPie(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) ([]*api.PieChart, error) {
	result := make([]*api.PieChart, 0)
	if err := db.CtxWith(ctx, "clickhouse").Model(&statistics_entity.StatisticsVisitor{}).Select(
		"device_type as key, count(*) as value",
	).Group("device_type").
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix()).
		Limit(10).
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (u *statisticsVisitorRepo) BrowserPie(ctx context.Context, scriptId int64, startTime time.Time, endTime time.Time) ([]*api.PieChart, error) {
	result := make([]*api.PieChart, 0)
	if err := db.CtxWith(ctx, "clickhouse").Model(&statistics_entity.StatisticsVisitor{}).Select(
		"browser_type as key, count(*) as value",
	).Group("browser_type").
		Where("script_id=? and visit_time >= ? and visit_time <= ?", scriptId, startTime.Unix(), endTime.Unix()).
		Limit(10).
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
