package notification_entity

import (
	"context"
	"net/http"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// Type 通知模板类型
type Type int

const (
	ScriptUpdateTemplate     Type = iota + 100 // 脚本更新
	IssueCreateTemplate                        // 问题创建
	CommentCreateTemplate                      // 评论创建
	ScriptScoreTemplate                        // 脚本评分
	AccessInviteTemplate                       // 访问邀请
	ScriptScoreReplyTemplate                   // 脚本评分回复
)

// 已读状态
const (
	StatusUnread int32 = 1 // 未读
	StatusRead   int32 = 2 // 已读
)

// Notification 通知实体
type Notification struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key;autoIncrement"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null;index:idx_user_id"`
	FromUserID int64  `gorm:"column:from_user_id;type:bigint(20);default:0"`
	Type       Type   `gorm:"column:type;type:int(11);not null;"`
	Title      string `gorm:"column:title;type:varchar(255);not null"`
	Content    string `gorm:"column:content;type:text"`
	Link       string `gorm:"column:link;type:varchar(512);default:''"` // 关联链接
	ReadStatus int32  `gorm:"column:read_status;type:tinyint(4);default:1;not null;index:idx_user_status,priority:2"`
	ReadTime   int64  `gorm:"column:read_time;type:bigint(20);default:0"`
	Params     string `gorm:"column:params;type:json"` // 额外信息
	Status     int32  `gorm:"column:status;type:tinyint(4);default:1;not null;index:idx_user_status,priority:1;index:idx_user_type,priority:1"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

// CheckOperate 检查是否可以操作
func (n *Notification) CheckOperate(ctx context.Context, userID int64) error {
	if n == nil {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.NotificationNotFound)
	}
	// 非激活状态
	if n.Status != consts.ACTIVE {
		return i18n.NewErrorWithStatus(ctx, http.StatusNotFound, code.NotificationNotFound)
	}
	// 检查是否属于当前用户
	if n.UserID != userID {
		return i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.NotificationPermissionDenied)
	}
	return nil
}

// MarkRead 标记为已读
func (n *Notification) MarkRead(readTime int64) {
	n.ReadStatus = StatusRead
	n.ReadTime = readTime
}

// MarkUnread 标记为未读
func (n *Notification) MarkUnread() {
	n.ReadStatus = StatusUnread
	n.ReadTime = 0
}
