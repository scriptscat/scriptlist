package script_svc

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type AccessSvc interface {
	// AccessList 访问控制列表
	AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error)
	// CreateAccess 创建访问控制
	CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error)
	// UpdateAccess 更新访问控制
	UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error)
	// DeleteAccess 删除访问控制
	DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error)
}

type accessSvc struct {
}

var defaultAccess = &accessSvc{}

func Access() AccessSvc {
	return defaultAccess
}

// AccessList 访问控制列表
func (a *accessSvc) AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error) {
	return nil, nil
}

// CreateAccess 创建访问控制
func (a *accessSvc) CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error) {
	return nil, nil
}

// UpdateAccess 更新访问控制
func (a *accessSvc) UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error) {
	return nil, nil
}

// DeleteAccess 删除访问控制
func (a *accessSvc) DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error) {
	return nil, nil
}
