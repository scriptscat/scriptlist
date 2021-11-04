package service

import (
	"crypto/tls"

	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendEmail(to, title, content, contextType string) error
}

type sender struct {
	config config.Email
}

func NewSender(config config.Email) Sender {
	return &sender{
		config: config,
	}
}

func (s *sender) SendEmail(to, title, content, contextType string) error {
	d := gomail.NewDialer(s.config.Smtp, s.config.Port, s.config.User, s.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[脚本猫]"+title)
	m.SetBody(contextType, content)

	return d.DialAndSend(m)
}
