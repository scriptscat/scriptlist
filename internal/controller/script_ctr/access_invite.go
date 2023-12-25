package script_ctr

import (
	"context"

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

}

// InviteCodeList 邀请码列表
func (a *AccessInvite) InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error) {
	return script_svc.AccessInvite().InviteCodeList(ctx, req)
}

// CreateInviteLink 创建邀请链接
func (a *AccessInvite) CreateInviteLink(ctx context.Context, req *api.CreateInviteLinkRequest) (*api.CreateInviteLinkResponse, error) {
	return script_svc.AccessInvite().CreateInviteLink(ctx, req)
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

// RejectInvite 拒绝邀请
func (a *AccessInvite) RejectInvite(ctx context.Context, req *api.RejectInviteRequest) (*api.RejectInviteResponse, error) {
	return script_svc.AccessInvite().RejectInvite(ctx, req)
}

// GroupInviteCode 群组邀请码列表
func (a *AccessInvite) GroupInviteCode(ctx context.Context, req *api.GroupInviteCodeRequest) (*api.GroupInviteCodeResponse, error) {
	return script_svc.AccessInvite().GroupInviteCode(ctx, req)
}

// CreateGroupInviteCode 创建群组邀请码
func (a *AccessInvite) CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error) {
	return script_svc.AccessInvite().CreateGroupInviteCode(ctx, req)
}
