package script

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type IScript interface {
	// List 获取脚本列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
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
