package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	_ "github.com/scriptscat/scriptlist/docs"
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"github.com/scriptscat/scriptlist/internal/infrastructure/logs"
	token2 "github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service5 "github.com/scriptscat/scriptlist/internal/service"
	application2 "github.com/scriptscat/scriptlist/internal/service/issue/application"
	api2 "github.com/scriptscat/scriptlist/internal/service/issue/interface/api"
	service7 "github.com/scriptscat/scriptlist/internal/service/notify/service"
	service6 "github.com/scriptscat/scriptlist/internal/service/resource/service"
	"github.com/scriptscat/scriptlist/internal/service/safe/service"
	"github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/interface/api"
	service4 "github.com/scriptscat/scriptlist/internal/service/statistics/service"
	api3 "github.com/scriptscat/scriptlist/internal/service/user/interface/api"
	service2 "github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/internal/subscriber"
	"github.com/scriptscat/scriptlist/pkg/httputils"
	"github.com/scriptscat/scriptlist/pkg/oauth"
	pkgValidator "github.com/scriptscat/scriptlist/pkg/utils/validator"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Service interface {
	Registry(ctx context.Context, r *gin.Engine)
}

type Subscribe interface {
	Subscribe(ctx context.Context) error
}

func Registry(ctx context.Context, r *gin.Engine, registry ...Service) {
	for _, v := range registry {
		v.Registry(ctx, r)
	}
}

func Subscriber(ctx context.Context, sub ...Subscribe) {
	for _, v := range sub {
		v.Subscribe(ctx)
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
		logrus.Errorf("%s - %s: %+v", ctx.Request.RequestURI, ctx.ClientIP(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1001, "msg": "系统错误",
		})
	case *httputils.List:
		list := resp.(*httputils.List)
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

// StartApi 初始化路由
// Swagger spec:
// @title                       脚本猫列表
// @version                     1.0
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @BasePath                    /api/v1
func StartApi(db *persistence.Repositories) error {
	token2.TokenAuth = func(enforce bool) func(ctx *gin.Context) {
		return token2.Middleware(db.Cache, enforce, token2.WithExpired(token2.TokenAuthMaxAge))
	}

	ctx := context.Background()
	binding.Validator = pkgValidator.NewValidator()
	c := cron.New()
	userSvc := service2.NewUser(db.User.User, db.User.Follow)
	scriptApp := application.NewScript(db, db.Script.Script, db.Script.Code,
		db.Script.Category, db.Script.Statistics, c)
	statisSvc := service4.NewStatistics(db.Statistics.Statistics)
	scoreSvc := application.NewScore(db.Script.Script, db.Script.Score)
	rateSvc := service.NewRate(db.Safe.Rate)
	notifySvc := service7.NewSender(config.AppConfig.Email, config.AppConfig.EmailNotify)
	issueSvc := application2.NewIssue(db.Issue.Issue, db.Issue.IssueComment)
	issueWatchSvc := application2.NewWatch(db.Issue.IssueWatch)
	scriptWatchSvc := application.NewWatch(db.Script.ScriptWatch)

	script := service5.NewScript(userSvc,
		scriptApp,
		scoreSvc,
		statisSvc,
		rateSvc,
		c,
	)

	token2.UserAuth = func(enforce bool) func(ctx *gin.Context) {
		authHandler := token2.TokenAuth(enforce)
		return func(ctx *gin.Context) {
			authHandler(ctx)
			if !ctx.IsAborted() {
				if uid, ok := token2.UserId(ctx); ok {
					// NOTE:用户信息可以写入context
					if user, err := userSvc.UserInfo(uid); err != nil {
						handelResp(ctx, err)
						ctx.Abort()
					} else {
						ctx.Set(token2.Userentity, user)
					}
				}
			}
		}
	}

	statis := service5.NewStatistics(statisSvc, scriptApp)

	r := gin.New()
	if config.AppConfig.Mode == "debug" {
		r.Use(cors.Default())
	}
	r.Use(logs.GinLogger()...)

	Registry(ctx, r,
		api.NewScript(script, scriptApp, scoreSvc, statis, userSvc, notifySvc, scriptWatchSvc, c),
		NewLogin(oauth.NewClient(&config.AppConfig.OAuth), db),
		NewResource(service6.NewResource(db.Resource.Resource), rateSvc),
		api.NewStatistics(db, statisSvc, scriptApp, c),
		api3.NewUser(db, userSvc, script),
		api2.NewScriptIssue(scriptApp, userSvc, notifySvc, issueSvc, issueWatchSvc),
	)

	Subscriber(ctx,
		subscriber.NewScriptSubscriber(notifySvc, scriptWatchSvc, issueWatchSvc, issueSvc, scriptApp, userSvc),
	)

	c.Start()

	if config.AppConfig.Mode == "debug" {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	return r.Run(":" + strconv.Itoa(config.AppConfig.WebPort))
}
