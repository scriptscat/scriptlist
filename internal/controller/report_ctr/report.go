package report_ctr

import (
	"context"
	"strconv"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/limit"
	"github.com/cago-frame/cago/pkg/utils/muxutils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/report"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/report_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Report struct {
	limit limit.Limit
}

func NewReport() *Report {
	return &Report{
		limit: limit.NewPeriodLimit(
			300, 3, redis.Default(), "limit:create:report",
		),
	}
}

func (r *Report) Router(root *mux.Router) {
	muxutils.BindTree(root, []*muxutils.RouterTree{{
		// 公开路由（无需登录）
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(false),
			script_svc.Script().RequireScript(),
		},
		Handler: []interface{}{
			r.List,
			muxutils.Use(report_svc.Report().RequireReport()).Append(r.GetReport),
		},
	}, {
		// 需要登录
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Script().RequireScript(),
		},
		Handler: []interface{}{
			r.CreateReport,
			&muxutils.RouterTree{
				Middleware: []gin.HandlerFunc{
					report_svc.Report().RequireReport(),
				},
				Handler: []interface{}{
					// 管理员才能解决举报
					muxutils.Use(
						script_svc.Access().CheckHandler("report", "manage"),
					).Append(
						r.Resolve,
					),
					// 管理员才能删除举报
					muxutils.Use(
						script_svc.Access().CheckHandler("report", "delete"),
					).Append(
						r.Delete,
					),
				},
			},
		},
	}})
}

// List 获取脚本举报列表
func (r *Report) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	return report_svc.Report().List(ctx, req)
}

// CreateReport 创建脚本举报
func (r *Report) CreateReport(ctx context.Context, req *api.CreateReportRequest) (*api.CreateReportResponse, error) {
	resp, err := r.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return report_svc.Report().CreateReport(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateReportResponse), nil
}

// GetReport 获取举报详情
func (r *Report) GetReport(ctx context.Context, req *api.GetReportRequest) (*api.GetReportResponse, error) {
	return report_svc.Report().GetReport(ctx, req)
}

// Resolve 解决/重新打开举报
func (r *Report) Resolve(ctx context.Context, req *api.ResolveRequest) (*api.ResolveResponse, error) {
	return report_svc.Report().Resolve(ctx, req)
}

// Delete 删除举报
func (r *Report) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return report_svc.Report().Delete(ctx, req)
}
