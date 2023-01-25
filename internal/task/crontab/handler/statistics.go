package handler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Statistics struct {
}

func (s *Statistics) Crontab(ctx context.Context, c *cron.Cron) error {
	_, err := c.AddFunc("0 3 * * *", s.saveScriptStatistics(ctx))
	if err != nil {
		return err
	}
	return nil
}

// 每日统计数据落库
func (s *Statistics) saveScriptStatistics(ctx context.Context) func() {

	return nil
}
