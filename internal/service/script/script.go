package script

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type IScript interface {
	// List 获取脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// Create 创建脚本
	Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error)
}

type script struct {
}

var defaultScript = &script{}

func Script() IScript {
	return defaultScript
}

// List 获取脚本列表
func (s *script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return nil, nil
}

// Create 创建脚本
func (s *script) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	return nil, nil
}
