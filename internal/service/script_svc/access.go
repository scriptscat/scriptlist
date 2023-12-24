package script_svc

import (
	"context"
	"strconv"
	"time"

	"github.com/codfrm/cago/pkg/consts"

	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
)

type AccessSvc interface {
	// AccessList 访问控制列表
	AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error)
	// CreateAccess 创建访问控制
	CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error)
	// UpdateAccess 更新访问控制
	UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error)
	// DeleteAccess 删除访问控制
	DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error)
	// Check 检查是否有访问权限
	Check(ctx context.Context, res, act string) (*CheckAccess, error)
	// CheckHandler 检查是否有访问权限中间件
	CheckHandler(res, act string, opts ...CheckOption) gin.HandlerFunc
	// RequireAccess 需要权限存在
	RequireAccess() gin.HandlerFunc
	// CtxAccess 获取访问权限
	CtxAccess(ctx context.Context) *script_entity.ScriptAccess
}

type CheckOption func(*CheckOptions)

type CheckOptions struct {
	Skip func(ctx *gin.Context) (bool, error)
}

func WithCheckSkip(f func(ctx *gin.Context) (bool, error)) CheckOption {
	return func(options *CheckOptions) {
		options.Skip = f
	}
}

type accessSvc struct {
}

var defaultAccess = &accessSvc{}

func Access() AccessSvc {
	return defaultAccess
}

// RequireAccess 检查访问权限
func (a *accessSvc) RequireAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		sSid := c.Param("aid")
		aid, err := strconv.ParseInt(sSid, 10, 64)
		if err != nil {
			httputils.HandleResp(c, i18n.NewError(c, code.AccessNotFound))
			return
		}
		access, err := script_repo.ScriptAccess().Find(c, Script().CtxScript(c).ID, aid)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		if err := access.Check(c); err != nil {
			httputils.HandleResp(c, err)
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), accessCtxKey, access))
	}
}

func (a *accessSvc) CtxAccess(ctx context.Context) *script_entity.ScriptAccess {
	return ctx.Value(accessCtxKey).(*script_entity.ScriptAccess)
}

// AccessList 访问控制列表
func (a *accessSvc) AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error) {
	list, total, err := script_repo.ScriptAccess().FindPage(ctx, Script().CtxScript(ctx).ID, req.PageRequest)
	if err != nil {
		return nil, err
	}
	ret := &api.AccessListResponse{
		PageResponse: httputils.PageResponse[*api.Access]{
			List:  make([]*api.Access, 0),
			Total: total,
		},
	}
	for _, v := range list {
		access := &api.Access{
			ID:         v.ID,
			LinkID:     v.LinkID,
			Type:       int32(v.Type),
			Role:       string(v.Role),
			IsExpire:   v.IsExpired(),
			Expiretime: v.Expiretime,
			Createtime: v.Createtime,
		}
		switch v.Type {
		case script_entity.AccessTypeUser:
			user, err := user_repo.User().Find(ctx, v.LinkID)
			if err != nil {
				return nil, err
			}
			if user == nil {
				continue
			}
			access.Name = user.Username
			access.Avatar = user.Avatar()
		case script_entity.AccessTypeGroup:
			group, err := script_repo.ScriptGroup().Find(ctx, v.ScriptID, v.LinkID)
			if err != nil {
				return nil, err
			}
			if group == nil {
				continue
			}
			access.Name = group.Name
			access.Avatar = ""
		}
		ret.List = append(ret.List, access)
	}
	return ret, nil
}

func (a *accessSvc) checkLinkExist(ctx context.Context, scriptId, linkId int64, accessType script_entity.AccessType) error {
	// 检查link id是否存在
	switch accessType {
	case script_entity.AccessTypeUser:
		user, err := user_repo.User().Find(ctx, linkId)
		if err != nil {
			return err
		}
		if err := user.IsBanned(ctx); err != nil {
			return err
		}
	case script_entity.AccessTypeGroup:
		group, err := script_repo.ScriptGroup().Find(ctx, scriptId, linkId)
		if err != nil {
			return err
		}
		if err := group.Check(ctx); err != nil {
			return err
		}
	}
	return nil
}

// CreateAccess 创建访问控制
func (a *accessSvc) CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error) {
	script := Script().CtxScript(ctx)
	// 检查是否重复
	list, err := script_repo.ScriptAccess().FindByLinkID(ctx, script.ID, req.LinkID, req.Type)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return nil, i18n.NewError(ctx, code.AccessAlreadyExist)
	}
	if err := a.checkLinkExist(ctx, script.ID, req.LinkID, req.Type); err != nil {
		return nil, err
	}
	// 创建
	access := &script_entity.ScriptAccess{
		ScriptID:   script.ID,
		LinkID:     req.LinkID,
		Type:       req.Type,
		Role:       req.Role,
		Status:     consts.ACTIVE,
		Expiretime: req.Expiretime,
		Createtime: time.Now().Unix(),
		Updatetime: time.Now().Unix(),
	}
	if err := script_repo.ScriptAccess().Create(ctx, access); err != nil {
		return nil, err
	}
	return &api.CreateAccessResponse{}, nil
}

// UpdateAccess 更新访问控制
func (a *accessSvc) UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error) {
	access := a.CtxAccess(ctx)
	access.Role = req.Role
	access.Expiretime = req.Expiretime
	if err := script_repo.ScriptAccess().Update(ctx, access); err != nil {
		return nil, err
	}
	return &api.UpdateAccessResponse{}, nil
}

// DeleteAccess 删除访问控制
func (a *accessSvc) DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error) {
	access := a.CtxAccess(ctx)
	if err := script_repo.ScriptAccess().Delete(ctx, access.ID); err != nil {
		return nil, err
	}
	return nil, nil
}

var roleAccessMap = map[script_entity.AccessRole]map[string]map[string]struct{}{
	"admin": {
		"script": {
			"delete:score": struct{}{},
			"delete":       struct{}{},
			"manage":       struct{}{},
			"read:info":    struct{}{},
		},
		"group": {
			"read":   struct{}{},
			"manage": struct{}{},
		},
		"access": {
			"read": struct{}{},
		},
		"issue": {
			"manage": struct{}{},
			"delete": struct{}{},
		},
		"statistics": {
			"manage": struct{}{},
		},
	},
	"owner": {
		"script": {
			"delete":    struct{}{},
			"manage":    struct{}{},
			"read:info": struct{}{},
		},
		"group": {
			"read":   struct{}{},
			"manage": struct{}{},
		},
		"access": {
			"read":   struct{}{},
			"manage": struct{}{},
		},
		"issue": {
			"manage": struct{}{},
			"delete": struct{}{},
		},
		"statistics": {
			"manage": struct{}{},
		},
	},
	"manager": {
		"script": {
			"manage":    struct{}{},
			"read:info": struct{}{},
		},
		"issue": {
			"manage": struct{}{},
			"delete": struct{}{},
		},
		"statistics": {
			"manage": struct{}{},
		},
	},
	"guest": {
		"script": {
			"read:info": struct{}{},
		},
	},
}

func (a *accessSvc) RoleToAccess(roles []script_entity.AccessRole) map[string]map[string]struct{} {
	ret := map[string]map[string]struct{}{}
	for _, role := range roles {
		for access := range roleAccessMap[role] {
			for act := range roleAccessMap[role][access] {
				if _, ok := ret[access]; !ok {
					ret[access] = map[string]struct{}{}
				}
				ret[access][act] = struct{}{}
			}
		}
	}
	return ret
}

func (a *accessSvc) GetUserAccess(ctx context.Context, scriptId, userId int64) ([]script_entity.AccessRole, error) {
	// 先在列表中查询
	list, err := script_repo.ScriptAccess().FindByLinkID(ctx, scriptId, userId, script_entity.AccessTypeUser)
	if err != nil {
		return nil, err
	}
	roles := make([]script_entity.AccessRole, 0)
	// 再在组中查询
	groups, err := script_repo.ScriptGroupMember().FindByUserId(ctx, scriptId, userId)
	if err != nil {
		return nil, err
	}
	// 再通过组id查出组的权限
	for _, v := range groups {
		// 检查是否过期
		if v.IsExpired() {
			continue
		}
		list, err := script_repo.ScriptAccess().FindByLinkID(ctx, scriptId, v.GroupID, script_entity.AccessTypeGroup)
		if err != nil {
			return nil, err
		}
		for _, v := range list {
			if !v.IsExpired() {
				roles = append(roles, v.Role)
			}
		}
	}
	for _, v := range list {
		if !v.IsExpired() {
			roles = append(roles, v.Role)
		}
	}
	if len(roles) == 0 {
		return nil, i18n.NewForbiddenError(ctx, code.UserNotPermission)
	}
	return roles, nil
}

type CheckAccess struct {
	Roles     []script_entity.AccessRole
	AccessMap map[string]map[string]struct{}
}

func (a *CheckAccess) Check(ctx context.Context, res, act string) error {
	if _, ok := a.AccessMap[res]; !ok {
		return i18n.NewForbiddenError(ctx, code.UserNotPermission)
	}
	if _, ok := a.AccessMap[res][act]; !ok {
		return i18n.NewForbiddenError(ctx, code.UserNotPermission)
	}
	return nil
}

func (a *accessSvc) Check(ctx context.Context, res, act string) (*CheckAccess, error) {
	// 获取用户对该脚本拥有的权限
	script := Script().CtxScript(ctx)
	user := auth_svc.Auth().Get(ctx)
	var (
		roles = make([]script_entity.AccessRole, 0)
		err   error
	)
	if user.AdminLevel.IsAdmin(model.Admin) {
		roles = append(roles, "admin")
	}
	if user.UID == script.UserID {
		roles = append(roles, "owner")
	}
	if len(roles) == 0 {
		roles, err = a.GetUserAccess(ctx, script.ID, user.UID)
		if err != nil {
			return nil, err
		}
	}
	accessMap := a.RoleToAccess(roles)
	access := &CheckAccess{
		Roles:     roles,
		AccessMap: accessMap,
	}
	if err := access.Check(ctx, res, act); err != nil {
		return nil, err
	}
	return access, nil
}

func (a *accessSvc) CheckHandler(res, act string, opts ...CheckOption) gin.HandlerFunc {
	options := &CheckOptions{}
	for _, o := range opts {
		o(options)
	}
	return func(ctx *gin.Context) {
		if options.Skip != nil {
			if ok, err := options.Skip(ctx); err != nil {
				httputils.HandleResp(ctx, err)
				return
			} else if ok {
				return
			}
		}
		access, ok := ctx.Value(checkAccessCtxKey).(*CheckAccess)
		if ok {
			if err := access.Check(ctx, res, act); err != nil {
				httputils.HandleResp(ctx, err)
				return
			}
		} else {
			if access, err := a.Check(ctx, res, act); err != nil {
				httputils.HandleResp(ctx, err)
				return
			} else {
				ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), checkAccessCtxKey, access))
			}
		}
	}
}
