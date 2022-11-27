package script

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/httputils"
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
	return &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Item]{
			List: []*api.Item{
				{
					ID: 1,
				},
			}, Total: 0,
		},
	}, nil
	return service.Script().List(ctx, req)
}
