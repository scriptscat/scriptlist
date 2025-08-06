package user_ctr

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/codfrm/cago/configs"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/user_svc"
	"github.com/scriptscat/scriptlist/pkg/oauth"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// CurrentUser 获取当前登录的用户信息
func (u *User) CurrentUser(ctx *gin.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	resp, err := user_svc.User().UserInfo(ctx, auth_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	return &api.CurrentUserResponse{InfoResponse: resp}, nil
}

// RefreshToken 刷新用户token
func (u *User) RefreshToken(ctx *gin.Context, req *api.RefreshTokenRequest) (*api.RefreshTokenResponse, error) {
	return user_svc.User().RefreshToken(ctx, req)
}

// Info 获取指定用户信息
func (u *User) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return user_svc.User().UserInfo(ctx, req.UID)
}

// Avatar 获取用户头像
func (u *User) Avatar() gin.HandlerFunc {
	config := &oauth.Config{}
	if err := configs.Default().Scan(context.Background(), "oauth.bbs", &config); err != nil {
		config.ServerUrl = "https://bbs.tampermonkey.net.cn"
	}
	// https://bbs.tampermonkey.net.cn/uc_server/avatar.php?uid=13895&size=middle
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		resp, err := http.Get(config.ServerUrl + "/uc_server/avatar.php?uid=" + uid + "&size=middle")
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ct := http.DetectContentType(b)
		if strings.Contains(ct, "image") {
			ctx.Writer.Header().Set("content-type", ct)
		} else {
			// svg图片
			ctx.Writer.Header().Set("content-type", "image/svg+xml")
		}
		// 缓存
		ctx.Writer.Header().Set("Cache-Control", "max-age=86400")
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

// Search 搜索用户
func (u *User) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResponse, error) {
	return user_svc.User().Search(ctx, req)
}

// Logout 登出账户
func (u *User) Logout(ctx *gin.Context, req *api.LogoutRequest) (*api.LogoutResponse, error) {
	return user_svc.User().Logout(ctx, req)
}
