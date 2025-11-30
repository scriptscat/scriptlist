package notification_entity

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"net/http"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
)

// 通知类型定义
const (
	TypeScriptUpdate     int32 = 1  // 脚本更新
	TypeIssueCreate      int32 = 2  // 反馈创建
	TypeCommentCreate    int32 = 3  // 评论创建
	TypeScriptScore      int32 = 4  // 脚本评分
	TypeAccessInvite     int32 = 5  // 协作邀请
	TypeScriptScoreReply int32 = 6  // 评分回复
	TypeSystem           int32 = 99 // 系统通知
)

// 已读状态
const (
	StatusUnread int32 = 0 // 未读
	StatusRead   int32 = 1 // 已读
)

// Extra 通知额外数据
type Extra struct {
	ScriptID   int64  `json:"script_id,omitempty"`
	IssueID    int64  `json:"issue_id,omitempty"`
	CommentID  int64  `json:"comment_id,omitempty"`
	ScoreID    int64  `json:"score_id,omitempty"`
	InviteID   int64  `json:"invite_id,omitempty"`
	CustomData string `json:"custom_data,omitempty"`
}

func (e *Extra) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, e)
}

func (e *Extra) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

// Notification 通知实体
type Notification struct {
	ID         int64  `gorm:"column:id;type:bigint(20);not null;primary_key;autoIncrement"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);not null;index:idx_user_id"`
	FromUserID int64  `gorm:"column:from_user_id;type:bigint(20);default:0"`
	Type       int32  `gorm:"column:type;type:int(11);not null;index:idx_user_type,priority:2"`
	Title      string `gorm:"column:title;type:varchar(255);not null"`
	Content    string `gorm:"column:content;type:text"`
	Link       string `gorm:"column:link;type:varchar(512);default:''"`
	ReadStatus int32  `gorm:"column:read_status;type:tinyint(4);default:0;not null;index:idx_user_status,priority:2"`
	ReadTime   int64  `gorm:"column:read_time;type:bigint(20);default:0"`
	Extra      *Extra `gorm:"column:extra;type:json"`
	Status     int32  `gorm:"column:status;type:tinyint(4);default:1;not null;index:idx_user_status,priority:1;index:idx_user_type,priority:1"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)"`
}

func (Notification) TableName() string {
	return "notification"
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
