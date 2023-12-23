package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type Access struct {
	ID         int64  `json:"id"`
	LinkID     int64  `json:"link_id"` // 关联id
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Type       int32  `json:"type"` // id类型 1=用户id 2=组id
	Role       string `json:"role"`
	IsExpire   bool   `json:"is_expire"`
	Expiretime int64  `json:"expiretime"`
	Createtime int64  `json:"createtime"`
}

// AccessListRequest 访问控制列表
type AccessListRequest struct {
	mux.Meta              `path:"/scripts/:id/access" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type AccessListResponse struct {
	httputils.PageResponse[*Access] `json:",inline"`
}

// CreateAccessRequest 创建访问控制
type CreateAccessRequest struct {
	mux.Meta   `path:"/scripts/:id/access" method:"POST"`
	LinkID     int64                    `form:"link_id" binding:"required" label:"关联id"`
	Type       script_entity.AccessType `form:"type" binding:"required,oneof=1 2" label:"id类型"`           // id类型 1=用户id 2=组id
	Role       script_entity.AccessRole `form:"role" binding:"required,oneof=guest manager" label:"访问权限"` // 访问权限 guest=访客 manager=管理员
	Expiretime int64                    `form:"expiretime,default=0" label:"过期时间"`                        // 0 为永久
}

type CreateAccessResponse struct {
}

// UpdateAccessRequest 更新访问控制
type UpdateAccessRequest struct {
	mux.Meta   `path:"/scripts/:id/access/:aid" method:"PUT"`
	ID         int64                    `uri:"aid" binding:"required" label:"id"`
	Role       script_entity.AccessRole `form:"role" binding:"required,oneof=guest manager" label:"访问权限"` // 访问权限 guest=访客 manager=管理员
	Expiretime int64                    `form:"expiretime,default=0" label:"过期时间"`                        // 0 为永久
}

type UpdateAccessResponse struct {
}

// DeleteAccessRequest 删除访问控制
type DeleteAccessRequest struct {
	mux.Meta `path:"/scripts/:id/access/:aid" method:"DELETE"`
	ID       int64 `uri:"aid" binding:"required" label:"id"`
}

type DeleteAccessResponse struct {
}
