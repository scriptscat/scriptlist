package apis

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt"
	"github.com/scriptscat/scriptweb/internal/application/service"
	"github.com/scriptscat/scriptweb/internal/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/oauth"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type User struct {
	svc       service.User
	scriptSvc service.Script
	client    *oauth.Client
	jwtToken  string
}

func NewUser(user service.User, scriptSvc service.Script) *User {
	return &User{
		svc:       user,
		scriptSvc: scriptSvc,
		jwtToken:  config.AppConfig.Jwt.Token,
	}
}

func (u *User) info(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		// 判断是否登录
		user, token, ok := jwttoken(ctx)
		resp := gin.H{}
		uid := ctx.Param("uid")
		if ok {
			if t, ok := token.Header["time"]; ok {
				if i, ok := t.(float64); !(ok && int64(i)+JwtAutoRenew > time.Now().Unix()) {
					// 刷新token
					tokenString, err := jwt3.GenJwt([]byte(u.jwtToken), jwt2.MapClaims{
						"uid":      user["uid"],
						"username": user["username"],
						"email":    user["email"],
					})
					if err != nil {
						return err
					}
					ctx.SetCookie("auth", tokenString, JwtAuthMaxAge, "/", "", false, true)
					resp["auth"] = tokenString
				}
			}
			if uid == "" {
				uid = user["uid"].(string)
			}
		}
		if uid == "" {
			return errs.NewError(http.StatusBadRequest, 1000, "请指定用户")
		}
		dUid := utils.StringToInt64(uid)
		userinfo, err := u.svc.UserInfo(dUid)
		if err != nil {
			return err
		}
		resp["user"] = userinfo

		return resp
	})
}

func (u *User) scripts(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		sUid := ctx.Param("uid")
		uid, ok := userinfo(ctx)
		self := false
		if sUid == "" {
			if !ok {
				return errs.NewError(http.StatusBadRequest, 1000, "请指定用户")
			}
			self = true
		} else {
			uid = utils.StringToInt64(sUid)
		}
		page := request.Pages{}
		if err := ctx.ShouldBind(&page); err != nil {
			return err
		}
		ret, err := u.scriptSvc.GetUserScript(uid, self, page)
		if err != nil {
			return err
		}
		return ret
	})
}

func (u *User) Registry(r *gin.Engine) {
	rg := r.Group("/api/v1/user", jwt3.Jwt([]byte(u.jwtToken), false, jwt3.WithExpired(JwtAuthMaxAge)))
	rg.GET("/info", u.info)
	rg.GET("/info/:uid", u.info)
	rg.GET("/scripts", u.scripts)
	rg.GET("/scripts/:uid", u.scripts)

}
