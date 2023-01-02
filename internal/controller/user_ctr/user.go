package user_ctr

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/user"
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
	resp, err := user_svc.User().UserInfo(ctx, user_svc.Auth().Get(ctx).UID)
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
	m, err := user_svc.Auth().GetLoginToken(ctx, user_svc.Auth().Get(ctx).UID, loginId, token)
	if err != nil {
		return nil, err
	}
	if m.Updatetime+user_svc.TokenAutoRegen < time.Now().Unix() {
		// 刷新token
		m, err = user_svc.Auth().Refresh(ctx, user_svc.Auth().Get(ctx).UID, loginId, token)
		if err != nil {
			return nil, err
		}
		// 设置cookie
		gctx.SetCookie("login_id", m.ID, user_svc.TokenAuthMaxAge, "/", "", false, true)
		gctx.SetCookie("token", m.Token, user_svc.TokenAuthMaxAge, "/", "", false, true)
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
