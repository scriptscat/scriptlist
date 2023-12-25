package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type Access struct {
	ID           int64                            `json:"id"`
	LinkID       int64                            `json:"link_id"` // 关联id
	Name         string                           `json:"name"`
	Avatar       string                           `json:"avatar"`
	Type         script_entity.AccessType         `json:"type"`         // id类型 1=用户id 2=组id
	InviteStatus script_entity.AccessInviteStatus `json:"invite_tatus"` // 邀请状态 1=已接受 2=已拒绝 3=待接受
	Role         script_entity.AccessRole         `json:"role"`
	IsExpire     bool                             `json:"is_expire"`
	Expiretime   int64                            `json:"expiretime"`
	Createtime   int64                            `json:"createtime"`
}

// AccessListRequest 访问控制列表
type AccessListRequest struct {
	mux.Meta              `path:"/scripts/:id/access" method:"GET"`
	ScriptID              int64 `uri:"id" binding:"required" label:"id"`
	httputils.PageRequest `form:",inline"`
}

type AccessListResponse struct {
	httputils.PageResponse[*Access] `json:",inline"`
}

// AddGroupAccessRequest 添加组权限
type AddGroupAccessRequest struct {
	mux.Meta   `path:"/scripts/:id/access/group" method:"POST"`
	ScriptID   int64                    `uri:"id" binding:"required" label:"id"`
	GroupID    int64                    `form:"group_id" json:"group_id"  binding:"required" label:"群组id"`
	Role       script_entity.AccessRole `form:"role" binding:"required,oneof=guest manager" label:"访问权限"` // 访问权限 guest=访客 manager=管理员
	Expiretime int64                    `form:"expiretime,default=0" label:"过期时间"`                        // 0 为永久
}

type AddGroupAccessResponse struct {
}

// AddUserAccessRequest 添加用户权限, 通过用户名进行邀请
type AddUserAccessRequest struct {
	mux.Meta   `path:"/scripts/:id/access/user" method:"POST"`
	ScriptID   int64                    `uri:"id" binding:"required" label:"id"`
	UserID     int64                    `form:"user_id" json:"user_id"  binding:"required" label:"用户id"`
	Role       script_entity.AccessRole `form:"role" binding:"required,oneof=guest manager" label:"访问权限"` // 访问权限 guest=访客 manager=管理员
	Expiretime int64                    `form:"expiretime,default=0" label:"过期时间"`                        // 0 为永久
}

type AddUserAccessResponse struct {
}

// UpdateAccessRequest 更新访问控制
type UpdateAccessRequest struct {
	mux.Meta   `path:"/scripts/:id/access/:aid" method:"PUT"`
	ScriptID   int64                    `uri:"id" binding:"required" label:"id"`
	AccessID   int64                    `uri:"aid" binding:"required" label:"id"`
	Role       script_entity.AccessRole `form:"role" binding:"required,oneof=guest manager" label:"访问权限"` // 访问权限 guest=访客 manager=管理员
	Expiretime int64                    `form:"expiretime,default=0" label:"过期时间"`                        // 0 为永久
}

type UpdateAccessResponse struct {
}

// DeleteAccessRequest 删除访问控制
type DeleteAccessRequest struct {
	mux.Meta `path:"/scripts/:id/access/:aid" method:"DELETE"`
	ScriptID int64 `uri:"id" binding:"required" label:"id"`
	AccessID int64 `uri:"aid" binding:"required" label:"id"`
}

type DeleteAccessResponse struct {
}
