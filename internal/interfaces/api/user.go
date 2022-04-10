package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptlist/internal/service"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/user/domain/vo"
	"github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/pkg/oauth"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"gorm.io/datatypes"
)

type User struct {
	db        *persistence.Repositories
	svc       service.User
	scriptSvc service2.Script
	client    *oauth.Client
}

func NewUser(db *persistence.Repositories, user service.User, scriptSvc service2.Script) *User {
	return &User{
		db:        db,
		svc:       user,
		scriptSvc: scriptSvc,
	}
}

func (u *User) info(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		// 判断是否登录
		tokenInfo, ok := token.Authtoken(ctx)
		resp := gin.H{}
		uid := ctx.Param("uid")
		self := false
		if ok {
			if tokenInfo.Createtime+TokenAutoRegen < time.Now().Unix() {
				// 刷新token
				tokenString, err := token.GenToken(u.db.Cache, gin.H{
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
		var userinfo *vo.User
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
		_, num, _ := u.svc.FollowList(userinfo.UID, &request.Pages{})
		_, num2, _ := u.svc.FollowerList(userinfo.UID, &request.Pages{})
		resp["follow"] = gin.H{
			"following": num,
			"followers": num2,
		}
		return resp
	})
}

func (u *User) scripts(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		sUid := ctx.Param("uid")
		currentUid, ok := token.UserId(ctx)
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

		list, err := u.scriptSvc.GetScriptList(&repository.SearchList{
			Uid: uid, Self: self,
			Domain:  ctx.Query("domain"),
			Sort:    ctx.Query("sort"),
			Status:  cnt.ACTIVE,
			Keyword: ctx.Query("keyword"),
		}, request.AllPage)
		if err != nil {
			return err
		}
		return list
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
		uid, _ := token.UserId(c)
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
		uid, _ := token.UserId(c)
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
		uid, _ := token.UserId(c)
		ret, err := u.svc.GetUserConfig(uid)
		if err != nil {
			return err
		}
		return ret
	})
}

func (u *User) notify(c *gin.Context) {
	handle(c, func() interface{} {
		uid, _ := token.UserId(c)
		notify := datatypes.JSONMap{
			service.UserNotifyScore:              c.PostForm(service.UserNotifyScore) == "true",
			service.UserNotifyCreateScript:       c.PostForm(service.UserNotifyCreateScript) == "true",
			service.UserNotifyScriptUpdate:       c.PostForm(service.UserNotifyScriptUpdate) == "true",
			service.UserNotifyScriptIssue:        c.PostForm(service.UserNotifyScriptIssue) == "true",
			service.UserNotifyScriptIssueComment: c.PostForm(service.UserNotifyScriptIssueComment) == "true",
			service.UserNotifyAt:                 c.PostForm(service.UserNotifyAt) == "true",
		}
		return u.svc.SetUserNotifyConfig(uid, notify)
	})
}

func (u *User) isfollow(c *gin.Context) {
	handle(c, func() interface{} {
		follow := utils.StringToInt64(c.Param("follow"))
		uid, _ := token.UserId(c)
		is, err := u.svc.IsFollow(uid, follow)
		if err != nil {
			return err
		}
		return is
	})
}

func (u *User) follow(c *gin.Context) {
	handle(c, func() interface{} {
		follow := utils.StringToInt64(c.Param("follow"))
		uid, _ := token.UserId(c)
		return u.svc.Follow(uid, follow)
	})
}

func (u *User) unfollow(c *gin.Context) {
	handle(c, func() interface{} {
		follow := utils.StringToInt64(c.Param("follow"))
		uid, _ := token.UserId(c)
		return u.svc.Unfollow(uid, follow)
	})
}

func (u *User) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/user")
	rgg := rg.Group("", token.TokenAuth(false))
	rgg.GET("/info", u.info)
	rgg.GET("/info/:uid", u.info)
	rgg.GET("/scripts", u.scripts)
	rgg.GET("/scripts/:uid", u.scripts)
	rgg.GET("/avatar/:uid", u.avatar)

	rgg = rg.Group("/follow/:follow", token.TokenAuth(true))
	rgg.GET("", u.isfollow)
	rgg.POST("", u.follow)
	rgg.DELETE("", u.unfollow)

	rgg = rg.Group("/config", token.UserAuth(true))
	rgg.GET("", u.config)
	rgg.PUT("/notify", u.notify)

	rgg = rg.Group("/webhook", token.UserAuth(true))
	rgg.GET("", u.getwebhook)
	rgg.PUT("", u.regenwebhook)
}
