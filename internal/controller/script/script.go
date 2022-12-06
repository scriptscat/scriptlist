package script

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/limit"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	service "github.com/scriptscat/scriptlist/internal/service/script"
	"github.com/scriptscat/scriptlist/internal/service/user"
)

type Script struct {
	limit *limit.PeriodLimit
}

func NewScript() *Script {
	return &Script{
		limit: limit.NewPeriodLimit(
			300, 10, redis.Default(), "limit:create:script",
		),
	}
}

// List 获取脚本列表
func (s *Script) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return service.Script().List(ctx, req)
}

// Create 创建脚本
func (s *Script) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	cancel, err := s.limit.Take(ctx, strconv.FormatInt(user.Auth().Get(ctx).UID, 10))
	if err != nil {
		return nil, err
	}
	resp, err := service.Script().Create(ctx, req)
	if err != nil {
		if err := cancel(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return resp, err
}
