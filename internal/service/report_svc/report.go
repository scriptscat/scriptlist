package report_svc

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/report"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/report_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

type contextKey int

const (
	reportCtxKey contextKey = iota
)

type ReportSvc interface {
	// List 获取脚本举报列表
	List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error)
	// CreateReport 创建脚本举报
	CreateReport(ctx context.Context, req *api.CreateReportRequest) (*api.CreateReportResponse, error)
	// GetReport 获取举报详情
	GetReport(ctx context.Context, req *api.GetReportRequest) (*api.GetReportResponse, error)
	// Resolve 解决/重新打开举报
	Resolve(ctx context.Context, req *api.ResolveRequest) (*api.ResolveResponse, error)
	// Delete 删除举报
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)
	// RequireReport 需要举报存在
	RequireReport() gin.HandlerFunc
	// CtxReport 获取举报
	CtxReport(ctx context.Context) *report_entity.ScriptReport
}

type reportSvc struct{}

var defaultReport = &reportSvc{}

func Report() ReportSvc {
	return defaultReport
}

// List 获取脚本举报列表
func (r *reportSvc) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	list, total, err := report_repo.Report().FindPage(ctx, req.ScriptID, req.Status, req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := &api.ListResponse{
		PageResponse: httputils.PageResponse[*api.Report]{
			Total: total,
			List:  make([]*api.Report, len(list)),
		},
	}
	for n, v := range list {
		item, err := r.toReport(ctx, v)
		if err != nil {
			return nil, err
		}
		resp.List[n] = item
	}
	return resp, nil
}

func (r *reportSvc) toReport(ctx context.Context, report *report_entity.ScriptReport) (*api.Report, error) {
	ret := &api.Report{
		ID:         report.ID,
		ScriptID:   report.ScriptID,
		Reason:     report.Reason,
		Status:     report.Status,
		Createtime: report.Createtime,
		Updatetime: report.Updatetime,
	}
	user, err := user_repo.User().Find(ctx, report.UserID)
	if err != nil {
		return nil, err
	}
	ret.UserInfo = user.UserInfo()
	commentCount, err := report_repo.Comment().CountByReport(ctx, report.ID)
	if err != nil {
		return nil, err
	}
	ret.CommentCount = commentCount
	return ret, nil
}

// CreateReport 创建脚本举报
func (r *reportSvc) CreateReport(ctx context.Context, req *api.CreateReportRequest) (*api.CreateReportResponse, error) {
	uid := auth_svc.Auth().Get(ctx).UID
	// 不允许举报自己的脚本
	if uid == script_svc.Script().CtxScript(ctx).UserID {
		return nil, i18n.NewError(ctx, code.ReportSelfReport)
	}
	report := &report_entity.ScriptReport{
		ScriptID:   req.ScriptID,
		UserID:     uid,
		Reason:     req.Reason,
		Content:    req.Content,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := report.ValidateReason(ctx); err != nil {
		return nil, err
	}
	if err := report_repo.Report().Create(ctx, report); err != nil {
		return nil, err
	}
	return &api.CreateReportResponse{ID: report.ID}, producer.PublishReportCreate(ctx, script_svc.Script().CtxScript(ctx), report)
}

// GetReport 获取举报详情
func (r *reportSvc) GetReport(ctx context.Context, req *api.GetReportRequest) (*api.GetReportResponse, error) {
	report := r.CtxReport(ctx)
	ret, err := r.toReport(ctx, report)
	if err != nil {
		return nil, err
	}
	return &api.GetReportResponse{
		Report:  ret,
		Content: report.Content,
	}, nil
}

// Resolve 解决/重新打开举报
func (r *reportSvc) Resolve(ctx context.Context, req *api.ResolveRequest) (*api.ResolveResponse, error) {
	report := r.CtxReport(ctx)

	// 幂等检查：已解决的不能再次解决，已打开的不能再次打开
	if req.Close && report.IsResolved() {
		return nil, i18n.NewError(ctx, code.ReportAlreadyResolved)
	}
	if !req.Close && !report.IsResolved() {
		return nil, i18n.NewError(ctx, code.ReportAlreadyResolved)
	}

	comments := make([]*api.Comment, 0)

	var commentType report_entity.CommentType
	var commentContent string

	if req.Close {
		report.Resolve(time.Now())
		commentContent = "解决举报"
		commentType = report_entity.CommentTypeResolve
	} else {
		report.Reopen(time.Now())
		commentContent = "重新打开举报"
		commentType = report_entity.CommentTypeReopen
	}

	if err := report_repo.Report().Update(ctx, report); err != nil {
		return nil, err
	}

	// 管理员附带的评论内容，直接创建评论记录而不通过 CreateComment，避免发送两条通知
	if req.Content != "" {
		userComment := &report_entity.ScriptReportComment{
			ReportID:   req.ReportID,
			UserID:     auth_svc.Auth().Get(ctx).UID,
			Content:    req.Content,
			Type:       report_entity.CommentTypeComment,
			Status:     consts.ACTIVE,
			Createtime: time.Now().Unix(),
		}
		if err := report_repo.Comment().Create(ctx, userComment); err != nil {
			return nil, err
		}
		resp, _ := ReportComment().toComment(ctx, userComment)
		if resp != nil {
			comments = append(comments, resp)
		}
	}

	comment := &report_entity.ScriptReportComment{
		ReportID:   req.ReportID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    commentContent,
		Type:       commentType,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := report_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	resp, _ := ReportComment().toComment(ctx, comment)
	if resp != nil {
		comments = append(comments, resp)
	}
	return &api.ResolveResponse{
		Comments: comments,
	}, producer.PublishReportCommentCreate(ctx, script_svc.Script().CtxScript(ctx), report, comment)
}

// Delete 删除举报
func (r *reportSvc) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	report := r.CtxReport(ctx)
	if err := report_repo.Report().Delete(ctx, script_svc.Script().CtxScript(ctx).ID, report.ID); err != nil {
		return nil, err
	}
	return &api.DeleteResponse{}, nil
}

func (r *reportSvc) RequireReport() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sReportId := ctx.Param("reportId")
		if sReportId == "" {
			httputils.HandleResp(ctx, httputils.NewError(http.StatusNotFound, -1, "举报ID不能为空"))
			return
		}
		reportId, err := strconv.ParseInt(sReportId, 10, 64)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		script := script_svc.Script().CtxScript(ctx)
		report, err := report_repo.Report().Find(ctx, script.ID, reportId)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		if err := report.CheckOperate(ctx); err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(
			ctx.Request.Context(), reportCtxKey, report,
		))
	}
}

func (r *reportSvc) CtxReport(ctx context.Context) *report_entity.ScriptReport {
	return ctx.Value(reportCtxKey).(*report_entity.ScriptReport)
}
