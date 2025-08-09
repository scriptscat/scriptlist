package handler

import (
	"context"

	"github.com/cago-frame/cago/server/cron"
)

type Crx struct {
	ctxList []string
}

func (s *Crx) Crontab(c cron.Crontab) error {
	s.ctxList = []string{}
	_, err := c.AddFunc("0 */6 * * *", s.checkCrxUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Crx) checkCrxUpdate(ctx context.Context) error {
	//
	return nil
}
