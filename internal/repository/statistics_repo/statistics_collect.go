package statistics_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/clickhouse"
	"github.com/codfrm/cago/database/redis"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
)

type StatisticsCollectRepo interface {
	Create(ctx context.Context, statistic *statistics_entity.StatisticsCollect) error
	CheckLimit(ctx context.Context, scriptId int64) (bool, error)
	GetLimit(ctx context.Context, scriptId int64) (int64, error)
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
	return redis.Ctx(ctx).Get(key).Int64()
}
