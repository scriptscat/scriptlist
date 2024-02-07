package issue_svc

import (
	"context"
	"time"

	"github.com/scriptscat/scriptlist/internal/service/script_svc"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

type CommentSvc interface {
	// ListComment 获取反馈评论列表
	ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error)
	// CreateComment 创建反馈评论
	CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error)
	// ToComment 转换为api.Comment
	ToComment(ctx context.Context, comment *issue_entity.ScriptIssueComment) (*api.Comment, error)
	// DeleteComment 删除反馈评论
	DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error)
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
	return resp, producer.PublishCommentCreate(ctx, script_svc.Script().CtxScript(ctx), Issue().CtxIssue(ctx), comment)
}

// DeleteComment 删除反馈评论
func (c *commentSvc) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error) {
	comment, err := issue_repo.Comment().Find(ctx, req.IssueID, req.CommentID)
	if err != nil {
		return nil, err
	}
	if err := comment.CheckOperate(ctx); err != nil {
		return nil, err
	}
	return &api.DeleteCommentResponse{}, issue_repo.Comment().Delete(ctx, req.CommentID)
}
