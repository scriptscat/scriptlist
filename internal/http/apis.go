package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	service7 "github.com/scriptscat/scriptlist/internal/domain/notify/service"
	repository5 "github.com/scriptscat/scriptlist/internal/domain/resource/repository"
	service6 "github.com/scriptscat/scriptlist/internal/domain/resource/service"
	repository4 "github.com/scriptscat/scriptlist/internal/domain/safe/repository"
	"github.com/scriptscat/scriptlist/internal/domain/safe/service"
	repository3 "github.com/scriptscat/scriptlist/internal/domain/script/repository"
	service3 "github.com/scriptscat/scriptlist/internal/domain/script/service"
	repository2 "github.com/scriptscat/scriptlist/internal/domain/statistics/repository"
	service4 "github.com/scriptscat/scriptlist/internal/domain/statistics/service"
	"github.com/scriptscat/scriptlist/internal/domain/user/repository"
	service2 "github.com/scriptscat/scriptlist/internal/domain/user/service"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service5 "github.com/scriptscat/scriptlist/internal/service"
	"github.com/scriptscat/scriptlist/pkg/middleware/token"
	"github.com/scriptscat/scriptlist/pkg/oauth"
	pkgValidator "github.com/scriptscat/scriptlist/pkg/utils/validator"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Registry(ctx context.Context, r *gin.Engine)
}

func Registry(ctx context.Context, r *gin.Engine, registry ...Service) {
	for _, v := range registry {
		v.Registry(ctx, r)
	}
}

func handle(ctx *gin.Context, f func() interface{}) {
	resp := f()
	if resp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok",
		})
		return
	}
	handelResp(ctx, resp)
}

func handelResp(ctx *gin.Context, resp interface{}) {
	switch resp.(type) {
	case *errs.JsonRespondError:
		err := resp.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := resp.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := resp.(error)
		logrus.Errorf("%s - %s: %v", ctx.Request.RequestURI, ctx.ClientIP(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	case *respond.List:
		list := resp.(*respond.List)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "list": list.List, "total": list.Total,
		})
	case string:

	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "data": resp,
		})
	}
}

var tokenAuth func(enforce bool) func(ctx *gin.Context)
var userAuth func(enforce bool) func(ctx *gin.Context)

func StartApi() error {
	tokenAuth = func(enforce bool) func(ctx *gin.Context) {
		return token.Middleware(db.Cache, enforce, token.WithExpired(TokenAuthMaxAge))
	}

	ctx := context.Background()
	binding.Validator = pkgValidator.NewValidator()
	c := cron.New()
	userSvc := service2.NewUser(repository.NewUser())
	scriptSvc := service3.NewScript(repository3.NewScript(), repository3.NewCode(), repository3.NewCategory(), repository3.NewStatistics())
	statisSvc := service4.NewStatistics(repository2.NewStatistics())
	scoreSvc := service3.NewScore(repository3.NewScore())
	rateSvc := service.NewRate(repository4.NewRate())
	notifySvc := service7.NewSender(config.AppConfig.Email)
	script := service5.NewScript(userSvc,
		scriptSvc,
		scoreSvc,
		statisSvc,
		rateSvc,
		c,
	)

	userAuth = func(enforce bool) func(ctx *gin.Context) {
		authHandler := tokenAuth(enforce)
		return func(ctx *gin.Context) {
			authHandler(ctx)
			if !ctx.IsAborted() {
				if uid, ok := userId(ctx); ok {
					// NOTE:用户信息可以写入context
					if user, err := userSvc.UserInfo(uid); err != nil {
						handelResp(ctx, err)
						ctx.Abort()
					} else {
						ctx.Set(SelfInfo, user)
					}
				}
			}
		}
	}

	statis := service5.NewStatistics(statisSvc, scriptSvc)

	r := gin.Default()
	Registry(ctx, r,
		NewScript(script, statis, userSvc, notifySvc, c),
		NewLogin(oauth.NewClient(&config.AppConfig.OAuth)),
		NewResource(service6.NewResource(repository5.NewResource()), rateSvc),
		NewStatistics(statisSvc, scriptSvc, c),
		NewUser(userSvc, script),
	)
	c.Start()
	return r.Run(":" + strconv.Itoa(config.AppConfig.WebPort))
}
