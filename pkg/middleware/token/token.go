package token

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const Userinfo = "userinfo"
const AuthToken = "auth_token"

func Middleware(cache cache.Cache, enforce bool, option ...Option) gin.HandlerFunc {
	opts := &options{}
	for _, o := range option {
		o(opts)
	}
	return func(ctx *gin.Context) {
		token, _ := ctx.Cookie("token")
		if token == "" {
			token = ctx.GetHeader("Authorization")
			if token == "" {
				token = ctx.PostForm("token")
			} else {
				auths := strings.Split(token, " ")
				if len(auths) != 2 {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code": 1000, "msg": "token string is empty",
					})
					return
				} else {
					token = auths[1]
				}
			}
		}
		if token == "" {
			if enforce {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": 1001, "msg": "token string is empty",
				})
			}
			return
		}
		tokenInfo := &Token{}
		err := cache.Get("token:token:"+token, tokenInfo)
		if err != nil {
			for _, v := range opts.authFailed {
				if err := v(tokenInfo); err != nil {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code": 1002, "msg": err.Error(),
					})
					return
				}
			}
			if tokenInfo.Info != nil {
				ctx.Set(Userinfo, tokenInfo.Info)
				ctx.Set(AuthToken, tokenInfo)
				return
			}
			if enforce {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": 1002, "msg": err.Error(),
				})
			}
			return
		}
		for _, v := range opts.tokenHandlerFunc {
			if err := v(tokenInfo); err != nil {
				if enforce {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code": 1002, "msg": err.Error(),
					})
				}
				return
			}
		}
		ctx.Set(Userinfo, tokenInfo.Info)
		ctx.Set(AuthToken, tokenInfo)
	}
}

func GenToken(c cache.Cache, info gin.H) (string, error) {
	tokenInfo := &Token{
		Info:       info,
		Token:      utils.RandString(16, 2),
		Createtime: time.Now().Unix(),
	}
	if err := c.Set("token:token:"+tokenInfo.Token, tokenInfo, cache.WithTTL(time.Hour*24*30)); err != nil {
		return "", err
	}
	return tokenInfo.Token, nil
}
