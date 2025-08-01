package user_svc

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/api/script"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type UserSvc interface {
	// UserInfo 获取用户信息
	UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error)
	// Script 用户脚本列表
	Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error)
	// GetFollow 获取用户关注信息
	GetFollow(ctx context.Context, req *api.GetFollowRequest) (*api.GetFollowResponse, error)
	// Follow 关注用户
	Follow(ctx context.Context, req *api.FollowRequest) (*api.FollowResponse, error)
	// GetWebhook 获取webhook配置
	GetWebhook(ctx context.Context, req *api.GetWebhookRequest) (*api.GetWebhookResponse, error)
	// RefreshWebhook 刷新webhook配置
	RefreshWebhook(ctx context.Context, req *api.RefreshWebhookRequest) (*api.RefreshWebhookResponse, error)
	// GetConfig 获取用户配置
	GetConfig(ctx context.Context, req *api.GetConfigRequest) (*api.GetConfigResponse, error)
	// UpdateConfig 更新用户配置
	UpdateConfig(ctx context.Context, req *api.UpdateConfigRequest) (*api.UpdateConfigResponse, error)
	// Search 搜索用户
	Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error)
	// RefreshToken 刷新用户token
	RefreshToken(ctx *gin.Context, req *api.RefreshTokenRequest) (*api.RefreshTokenResponse, error)
	// Logout TODO
	Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error)
}

type userSvc struct {
}

var defaultUser = &userSvc{}

func User() UserSvc {
	return defaultUser
}

// UserInfo 获取用户信息
func (u *userSvc) UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error) {
	user, err := user_repo.User().Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, i18n.NewError(ctx, code.UserNotFound)
	}
	return &api.InfoResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Avatar:      user.Avatar(),
		IsAdmin:     model.AdminLevel(user.Adminid),
		EmailStatus: user.Emailstatus,
	}, nil
}

// Script 用户脚本列表
func (u *userSvc) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	self := false
	user := auth_svc.Auth().Get(ctx)
	if user != nil {
		self = user.UID == req.UID
	}
	resp, total, err := script_repo.Script().Search(ctx, &script_repo.SearchOptions{
		UserID:   req.UID,
		Keyword:  req.Keyword,
		Type:     req.ScriptType,
		Sort:     req.Sort,
		Category: make([]int64, 0),
		Self:     self,
	}, req.PageRequest)
	if err != nil {
		return nil, err
	}
	list := make([]*script.Script, 0)
	for _, item := range resp {
		data, err := script_svc.Script().ToScript(ctx, item, false, "")
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return &api.ScriptResponse{
		PageResponse: httputils.PageResponse[*script.Script]{
			List:  list,
			Total: total,
		},
	}, nil
}

// GetFollow 获取用户关注信息
func (u *userSvc) GetFollow(ctx context.Context, req *api.GetFollowRequest) (*api.GetFollowResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	isFollow := false
	if user != nil {
		record, err := user_repo.Follow().Find(ctx, user.UID, req.UID)
		if err != nil {
			return nil, err
		}
		if record != nil {
			isFollow = true
		}
	}
	_, follower, err := user_repo.Follow().FollowerList(ctx, req.UID, httputils.PageRequest{})
	if err != nil {
		return nil, err
	}
	_, following, err := user_repo.Follow().List(ctx, req.UID, httputils.PageRequest{})
	if err != nil {
		return nil, err
	}
	return &api.GetFollowResponse{
		IsFollow:  isFollow,
		Followers: follower,
		Following: following,
	}, nil
}

// Follow 关注用户
func (u *userSvc) Follow(ctx context.Context, req *api.FollowRequest) (*api.FollowResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	uid := user.UID
	if uid == req.UID {
		return nil, i18n.NewError(ctx, code.UserNotFollowSelf)
	}
	ok, err := user_repo.Follow().Find(ctx, uid, req.UID)
	if err != nil {
		return nil, err
	}
	if req.Unfollow {
		if ok == nil {
			return nil, i18n.NewError(ctx, code.UserNotFollow)
		}
		mutual, err := user_repo.Follow().Find(ctx, req.UID, uid)
		if err != nil {
			return nil, err
		}
		if mutual != nil {
			if err := user_repo.Follow().UpdateMutual(ctx, req.UID, uid, 0); err != nil {
				return nil, err
			}
		}
		return &api.FollowResponse{}, user_repo.Follow().Delete(ctx, uid, req.UID)
	}
	if ok != nil {
		return nil, i18n.NewError(ctx, code.UserExistFollow)
	}
	mutual, err := user_repo.Follow().Find(ctx, req.UID, uid)
	if err != nil {
		return nil, err
	}
	fo, err := u.UserInfo(ctx, req.UID)
	if err != nil {
		return nil, err
	}
	hf := &user_entity.HomeFollow{
		Uid:       uid,
		Username:  user.Username,
		Followuid: fo.UserID,
		Fusername: fo.Username,
		Bkname:    "",
		Status:    0,
		Dateline:  time.Now().Unix(),
	}
	if mutual != nil {
		hf.Mutual = 1
		if err := user_repo.Follow().UpdateMutual(ctx, req.UID, uid, 1); err != nil {
			return nil, err
		}
	}
	return &api.FollowResponse{}, user_repo.Follow().Save(ctx, hf)
}

// GetWebhook 获取webhook配置
func (u *userSvc) GetWebhook(ctx context.Context, req *api.GetWebhookRequest) (*api.GetWebhookResponse, error) {
	cfg, err := u.getConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetWebhookResponse{Token: cfg.Token}, nil
}

// RefreshWebhook 刷新webhook配置
func (u *userSvc) RefreshWebhook(ctx context.Context, req *api.RefreshWebhookRequest) (*api.RefreshWebhookResponse, error) {
	cfg, err := u.getConfig(ctx)
	if err != nil {
		return nil, err
	}
	cfg.GenToken()
	if err := user_repo.UserConfig().Update(ctx, cfg); err != nil {
		return nil, err
	}
	return &api.RefreshWebhookResponse{Token: cfg.Token}, nil
}

func (u *userSvc) getConfig(ctx context.Context) (*user_entity.UserConfig, error) {
	cfg, err := user_repo.UserConfig().FindByUserID(ctx, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &user_entity.UserConfig{
			Uid:        auth_svc.Auth().Get(ctx).UID,
			Notify:     &user_entity.Notify{},
			Createtime: time.Now().Unix(),
		}
		cfg.GenToken()
		if err := user_repo.UserConfig().Create(ctx, cfg); err != nil {
			return nil, err
		}
	} else if cfg.Token == "" {
		cfg.GenToken()
		if err := user_repo.UserConfig().Update(ctx, cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// GetConfig 获取用户配置
func (u *userSvc) GetConfig(ctx context.Context, req *api.GetConfigRequest) (*api.GetConfigResponse, error) {
	cfg, err := u.getConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &api.GetConfigResponse{Notify: cfg.Notify}, nil
}

// UpdateConfig 更新用户配置
func (u *userSvc) UpdateConfig(ctx context.Context, req *api.UpdateConfigRequest) (*api.UpdateConfigResponse, error) {
	cfg, err := u.getConfig(ctx)
	if err != nil {
		return nil, err
	}
	cfg.Notify = req.Notify
	if err := user_repo.UserConfig().Update(ctx, cfg); err != nil {
		return nil, err
	}
	return &api.UpdateConfigResponse{}, nil
}

// Search 搜索用户
func (u *userSvc) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	resp, err := user_repo.User().FindByPrefix(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	ret := make([]*api.InfoResponse, 0)
	for _, item := range resp {
		ret = append(ret, &api.InfoResponse{
			UserID:   item.ID,
			Username: item.Username,
			Avatar:   item.Avatar(),
		})
	}
	return &api.SearchResponse{
		Users: ret,
	}, nil
}

func (u *userSvc) Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	loginId, err := ctx.Cookie("login_id")
	if err != nil {
		return nil, err
	}
	token, err := ctx.Cookie("token")
	if err != nil {
		return nil, err
	}
	_, err = auth_svc.Auth().Logout(ctx, auth_svc.Auth().Get(ctx).UID, loginId, token)
	if err != nil {
		return nil, err
	}
	ctx.SetCookie("login_id", "", -1, "/", "", false, true)
	ctx.SetCookie("token", "", -1, "/", "", false, true)
	return &api.LogoutResponse{}, nil
}

// RefreshToken 刷新用户token
func (u *userSvc) RefreshToken(ctx *gin.Context, req *api.RefreshTokenRequest) (*api.RefreshTokenResponse, error) {
	loginId, err := ctx.Cookie("login_id")
	if err != nil {
		return nil, err
	}
	token, err := ctx.Cookie("token")
	if err != nil {
		return nil, err
	}

	// 获取token信息, 判断是否需要刷新
	m, err := auth_svc.Auth().GetLoginToken(ctx, auth_svc.Auth().Get(ctx).UID, loginId, token)
	if err != nil {
		return nil, err
	}

	if m.Updatetime+auth_svc.TokenAutoRegen < time.Now().Unix() {
		// 刷新token
		m, err = auth_svc.Auth().Refresh(ctx, auth_svc.Auth().Get(ctx).UID, loginId, token)
		if err != nil {
			return nil, err
		}
		// 设置cookie
		ctx.SetCookie("login_id", m.ID, auth_svc.TokenAuthMaxAge, "/", "", false, true)
		ctx.SetCookie("token", m.Token, auth_svc.TokenAuthMaxAge, "/", "", false, true)
	}

	return &api.RefreshTokenResponse{}, nil
}
