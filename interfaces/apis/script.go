package apis

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	request2 "github.com/scriptscat/scriptweb/interfaces/dto/request"
	"github.com/scriptscat/scriptweb/interfaces/dto/respond"
	"github.com/scriptscat/scriptweb/internal/application/service"
	"github.com/scriptscat/scriptweb/internal/pkg/config"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	jwt3 "github.com/scriptscat/scriptweb/pkg/middleware/jwt"
	"github.com/scriptscat/scriptweb/pkg/utils"
)

type Script struct {
	svc       service.Script
	statisSvc service.Statistics
}

func NewScript(svc service.Script, statisSvc service.Statistics) *Script {
	return &Script{
		svc:       svc,
		statisSvc: statisSvc,
	}
}

func (s *Script) Registry(r *gin.Engine) {
	jwtAuth := jwt3.Jwt([]byte(config.AppConfig.Jwt.Token), false, jwt3.WithExpired(JwtAuthMaxAge))
	r.Use(func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Writer.Status() != http.StatusNotFound {
			return
		}
		if strings.HasSuffix(ctx.Request.RequestURI, ".user.js") {
			jwtAuth(ctx)
			if !ctx.IsAborted() {
				s.downloadScript(ctx)
			}
		} else if strings.HasSuffix(ctx.Request.RequestURI, ".meta.js") {
			jwtAuth(ctx)
			if !ctx.IsAborted() {
				s.getScriptMeta(ctx)
			}
		}
	})
	rg := r.Group("/api/v1/scripts")
	rg.GET("", s.list)
	rg.GET("/:id", s.get(false))
	rg.GET("/:id/code", s.get(true))
	rg.GET("/:id/versions", s.versions)
	rg.GET("/:id/versions/:version", s.versionsGet(false))
	rg.GET("/:id/versions/:version/code", s.versionsGet(true))

	rgg := rg.Group("/:id/score", jwtAuth)
	rgg.GET("", s.scoreList)
	rgg.PUT("", s.putScore)
	rgg.GET("/self", s.selfScore)

	rg = r.Group("/api/v1/category")
	rg.GET("", s.category)

}

func (s *Script) parseScriptInfo(url string) (int64, string) {
	path := url[strings.LastIndex(url, "/")+1:]
	id, _ := strconv.ParseInt(strings.Split(path, ".")[0], 10, 64)
	if id <= 0 {
		return 0, ""
	}
	version := ""
	if strings.Index(url, "/version/") != -1 {
		version = url[strings.LastIndex(url, "/version/")+9:]
		version = version[:strings.Index(version, "/")]
	}
	return id, version
}

func (s *Script) downloadScript(ctx *gin.Context) {
	uid, _ := userinfo(ctx)
	id, version := s.parseScriptInfo(ctx.Request.RequestURI)
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *respond.ScriptCode
	var err error
	if version != "" {
		code, err = s.svc.GetScriptCodeByVersion(id, version, true)
	} else {
		code, err = s.svc.GetLatestScriptCode(id, true)
	}
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Code)
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, true)
}

func (s *Script) getScriptMeta(ctx *gin.Context) {
	uid, _ := userinfo(ctx)
	id, version := s.parseScriptInfo(ctx.Request.RequestURI)
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *respond.ScriptCode
	var err error
	if version != "" {
		code, err = s.svc.GetScriptCodeByVersion(id, version, false)
	} else {
		code, err = s.svc.GetLatestScriptCode(id, false)
	}
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Meta)
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, false)
}

func (s *Script) list(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		req := request2.Pages{}
		if err := ctx.ShouldBind(&req); err != nil {
			return err
		}
		categorys := make([]int64, 0)
		for _, v := range strings.Split(ctx.Query("category"), ",") {
			if v != "" {
				categorys = append(categorys, utils.StringToInt64(v))
			}
		}
		list, err := s.svc.GetScriptList(categorys, ctx.Query("domain"),
			ctx.Query("keyword"), ctx.Query("sort"), req)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) get(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handle(ctx, func() interface{} {
			id := utils.StringToInt64(ctx.Param("id"))
			ret, err := s.svc.GetScript(id, "", withcode)
			if err != nil {
				return err
			}
			return ret
		})
	}
}

func (s *Script) versions(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("id"))
		list, err := s.svc.GetScriptCodeList(id)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) versionsGet(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handle(ctx, func() interface{} {
			id := utils.StringToInt64(ctx.Param("id"))
			version := ctx.Param("version")
			code, err := s.svc.GetScript(id, version, withcode)
			if err != nil {
				return err
			}
			return code
		})
	}
}

func (s *Script) category(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		list, err := s.svc.GetCategory()
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) putScore(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		uid, ok := userinfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		id := utils.StringToInt64(ctx.Param("id"))
		score := &request2.Score{}
		if err := ctx.ShouldBind(score); err != nil {
			return err
		}
		return s.svc.AddScore(uid, id, score)
	})
}

func (s *Script) scoreList(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("id"))
		page := &request2.Pages{}
		if err := ctx.ShouldBind(page); err != nil {
			return err
		}
		list, err := s.svc.ScoreList(id, page)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) selfScore(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		uid, ok := userinfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		id := utils.StringToInt64(ctx.Param("id"))
		ret, err := s.svc.UserScore(uid, id)
		if err != nil {
			return err
		}
		return ret
	})
}
