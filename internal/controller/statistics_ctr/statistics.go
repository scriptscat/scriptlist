package statistics_ctr

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/service/statistics_svc"
)

type Statistics struct {
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

// Script 脚本统计数据
func (s *Statistics) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	return statistics_svc.Statistics().Script(ctx, req)
}

// ScriptRealtime 脚本实时统计数据
func (s *Statistics) ScriptRealtime(ctx context.Context, req *api.ScriptRealtimeRequest) (*api.ScriptRealtimeResponse, error) {
	return statistics_svc.Statistics().ScriptRealtime(ctx, req)
}

// Collect 统计数据收集
func (s *Statistics) Collect(ctx *gin.Context, req *api.CollectRequest) (*api.CollectResponse, error) {
	req.UA = ctx.Request.UserAgent()
	if req.UA == "" {
		return nil, errors.New("ua is empty")
	}
	req.IP = ctx.ClientIP()
	return statistics_svc.Statistics().Collect(ctx.Request.Context(), req)
}
