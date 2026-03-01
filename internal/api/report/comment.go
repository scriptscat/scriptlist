package report

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Comment struct {
	user_entity.UserInfo `json:",inline"`
	ID                   int64                     `json:"id"`
	ReportID             int64                     `json:"report_id"`
	Content              string                    `json:"content"`
	Type                 report_entity.CommentType `json:"type"`
	Status               int32                     `json:"status"`
	Createtime           int64                     `json:"createtime"`
	Updatetime           int64                     `json:"updatetime"`
}

// ListCommentRequest 获取举报评论列表
type ListCommentRequest struct {
	mux.Meta `path:"/scripts/:id/reports/:reportId/comments" method:"GET"`
	ScriptID int64 `uri:"id" binding:"required"`
	ReportID int64 `uri:"reportId" binding:"required"`
}

type ListCommentResponse struct {
	httputils.PageResponse[*Comment] `json:",inline"`
}

// CreateCommentRequest 创建举报评论
type CreateCommentRequest struct {
	mux.Meta `path:"/scripts/:id/reports/:reportId/comments" method:"POST"`
	ScriptID int64  `uri:"id" binding:"required"`
	ReportID int64  `uri:"reportId" binding:"required"`
	Content  string `json:"content" binding:"required,max=10485760" label:"评论内容"`
}

type CreateCommentResponse struct {
	*Comment `json:",inline"`
}

// DeleteCommentRequest 删除举报评论
type DeleteCommentRequest struct {
	mux.Meta  `path:"/scripts/:id/reports/:reportId/comments/:commentId" method:"DELETE"`
	ScriptID  int64 `uri:"id" binding:"required"`
	ReportID  int64 `uri:"reportId" binding:"required"`
	CommentID int64 `uri:"commentId" binding:"required"`
}

type DeleteCommentResponse struct {
}
