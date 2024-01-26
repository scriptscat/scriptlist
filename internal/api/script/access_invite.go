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
	InviteStatus script_entity.InviteStatus `json:"invite_status"`
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

// GroupInviteCodeListRequest 群组邀请码列表
type GroupInviteCodeListRequest struct {
	mux.Meta              `path:"/scripts/:id/invite/group/:gid/code" method:"GET"`
	ScriptID              int64 `uri:"id" binding:"required" label:"id"`
	GroupID               int64 `uri:"gid" binding:"required" label:"gid"`
	httputils.PageRequest `form:",inline"`
}

type GroupInviteCodeListResponse struct {
	httputils.PageResponse[*InviteCode] `json:",inline"`
}

// CreateInviteCodeRequest 创建邀请码
type CreateInviteCodeRequest struct {
	mux.Meta `path:"/scripts/:id/invite/code" method:"POST"`
	ScriptID int64 `uri:"id" binding:"required" label:"id"`
	Count    int32 `form:"count,default=1" json:"count" label:"数量"`
	Audit    bool  `form:"audit" json:"audit" label:"是否需要审核"`
	Days     int32 `form:"days,default=0" json:"days" label:"有效天数"` // 0 为永久
}

type CreateInviteCodeResponse struct {
	Code []string `json:"code"`
}

// CreateGroupInviteCodeRequest 创建群组邀请码
type CreateGroupInviteCodeRequest struct {
	mux.Meta `path:"/scripts/:id/invite/group/:gid/code" method:"POST"`
	ScriptID int64 `uri:"id" binding:"required" label:"id"`
	GroupID  int64 `uri:"gid" binding:"required" label:"gid"`
	Count    int32 `form:"count,default=1" json:"count" label:"数量"`
	Audit    bool  `form:"audit" json:"audit" label:"是否需要审核"`
	Days     int32 `form:"days,default=0" json:"days" label:"有效天数"` // 0 为永久
}

type CreateGroupInviteCodeResponse struct {
	Code []string `json:"code"`
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
	ScriptID int64
	LinkID   int64
	Type     script_entity.InviteType // 1=权限邀请码 2=群组邀请码
}

// AcceptInviteRequest 接受邀请
type AcceptInviteRequest struct {
	mux.Meta `path:"/scripts/invite/:code/accept" method:"PUT"`
	Code     string `uri:"code" binding:"required" label:"code"`
	// 邀请码类型不能拒绝
	Accept bool `form:"accept" json:"accept" binding:"required" label:"是否接受"`
}

type AcceptInviteResponse struct {
}

// InviteCodeInfoRequest 邀请码信息
type InviteCodeInfoRequest struct {
	mux.Meta `path:"/scripts/invite/:code" method:"GET"`
	Code     string `uri:"code" binding:"required" label:"code"`
}

type InviteCodeInfoAccess struct {
	Role script_entity.AccessRole `json:"role"` // 访问权限 guest=访客 manager=管理员
}

type InviteCodeInfoGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InviteCodeInfoResponse struct {
	CodeType     script_entity.InviteCodeType `json:"code_type"`     // 邀请码类型 1=邀请码 2=邀请链接
	InviteStatus script_entity.InviteStatus   `json:"invite_status"` // 使用状态
	Type         script_entity.InviteType     `json:"type"`          // 邀请类型 1=权限邀请码 2=群组邀请码
	IsAudit      bool                         `json:"is_audit"`      // 是否需要审核 邀请码类型为邀请链接时，该字段固定为false
	Script       *Script                      `json:"script"`
	Access       *InviteCodeInfoAccess        `json:"access,omitempty"` // 如果type=1, 则返回权限信息
	Group        *InviteCodeInfoGroup         `json:"group,omitempty"`  // 如果type=2, 则返回群组信息
}
