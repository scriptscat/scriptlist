package user_ctr

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/user_svc"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// CurrentUser 获取当前登录的用户信息
func (u *User) CurrentUser(gctx *gin.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	ctx := gctx.Request.Context()
	resp, err := user_svc.User().UserInfo(ctx, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	loginId, err := gctx.Cookie("login_id")
	if err != nil {
		return nil, err
	}
	token, err := gctx.Cookie("token")
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
		gctx.SetCookie("login_id", m.ID, auth_svc.TokenAuthMaxAge, "/", "", false, true)
		gctx.SetCookie("token", m.Token, auth_svc.TokenAuthMaxAge, "/", "", false, true)
	}
	return &api.CurrentUserResponse{InfoResponse: resp}, nil
}

// Info 获取指定用户信息
func (u *User) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return user_svc.User().UserInfo(ctx, req.UID)
}

// Avatar 获取用户头像
func (u *User) Avatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		resp, err := http.Get("https://bbs.tampermonkey.net.cn/uc_server/avatar.php?uid=" + uid + "&size=middle")
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()
		ctx.Writer.Header().Set("content-type", "image/jpeg")
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		_, _ = ctx.Writer.Write(b)
	}
}

// Script 用户脚本列表
func (u *User) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	return user_svc.User().Script(ctx, req)
}

// GetFollow 获取用户关注信息
func (u *User) GetFollow(ctx context.Context, req *api.GetFollowRequest) (*api.GetFollowResponse, error) {
	return user_svc.User().GetFollow(ctx, req)
}

// Follow 关注用户
func (u *User) Follow(ctx context.Context, req *api.FollowRequest) (*api.FollowResponse, error) {
	return user_svc.User().Follow(ctx, req)
}

// GetWebhook 获取webhook配置
func (u *User) GetWebhook(ctx context.Context, req *api.GetWebhookRequest) (*api.GetWebhookResponse, error) {
	return user_svc.User().GetWebhook(ctx, req)
}

// RefreshWebhook 刷新webhook配置
func (u *User) RefreshWebhook(ctx context.Context, req *api.RefreshWebhookRequest) (*api.RefreshWebhookResponse, error) {
	return user_svc.User().RefreshWebhook(ctx, req)
}

// GetConfig 获取用户配置
func (u *User) GetConfig(ctx context.Context, req *api.GetConfigRequest) (*api.GetConfigResponse, error) {
	return user_svc.User().GetConfig(ctx, req)
}

// UpdateConfig 更新用户配置
func (u *User) UpdateConfig(ctx context.Context, req *api.UpdateConfigRequest) (*api.UpdateConfigResponse, error) {
	return user_svc.User().UpdateConfig(ctx, req)
}
