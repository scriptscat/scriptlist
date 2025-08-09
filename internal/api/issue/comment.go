package issue

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Comment struct {
	user_entity.UserInfo `json:",inline"`
	ID                   int64                    `json:"id"`
	IssueID              int64                    `json:"issue_id"`
	Content              string                   `json:"content"`
	Type                 issue_entity.CommentType `json:"type"`
	Status               int32                    `json:"status"`
	Createtime           int64                    `json:"createtime"`
	Updatetime           int64                    `json:"updatetime"`
}

// ListCommentRequest 获取反馈评论列表
type ListCommentRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/comment" method:"GET"`
	ScriptID int64 `uri:"id" binding:"required"`
	IssueID  int64 `uri:"issueId" binding:"required"`
}

type ListCommentResponse struct {
	httputils.PageResponse[*Comment] `json:",inline"`
}

// CreateCommentRequest 创建反馈评论
type CreateCommentRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/comment" method:"POST"`
	ScriptID int64  `uri:"id" binding:"required"`
	IssueID  int64  `uri:"issueId" binding:"required"`
	Content  string `json:"content" binding:"required,max=10485760" label:"评论内容"`
}

type CreateCommentResponse struct {
	*Comment `json:",inline"`
}

// DeleteCommentRequest 删除反馈评论
type DeleteCommentRequest struct {
	mux.Meta  `path:"/scripts/:id/issues/:issueId/comment/:commentId" method:"DELETE"`
	ScriptID  int64 `uri:"id" binding:"required"`
	IssueID   int64 `uri:"issueId" binding:"required"`
	CommentID int64 `uri:"commentId" binding:"required"`
}

type DeleteCommentResponse struct {
}
