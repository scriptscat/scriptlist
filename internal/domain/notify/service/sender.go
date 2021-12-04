package service

import (
	"crypto/tls"
	"mime"

	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendEmail(to, title, content, contextType string) error
	SendEmailFrom(from, to, title, content, contextType string) error
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
	m.SetHeader("From", mime.QEncoding.Encode("utf-8", "ScriptCat")+" <"+s.config.User+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[脚本猫]"+title)
	m.SetBody(contextType, content)

	err := d.DialAndSend(m)
	if err != nil {
		logrus.Errorf("send %v %v email: %v", to, title, err)
	}
	return err
}

func (s *sender) SendEmailFrom(from, to, title, content, contextType string) error {
	d := gomail.NewDialer(s.config.Smtp, s.config.Port, s.config.User, s.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("utf-8", from)+" <"+s.config.User+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "[脚本猫]"+title)
	m.SetBody(contextType, content)

	err := d.DialAndSend(m)
	if err != nil {
		logrus.Errorf("send %v %v email: %v", to, title, err)
	}
	return err
}
