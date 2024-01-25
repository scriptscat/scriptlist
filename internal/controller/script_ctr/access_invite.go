package script_ctr

import (
	"context"
	"github.com/codfrm/cago/pkg/utils/muxutils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"

	"github.com/codfrm/cago/server/mux"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type AccessInvite struct {
}

func NewAccessInvite() *AccessInvite {
	return &AccessInvite{}
}

func (a *AccessInvite) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
		},
		Handler: []interface{}{
			a.AcceptInvite,
			a.InviteCodeInfo,
		},
	}, {
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Script().RequireScript(),
		},
		Handler: []interface{}{
			muxutils.Use(script_svc.Access().CheckHandler("access", "read")).Append(
				a.InviteCodeList,
				muxutils.Use(script_svc.Group().RequireGroup()).Append(a.GroupInviteCodeList),
			),
			muxutils.Use(script_svc.Access().CheckHandler("access", "manage")).Append(
				a.CreateInviteCode,
				muxutils.Use(script_svc.Group().RequireGroup()).Append(a.CreateGroupInviteCode),
				a.DeleteInviteCode,
				a.AuditInviteCode,
			),
		},
	}})
}

// InviteCodeList 邀请码列表
func (a *AccessInvite) InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error) {
	return script_svc.AccessInvite().InviteCodeList(ctx, req)
}

// CreateInviteCode 创建邀请码
func (a *AccessInvite) CreateInviteCode(ctx context.Context, req *api.CreateInviteCodeRequest) (*api.CreateInviteCodeResponse, error) {
	return script_svc.AccessInvite().CreateInviteCode(ctx, req)
}

// DeleteInviteCode 删除邀请码
func (a *AccessInvite) DeleteInviteCode(ctx context.Context, req *api.DeleteInviteCodeRequest) (*api.DeleteInviteCodeResponse, error) {
	return script_svc.AccessInvite().DeleteInviteCode(ctx, req)
}

// AuditInviteCode 审核邀请码
func (a *AccessInvite) AuditInviteCode(ctx context.Context, req *api.AuditInviteCodeRequest) (*api.AuditInviteCodeResponse, error) {
	return script_svc.AccessInvite().AuditInviteCode(ctx, req)
}

// AcceptInvite 接受邀请
func (a *AccessInvite) AcceptInvite(ctx context.Context, req *api.AcceptInviteRequest) (*api.AcceptInviteResponse, error) {
	return script_svc.AccessInvite().AcceptInvite(ctx, req)
}

// GroupInviteCodeList 群组邀请码列表
func (a *AccessInvite) GroupInviteCodeList(ctx context.Context, req *api.GroupInviteCodeListRequest) (*api.GroupInviteCodeListResponse, error) {
	return script_svc.AccessInvite().GroupInviteCodeList(ctx, req)
}

// CreateGroupInviteCode 创建群组邀请码
func (a *AccessInvite) CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error) {
	return script_svc.AccessInvite().CreateGroupInviteCode(ctx, req)
}

// InviteCodeInfo 邀请码信息
func (a *AccessInvite) InviteCodeInfo(ctx context.Context, req *api.InviteCodeInfoRequest) (*api.InviteCodeInfoResponse, error) {
	return script_svc.AccessInvite().InviteCodeInfo(ctx, req)
}
