package notification_ctr

import (
	"github.com/cago-frame/cago/pkg/utils/muxutils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/notification"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc"
)

type Notification struct {
}

func NewNotification() *Notification {
	return &Notification{}
}

func (n *Notification) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{auth_svc.Auth().RequireLogin(true)},
		Handler: []interface{}{
			n.List,
			n.GetUnreadCount,
			n.Get,
			n.MarkRead,
			n.BatchMarkRead,
		},
	}})
}

// List 获取通知列表
func (n *Notification) List(ctx *gin.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return notification_svc.Notification().List(ctx, req)
}

// GetUnreadCount 获取未读通知数量
func (n *Notification) GetUnreadCount(ctx *gin.Context, req *api.GetUnreadCountRequest) (*api.GetUnreadCountResponse, error) {
	return notification_svc.Notification().GetUnreadCount(ctx, req)
}

// Get 获取通知详情
func (n *Notification) Get(ctx *gin.Context, req *api.GetRequest) (*api.GetResponse, error) {
	return notification_svc.Notification().Get(ctx, req)
}

// MarkRead 标记通知为已读
func (n *Notification) MarkRead(ctx *gin.Context, req *api.MarkReadRequest) (*api.MarkReadResponse, error) {
	err := notification_svc.Notification().MarkRead(ctx, req)
	return &api.MarkReadResponse{}, err
}

// BatchMarkRead 批量标记已读
func (n *Notification) BatchMarkRead(ctx *gin.Context, req *api.BatchMarkReadRequest) (*api.BatchMarkReadResponse, error) {
	count, err := notification_svc.Notification().BatchMarkRead(ctx, req)
	if err != nil {
		return nil, err
	}
	return &api.BatchMarkReadResponse{Count: count}, nil
}
