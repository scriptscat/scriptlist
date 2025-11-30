package notification_svc

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
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
	"go.uber.org/zap"
)

type NotificationSvc interface {
	// 查询和管理通知
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	GetUnreadCount(ctx context.Context, req *api.GetUnreadCountRequest) (*api.GetUnreadCountResponse, error)
	Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error)
	MarkRead(ctx context.Context, req *api.MarkReadRequest) error
	BatchMarkRead(ctx context.Context, req *api.BatchMarkReadRequest) (int64, error)

	// 发送通知（原 notice_svc 功能）
	Send(ctx context.Context, toUser int64, template TemplateID, options ...Option) error
	MultipleSend(ctx context.Context, toUsers []int64, template TemplateID, options ...Option) error
}

type notificationSvc struct {
	senderMap map[sender.Type]sender.Sender
}

var defaultNotification = &notificationSvc{
	senderMap: map[sender.Type]sender.Sender{
		sender.MailSender: sender.NewMail(),
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
		items[i] = n.toAPINotification(item, userMap)
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
	total, err := notification_repo.Notification().CountUnread(ctx, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}

	countMap, err := notification_repo.Notification().CountUnreadByType(ctx, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}

	items := make([]*api.UnreadCountItem, 0, len(countMap))
	for t, count := range countMap {
		items = append(items, &api.UnreadCountItem{
			Type:  t,
			Count: count,
		})
	}

	return &api.GetUnreadCountResponse{
		Total: total,
		Items: items,
	}, nil
}

// Get 获取通知详情
func (n *notificationSvc) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	notification, err := notification_repo.Notification().FindByUserID(ctx, auth_svc.Auth().Get(ctx).UID, req.NotificationID)
	if err != nil {
		return nil, err
	}
	if err := notification.CheckOperate(ctx, auth_svc.Auth().Get(ctx).UID); err != nil {
		return nil, err
	}

	// 自动标记为已读
	if notification.ReadStatus == notification_entity.StatusUnread {
		_ = notification_repo.Notification().MarkRead(ctx, auth_svc.Auth().Get(ctx).UID, req.NotificationID, time.Now().Unix())
	}

	// 获取发起用户信息
	userMap := make(map[int64]*user_entity.User)
	if notification.FromUserID > 0 {
		fromUser, err := user_repo.User().Find(ctx, notification.FromUserID)
		if err == nil && fromUser != nil {
			userMap[notification.FromUserID] = fromUser
		}
	}

	return &api.GetResponse{
		Notification: n.toAPINotification(notification, userMap),
	}, nil
}

// MarkRead 标记通知为已读
func (n *notificationSvc) MarkRead(ctx context.Context, req *api.MarkReadRequest) error {
	notification, err := notification_repo.Notification().FindByUserID(ctx, auth_svc.Auth().Get(ctx).UID, req.NotificationID)
	if err != nil {
		return err
	}
	if err := notification.CheckOperate(ctx, auth_svc.Auth().Get(ctx).UID); err != nil {
		return err
	}

	return notification_repo.Notification().MarkRead(ctx, auth_svc.Auth().Get(ctx).UID, req.NotificationID, time.Now().Unix())
}

// BatchMarkRead 批量标记已读
func (n *notificationSvc) BatchMarkRead(ctx context.Context, req *api.BatchMarkReadRequest) (int64, error) {
	if len(req.IDs) == 0 {
		// 全部已读
		return notification_repo.Notification().MarkAllRead(ctx, auth_svc.Auth().Get(ctx).UID, req.Type, time.Now().Unix())
	}
	return notification_repo.Notification().BatchMarkRead(ctx, auth_svc.Auth().Get(ctx).UID, req.IDs, time.Now().Unix())
}

// Send 根据模板id发送通知给指定用户
func (n *notificationSvc) Send(ctx context.Context, toUser int64, template TemplateID, options ...Option) error {
	return n.MultipleSend(ctx, []int64{toUser}, template, options...)
}

// MultipleSend 根据模板id发送通知给多个用户
func (n *notificationSvc) MultipleSend(ctx context.Context, toUsers []int64, template TemplateID, options ...Option) error {
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

	// 协程去发送邮件和保存通知记录
	return gogo.Go(ctx, func(ctx context.Context) error {
		now := time.Now().Unix()
		notifications := make([]*notification_entity.Notification, 0, len(toUsers))

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

			// 获取通知内容用于保存到数据库
			var notifyTitle, notifyContent, notifyLink string
			var notifyType int32
			var extra *notification_entity.Extra

			if content, ok := tplContent[sender.MailSender]; ok {
				notifyTitle = content.Title
				notifyContent = content.Template
			}

			// 根据模板类型设置通知类型和额外数据
			notifyType = n.getNotifyType(template)
			extra = n.extractExtra(opts.params)
			notifyLink = n.buildLink(template, opts.params)

			// 保存通知记录到数据库
			notification := &notification_entity.Notification{
				UserID:     toUser,
				FromUserID: opts.from,
				Type:       notifyType,
				Title:      notifyTitle,
				Content:    notifyContent,
				Link:       notifyLink,
				ReadStatus: notification_entity.StatusUnread,
				Extra:      extra,
				Status:     consts.ACTIVE,
				Createtime: now,
				Updatetime: now,
			}
			notifications = append(notifications, notification)

			// 发送邮件
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
					logger.Ctx(ctx).Error("send notification error", zap.Error(err), zap.Int64("user_id", toUser))
				}
			}
		}

		// 批量保存通知记录
		if len(notifications) > 0 {
			if err := notification_repo.Notification().BatchCreate(ctx, notifications); err != nil {
				logger.Ctx(ctx).Error("batch create notifications error", zap.Error(err))
				return err
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

// IsNotify 判断用户是否允许通知
func (n *notificationSvc) IsNotify(ctx context.Context, userConfig *user_entity.UserConfig, senderType sender.Type, tpl TemplateID) (bool, error) {
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
	default:
		return true, nil
	}
}

// toAPINotification 转换为API通知对象
func (n *notificationSvc) toAPINotification(entity *notification_entity.Notification, userMap map[int64]*user_entity.User) *api.Notification {
	notification := &api.Notification{
		ID:         entity.ID,
		UserID:     entity.UserID,
		FromUserID: entity.FromUserID,
		Type:       entity.Type,
		Title:      entity.Title,
		Content:    entity.Content,
		Link:       entity.Link,
		ReadStatus: entity.ReadStatus,
		ReadTime:   entity.ReadTime,
		Createtime: entity.Createtime,
		Updatetime: entity.Updatetime,
	}

	if entity.FromUserID > 0 {
		if user, ok := userMap[entity.FromUserID]; ok {
			notification.FromUser = user.UserInfo()
		}
	}

	return notification
}

// getNotifyType 根据模板获取通知类型
func (n *notificationSvc) getNotifyType(template TemplateID) int32 {
	switch template {
	case ScriptUpdateTemplate:
		return notification_entity.TypeScriptUpdate
	case IssueCreateTemplate:
		return notification_entity.TypeIssueCreate
	case CommentCreateTemplate:
		return notification_entity.TypeCommentCreate
	case ScriptScoreTemplate:
		return notification_entity.TypeScriptScore
	case AccessInviteTemplate:
		return notification_entity.TypeAccessInvite
	case ScriptScoreReplyTemplate:
		return notification_entity.TypeScriptScoreReply
	default:
		return notification_entity.TypeSystem
	}
}

// extractExtra 从参数中提取额外数据
func (n *notificationSvc) extractExtra(params interface{}) *notification_entity.Extra {
	extra := &notification_entity.Extra{}
	// 这里可以根据不同的参数类型提取相应的数据
	// 由于params是interface{}，需要使用类型断言
	// 这部分可以根据实际的template参数结构来完善
	return extra
}

// buildLink 构建通知链接
func (n *notificationSvc) buildLink(template TemplateID, params interface{}) string {
	// 根据模板类型和参数构建链接
	// 这部分可以根据实际需求来完善
	return ""
}
