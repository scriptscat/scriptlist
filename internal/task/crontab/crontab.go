package crontab

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptlist/internal/task/crontab/handler"
)

type Cron interface {
	Crontab(ctx context.Context, c *cron.Cron) error
}

// Crontab 定时任务
func Crontab(ctx context.Context, config *configs.Config) error {
	c := cron.New()
	crontab := []Cron{&handler.Statistics{}, &handler.Script{}}
	for _, v := range crontab {
		if err := v.Crontab(ctx, c); err != nil {
			return err
		}
	}
	c.Start()
	return nil
}
