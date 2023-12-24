package script_ctr

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/muxutils"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Access struct {
}

func NewAccess() *Access {
	return &Access{}
}

func (a *Access) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Script().RequireScript(),
		},
		Handler: []interface{}{
			muxutils.Use(script_svc.Access().CheckHandler("access", "read")).Append(
				a.AccessList,
			),
			muxutils.Use(script_svc.Access().CheckHandler("access", "manage")).Append(
				a.CreateAccess,
				muxutils.Use(script_svc.Access().RequireAccess()).Append(
					a.UpdateAccess,
					a.DeleteAccess,
				),
			),
		},
	}})
}

// AccessList 访问控制列表
func (a *Access) AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error) {
	return script_svc.Access().AccessList(ctx, req)
}

// CreateAccess 创建访问控制
func (a *Access) CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error) {
	return script_svc.Access().CreateAccess(ctx, req)
}

// UpdateAccess 更新访问控制
func (a *Access) UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error) {
	return script_svc.Access().UpdateAccess(ctx, req)
}

// DeleteAccess 删除访问控制
func (a *Access) DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error) {
	return script_svc.Access().DeleteAccess(ctx, req)
}
