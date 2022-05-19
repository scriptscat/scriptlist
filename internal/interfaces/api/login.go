package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/oauth"
)

type Login struct {
	client *oauth.Client
	db     *persistence.Repositories
}

func NewLogin(client *oauth.Client, db *persistence.Repositories) *Login {
	return &Login{
		client: client,
		db:     db,
	}
}

func (l *Login) oauth(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		code := ctx.Query("code")
		if code == "" {
			return errs.NewError(http.StatusBadRequest, 1001, "code不能为空")
		}
		resp, err := l.client.RequestAccessToken(code)
		if err != nil {
			return err
		}
		//请求用户信息,写cookie
		userResp, err := l.client.RequestUser(resp.AccessToken)
		if err != nil {
			return err
		}
		tokenString, err := token.GenToken(l.db.Cache, gin.H{
			"uid":      userResp.User.Uid,
			"username": userResp.User.Username,
			"email":    userResp.User.Email,
		})
		if err != nil {
			return err
		}
		ctx.SetCookie("token", tokenString, token.TokenAuthMaxAge, "/", "", false, true)
		if uri := ctx.Query("redirect_uri"); uri != "" {
			ctx.Redirect(http.StatusFound, uri)
			return nil
		}
		return gin.H{
			"token": tokenString,
		}
	})
}

func (l *Login) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/login")
	rg.GET("/oauth", l.oauth)
}
