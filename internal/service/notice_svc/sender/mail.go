package sender

import (
	"context"
	"crypto/tls"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type mailConfig struct {
	SMTP     string
	Port     int
	User     string
	Password string
}

type mail struct {
}

func NewMail() Sender {
	return &mail{}
}

func (m *mail) Send(ctx context.Context, user *user_entity.User, content string, options *SendOptions) error {
	config := &mailConfig{}
	if err := configs.Default().Scan("mail", config); err != nil {
		return err
	}
	d := gomail.NewDialer(config.SMTP, config.Port, config.User, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	from := "ScriptCat"
	if options.From != nil {
		from = options.From.Username
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(config.User, from))
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", "[脚本猫]"+options.Title)
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
