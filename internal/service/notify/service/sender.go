package service

import (
	"crypto/tls"

	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendEmail(to, title, content, contextType string) error
	SendEmailFrom(from, to, title, content, contextType string) error
	NotifyEmail(to, title, content, contextType string) error
	NotifyEmailFrom(from, to, title, content, contextType string) error
}

type sender struct {
	config config.Email
	notify config.Email
}

func NewSender(config config.Email, notify config.Email) Sender {
	return &sender{
		config: config,
		notify: notify,
	}
}

func (s *sender) SendEmail(to, title, content, contextType string) error {
	return s.SendEmailFrom("ScriptCat", to, title, content, contextType)
}

func (s *sender) SendEmailFrom(from, to, title, content, contextType string) error {
	d := gomail.NewDialer(s.config.Smtp, s.config.Port, s.config.User, s.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.notify.User, from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[脚本猫]"+title)
	m.SetBody(contextType, content)

	err := d.DialAndSend(m)
	if err != nil {
		logrus.Errorf("send %v %v email: %v", to, title, err)
	}
	return err
}

func (s *sender) NotifyEmail(to, title, content, contextType string) error {
	return s.NotifyEmailFrom("ScriptCat", to, title, content, contextType)
}

func (s *sender) NotifyEmailFrom(from, to, title, content, contextType string) error {
	d := gomail.NewDialer(s.notify.Smtp, s.notify.Port, s.notify.User, s.notify.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.notify.User, from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[脚本猫]"+title)
	m.SetBody(contextType, content)

	err := d.DialAndSend(m)
	if err != nil {
		logrus.Errorf("send %v %v email: %v", to, title, err)
	}
	return err
}