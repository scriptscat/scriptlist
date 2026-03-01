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

type ReportComment struct {
	limit *limit.PeriodLimit
}

func NewReportComment() *ReportComment {
	return &ReportComment{
		limit: limit.NewPeriodLimit(
			300, 5, redis.Default(), "limit:create:report_comment",
		),
	}
}

func (c *ReportComment) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(false),
			script_svc.Script().RequireScript(),
			report_svc.Report().RequireReport(),
		},
		Handler: []interface{}{
			c.ListComment,
		},
	}, {
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Script().RequireScript(),
			report_svc.Report().RequireReport(),
		},
		Handler: []interface{}{
			c.CreateComment,
			muxutils.Use(script_svc.Access().CheckHandler("report", "delete")).Append(
				c.DeleteComment,
			),
		},
	}})
}

// ListComment 获取举报评论列表
func (c *ReportComment) ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error) {
	return report_svc.ReportComment().ListComment(ctx, req)
}

// CreateComment 创建举报评论
func (c *ReportComment) CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error) {
	resp, err := c.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return report_svc.ReportComment().CreateComment(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateCommentResponse), nil
}

// DeleteComment 删除举报评论
func (c *ReportComment) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error) {
	return report_svc.ReportComment().DeleteComment(ctx, req)
}
