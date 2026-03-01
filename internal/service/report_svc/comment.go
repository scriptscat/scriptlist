package report_svc

import (
	"context"
	"time"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/report"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/repository/report_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

type ReportCommentSvc interface {
	// ListComment 获取举报评论列表
	ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error)
	// CreateComment 创建举报评论
	CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error)
	// DeleteComment 删除举报评论
	DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error)
	// toComment 转换为api.Comment
	toComment(ctx context.Context, comment *report_entity.ScriptReportComment) (*api.Comment, error)
}

type reportCommentSvc struct{}

var defaultReportComment = &reportCommentSvc{}

func ReportComment() ReportCommentSvc {
	return defaultReportComment
}

// ListComment 获取举报评论列表
func (c *reportCommentSvc) ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error) {
	list, err := report_repo.Comment().FindAll(ctx, req.ReportID)
	if err != nil {
		return nil, err
	}
	ret := &api.ListCommentResponse{
		PageResponse: httputils.PageResponse[*api.Comment]{
			List:  make([]*api.Comment, len(list)),
			Total: int64(len(list)),
		},
	}
	for n, v := range list {
		item, err := c.toComment(ctx, v)
		if err != nil {
			return nil, err
		}
		ret.List[n] = item
	}
	return ret, nil
}

func (c *reportCommentSvc) toComment(ctx context.Context, comment *report_entity.ScriptReportComment) (*api.Comment, error) {
	ret := &api.Comment{
		ID:         comment.ID,
		ReportID:   comment.ReportID,
		Content:    comment.Content,
		Type:       comment.Type,
		Status:     comment.Status,
		Createtime: comment.Createtime,
		Updatetime: comment.Updatetime,
	}
	user, err := user_repo.User().Find(ctx, comment.UserID)
	if err != nil {
		return nil, err
	}
	ret.UserInfo = user.UserInfo()
	return ret, nil
}

// CreateComment 创建举报评论
func (c *reportCommentSvc) CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error) {
	comment := &report_entity.ScriptReportComment{
		ReportID:   req.ReportID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    req.Content,
		Type:       report_entity.CommentTypeComment,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := report_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	resp := &api.CreateCommentResponse{}
	apiComment, err := c.toComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	resp.Comment = apiComment
	return resp, producer.PublishReportCommentCreate(ctx, script_svc.Script().CtxScript(ctx), Report().CtxReport(ctx), comment)
}

// DeleteComment 删除举报评论
func (c *reportCommentSvc) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error) {
	comment, err := report_repo.Comment().Find(ctx, req.ReportID, req.CommentID)
	if err != nil {
		return nil, err
	}
	if err := comment.CheckOperate(ctx); err != nil {
		return nil, err
	}
	return &api.DeleteCommentResponse{}, report_repo.Comment().Delete(ctx, req.CommentID)
}
