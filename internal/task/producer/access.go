package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type AccessInviteMsg struct {
	UserID       int64 `json:"user_id"`
	InviteUserID int64 `json:"invite_user_id"`
	Invite       *script_entity.ScriptInvite
}

func PublishAccessInvite(ctx context.Context, user, inviteUserID int64, invite *script_entity.ScriptInvite) error {
	body, err := json.Marshal(&AccessInviteMsg{
		UserID:       user,
		InviteUserID: inviteUserID,
		Invite:       invite,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptAccessInviteTopic, &broker2.Message{
		Body: body,
	})
}

func SubscribeAccessInvite(ctx context.Context, fn func(ctx context.Context,
	userId, inviteUserId int64, invite *script_entity.ScriptInvite) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptAccessInviteTopic, func(ctx context.Context, ev broker2.Event) error {
		data := &AccessInviteMsg{}
		err := json.Unmarshal(ev.Message().Body, data)
		if err != nil {
			return err
		}
		return fn(ctx, data.UserID, data.InviteUserID, data.Invite)
	}, opts...)
	return err
}
