package notice_svc

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/codfrm/cago/pkg/gogo"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/scriptscat/scriptlist/configs"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/sender"
	"go.uber.org/zap"
)

type NoticeSvc interface {
	// Send 根据模板id发送通知给指定用户
	Send(ctx context.Context, toUser int64, template TemplateID, options ...Option) error
	// MultipleSend 根据模板id发送通知给多个用户
	MultipleSend(ctx context.Context, toUser []int64, template TemplateID, options ...Option) error
}

type noticeSvc struct {
	senderMap map[sender.Type]sender.Sender
}

var defaultNotice = &noticeSvc{
	senderMap: map[sender.Type]sender.Sender{
		sender.MailSender: sender.NewMail(),
	},
}

func Notice() NoticeSvc {
	return defaultNotice
}

// Send 根据模板id发送通知给指定用户
func (n *noticeSvc) Send(ctx context.Context, toUser int64, template TemplateID, options ...Option) error {
	return n.MultipleSend(ctx, []int64{toUser}, template, options...)
}

func (n *noticeSvc) MultipleSend(ctx context.Context, toUsers []int64, template TemplateID, options ...Option) error {
	opts := newOptions(options...)
	tpl, ok := templateMap[template]
	if !ok {
		return errors.New("template not found")
	}
	var err error
	var from *user_entity.User
	if opts.from != 0 {
		from, err = user_repo.User().Find(ctx, opts.from)
		if err != nil {
			return err
		}
	}
	url := configs.Url()
	tplContent := make(map[sender.Type]Template)
	for senderType, tpl := range tpl {
		content, err := n.parseTpl(tpl.Template, map[string]interface{}{
			"Config": map[string]interface{}{
				"Url": url,
			},
			"Value": opts.params,
		})
		if err != nil {
			return err
		}
		title, err := n.parseTpl(tpl.Title, map[string]interface{}{
			"Value": opts.params,
		})
		if err != nil {
			return err
		}
		tplContent[senderType] = Template{
			Title:    title,
			Template: content,
		}
	}
	// 协程去发送邮件
	return gogo.Go(func(ctx context.Context) error {
		for _, toUser := range toUsers {
			to, err := user_repo.User().Find(ctx, toUser)
			if err != nil {
				logger.Ctx(ctx).Error("find error", zap.Error(err), zap.Int64("user_id", toUser))
				continue
			}
			if to == nil {
				logger.Ctx(ctx).Error("user not found", zap.Int64("user_id", toUser))
				continue
			}
			userConfig, err := user_repo.UserConfig().FindByUserID(ctx, toUser)
			if err != nil {
				logger.Ctx(ctx).Error("find user config error", zap.Error(err), zap.Int64("user_id", toUser))
				userConfig = &user_entity.UserConfig{
					Uid:    toUser,
					Notify: &user_entity.Notify{},
				}
				userConfig.Notify.DefaultValue()
			}
			if userConfig == nil {
				userConfig = &user_entity.UserConfig{
					Uid:    toUser,
					Notify: &user_entity.Notify{},
				}
				userConfig.Notify.DefaultValue()
			}
			for senderType, content := range tplContent {
				s, ok := n.senderMap[senderType]
				if !ok {
					return errors.New("sender not found")
				}
				if ok, err := n.IsNotify(ctx, userConfig, senderType, template); err != nil {
					logger.Ctx(ctx).Error("IsNotify error", zap.Error(err), zap.Int64("user_id", toUser))
					continue
				} else if !ok {
					continue
				}
				if err := s.Send(ctx, to, content.Template, &sender.SendOptions{
					From:  from,
					Title: content.Title,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	}, gogo.WithContext(utils.BaseContext(ctx)))
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

// IsNotify 判断用户是否允许通知
func (n *noticeSvc) IsNotify(ctx context.Context, userConfig *user_entity.UserConfig, senderType sender.Type, tpl TemplateID) (bool, error) {
	if senderType != sender.MailSender {
		return true, nil
	}
	switch tpl {
	case ScriptUpdateTemplate:
		return *userConfig.Notify.ScriptUpdate, nil
	case CommentCreateTemplate:
		return *userConfig.Notify.ScriptIssueComment, nil
	case IssueCreateTemplate:
		return *userConfig.Notify.ScriptIssue, nil
	}
	return true, nil
}
