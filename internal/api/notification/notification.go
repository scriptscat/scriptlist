package notification

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

// Notification 通知信息
type Notification struct {
	ID         int64                    `json:"id"`
	UserID     int64                    `json:"user_id"`
	FromUser   user_entity.UserInfo     `json:"from_user,omitempty"` // 发起用户信息
	Type       notification_entity.Type `json:"type"`                // 通知类型
	Title      string                   `json:"title"`               // 通知标题
	Content    string                   `json:"content"`             // 通知内容
	Params     interface{}              `json:"params,omitempty"`    // 额外参数
	ReadStatus int32                    `json:"read_status"`         // 0:未读 1:已读
	ReadTime   int64                    `json:"read_time,omitempty"` // 阅读时间
	Createtime int64                    `json:"createtime"`
	Updatetime int64                    `json:"updatetime"`
}

// ListRequest 获取通知列表
type ListRequest struct {
	mux.Meta              `path:"/notifications" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ReadStatus            int32 `form:"read_status" binding:"omitempty,oneof=0 1 2"` // 0:全部 1:未读 2:已读
}

type ListResponse struct {
	httputils.PageResponse[*Notification] `json:",inline"`
}

// GetUnreadCountRequest 获取未读通知数量
type GetUnreadCountRequest struct {
	mux.Meta `path:"/notifications/unread-count" method:"GET"`
}

type GetUnreadCountResponse struct {
	Total int64 `json:"total"` // 总未读数
}

// GetRequest 获取通知详情
type GetRequest struct {
	mux.Meta       `path:"/notifications/:id" method:"GET"`
	NotificationID int64 `uri:"id" binding:"required"`
}

type GetResponse struct {
	*Notification `json:",inline"`
}

// MarkReadRequest 标记通知为已读
type MarkReadRequest struct {
	mux.Meta       `path:"/notifications/:id/read" method:"PUT"`
	NotificationID int64 `uri:"id" binding:"required"`
	Unread         int32 `json:"unread" binding:"omitempty"`
}

type MarkReadResponse struct{}

// BatchMarkReadRequest 批量标记已读
type BatchMarkReadRequest struct {
	mux.Meta `path:"/notifications/read" method:"PUT"`
	IDs      []int64 `json:"ids" binding:"omitempty"` // 通知ID列表，为空则全部标记已读
}

type BatchMarkReadResponse struct {
}
