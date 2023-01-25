package handler

import (
	"context"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"go.uber.org/zap"
)

type Script struct {
}

func (s *Script) Crontab(ctx context.Context, c *cron.Cron) error {
	_, err := c.AddFunc("0 */6 * * *", s.checkSyncUpdate(ctx))
	if err != nil {
		return err
	}
	return nil
}

// 检查设置的同步更新
func (s *Script) checkSyncUpdate(ctx context.Context) func() {
	return func() {
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
				return
			}
			if len(list) == 0 {
				logger.Info("检查同步更新结束")
				return
			}
			for _, v := range list {
				if err := script_svc.Script().SyncOnce(ctx, v); err != nil {
					logger.Error("脚本检查更新失败", zap.Int64("script_id", v.ID),
						zap.String("sync_url", v.SyncUrl), zap.Error(err))
				}
			}
			page++
		}
	}
}
