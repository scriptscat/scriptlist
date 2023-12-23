package script

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
)

type Group struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Member      []*GroupMember `json:"member"`
	MemberCount int64          `json:"member_count"`
	Createtime  int64          `json:"createtime"`
}

// GroupListRequest 群组列表
type GroupListRequest struct {
	mux.Meta              `path:"/scripts/:id/group" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type GroupListResponse struct {
	httputils.PageResponse[*Group] `json:",inline"`
}

// CreateGroupRequest 创建群组
type CreateGroupRequest struct {
	mux.Meta    `path:"/scripts/:id/group" method:"POST"`
	Name        string `json:"name" binding:"required,max=20" label:"群组名"`
	Description string `json:"description" binding:"required,max=200" label:"群组描述"`
}

type CreateGroupResponse struct {
}

// UpdateGroupRequest 更新群组
type UpdateGroupRequest struct {
	mux.Meta    `path:"/scripts/:id/group/:gid" method:"PUT"`
	Name        string `json:"name" binding:"required,max=20" label:"群组名"`
	Description string `json:"description" binding:"required,max=200" label:"群组描述"`
}

type UpdateGroupResponse struct {
}

// DeleteGroupRequest 删除群组
type DeleteGroupRequest struct {
	mux.Meta `path:"/scripts/:id/group/:gid" method:"DELETE"`
}

type DeleteGroupResponse struct {
}

// GroupMemberListRequest 群组成员列表
type GroupMemberListRequest struct {
	mux.Meta              `path:"/scripts/:id/group/:gid/member" method:"GET"`
	httputils.PageRequest `form:",inline"`
}

type GroupMember struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	IsExpire   bool   `json:"is_expire"`
	Expiretime int64  `json:"expiretime"`
	Createtime int64  `json:"createtime"`
}

type GroupMemberListResponse struct {
	httputils.PageResponse[*GroupMember] `json:",inline"`
}

// AddMemberRequest 添加成员
type AddMemberRequest struct {
	mux.Meta   `path:"/scripts/:id/group/:gid/member" method:"POST"`
	UserID     int64 `json:"user_id" binding:"required"`
	Expiretime int64 `json:"expiretime" binding:"required"`
}

type AddMemberResponse struct {
}

// UpdateMemberRequest 更新成员
type UpdateMemberRequest struct {
	mux.Meta   `path:"/scripts/:id/group/:gid/member/:mid" method:"PUT"`
	ID         int64 `uri:"mid" binding:"required"`
	Expiretime int64 `json:"expiretime" binding:"required"`
}

type UpdateMemberResponse struct {
}

// RemoveMemberRequest 移除成员
type RemoveMemberRequest struct {
	mux.Meta `path:"/scripts/:id/group/:gid/member/:mid" method:"DELETE"`
	ID       int64 `uri:"mid" binding:"required"`
}

type RemoveMemberResponse struct {
}
