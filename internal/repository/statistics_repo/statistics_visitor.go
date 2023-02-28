package statistics_repo

import (
	"context"

	"github.com/codfrm/cago/database/clickhouse"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsVisitorRepo interface {
	Create(ctx context.Context, statistic *statistics_entity.StatisticsVisitor) error
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
