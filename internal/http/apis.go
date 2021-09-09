package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	jwt2 "github.com/golang-jwt/jwt"
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	repository5 "github.com/scriptscat/scriptweb/internal/domain/resource/repository"
	service6 "github.com/scriptscat/scriptweb/internal/domain/resource/service"
	repository4 "github.com/scriptscat/scriptweb/internal/domain/safe/repository"
	"github.com/scriptscat/scriptweb/internal/domain/safe/service"
	repository3 "github.com/scriptscat/scriptweb/internal/domain/script/repository"
	service3 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	repository2 "github.com/scriptscat/scriptweb/internal/domain/statistics/repository"
	service4 "github.com/scriptscat/scriptweb/internal/domain/statistics/service"
	"github.com/scriptscat/scriptweb/internal/domain/user/repository"
	service2 "github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/http/dto/respond"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	service5 "github.com/scriptscat/scriptweb/internal/service"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/oauth"
	"github.com/scriptscat/scriptweb/pkg/utils"
	pkgValidator "github.com/scriptscat/scriptweb/pkg/utils/validator"
)

const CheckUserInfo = "user.CheckUserInfo"

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
		glog.Errorf("%s - %s: %v", ctx.Request.RequestURI, ctx.ClientIP(), err)
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

func userId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(jwt3.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(jwt2.MapClaims)["uid"].(string)), true
}

func isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(jwt3.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(jwt2.MapClaims)["uid"].(string)), false
}

func jwttoken(ctx *gin.Context) (jwt2.MapClaims, *jwt2.Token, bool) {
	u, ok := ctx.Get(jwt3.Userinfo)
	if !ok {
		return nil, nil, false
	}
	t, ok := ctx.Get(jwt3.JwtToken)
	if !ok {
		return nil, nil, false
	}
	return u.(jwt2.MapClaims), t.(*jwt2.Token), true
}

func StartApi() error {
	ctx := context.Background()
	binding.Validator = pkgValidator.NewValidator()
	c := cron.New()
	userSvc := service2.NewUser(repository.NewUser())
	scriptSvc := service3.NewScript(repository3.NewScript(), repository3.NewCode(), repository3.NewCategory(), repository3.NewStatistics(), c)
	rateSvc := service.NewRate(repository4.NewRate())
	script := service5.NewScript(userSvc,
		scriptSvc,
		service3.NewScore(repository3.NewScore()),
		service4.NewStatistics(repository2.NewStatistics()),
		rateSvc,
	)

	statis := service5.NewStatistical(service4.NewStatistics(repository2.NewStatistics()), scriptSvc)
	user := service5.NewUser(userSvc)
	userApi := NewUser(user, script)
	ctx = context.WithValue(ctx, CheckUserInfo, userApi.CheckUserInfo())

	r := gin.Default()
	Registry(ctx, r,
		NewScript(script, statis),
		NewLogin(oauth.NewClient(&config.AppConfig.OAuth), config.AppConfig.Jwt.Token),
		NewResource(service6.NewResource(repository5.NewResource()), rateSvc),
		userApi,
	)
	c.Start()
	return r.Run(":" + strconv.Itoa(config.AppConfig.WebPort))
}
