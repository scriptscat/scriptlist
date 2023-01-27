package issue_ctr

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/limit"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/issue_svc"
)

type Comment struct {
	limit *limit.PeriodLimit
}

func NewComment() *Comment {
	return &Comment{
		limit: limit.NewPeriodLimit(
			300, 10, redis.Default(), "limit:create:issue",
		),
	}
}

// ListComment 获取反馈评论列表
func (c *Comment) ListComment(ctx context.Context, req *api.ListCommentRequest) (*api.ListCommentResponse, error) {
	return issue_svc.Comment().ListComment(ctx, req)
}

// CreateComment 创建反馈评论
func (c *Comment) CreateComment(ctx context.Context, req *api.CreateCommentRequest) (*api.CreateCommentResponse, error) {
	resp, err := c.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return issue_svc.Comment().CreateComment(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateCommentResponse), nil
}

func (c *Comment) Middleware() gin.HandlerFunc {
	return issue_svc.Comment().Middleware()
}

// DeleteComment 删除反馈评论
func (c *Comment) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.DeleteCommentResponse, error) {
	return issue_svc.Comment().DeleteComment(ctx, req)
}
