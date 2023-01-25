package user_ctr

import (
	"net/http"
	"strings"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/auth"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Debug() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := auth_svc.Auth().Login(ctx.Request.Context(), 1)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		// 设置cookie
		ctx.SetCookie("login_id", token.ID, auth_svc.TokenAuthMaxAge, "/", "", false, true)
		ctx.SetCookie("token", token.Token, auth_svc.TokenAuthMaxAge, "/", "", false, true)
	}
}

// OAuthCallback 第三方登录回调
func (a *Auth) OAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &api.OAuthCallbackRequest{}
		if err := c.BindQuery(req); err != nil {
			httputils.HandleResp(c, err)
			return
		}
		resp, err := auth_svc.Auth().OAuthCallback(c.Request.Context(), req)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		token, err := auth_svc.Auth().Login(c.Request.Context(), resp.UID)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		// 设置cookie
		c.SetCookie("login_id", token.ID, auth_svc.TokenAuthMaxAge, "/", "", false, true)
		c.SetCookie("token", token.Token, auth_svc.TokenAuthMaxAge, "/", "", false, true)
		// 重定向
		path := configs.Default().String("website.url")
		if strings.HasPrefix(resp.RedirectUri, "/") {
			path = path + resp.RedirectUri
		} else {
			path = path + "/" + resp.RedirectUri
		}
		c.Redirect(http.StatusFound, path)
	}
}

func (a *Auth) Middleware(force bool) gin.HandlerFunc {
	return auth_svc.Auth().Middleware(force)
}
