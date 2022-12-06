package user

import (
	"net/http"
	"strings"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/middleware/sessions"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/service/user"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Router(r *mux.Router) error {
	if configs.Default().Debug && configs.Default().Env == configs.DEV {
		r.GET("/login/debug", a.Debug())
	}
	r.GET("/login/oauth", a.OAuthCallback())
	return nil
}

func (a *Auth) Debug() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置session
		session := sessions.Ctx(ctx)
		session.Set("uid", int64(1))
		_ = session.Save()
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
		resp, err := user.Auth().OAuthCallback(c.Request.Context(), req)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		// 设置session
		session := sessions.Ctx(c)
		session.Set("uid", resp.UID)
		if err := session.Save(); err != nil {
			httputils.HandleResp(c, err)
			return
		}
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
	return user.Auth().Middleware(force)
}
