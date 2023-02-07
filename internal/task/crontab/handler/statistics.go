package handler

import (
	"context"

	"github.com/codfrm/cago/server/cron"
)

type Statistics struct {
}

func (s *Statistics) Crontab(c cron.Crontab) error {
	//_, err := c.AddFunc("0 3 * * *", s.saveScriptStatistics)
	//if err != nil {
	//	return err
	//}
	return nil
}

// 每日统计数据落库
func (s *Statistics) saveScriptStatistics(ctx context.Context) error {
	return nil
}
