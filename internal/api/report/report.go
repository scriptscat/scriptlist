package report

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

type Report struct {
	user_entity.UserInfo `json:",inline"`
	ID                   int64  `json:"id"`
	ScriptID             int64  `json:"script_id"`
	Reason               string `json:"reason"`
	CommentCount         int64  `json:"comment_count"`
	Status               int32  `json:"status"`
	Createtime           int64  `json:"createtime"`
	Updatetime           int64  `json:"updatetime"`
}

// ListRequest 获取脚本举报列表
type ListRequest struct {
	mux.Meta              `path:"/scripts/:id/reports" method:"GET"`
	httputils.PageRequest `form:",inline"`
	ScriptID              int64 `uri:"id" binding:"required"`
	Status                int32 `form:"status,default=0" binding:"oneof=0 1 3"` // 0:全部 1:待处理 3:已解决
}

type ListResponse struct {
	httputils.PageResponse[*Report] `json:",inline"`
}

// CreateReportRequest 创建脚本举报
type CreateReportRequest struct {
	mux.Meta `path:"/scripts/:id/reports" method:"POST"`
	ScriptID int64  `uri:"id" binding:"required"`
	Reason   string `json:"reason" binding:"required,max=64" label:"举报原因"`
	Content  string `json:"content" binding:"required,max=10485760" label:"举报内容"`
}

type CreateReportResponse struct {
	ID int64 `json:"id"`
}

// GetReportRequest 获取举报详情
type GetReportRequest struct {
	mux.Meta `path:"/scripts/:id/reports/:reportId" method:"GET"`
	ScriptID int64 `uri:"id" binding:"required"`
	ReportID int64 `uri:"reportId" binding:"required"`
}

type GetReportResponse struct {
	*Report `json:",inline"`
	Content string `json:"content"`
}

// ResolveRequest 解决/重新打开举报
type ResolveRequest struct {
	mux.Meta `path:"/scripts/:id/reports/:reportId/resolve" method:"PUT"`
	ScriptID int64  `uri:"id" binding:"required"`
	ReportID int64  `uri:"reportId" binding:"required"`
	Content  string `json:"content" binding:"max=10485760" label:"评论内容"`
	Close    bool   `json:"close" binding:"omitempty" label:"关闭状态"` // true:解决 false:重新打开
}

type ResolveResponse struct {
	Comments []*Comment `json:"comments"`
}

// DeleteRequest 删除举报
type DeleteRequest struct {
	mux.Meta `path:"/scripts/:id/reports/:reportId" method:"DELETE"`
	ScriptID int64 `uri:"id" binding:"required"`
	ReportID int64 `uri:"reportId" binding:"required"`
}

type DeleteResponse struct {
}
