package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Access struct {
	ID               int64  `json:"id"`
	LinkID           int64  `json:"link_id"` // 关联id
	Name             string `json:"name"`
	Avatar           string `json:"avatar"`
	Type             int32  `json:"type"` // id类型 1=用户id 2=组id
	AccessPermission string `json:"access_permission"`
	Expiretime       int64  `json:"expiretime"`
	Createtime       int64  `json:"createtime"`
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
	mux.Meta         `path:"/scripts/:id/access" method:"POST"`
	LinkID           int64  `json:"link_id" binding:"required" label:"关联id"`
	Type             int32  `json:"type" binding:"required,oneof=1 2" label:"id类型"` // id类型 1=用户id 2=组id
	AccessPermission string `json:"access_permission" binding:"required" label:"访问权限"`
	Expiretime       int64  `json:"expiretime" binding:"required" label:"过期时间"`
}

type CreateAccessResponse struct {
}

// UpdateAccessRequest 更新访问控制
type UpdateAccessRequest struct {
	mux.Meta         `path:"/scripts/:id/access/:aid" method:"PUT"`
	LinkID           int64  `json:"link_id" binding:"required" label:"关联id"`
	Type             int32  `json:"type" binding:"required,oneof=1 2" label:"id类型"` // id类型 1=用户id 2=组id
	AccessPermission string `json:"access_permission" binding:"required" label:"访问权限"`
	Expiretime       int64  `json:"expiretime" binding:"required" label:"过期时间"`
}

type UpdateAccessResponse struct {
}

// DeleteAccessRequest 删除访问控制
type DeleteAccessRequest struct {
	mux.Meta `path:"/scripts/:id/access/:aid" method:"DELETE"`
}

type DeleteAccessResponse struct {
}
