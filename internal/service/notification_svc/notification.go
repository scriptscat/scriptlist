package notification_svc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"time"

	"github.com/cago-frame/cago/pkg/gogo"
	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/configs"
	api "github.com/scriptscat/scriptlist/internal/api/notification"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/repository/notification_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc/sender"
	template2 "github.com/scriptscat/scriptlist/internal/service/notification_svc/template"
	"go.uber.org/zap"
)

type NotificationSvc interface {
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	GetUnreadCount(ctx context.Context, req *api.GetUnreadCountRequest) (*api.GetUnreadCountResponse, error)
	MarkRead(ctx context.Context, req *api.MarkReadRequest) error
	BatchMarkRead(ctx context.Context, req *api.BatchMarkReadRequest) error

	// Send 发送通知
	Send(ctx context.Context, toUser int64, notificationType notification_entity.Type, options ...Option) error
	// MultipleSend 发送通知给多个用户
	MultipleSend(ctx context.Context, toUsers []int64, notificationType notification_entity.Type, options ...Option) error
}

type notificationSvc struct {
	senderMap map[sender.Type]sender.Sender
}

var defaultNotification = &notificationSvc{
	senderMap: map[sender.Type]sender.Sender{
		sender.InAppSender: sender.NewApp(),
		sender.MailSender:  sender.NewMail(),
	},
}

func Notification() NotificationSvc {
	return defaultNotification
}

// List 获取通知列表
func (n *notificationSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := notification_repo.Notification().FindPage(ctx, auth_svc.Auth().Get(ctx).UID, req)
	if err != nil {
		return nil, err
	}

	// 收集发起用户ID
	fromUserIDs := make(map[int64]bool)
	for _, item := range list {
		if item.FromUserID > 0 {
			fromUserIDs[item.FromUserID] = true
		}
	}

	// 批量查询用户信息
	userMap := make(map[int64]*user_entity.User)
	for userId := range fromUserIDs {
		user, err := user_repo.User().Find(ctx, userId)
		if err != nil {
			logger.Ctx(ctx).Error("find user error", zap.Error(err), zap.Int64("user_id", userId))
			continue
		}
		if user != nil {
			userMap[userId] = user
		}
	}

	// 组装响应
	items := make([]*api.Notification, len(list))
	for i, item := range list {
		items[i] = n.toAPINotification(ctx, item, userMap)
	}

	return &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Notification]{
			Total: total,
			List:  items,
		},
	}, nil
}

// GetUnreadCount 获取未读通知数量
func (n *notificationSvc) GetUnreadCount(ctx context.Context, req *api.GetUnreadCountRequest) (*api.GetUnreadCountResponse, error) {
	total, err := notification_repo.Notification().CountUnread(ctx, auth_svc.Auth().Get(ctx).UID, 0)
	if err != nil {
		return nil, err
	}

	return &api.GetUnreadCountResponse{
		Total: total,
	}, nil
}

// MarkRead 标记通知为已读或未读
func (n *notificationSvc) MarkRead(ctx context.Context, req *api.MarkReadRequest) error {
	notification, err := notification_repo.Notification().Find(ctx, auth_svc.Auth().Get(ctx).UID, req.NotificationID)
	if err != nil {
		return err
	}
	if err := notification.CheckOperate(ctx, auth_svc.Auth().Get(ctx).UID); err != nil {
		return err
	}

	if req.Unread == 0 {
		notification.MarkRead(time.Now().Unix())
	} else {
		notification.MarkUnread()
	}

	if err := notification_repo.Notification().Update(ctx, notification); err != nil {
		return err
	}
	return nil
}

// BatchMarkRead 批量标记已读
func (n *notificationSvc) BatchMarkRead(ctx context.Context, req *api.BatchMarkReadRequest) error {
	return notification_repo.Notification().BatchMarkRead(ctx, auth_svc.Auth().Get(ctx).UID, req.IDs)
}

// Send 根据模板id发送通知给指定用户
func (n *notificationSvc) Send(ctx context.Context, toUser int64, notificationType notification_entity.Type, options ...Option) error {
	return n.MultipleSend(ctx, []int64{toUser}, notificationType, options...)
}

type Link interface {
	Link() string
}

// MultipleSend 根据模板id发送通知给多个用户
func (n *notificationSvc) MultipleSend(ctx context.Context, toUsers []int64, notificationType notification_entity.Type, options ...Option) error {
	opts := newOptions(options...)
	tpl, ok := template2.TplMap[notificationType]
	if !ok {
		return errors.New("notificationType not found")
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
	tplContent := make(map[sender.Type]template2.Template)
	for senderType, tpl := range tpl {
		content, err := n.parseTpl(tpl.Content, map[string]interface{}{
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
		t := template2.Template{
			Title:   title,
			Content: content,
			Link:    "",
		}
		if link, ok := opts.params.(Link); ok {
			t.Link = url + link.Link()
		}
		tplContent[senderType] = t
	}

	// 协程去发送邮件和保存通知记录
	return gogo.Go(ctx, func(ctx context.Context) error {
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

			// 发送邮件
			for senderType, content := range tplContent {
				s, ok := n.senderMap[senderType]
				if !ok {
					return errors.New("sender not found")
				}
				if ok, err := n.IsNotify(ctx, userConfig, senderType, notificationType); err != nil {
					logger.Ctx(ctx).Error("IsNotify error", zap.Error(err), zap.Int64("user_id", toUser))
					continue
				} else if !ok {
					continue
				}
				if err := s.Send(ctx, to, content.Content, &sender.SendOptions{
					From:   from,
					Title:  content.Title,
					Type:   notificationType,
					Link:   content.Link,
					Params: opts.params,
				}); err != nil {
					logger.Ctx(ctx).Error("send notification error", zap.Error(err), zap.Int64("user_id", toUser))
				}
			}
		}

		return nil
	})
}

func (n *notificationSvc) parseTpl(tpl string, data interface{}) (string, error) {
	t := template.New("tpl")
	t = template.Must(t.Parse(tpl))
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// IsNotify 判断用户是否允许邮件通知
func (n *notificationSvc) IsNotify(ctx context.Context, userConfig *user_entity.UserConfig,
	senderType sender.Type, tpl notification_entity.Type) (bool, error) {
	if senderType != sender.MailSender {
		return true, nil
	}
	switch tpl {
	case notification_entity.ScriptUpdateTemplate:
		return *userConfig.Notify.ScriptUpdate, nil
	case notification_entity.CommentCreateTemplate:
		return *userConfig.Notify.ScriptIssueComment, nil
	case notification_entity.IssueCreateTemplate:
		return *userConfig.Notify.ScriptIssue, nil
	default:
		return true, nil
	}
}

// toAPINotification 转换为API通知对象
func (n *notificationSvc) toAPINotification(ctx context.Context, entity *notification_entity.Notification, userMap map[int64]*user_entity.User) *api.Notification {
	notification := &api.Notification{
		ID:         entity.ID,
		UserID:     entity.UserID,
		Type:       entity.Type,
		Title:      entity.Title,
		Content:    entity.Content,
		Link:       entity.Link,
		ReadStatus: entity.ReadStatus,
		ReadTime:   entity.ReadTime,
		Createtime: entity.Createtime,
		Updatetime: entity.Updatetime,
	}

	var params map[string]interface{}
	if err := json.Unmarshal([]byte(entity.Params), &params); err != nil {
		logger.Ctx(ctx).Error("unmarshal notification params error", zap.Error(err), zap.Int64("notification_id", entity.ID))
	} else {
		notification.Params = params
	}

	if entity.FromUserID > 0 {
		if user, ok := userMap[entity.FromUserID]; ok {
			notification.FromUser = user.UserInfo()
		}
	}

	return notification
}
