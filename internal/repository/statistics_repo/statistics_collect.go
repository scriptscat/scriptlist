package statistics_repo

import (
	"context"

	"github.com/codfrm/cago/database/clickhouse"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsCollectRepo interface {
	Create(ctx context.Context, statistic *statistics_entity.StatisticsCollect) error
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
