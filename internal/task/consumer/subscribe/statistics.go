package subscribe

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

// Statistics 处理统计平台数据
type Statistics struct {
}

func (e *Statistics) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptStatistics(ctx, e.scriptStatistics); err != nil {
		return err

	}
	return nil
}

func (e *Statistics) scriptStatistics(ctx context.Context, msg *producer.ScriptStatisticsMsg) error {
	switch msg.Download {
	case statistics_repo.DownloadStatistics:
		if ok, err := statistics_repo.Statistics().IncrDownload(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			return err
		} else if ok {
			// 统计总量
			if err := script_repo.ScriptStatistics().IncrDownload(ctx, msg.ScriptID); err != nil {
				return err
			}
			// 统计当日
			if err := script_repo.ScriptDateStatistics().IncrDownload(ctx, msg.ScriptID, msg.Time); err != nil {
				return err
			}
		}
	case statistics_repo.UpdateStatistics:
		if ok, err := statistics_repo.Statistics().IncrUpdate(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			return err
		} else if ok {
			// 统计总量
			if err := script_repo.ScriptStatistics().IncrUpdate(ctx, msg.ScriptID); err != nil {
				return err
			}
			// 统计当日
			if err := script_repo.ScriptDateStatistics().IncrUpdate(ctx, msg.ScriptID, msg.Time); err != nil {
				return err
			}
		}
	case statistics_repo.ViewStatistics:
		if _, err := statistics_repo.Statistics().IncrPageView(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			return err
		}
	}

	return nil
}