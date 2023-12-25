package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type InviteCode struct {
	ID           int64                      `json:"id"`
	Code         string                     `json:"code"`     // 邀请码
	UserID       int64                      `json:"used"`     // 使用用户
	Username     string                     `json:"username"` // 使用用户名
	IsAudit      bool                       `json:"is_audit"` // 是否需要审核
	InviteStatus script_entity.InviteStatus `json:"inv_status"`
	Expiretime   int64                      `json:"expiretime"` // 到期时间
	Createtime   int64                      `json:"createtime"`
}

// InviteCodeListRequest 邀请码列表
type InviteCodeListRequest struct {
	mux.Meta              `path:"/scripts/:id/invite/code" method:"GET"`
	ScriptID              int64 `uri:"id" binding:"required" label:"id"`
	httputils.PageRequest `form:",inline"`
}

type InviteCodeListResponse struct {
	httputils.PageResponse[*InviteCode] `json:",inline"`
}

// GroupInviteCodeRequest 群组邀请码列表
type GroupInviteCodeRequest struct {
	mux.Meta              `path:"/scripts/:id/invite/group/:gid/code" method:"GET"`
	ScriptID              int64 `uri:"id" binding:"required" label:"id"`
	GroupID               int64 `uri:"gid" binding:"required" label:"gid"`
	httputils.PageRequest `form:",inline"`
}

type GroupInviteCodeResponse struct {
	httputils.PageResponse[*InviteCode] `json:",inline"`
}

// CreateInviteCodeRequest 创建邀请码
type CreateInviteCodeRequest struct {
	mux.Meta   `path:"/scripts/:id/invite/code" method:"POST"`
	ScriptID   int64 `uri:"id" binding:"required" label:"id"`
	Count      int32 `form:"count,default=1" label:"数量"`
	Audit      bool  `form:"audit" label:"是否需要审核"`
	Expiretime int64 `form:"expiretime,default=0" label:"过期时间"` // 0 为永久
}

type CreateInviteCodeResponse struct {
}

// CreateGroupInviteCodeRequest 创建群组邀请码
type CreateGroupInviteCodeRequest struct {
	mux.Meta   `path:"/scripts/:id/invite/group/:gid/code" method:"POST"`
	ScriptID   int64 `uri:"id" binding:"required" label:"id"`
	GroupID    int64 `uri:"gid" binding:"required" label:"gid"`
	Count      int32 `form:"count,default=1" label:"数量"`
	Audit      bool  `form:"audit" label:"是否需要审核"`
	Expiretime int64 `form:"expiretime,default=0" label:"过期时间"` // 0 为永久
}

type CreateGroupInviteCodeResponse struct {
}

// DeleteInviteCodeRequest 删除邀请码
type DeleteInviteCodeRequest struct {
	mux.Meta `path:"/scripts/:id/invite/code/:code_id" method:"DELETE"`
	ScriptID int64 `uri:"id" binding:"required" label:"id"`
	CodeID   int64 `uri:"code_id" binding:"required" label:"code_id"`
}

type DeleteInviteCodeResponse struct {
}

// AuditInviteCodeRequest 审核邀请码
type AuditInviteCodeRequest struct {
	mux.Meta `path:"/scripts/:id/invite/code/:code_id/audit" method:"PUT"`
	ScriptID int64 `uri:"id" binding:"required" label:"id"`
	CodeID   int64 `uri:"code_id" binding:"required" label:"code_id"`
	Status   int32 `form:"status" binding:"required,oneof=1 2" label:"状态"` // 1=通过 2=拒绝
}

type AuditInviteCodeResponse struct {
}

// CreateInviteLinkRequest 创建邀请链接
type CreateInviteLinkRequest struct {
	ScriptID int64                    `json:"script_id" label:"脚本ID"`
	Type     script_entity.InviteType `json:"type" label:"类型"` // 1=权限邀请码 2=群组邀请码
}

type CreateInviteLinkResponse struct {
}

// AcceptInviteRequest 接受邀请
type AcceptInviteRequest struct {
	mux.Meta `path:"/scripts/invite/:code/accept" method:"PUT"`
	Code     string `uri:"code" binding:"required" label:"code"`
}

type AcceptInviteResponse struct {
}

// RejectInviteRequest 拒绝邀请
type RejectInviteRequest struct {
	mux.Meta `path:"/scripts/invite/:code/reject" method:"PUT"`
	Code     string `uri:"code" binding:"required" label:"code"`
}

type RejectInviteResponse struct {
}
