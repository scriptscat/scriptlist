package statistics_repo

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/clickhouse"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsVisitorRepo interface {
	Create(ctx context.Context, statistic *statistics_entity.StatisticsVisitor) error
	FirstUserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error)
	UserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error)
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

func (u *statisticsVisitorRepo) Create(ctx context.Context, statistic *statistics_entity.StatisticsVisitor) error {
	return clickhouse.Ctx(ctx).Create(statistic).Error
}

func (u *statisticsVisitorRepo) FirstUserNumber(ctx context.Context, scriptId int64, startTime, endTime time.Time) (int64, error) {
	var count int64
	err := clickhouse.Ctx(ctx).
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
	err := clickhouse.Ctx(ctx).
		Model(&statistics_entity.StatisticsVisitor{}).
		Where("script_id = ? and visit_time >= ? and visit_time <= ?", scriptId, startTime, endTime).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
