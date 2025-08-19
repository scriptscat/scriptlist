package subscribe

import (
	"context"
	"fmt"
	"time"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// Statistics 处理统计平台数据
type Statistics struct {
}

func (s *Statistics) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptStatistics(ctx, s.scriptStatistics); err != nil {
		return err
	}
	// TODO: 每天扫描并同步一次数据
	return nil
}

func (s *Statistics) statisticSyncKey(scriptId int64, key string) string {
	return fmt.Sprintf("script:statistic:sync:statistic:%d:%s", scriptId, key)
}

func SyncIncr(ctx context.Context, key, field string, update func(ctx context.Context, num int64) error) error {
	num, err := redis.Ctx(ctx).HIncrBy(key, field+"_num", 1).Result()
	if err != nil {
		return err
	}
	// 当囤了1000条记录或者时间超过了5分钟, 进行更新
	if num < 1000 {
		t, err := redis.Ctx(ctx).HGet(key, field+"_time").Int64()
		if err != nil {
			if !redis.Nil(err) {
				return err
			}
		}
		if time.Now().Unix()-t < 300 {
			return nil
		}
	}
	// 统计总量
	if err := update(ctx, num); err != nil {
		logger.Ctx(ctx).Error("更新失败", zap.Error(err))
		return err
	}
	// 设置时间
	if err := redis.Ctx(ctx).HSet(key, field+"_time", time.Now().Unix()).Err(); err != nil {
		logger.Ctx(ctx).Error("设置时间失败", zap.Error(err))
		return err
	}
	// 重置数量
	if err := redis.Ctx(ctx).HSet(key, field+"_num", 0).Err(); err != nil {
		logger.Ctx(ctx).Error("重置数量失败", zap.Error(err))
		return err
	}
	return nil
}

func (s *Statistics) scriptStatistics(ctx context.Context, msg *producer.ScriptStatisticsMsg) error {
	switch msg.Download {
	case statistics_repo.DownloadScriptStatistics:
		if ok, err := statistics_repo.ScriptStatistics().IncrDownload(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计下载量失败", zap.Error(err))
			return err
		} else if ok {
			// 统计总量
			if err := SyncIncr(ctx, s.statisticSyncKey(msg.ScriptID, "download"), "total",
				func(ctx context.Context, num int64) error {
					return script_repo.ScriptStatistics().IncrDownload(ctx, msg.ScriptID, num)
				}); err != nil {
				logger.Ctx(ctx).Error("统计总更新量失败", zap.Error(err))
			}
			// 统计当日
			if err := SyncIncr(ctx, s.statisticSyncKey(msg.ScriptID, "download"), msg.Time.Format("2006-01-02"),
				func(ctx context.Context, num int64) error {
					return script_repo.ScriptDateStatistics().IncrDownload(ctx, msg.ScriptID, msg.Time, num)
				}); err != nil {
				logger.Ctx(ctx).Error("统计总更新量失败", zap.Error(err))
			}
		}
	case statistics_repo.UpdateScriptStatistics:
		if ok, err := statistics_repo.ScriptStatistics().IncrUpdate(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计更新量失败", zap.Error(err))
			return err
		} else if ok {
			// 统计总量
			if err := SyncIncr(ctx, s.statisticSyncKey(msg.ScriptID, "update"), "total",
				func(ctx context.Context, num int64) error {
					return script_repo.ScriptStatistics().IncrUpdate(ctx, msg.ScriptID, num)
				}); err != nil {
				logger.Ctx(ctx).Error("统计总更新量失败", zap.Error(err))
			}
			// 统计当日
			if err := SyncIncr(ctx, s.statisticSyncKey(msg.ScriptID, "update"), msg.Time.Format("2006-01-02"),
				func(ctx context.Context, num int64) error {
					return script_repo.ScriptDateStatistics().IncrUpdate(ctx, msg.ScriptID, msg.Time, num)
				}); err != nil {
				logger.Ctx(ctx).Error("统计当日更新量失败", zap.Error(err))
			}
		}
	case statistics_repo.ViewScriptStatistics:
		if _, err := statistics_repo.ScriptStatistics().IncrPageView(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计浏览量失败", zap.Error(err))
			return err
		}
	}

	return nil
}
