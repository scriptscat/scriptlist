package script_svc

import (
	"context"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"time"
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

// AccessList 访问控制列表
func (a *accessSvc) AccessList(ctx context.Context, req *api.AccessListRequest) (*api.AccessListResponse, error) {
	return nil, nil
}

// CreateAccess 创建访问控制
func (a *accessSvc) CreateAccess(ctx context.Context, req *api.CreateAccessRequest) (*api.CreateAccessResponse, error) {
	return nil, nil
}

// UpdateAccess 更新访问控制
func (a *accessSvc) UpdateAccess(ctx context.Context, req *api.UpdateAccessRequest) (*api.UpdateAccessResponse, error) {
	return nil, nil
}

// DeleteAccess 删除访问控制
func (a *accessSvc) DeleteAccess(ctx context.Context, req *api.DeleteAccessRequest) (*api.DeleteAccessResponse, error) {
	return nil, nil
}

var roleAccessMap = map[string]map[string]map[string]struct{}{
	"admin": {
		"script": {
			"delete:score": struct{}{},
			"delete":       struct{}{},
			"manage":       struct{}{},
			"read:info":    struct{}{},
		},
		"group": {
			"read": struct{}{},
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
			"read":    struct{}{},
			"manager": struct{}{},
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

func (a *accessSvc) RoleToAccess(roles []string) map[string]map[string]struct{} {
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

func (a *accessSvc) GetUserAccess(ctx context.Context, scriptId, userId int64) ([]string, error) {
	// 先在列表中查询
	access, err := script_repo.ScriptAccess().FindByUserId(ctx, scriptId, userId)
	if err != nil {
		return []string{}, err
	}
	if access == nil {
		// 再在组中查询
		groups, err := script_repo.ScriptGroupMember().FindByUserId(ctx, scriptId, userId)
		if err != nil {
			return nil, err
		}
		if len(groups) == 0 {
			return nil, i18n.NewForbiddenError(ctx, code.UserNotPermission)
		}
		// 再通过组id查出组的权限
		access, err = script_repo.ScriptAccess().FindByGroupId(ctx, scriptId, groups[0].GroupID)
		if err != nil {
			return nil, err
		}
		if access == nil {
			return nil, i18n.NewForbiddenError(ctx, code.UserNotPermission)
		}
	}
	if time.Unix(access.Expiretime, 0).Before(time.Now()) {
		return nil, i18n.NewForbiddenError(ctx, code.UserNotPermission)
	}
	return []string{access.Role}, nil
}

type CheckAccess struct {
	Roles     []string
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
		roles []string
		err   error
	)
	if user.AdminLevel.IsAdmin(model.Admin) {
		roles = []string{"admin"}
	} else if user.UID == script.UserID {
		roles = []string{"owner"}
	} else {
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
		access, ok := ctx.Value(accessCtxKey).(*CheckAccess)
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
				ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), accessCtxKey, access))
			}
		}
	}
}
