package script_ctr

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Access struct {
}

func NewAccess() *Access {
	return &Access{}
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
