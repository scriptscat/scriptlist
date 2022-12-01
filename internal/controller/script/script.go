package script

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	service "github.com/scriptscat/scriptlist/internal/service/script"
)

type Script struct {
}

func NewScript() *Script {
	return &Script{}
}

// List 获取脚本列表
func (s *Script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return service.Script().List(ctx, req)
}

// Create 创建脚本
func (s *Script) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return service.Script().Create(ctx, req)
}
