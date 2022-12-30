package notice_svc

import (
	"context"
	"crypto/tls"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type SenderType int

const (
	MailSender SenderType = iota + 1
)

type Sender interface {
	Send(ctx context.Context, user *entity.User, content string, options *sendOptions) error
}

type mail struct {
}

type mailConfig struct {
	SMTP     string
	Port     int
	User     string
	Password string
}

func (m *mail) Send(ctx context.Context, user *entity.User, content string, options *sendOptions) error {
	config := &mailConfig{}
	if err := configs.Default().Scan("mail", config); err != nil {
		return err
	}
	d := gomail.NewDialer(config.SMTP, config.Port, config.User, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	from := "ScriptCat"
	if options.from != nil {
		from = options.from.Username
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(config.User, from))
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "[脚本猫]"+options.title)
	msg.SetBody("text/html", content)

	err := d.DialAndSend(msg)
	if err != nil {
		logger.Ctx(ctx).Error("send email failed",
			zap.Int64("to", user.UID), zap.String("toAddress", user.Email),
			zap.String("from", from),
			zap.Error(err))
		return err
	}
	return nil

}
