package issue

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Issue struct {
	user_entity.UserInfo `json:",inline"`
	ID                   int64    `json:"id"`
	ScriptID             int64    `json:"script_id"`
	Title                string   `json:"title"`
	Labels               []string `json:"labels"`
	CommentCount         int64    `json:"comment_count"`
	Status               int32    `json:"status"`
	Createtime           int64    `json:"createtime"`
	Updatetime           int64    `json:"updatetime"`
}

// ListRequest 获取脚本反馈列表
type ListRequest struct {
	mux.Meta              `path:"/scripts/:id/issues" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ScriptID              int64  `uri:"id" binding:"required"`
	Keyword               string `form:"keyword"`
	Status                int32  `form:"status,default=0" binding:"oneof=0 1 3"` // 0:全部 1:待解决 3:已关闭
}

type ListResponse struct {
	httputils.PageResponse[*Issue] `json:",inline"`
}

// CreateIssueRequest 创建脚本反馈
type CreateIssueRequest struct {
	mux.Meta `path:"/scripts/:id/issues" method:"POST"`
	ScriptID int64    `uri:"id" binding:"required"`
	Title    string   `json:"title" binding:"required,max=128" label:"标题"`
	Content  string   `json:"content" binding:"max=10485760" label:"反馈内容"`
	Labels   []string `json:"labels" binding:"max=128" label:"标签"`
}

type CreateIssueResponse struct {
	ID int64 `json:"id"`
}

// GetIssueRequest 获取issue信息
type GetIssueRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId" method:"GET"`
	ScriptID int64 `uri:"id" binding:"required"`
	IssueID  int64 `uri:"issueId" binding:"required"`
}

type GetIssueResponse struct {
	*Issue  `json:",inline"`
	Content string `json:"content"`
}

// GetWatchRequest 获取issue关注状态
type GetWatchRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/watch" method:"GET"`
	ScriptID int64 `uri:"id" binding:"required"`
	IssueID  int64 `uri:"issueId" binding:"required"`
}

type GetWatchResponse struct {
	Watch bool `json:"watch"`
}

// WatchRequest 关注issue
type WatchRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/watch" method:"PUT"`
	ScriptID int64 `uri:"id" binding:"required"`
	IssueID  int64 `uri:"issueId" binding:"required"`
	Watch    bool  `form:"watch" binding:"omitempty" label:"关注状态"`
}

type WatchResponse struct {
}

// OpenRequest 打开/关闭issue
type OpenRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/open" method:"PUT"`
	ScriptID int64  `uri:"id" binding:"required"`
	IssueID  int64  `uri:"issueId" binding:"required"`
	Content  string `json:"content" binding:"max=10485760" label:"评论内容"`
	Close    bool   `json:"close" binding:"omitempty" label:"关闭状态"` // true:关闭 false:打开
}

type OpenResponse struct {
	Comments []*Comment `json:"comments"`
}

// DeleteRequest 删除issue
type DeleteRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId" method:"DELETE"`
	ScriptID int64 `uri:"id" binding:"required"`
	IssueID  int64 `uri:"issueId" binding:"required"`
}

type DeleteResponse struct {
}

// UpdateLabelsRequest 更新issue标签
type UpdateLabelsRequest struct {
	mux.Meta `path:"/scripts/:id/issues/:issueId/labels" method:"PUT"`
	ScriptID int64    `uri:"id" binding:"required"`
	IssueID  int64    `uri:"issueId" binding:"required"`
	Labels   []string `form:"labels" binding:"max=128" label:"标签"`
}

type UpdateLabelsResponse struct {
	*Comment `json:",inline"`
}
