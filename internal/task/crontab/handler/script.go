package handler

import (
	"context"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/cron"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"go.uber.org/zap"
)

type Script struct {
}

func (s *Script) Crontab(c cron.Crontab) error {
	_, err := c.AddFunc("0 */6 * * *", s.checkSyncUpdate)
	if err != nil {
		return err
	}
	return nil
}

// 检查设置的同步更新
func (s *Script) checkSyncUpdate(ctx context.Context) error {
	if ok, err := redis.Ctx(ctx).SetNX("checkSyncUpdate", "1", time.Minute).Result(); err != nil {
		logger.Ctx(ctx).Error("检查同步失败", zap.Error(err))
		return err
	} else if !ok {
		logger.Ctx(ctx).Info("其他机器检查同步更新中")
		return nil
	}
	page := 1
	logger.Ctx(ctx).Info("检查同步更新开始")
	for {
		logger := logger.Ctx(ctx).With(zap.Int("page", page))
		list, err := script_repo.Script().FindSyncScript(ctx, httputils.PageRequest{
			Page: page,
			Size: 20,
		})
		if err != nil {
			logger.Error("checkSyncUpdate", zap.Error(err))
			return err
		}
		if len(list) == 0 {
			logger.Info("检查同步更新结束")
			return nil
		}
		for _, v := range list {
			if err := script_svc.Script().SyncOnce(ctx, v); err != nil {
				logger.Error("脚本检查更新失败", zap.Int64("script_id", v.ID),
					zap.String("sync_url", v.SyncUrl), zap.Error(err))
			} else {
				logger.Info("脚本检查更新成功", zap.Int64("script_id", v.ID),
					zap.String("sync_url", v.SyncUrl))
			}
		}
		page++
	}
}
