package notice_svc

import (
	"context"

	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type SenderType int

const (
	MailSender SenderType = iota + 1
)

type Sender interface {
	Send(ctx context.Context, user *user_entity.User, content string, options *SendOptions) error
}
