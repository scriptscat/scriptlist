package script_ctr

import (
	"context"
	"strconv"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/limit"
	"github.com/cago-frame/cago/pkg/utils/muxutils"
	"github.com/cago-frame/cago/server/mux"
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
				a.AddGroupAccess,
				a.AddUserAccess,
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

// UpdateAccess 更新访问控制
func (a *Access) UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error) {
	return script_svc.Access().UpdateAccess(ctx, req)
}

// DeleteAccess 删除访问控制
func (a *Access) DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error) {
	return script_svc.Access().DeleteAccess(ctx, req)
}

// AddGroupAccess 添加组权限
func (a *Access) AddGroupAccess(ctx context.Context, req *api.AddGroupAccessRequest) (*api.AddGroupAccessResponse, error) {
	return script_svc.Access().AddGroupAccess(ctx, req)
}

// AddUserAccess 添加用户权限, 通过用户名进行邀请
func (a *Access) AddUserAccess(ctx context.Context, req *api.AddUserAccessRequest) (*api.AddUserAccessResponse, error) {
	ret, err := limit.NewPeriodLimit(1, 1, redis.Default(), "script:access:").
		FuncTake(ctx, strconv.FormatInt(req.ScriptID, 10), func() (interface{}, error) {
			return script_svc.Access().AddUserAccess(ctx, req)
		})
	if err != nil {
		return nil, err
	}
	return ret.(*api.AddUserAccessResponse), nil
}
