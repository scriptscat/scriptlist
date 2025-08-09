package script_svc

import (
	"context"
	"strconv"
	"time"

	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/task/producer"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"

	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"

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
	// AddMemberInternal 添加成员 内部调用
	AddMemberInternal(ctx context.Context, entity *script_entity.ScriptGroupMember) error
	// RemoveMember 移除成员
	RemoveMember(ctx context.Context, req *api.RemoveMemberRequest) (*api.RemoveMemberResponse, error)
	// RequireGroup 需要群组存在
	RequireGroup() gin.HandlerFunc
	// CtxGroup 获取群组
	CtxGroup(c context.Context) *script_entity.ScriptGroup
	// RequireMember 需要成员存在
	RequireMember() gin.HandlerFunc
	// CtxMember 获取成员
	CtxMember(c context.Context) *script_entity.ScriptGroupMember
	// UpdateMember 更新成员
	UpdateMember(ctx context.Context, req *api.UpdateMemberRequest) (*api.UpdateMemberResponse, error)
}

type groupSvc struct {
}

var defaultGroup = &groupSvc{}

func Group() GroupSvc {
	return defaultGroup
}

func (g *groupSvc) toMembers(ctx context.Context, m []*script_entity.ScriptGroupMember) ([]*api.GroupMember, error) {
	list := make([]*api.GroupMember, 0, len(m))
	for _, v := range m {
		user, err := user_repo.User().Find(ctx, v.UserID)
		if err != nil {
			return nil, err
		}
		list = append(list, &api.GroupMember{
			ID:           v.ID,
			UserID:       v.UserID,
			Username:     user.Username,
			Avatar:       user.Avatar(),
			InviteStatus: v.InviteStatus,
			IsExpire:     v.IsExpired(),
			Expiretime:   v.Expiretime,
			Createtime:   v.Createtime,
		})
	}
	return list, nil
}

// GroupList 群组列表
func (g *groupSvc) GroupList(ctx context.Context, req *api.GroupListRequest) (*api.GroupListResponse, error) {
	list, total, err := script_repo.ScriptGroup().FindPage(ctx, Script().CtxScript(ctx).ID, req)
	if err != nil {
		return nil, err
	}
	ret := make([]*api.Group, 0, len(list))
	for _, v := range list {
		// 获取群组前10人与群组总人数
		memberList, memberTotal, err := script_repo.ScriptGroupMember().FindPage(ctx, Script().CtxScript(ctx).ID, v.ID, httputils.PageRequest{Page: 1, Size: 10})
		if err != nil {
			return nil, err
		}
		member, err := g.toMembers(ctx, memberList)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &api.Group{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Member:      member,
			MemberCount: memberTotal,
			Createtime:  v.Createtime,
		})
	}
	return &api.GroupListResponse{
		PageResponse: httputils.PageResponse[*api.Group]{
			List:  ret,
			Total: total,
		},
	}, nil
}

// CreateGroup 创建群组
func (g *groupSvc) CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.CreateGroupResponse, error) {
	group := &script_entity.ScriptGroup{
		ScriptID:    Script().CtxScript(ctx).ID,
		Name:        req.Name,
		Description: req.Description,
		Status:      consts.ACTIVE,
		Createtime:  time.Now().Unix(),
		Updatetime:  time.Now().Unix(),
	}
	if err := script_repo.ScriptGroup().Create(ctx, group); err != nil {
		return nil, err
	}
	return &api.CreateGroupResponse{}, nil
}

// UpdateGroup 更新群组
func (g *groupSvc) UpdateGroup(ctx context.Context, req *api.UpdateGroupRequest) (*api.UpdateGroupResponse, error) {
	group := g.CtxGroup(ctx)
	group.Name = req.Name
	group.Description = req.Description
	group.Updatetime = time.Now().Unix()
	if err := script_repo.ScriptGroup().Update(ctx, group); err != nil {
		return nil, err
	}
	return &api.UpdateGroupResponse{}, nil
}

// DeleteGroup 删除群组
func (g *groupSvc) DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) (*api.DeleteGroupResponse, error) {
	group := g.CtxGroup(ctx)
	if err := script_repo.ScriptGroup().Delete(ctx, group.ID); err != nil {
		return nil, err
	}
	return &api.DeleteGroupResponse{}, nil
}

// GroupMemberList 群组成员列表
func (g *groupSvc) GroupMemberList(ctx context.Context, req *api.GroupMemberListRequest) (*api.GroupMemberListResponse, error) {
	list, total, err := script_repo.ScriptGroupMember().FindPage(ctx, Script().CtxScript(ctx).ID, g.CtxGroup(ctx).ID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret, err := g.toMembers(ctx, list)
	if err != nil {
		return nil, err
	}
	return &api.GroupMemberListResponse{
		PageResponse: httputils.PageResponse[*api.GroupMember]{
			List:  ret,
			Total: total,
		},
	}, nil
}

// AddMember 添加成员
func (g *groupSvc) AddMember(ctx context.Context, req *api.AddMemberRequest) (*api.AddMemberResponse, error) {
	script := Script().CtxScript(ctx)
	group := g.CtxGroup(ctx)
	// 添加成员
	member := &script_entity.ScriptGroupMember{
		ScriptID:     script.ID,
		GroupID:      group.ID,
		UserID:       req.UserID,
		InviteStatus: script_entity.AccessInviteStatusPending,
		Status:       consts.ACTIVE,
		Expiretime:   req.Expiretime,
		Createtime:   time.Now().Unix(),
		Updatetime:   time.Now().Unix(),
	}
	if err := g.AddMemberInternal(ctx, member); err != nil {
		return nil, err
	}
	// 创建邀请链接
	invResp, err := AccessInvite().CreateInviteLink(ctx, &script_entity.ScriptInvite{
		ScriptID: script.ID,
		GroupID:  group.ID,
		Type:     script_entity.InviteTypeGroup,
		UserID:   member.ID,
	})
	if err != nil {
		return nil, err
	}
	// 发送邮件通知
	if err := producer.PublishAccessInvite(ctx, auth_svc.Auth().Get(ctx).UID, req.UserID, invResp); err != nil {
		return nil, err
	}
	return &api.AddMemberResponse{}, nil
}

// AddMemberInternal 添加成员 内部调用
func (g *groupSvc) AddMemberInternal(ctx context.Context, entity *script_entity.ScriptGroupMember) error {
	// 检查用户
	user, err := user_repo.User().Find(ctx, entity.UserID)
	if err != nil {
		return err
	}
	if err := user.IsBanned(ctx); err != nil {
		return err
	}
	// 添加成员
	entity.Updatetime = time.Now().Unix()
	if entity.ID == 0 {
		// 检查是否已经在群组中
		if list, err := script_repo.ScriptGroupMember().FindByUserId(ctx, entity.ScriptID, entity.UserID); err != nil {
			return err
		} else if len(list) > 0 {
			return i18n.NewError(ctx, code.GroupMemberExist)
		}
		entity.Createtime = time.Now().Unix()
		if err := script_repo.ScriptGroupMember().Create(ctx, entity); err != nil {
			return err
		}
	} else {
		if err := script_repo.ScriptGroupMember().Update(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMember 更新成员
func (g *groupSvc) UpdateMember(ctx context.Context, req *api.UpdateMemberRequest) (*api.UpdateMemberResponse, error) {
	member := g.CtxMember(ctx)
	member.Expiretime = req.Expiretime
	member.Updatetime = time.Now().Unix()
	if err := script_repo.ScriptGroupMember().Update(ctx, member); err != nil {
		return nil, err
	}
	return &api.UpdateMemberResponse{}, nil
}

// RemoveMember 移除成员
func (g *groupSvc) RemoveMember(ctx context.Context, req *api.RemoveMemberRequest) (*api.RemoveMemberResponse, error) {
	member := g.CtxMember(ctx)
	if err := script_repo.ScriptGroupMember().Delete(ctx, member.ID); err != nil {
		return nil, err
	}
	return &api.RemoveMemberResponse{}, nil
}

func (g *groupSvc) RequireGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		gid := c.Param("gid")
		id, err := strconv.ParseInt(gid, 10, 64)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		group, err := script_repo.ScriptGroup().Find(c, Script().CtxScript(c).ID, id)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		if err := group.Check(c); err != nil {
			httputils.HandleResp(c, err)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), groupCtxKey, group))

	}
}

func (g *groupSvc) CtxGroup(c context.Context) *script_entity.ScriptGroup {
	return c.Value(groupCtxKey).(*script_entity.ScriptGroup)
}

func (g *groupSvc) RequireMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		gid := c.Param("mid")
		id, err := strconv.ParseInt(gid, 10, 64)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		member, err := script_repo.ScriptGroupMember().Find(c, Script().CtxScript(c).ID, id)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		if err := member.Check(c); err != nil {
			httputils.HandleResp(c, err)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), memberCtxKey, member))

	}
}

func (g *groupSvc) CtxMember(c context.Context) *script_entity.ScriptGroupMember {
	return c.Value(memberCtxKey).(*script_entity.ScriptGroupMember)
}
