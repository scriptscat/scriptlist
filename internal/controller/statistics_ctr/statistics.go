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

// RealtimeChart 实时统计数据图表
func (s *Statistics) RealtimeChart(ctx context.Context, req *api.RealtimeChartRequest) (*api.RealtimeChartResponse, error) {
	return statistics_svc.Statistics().RealtimeChart(ctx, req)
}

// VisitList 访问列表
func (s *Statistics) VisitList(ctx context.Context, req *api.VisitListRequest) (*api.VisitResponse, error) {
	req.PageRequest.Size = 10
	return statistics_svc.Statistics().VisitList(ctx, req)
}

// AdvancedInfo 高级统计信息
func (s *Statistics) AdvancedInfo(ctx context.Context, req *api.AdvancedInfoRequest) (*api.AdvancedInfoResponse, error) {
	return statistics_svc.Statistics().AdvancedInfo(ctx, req)
}

// UserOrigin 用户来源统计
func (s *Statistics) UserOrigin(ctx context.Context, req *api.UserOriginRequest) (*api.UserOriginResponse, error) {
	req.PageRequest.Size = 10
	return statistics_svc.Statistics().UserOrigin(ctx, req)
}

func (s *Statistics) Middleware() gin.HandlerFunc {
	return statistics_svc.Statistics().Middleware()
}

// VisitDomain 访问域名统计
func (s *Statistics) VisitDomain(ctx context.Context, req *api.VisitDomainRequest) (*api.VisitDomainResponse, error) {
	req.PageRequest.Size = 10
	return statistics_svc.Statistics().VisitDomain(ctx, req)
}

// UpdateWhitelist 更新统计白名单
func (s *Statistics) UpdateWhitelist(ctx context.Context, req *api.UpdateWhitelistRequest) (*api.UpdateWhitelistResponse, error) {
	return statistics_svc.Statistics().UpdateWhitelist(ctx, req)
}

// CollectWhitelist 获取统计收集白名单
func (s *Statistics) CollectWhitelist(ctx context.Context, req *api.CollectWhitelistRequest) (*api.CollectWhitelistResponse, error) {
	return statistics_svc.Statistics().CollectWhitelist(ctx, req)
}
