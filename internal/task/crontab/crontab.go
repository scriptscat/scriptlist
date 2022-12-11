package crontab

import (
	"context"

	"github.com/codfrm/cago/configs"
	cron "github.com/robfig/cron/v3"
)

// 定时任务

func Crontab(ctx context.Context, config *configs.Config) error {
	c := cron.New()
	// 定时同步数据到es
	_, err := c.AddFunc("0 3 * * *", syncScriptToEs())
	if err != nil {
		return err
	}
	// 定时检查更新

	return nil
}
