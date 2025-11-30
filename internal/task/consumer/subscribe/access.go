package subscribe

import (
	"context"

	"github.com/cago-frame/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type Access struct {
}

func (s *Access) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeAccessInvite(ctx, s.Invite); err != nil {
		return err
	}
	return nil
}

func (s *Access) Invite(ctx context.Context, userId, inviteUserId int64, invite *script_entity.ScriptInvite) error {
	script, err := script_repo.Script().Find(ctx, invite.ScriptID)
	if err != nil {
		return err
	}
	if script == nil {
		logger.Ctx(ctx).Error("脚本不存在", zap.Int64("invite", invite.ID), zap.Int64("script", invite.ScriptID))
		return nil
	}
	user, err := user_repo.User().Find(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		logger.Ctx(ctx).Error("用户不存在", zap.Int64("invite", invite.ID), zap.Int64("user", userId))
		return nil
	}
	return notification_svc.Notification().Send(ctx, inviteUserId, notification_svc.AccessInviteTemplate,
		notification_svc.WithParams(&template.AccessInvite{
			Code:     invite.Code,
			Name:     script.Name,
			Username: user.Username,
		}), notification_svc.WithFrom(userId))
}
