package issue_svc

import (
	"context"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

type ctxScript string

type CommentSvc interface {
	// ListComment 获取反馈评论列表
	ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error)
	// CreateComment 创建反馈评论
	CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error)
	// Middleware 中间件
	Middleware() gin.HandlerFunc
	// CheckOperate 检查脚本和issue状态是否正确
	CheckOperate(ctx context.Context, scriptId, issueId int64) (*script_entity.Script, *issue_entity.ScriptIssue, error)
	// ToComment 转换为api.Comment
	ToComment(ctx context.Context, comment *issue_entity.ScriptIssueComment) (*api.Comment, error)
	// DeleteComment 删除反馈评论
	DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error)
	// RequireComment 需要存在评论
	RequireComment() gin.HandlerFunc
	// CtxComment 获取评论
	CtxComment(ctx context.Context) *issue_entity.ScriptIssueComment
}

type commentSvc struct {
}

var defaultComment = &commentSvc{}

func Comment() CommentSvc {
	return defaultComment
}

// ListComment 获取反馈评论列表
func (c *commentSvc) ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error) {
	list, err := issue_repo.Comment().FindAll(ctx, req.IssueID)
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
		ret.List[n], _ = c.ToComment(ctx, v)
	}
	return ret, nil
}

func (c *commentSvc) ToComment(ctx context.Context, comment *issue_entity.ScriptIssueComment) (*api.Comment, error) {
	ret := &api.Comment{
		ID:         comment.ID,
		IssueID:    comment.IssueID,
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

func (c *commentSvc) CheckOperate(ctx context.Context, scriptId, issueId int64) (*script_entity.Script, *issue_entity.ScriptIssue, error) {
	script, err := script_repo.Script().Find(ctx, scriptId)
	if err != nil {
		return nil, nil, err
	}
	issue, err := issue_repo.Issue().Find(ctx, scriptId, issueId)
	if err != nil {
		return nil, nil, err
	}
	return script, issue, issue.CheckOperate(ctx, script)
}

// CreateComment 创建反馈评论
func (c *commentSvc) CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error) {
	comment := &issue_entity.ScriptIssueComment{
		IssueID:    req.IssueID,
		UserID:     auth_svc.Auth().Get(ctx).UID,
		Content:    req.Content,
		Type:       issue_entity.CommentTypeComment,
		Status:     consts.ACTIVE,
		Createtime: time.Now().Unix(),
	}
	if err := issue_repo.Comment().Create(ctx, comment); err != nil {
		return nil, err
	}
	resp := &api.CreateCommentResponse{}
	resp.Comment, _ = c.ToComment(ctx, comment)
	return resp, producer.PublishCommentCreate(ctx, script_svc.Script().CtxScript(ctx), c.CtxIssue(ctx), comment)
}

// Middleware 中间件,检查是否可以访问
func (c *commentSvc) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var script *script_entity.Script
		var issue *issue_entity.ScriptIssue
		var err error
		id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
		// 非GET请求,需要验证邮箱,判断是否归档
		if ctx.Request.Method != http.MethodGet {
			if !auth_svc.Auth().Get(ctx).EmailVerified {
				httputils.HandleResp(ctx, i18n.NewErrorWithStatus(ctx, http.StatusForbidden, code.UserEmailNotVerified))
				return
			}
			script, err = script_repo.Script().Find(ctx, id)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
			err = script.IsArchive(ctx)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
		}
		issueId, _ := strconv.ParseInt(ctx.Param("issueId"), 10, 64)
		if issueId != 0 {
			script, issue, err = c.CheckOperate(ctx, id, issueId)
			if err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(context.WithValue(
			ctx.Request.Context(),
			issue_entity.ScriptIssue{}, issue),
			ctxScript("ctxScript"), script))
		ctx.Next()
	}
}

// CtxIssue 获取issue
func (c *commentSvc) CtxIssue(ctx context.Context) *issue_entity.ScriptIssue {
	return ctx.Value(issue_entity.ScriptIssue{}).(*issue_entity.ScriptIssue)
}

// DeleteComment 删除反馈评论
func (c *commentSvc) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error) {
	comment, err := issue_repo.Comment().Find(ctx, req.IssueID, req.CommentID)
	if err != nil {
		return nil, err
	}
	if err := comment.CheckPermission(ctx, script_svc.Script().CtxScript(ctx), c.CtxIssue(ctx)); err != nil {
		return nil, err
	}
	return &api.DeleteCommentResponse{}, issue_repo.Comment().Delete(ctx, req.CommentID)
}

func (c *commentSvc) RequireComment() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (c *commentSvc) CtxComment(ctx context.Context) *issue_entity.ScriptIssueComment {
	//TODO implement me
	panic("implement me")
}
