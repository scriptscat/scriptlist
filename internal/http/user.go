package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt"
	"github.com/scriptscat/scriptweb/internal/http/dto/request"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptweb/internal/service"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/oauth"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type User struct {
	svc       service2.User
	scriptSvc service2.Script
	client    *oauth.Client
	jwtToken  string
}

func NewUser(user service2.User, scriptSvc service2.Script) *User {
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
		currentUid, ok := userId(ctx)
		self := false
		uid := utils.StringToInt64(sUid)
		if uid == 0 {
			if !ok {
				return errs.NewError(http.StatusBadRequest, 1000, "请指定用户")
			}
			uid = currentUid
			self = true
		}
		//page := request.Pages{}
		//if err := ctx.ShouldBind(&page); err != nil {
		//	return err
		//}
		ret, err := u.scriptSvc.GetUserScript(uid, self, request.AllPage)
		if err != nil {
			return err
		}
		return ret
	})
}

func (u *User) avatar(ctx *gin.Context) {
	uid := ctx.Param("uid")
	resp, err := http.Get("https://bbs.tampermonkey.net.cn/uc_server/avatar.php?uid=" + uid + "&size=middle")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	ctx.Writer.Header().Set("content-type", "image/jpeg")
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Writer.Write(b)
}

// CheckUserInfo 校验用户信息
func (u *User) CheckUserInfo() gin.HandlerFunc {
	jwtAuth := jwt3.Jwt([]byte(config.AppConfig.Jwt.Token), true, jwt3.WithExpired(JwtAuthMaxAge))
	return func(c *gin.Context) {
		jwtAuth(c)
		if c.IsAborted() {
			return
		}
		uid, ok := userId(c)
		if !ok {
			c.JSON(http.StatusForbidden, errs.ErrNotLogin)
			c.Abort()
			return
		}
		_, err := u.svc.UserInfo(uid)
		if err != nil {
			c.JSON(http.StatusForbidden, err)
			c.Abort()
			return
		}
	}
}

func (u *User) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/user", jwt3.Jwt([]byte(u.jwtToken), false, jwt3.WithExpired(JwtAuthMaxAge)))
	rg.GET("/info", u.info)
	rg.GET("/info/:uid", u.info)
	rg.GET("/scripts", u.scripts)
	rg.GET("/scripts/:uid", u.scripts)
	rg.GET("/avatar/:uid", u.avatar)
}
