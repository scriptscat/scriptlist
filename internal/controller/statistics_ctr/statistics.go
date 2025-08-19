package statistics_ctr

import (
	"context"

	"github.com/cago-frame/cago/pkg/utils/muxutils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"github.com/scriptscat/scriptlist/internal/service/statistics_svc"
)

type Statistics struct {
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

func (s *Statistics) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{muxutils.
		Use(
			auth_svc.Auth().RequireLogin(true),
			script_svc.Script().RequireScript(script_svc.WithRequireScriptAccess("statistics", "manage")),
		).Append(
		s.Script,
		s.ScriptRealtime,
	)})
}

// Script 脚本统计数据
func (s *Statistics) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	return statistics_svc.Statistics().Script(ctx, req)
}

// ScriptRealtime 脚本实时统计数据
func (s *Statistics) ScriptRealtime(ctx context.Context, req *api.ScriptRealtimeRequest) (*api.ScriptRealtimeResponse, error) {
	return statistics_svc.Statistics().ScriptRealtime(ctx, req)
}

func (s *Statistics) Middleware() gin.HandlerFunc {
	return statistics_svc.Statistics().Middleware()
}
