package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	jwt2 "github.com/golang-jwt/jwt"
	"github.com/golang/glog"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptweb/interfaces/dto/respond"
	"github.com/scriptscat/scriptweb/internal/application/service"
	repository3 "github.com/scriptscat/scriptweb/internal/domain/script/repository"
	service3 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	repository2 "github.com/scriptscat/scriptweb/internal/domain/statistics/repository"
	service4 "github.com/scriptscat/scriptweb/internal/domain/statistics/service"
	"github.com/scriptscat/scriptweb/internal/domain/user/repository"
	service2 "github.com/scriptscat/scriptweb/internal/domain/user/service"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/oauth"
	"github.com/scriptscat/scriptweb/pkg/utils"
	pkgValidator "github.com/scriptscat/scriptweb/pkg/utils/validator"
)

type Service interface {
	Registry(r *gin.Engine)
}

func Registry(r *gin.Engine, registry ...Service) {
	for _, v := range registry {
		v.Registry(r)
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
	switch resp.(type) {
	case *errs.RespondError:
		err := resp.(*errs.RespondError)
		ctx.String(err.Status, err.Msg)
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

func userinfo(ctx *gin.Context) (int64, bool) {
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
	binding.Validator = pkgValidator.NewValidator()
	c := cron.New()
	userSvc := service2.NewUser(repository.NewUser())
	scriptSvc := service3.NewScript(repository3.NewScript(), repository3.NewCode(), repository3.NewCategory(), repository3.NewStatistics(), c)
	script := service.NewScript(userSvc,
		scriptSvc,
		service3.NewScore(repository3.NewScore()),
		service4.NewStatistics(repository2.NewStatistics()),
	)

	statis := service.NewStatistical(service4.NewStatistics(repository2.NewStatistics()), scriptSvc)
	user := service.NewUser(userSvc)

	r := gin.Default()
	Registry(r,
		NewScript(script, statis),
		NewLogin(oauth.NewClient(&config.AppConfig.OAuth), config.AppConfig.Jwt.Token),
		NewUser(user, script),
	)
	c.Start()
	return r.Run(":" + strconv.Itoa(config.AppConfig.WebPort))
}
