package script_svc

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type AccessInviteSvc interface {
	// InviteCodeList 邀请码列表
	InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error)
	// CreateInviteLink 创建邀请链接
	CreateInviteLink(ctx context.Context, req *api.CreateInviteLinkRequest) (*api.CreateInviteLinkResponse, error)
	// CreateInviteCode 创建邀请码
	CreateInviteCode(ctx context.Context, req *api.CreateInviteCodeRequest) (*api.CreateInviteCodeResponse, error)
	// DeleteInviteCode 删除邀请码
	DeleteInviteCode(ctx context.Context, req *api.DeleteInviteCodeRequest) (*api.DeleteInviteCodeResponse, error)
	// AuditInviteCode 审核邀请码
	AuditInviteCode(ctx context.Context, req *api.AuditInviteCodeRequest) (*api.AuditInviteCodeResponse, error)
	// AcceptInvite 接受邀请
	AcceptInvite(ctx context.Context, req *api.AcceptInviteRequest) (*api.AcceptInviteResponse, error)
	// RejectInvite 拒绝邀请
	RejectInvite(ctx context.Context, req *api.RejectInviteRequest) (*api.RejectInviteResponse, error)
	// GroupInviteCode 群组邀请码列表
	GroupInviteCode(ctx context.Context, req *api.GroupInviteCodeRequest) (*api.GroupInviteCodeResponse, error)
	// CreateGroupInviteCode 创建群组邀请码
	CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error)
}

type accessInviteSvc struct {
}

var defaultAccessInvite = &accessInviteSvc{}

func AccessInvite() AccessInviteSvc {
	return defaultAccessInvite
}

// InviteCodeList 邀请码列表
func (a *accessInviteSvc) InviteCodeList(ctx context.Context, req *api.InviteCodeListRequest) (*api.InviteCodeListResponse, error) {
	return nil, nil
}

// CreateInviteLink 创建邀请链接
func (a *accessInviteSvc) CreateInviteLink(ctx context.Context, req *api.CreateInviteLinkRequest) (*api.CreateInviteLinkResponse, error) {

	return nil, nil
}

// CreateInviteCode 创建邀请码
func (a *accessInviteSvc) CreateInviteCode(ctx context.Context, req *api.CreateInviteCodeRequest) (*api.CreateInviteCodeResponse, error) {
	return nil, nil
}

// DeleteInviteCode 删除邀请码
func (a *accessInviteSvc) DeleteInviteCode(ctx context.Context, req *api.DeleteInviteCodeRequest) (*api.DeleteInviteCodeResponse, error) {
	return nil, nil
}

// AuditInviteCode 审核邀请码
func (a *accessInviteSvc) AuditInviteCode(ctx context.Context, req *api.AuditInviteCodeRequest) (*api.AuditInviteCodeResponse, error) {
	return nil, nil
}

// AcceptInvite 接受邀请
func (a *accessInviteSvc) AcceptInvite(ctx context.Context, req *api.AcceptInviteRequest) (*api.AcceptInviteResponse, error) {
	return nil, nil
}

// RejectInvite 拒绝邀请
func (a *accessInviteSvc) RejectInvite(ctx context.Context, req *api.RejectInviteRequest) (*api.RejectInviteResponse, error) {
	return nil, nil
}

// GroupInviteCode 群组邀请码列表
func (a *accessInviteSvc) GroupInviteCode(ctx context.Context, req *api.GroupInviteCodeRequest) (*api.GroupInviteCodeResponse, error) {
	return nil, nil
}

// CreateGroupInviteCode 创建群组邀请码
func (a *accessInviteSvc) CreateGroupInviteCode(ctx context.Context, req *api.CreateGroupInviteCodeRequest) (*api.CreateGroupInviteCodeResponse, error) {
	return nil, nil
}
