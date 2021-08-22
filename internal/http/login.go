package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/oauth"
)

const JwtAuthMaxAge = 432000
const JwtAutoRenew = 259200

type Login struct {
	client   *oauth.Client
	jwtToken string
}

func NewLogin(client *oauth.Client, jwtToken string) *Login {
	return &Login{
		client:   client,
		jwtToken: jwtToken,
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
		tokenString, err := jwt3.GenJwt([]byte(l.jwtToken), jwt.MapClaims{
			"uid":      userResp.User.Uid,
			"username": userResp.User.Username,
			"email":    userResp.User.Email,
		})
		if err != nil {
			return err
		}
		ctx.SetCookie("auth", tokenString, JwtAuthMaxAge, "/", "", false, true)
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
