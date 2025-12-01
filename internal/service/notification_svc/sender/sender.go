package sender

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Type int

const (
	MailSender  Type = iota + 1
	InAppSender      // 应用内通知
)

type Sender interface {
	Send(ctx context.Context, user *user_entity.User, content string, options *SendOptions) error
}
