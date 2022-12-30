package notice_svc

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/scriptscat/scriptlist/internal/model/entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
)

type NoticeSvc interface {
	// Send 根据模板id发送通知给指定用户
	Send(ctx context.Context, toUser int64, template int, options ...Option) error
	// MultipleSend 根据模板id发送通知给多个用户
	MultipleSend(ctx context.Context, toUser []int64, template int, options ...Option) error
}

type noticeSvc struct {
	senderMap map[SenderType]Sender
}

var defaultNotice = &noticeSvc{
	senderMap: map[SenderType]Sender{
		MailSender: &mail{},
	},
}

func Notice() NoticeSvc {
	return defaultNotice
}

// Send 根据模板id发送通知给指定用户
func (n *noticeSvc) Send(ctx context.Context, toUser int64, template int, options ...Option) error {
	return n.MultipleSend(ctx, []int64{toUser}, template, options...)
}

func (n *noticeSvc) MultipleSend(ctx context.Context, toUsers []int64, template int, options ...Option) error {
	opts := newOptions(options...)
	tpl, ok := templateMap[template]
	if !ok {
		return errors.New("template not found")
	}
	var err error
	var from *entity.User
	if opts.from != 0 {
		from, err = user_repo.User().Find(ctx, opts.from)
		if err != nil {
			return err
		}
	}
	senderOptions := &sendOptions{
		from:  from,
		title: opts.title,
	}
	tplContent := make(map[SenderType]string)
	for senderType, tpl := range tpl {
		content, err := n.parseTpl(tpl, map[string]interface{}{
			"Value": opts.params,
		})
		if err != nil {
			return err
		}
		tplContent[senderType] = content
	}
	for _, toUser := range toUsers {
		to, err := user_repo.User().Find(ctx, toUser)
		if err != nil {
			return err
		}
		if to == nil {
			return errors.New("user not found")
		}
		for senderType, content := range tplContent {
			sender, ok := n.senderMap[senderType]
			if !ok {
				return errors.New("sender not found")
			}
			if err := sender.Send(ctx, to, content, senderOptions); err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *noticeSvc) parseTpl(tpl string, data interface{}) (string, error) {
	t := template.New("tpl")
	t = template.Must(t.Parse(tpl))
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
