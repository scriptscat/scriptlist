package script_svc

import (
	"context"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"net/http"
	"strconv"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type GroupSvc interface {
	// GroupList 群组列表
	GroupList(ctx context.Context, req *api.GroupListRequest) (*api.GroupListResponse, error)
	// CreateGroup 创建群组
	CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.CreateGroupResponse, error)
	// UpdateGroup 更新群组
	UpdateGroup(ctx context.Context, req *api.UpdateGroupRequest) (*api.UpdateGroupResponse, error)
	// DeleteGroup 删除群组
	DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) (*api.DeleteGroupResponse, error)
	// GroupMemberList 群组成员列表
	GroupMemberList(ctx context.Context, req *api.GroupMemberListRequest) (*api.GroupMemberListResponse, error)
	// AddMember 添加成员
	AddMember(ctx context.Context, req *api.AddMemberRequest) (*api.AddMemberResponse, error)
	// RemoveMember 移除成员
	RemoveMember(ctx context.Context, req *api.RemoveMemberRequest) (*api.RemoveMemberResponse, error)
	// Middleware 中间件
	Middleware() gin.HandlerFunc
}

type groupSvc struct {
}

var defaultGroup = &groupSvc{}

func Group() GroupSvc {
	return defaultGroup
}

// GroupList 群组列表
func (g *groupSvc) GroupList(ctx context.Context, req *api.GroupListRequest) (*api.GroupListResponse, error) {
	return nil, nil
}

// CreateGroup 创建群组
func (g *groupSvc) CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.CreateGroupResponse, error) {
	return nil, nil
}

// UpdateGroup 更新群组
func (g *groupSvc) UpdateGroup(ctx context.Context, req *api.UpdateGroupRequest) (*api.UpdateGroupResponse, error) {
	return nil, nil
}

// DeleteGroup 删除群组
func (g *groupSvc) DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) (*api.DeleteGroupResponse, error) {
	return nil, nil
}

// GroupMemberList 群组成员列表
func (g *groupSvc) GroupMemberList(ctx context.Context, req *api.GroupMemberListRequest) (*api.GroupMemberListResponse, error) {
	return nil, nil
}

// AddMember 添加成员
func (g *groupSvc) AddMember(ctx context.Context, req *api.AddMemberRequest) (*api.AddMemberResponse, error) {
	return nil, nil
}

// RemoveMember 移除成员
func (g *groupSvc) RemoveMember(ctx context.Context, req *api.RemoveMemberRequest) (*api.RemoveMemberResponse, error) {
	return nil, nil
}

func (g *groupSvc) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("gid")
		if sid == "" {
			httputils.HandleResp(c, httputils.NewError(http.StatusNotFound, -1, "群组ID不能为空"))
			return
		}
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		script, err := script_repo.Script().Find(c, id)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		if c.Request.Method == http.MethodGet {
			if err := script.CheckOperate(c); err != nil {
				httputils.HandleResp(c, err)
				return
			}
		} else {
			if err := script.CheckPermission(c, model.Moderator); err != nil {
				httputils.HandleResp(c, err)
				return
			}
		}
		c.Request = c.Request.WithContext(context.WithValue(
			c.Request.Context(), ctxScript("ctxScript"), script,
		))
		c.Next()
	}
}
