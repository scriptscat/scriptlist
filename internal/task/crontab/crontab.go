package crontab

import (
	"context"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/server/cron"
	"github.com/scriptscat/scriptlist/internal/task/crontab/handler"
)

type Cron interface {
	Crontab(c cron.Crontab) error
}

// Crontab 定时任务
func Crontab(ctx context.Context, cfg *configs.Config) error {
	// pre环境不执行定时任务,避免冲突
	if configs.Default().Env == configs.PRE {
		return nil
	}
	crontab := []Cron{&handler.Statistics{}, &handler.Script{}}
	for _, v := range crontab {
		if err := v.Crontab(cron.Default()); err != nil {
			return err
		}
	}
	return nil
}
