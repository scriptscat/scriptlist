package handler

import (
	"context"

	"github.com/cago-frame/cago/database/db"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/server/cron"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"go.uber.org/zap"
)

type Invite struct {
	ctxList []string
}

func (s *Invite) Crontab(c cron.Crontab) error {
	s.ctxList = []string{}
	_, err := c.AddFunc("0 0 * * *", s.checkInviteUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Invite) checkInviteUpdate(ctx context.Context) error {

	err := db.Ctx(ctx).Model(&script_entity.ScriptInvite{}).Where("expiretime != 0 AND expiretime < unix_timestamp() AND invite_status = ?", script_entity.InviteStatusUnused).Update("invite_status", script_entity.InviteStatusExpired).Error
	if err != nil {
		logger.Ctx(ctx).Error("Invite Expiretime Check Error!", zap.Error(err))
		return err
	}
	return err
}
