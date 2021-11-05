package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/http/dto/request"
	"github.com/scriptscat/scriptweb/internal/http/dto/respond"
	"github.com/scriptscat/scriptweb/internal/pkg/db"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptweb/internal/service"
	"github.com/scriptscat/scriptweb/pkg/middleware/token"
	"github.com/scriptscat/scriptweb/pkg/oauth"
	"github.com/scriptscat/scriptweb/pkg/utils"
	"gorm.io/datatypes"
)

type User struct {
	svc       service.User
	scriptSvc service2.Script
	client    *oauth.Client
}

func NewUser(user service.User, scriptSvc service2.Script) *User {
	return &User{
		svc:       user,
		scriptSvc: scriptSvc,
	}
}

func (u *User) info(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		// 判断是否登录
		tokenInfo, ok := authtoken(ctx)
		resp := gin.H{}
		uid := ctx.Param("uid")
		self := false
		if ok {
			if tokenInfo.Createtime+TokenAutoRegen < time.Now().Unix() {
				// 刷新token
				tokenString, err := token.GenToken(db.Cache, gin.H{
					"uid":      tokenInfo.Info["uid"],
					"username": tokenInfo.Info["username"],
					"email":    tokenInfo.Info["email"],
				})
				if err != nil {
					return err
				}
				ctx.SetCookie("token", tokenString, TokenAuthMaxAge, "/", "", false, true)
				resp["token"] = tokenString
			}
			if uid == "" {
				uid = tokenInfo.Info["uid"].(string)
				self = true
			}
		}
		if uid == "" {
			return errs.NewError(http.StatusBadRequest, 1000, "请指定用户")
		}
		dUid := utils.StringToInt64(uid)
		var userinfo *respond.User
		var err error
		if self {
			userinfo, err = u.svc.SelfInfo(dUid)
		} else {
			userinfo, err = u.svc.UserInfo(dUid)
		}
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

func (u *User) getwebhook(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
		ret, err := u.svc.GetUserWebhook(uid)
		if err != nil {
			return err
		}
		return gin.H{
			"token": ret,
		}
	})
}

func (u *User) regenwebhook(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
		ret, err := u.svc.RegenWebhook(uid)
		if err != nil {
			return err
		}
		return gin.H{
			"token": ret,
		}
	})
}

func (u *User) config(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
		ret, err := u.svc.GetUserConfig(uid)
		if err != nil {
			return err
		}
		return ret
	})
}

func (u *User) notify(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := userId(c)
		notify := datatypes.JSONMap{
			"score": c.PostForm("score") == "true",
		}
		return u.svc.SetUserNotifyConfig(uid, notify)
	})
}

func (u *User) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/user")
	rgg := rg.Group("", tokenAuth(false))
	rgg.GET("/info", u.info)
	rgg.GET("/info/:uid", u.info)
	rgg.GET("/scripts", u.scripts)
	rgg.GET("/scripts/:uid", u.scripts)
	rgg.GET("/avatar/:uid", u.avatar)

	rgg = rg.Group("/config", userAuth(true))
	rgg.GET("", u.config)
	rgg.PUT("/notifySvc", u.notify)

	rgg = rg.Group("/webhook", userAuth(true))
	rgg.GET("", u.getwebhook)
	rgg.PUT("", u.regenwebhook)
}
