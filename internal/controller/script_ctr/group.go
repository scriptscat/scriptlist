package script_ctr

import (
	"context"

	"github.com/codfrm/cago/pkg/utils/muxutils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"

	"github.com/gin-gonic/gin"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Group struct {
}

func NewGroup() *Group {
	return &Group{}
}

// GroupList 群组列表
func (g *Group) GroupList(ctx context.Context, req *api.GroupListRequest) (*api.GroupListResponse, error) {
	return script_svc.Group().GroupList(ctx, req)
}

// CreateGroup 创建群组
func (g *Group) CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.CreateGroupResponse, error) {
	return script_svc.Group().CreateGroup(ctx, req)
}

// UpdateGroup 更新群组
func (g *Group) UpdateGroup(ctx context.Context, req *api.UpdateGroupRequest) (*api.UpdateGroupResponse, error) {
	return script_svc.Group().UpdateGroup(ctx, req)
}

// DeleteGroup 删除群组
func (g *Group) DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) (*api.DeleteGroupResponse, error) {
	return script_svc.Group().DeleteGroup(ctx, req)
}

// GroupMemberList 群组成员列表
func (g *Group) GroupMemberList(ctx context.Context, req *api.GroupMemberListRequest) (*api.GroupMemberListResponse, error) {
	return script_svc.Group().GroupMemberList(ctx, req)
}

// AddMember 添加成员
func (g *Group) AddMember(ctx context.Context, req *api.AddMemberRequest) (*api.AddMemberResponse, error) {
	return script_svc.Group().AddMember(ctx, req)
}

// RemoveMember 移除成员
func (g *Group) RemoveMember(ctx context.Context, req *api.RemoveMemberRequest) (*api.RemoveMemberResponse, error) {
	return script_svc.Group().RemoveMember(ctx, req)
}

func (g *Group) Middleware() gin.HandlerFunc {
	return script_svc.Group().Middleware()
}

func (g *Group) Router(r *mux.Router) {
	muxutils.BindTree(r, []*muxutils.RouterTree{{
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Access().CheckHandler("group", "read"),
		},
		Handler: []interface{}{
			g.GroupList,
			g.GroupMemberList,
		},
	}, {
		Middleware: []gin.HandlerFunc{
			auth_svc.Auth().RequireLogin(true),
			script_svc.Access().CheckHandler("group", "manage"),
		},
		Handler: []interface{}{
			g.CreateGroup,
			g.UpdateGroup,
			g.DeleteGroup,
			g.AddMember,
			g.RemoveMember,
		},
	}})
}
