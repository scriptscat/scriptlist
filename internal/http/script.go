package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/scriptweb/internal/domain/script/repository"
	service3 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	"github.com/scriptscat/scriptweb/internal/domain/user/service"
	request2 "github.com/scriptscat/scriptweb/internal/http/dto/request"
	"github.com/scriptscat/scriptweb/internal/http/dto/respond"
	"github.com/scriptscat/scriptweb/internal/pkg/cnt"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptweb/internal/service"
	"github.com/scriptscat/scriptweb/pkg/utils"
	"github.com/scriptscat/scriptweb/pkg/utils/diff"
)

type Script struct {
	svc       service2.Script
	statisSvc service2.Statistics
	userSvc   service.User
}

func NewScript(svc service2.Script, statisSvc service2.Statistics, userSvc service.User) *Script {
	return &Script{
		svc:       svc,
		statisSvc: statisSvc,
		userSvc:   userSvc,
	}
}

func (s *Script) Registry(ctx context.Context, r *gin.Engine) {
	tokenAuth := tokenAuth(false)
	r.Use(func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Writer.Status() != http.StatusNotFound {
			return
		}
		if strings.HasSuffix(ctx.Request.RequestURI, ".user.js") || strings.HasSuffix(ctx.Request.RequestURI, ".user.sub.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.downloadScript(ctx)
			}
		} else if strings.HasSuffix(ctx.Request.RequestURI, ".meta.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.getScriptMeta(ctx)
			}
		}
	})
	rg := r.Group("/api/v1/scripts")
	rg.GET("", s.list)
	rg.POST("", userAuth(true), s.add)
	rgg := rg.Group("/:id", userAuth(true))
	rgg.PUT("", s.update)
	rgg.PUT("/code", s.updatecode)
	rgg.POST("/sync", s.sync)

	rgg = rg.Group("/:id", tokenAuth)
	rgg.GET("", s.get(false))
	rgg.GET("/code", s.get(true))
	rgg.GET("/diff/:v1/:v2", s.diff)
	rggg := rgg.Group("/versions")
	rggg.GET("", s.versions)
	rggg.GET("/:version", s.versionsGet(false))
	rggg.GET("/:version/code", s.versionsGet(true))

	rggg = rgg.Group("/score")
	rggg.GET("", s.scoreList)
	rggg.PUT("", s.putScore)
	rggg.GET("/self", s.selfScore)

	rg = r.Group("/api/v1/category")
	rg.GET("", s.category)

	r.Any("/api/v1/webhook/:uid", s.webhook)

}

type githubWebhook struct {
	Hook struct {
		Type string `json:"type"`
	} `json:"hook"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func (s *Script) webhook(c *gin.Context) {
	handle(c, func() interface{} {
		uid := utils.StringToInt64(c.Param("uid"))
		secret, err := s.userSvc.GetUserWebhook(uid)
		if err != nil {
			return err
		}
		ua := c.GetHeader("User-Agent")
		if strings.Index(ua, "GitHub") != -1 {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				return err
			}
			hash := hmac.New(sha256.New, []byte(secret))
			if _, err := hash.Write(b); err != nil {
				return err
			}
			if fmt.Sprintf("sha256=%x", hash.Sum(nil)) != c.GetHeader("X-Hub-Signature-256") {
				return errs.NewBadRequestError(1000, "密钥校验错误")
			}
			// 处理github
			data := &githubWebhook{}
			if err := json.Unmarshal(b, data); err != nil {
				return err
			}
			if data.Hook.Type != "Repository" {
				return errs.NewBadRequestError(10001, "只能识别data.hook.type=repository")
			}
			list, err := s.svc.FindSyncPrefix(uid, "https://raw.githubusercontent.com/"+data.Repository.FullName)
			if err != nil {
				return gin.H{
					"success": nil,
					"error":   nil,
				}
			}
			success := []gin.H{}
			error := []gin.H{}
			for _, v := range list {
				if v.SyncMode != service3.SYNC_MODE_AUTO {
					continue
				}
				if err := s.svc.SyncScript(uid, v.ID); err != nil {
					error = append(error, gin.H{"id": v.ID, "name": v.Name, "err": err.Error()})
				} else {
					success = append(success, gin.H{"id": v.ID, "name": v.Name, "err": err.Error()})
				}
			}
			return gin.H{
				"success": success,
				"error":   error,
			}
		}
		return nil
	})
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
	uid, _ := userId(ctx)
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
	uid, _ := userId(ctx)
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
		list, err := s.svc.GetScriptList(&repository.SearchList{
			Category: categorys,
			Domain:   ctx.Query("domain"),
			Sort:     ctx.Query("sort"),
			Status:   cnt.ACTIVE,
			Keyword:  ctx.Query("keyword"),
		}, req)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) get(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handle(ctx, func() interface{} {
			uid, _ := userId(ctx)
			id := utils.StringToInt64(ctx.Param("id"))
			ret, err := s.svc.GetScript(id, "", withcode)
			if err != nil {
				return err
			}
			ret.IsManager = uid == ret.UID
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
		uid, ok := userId(ctx)
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
		uid, ok := userId(ctx)
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

func (s *Script) add(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request2.CreateScript{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		ret, err := s.svc.CreateScript(uid, script)
		if err != nil {
			return err
		}
		return ret
	})
}

func (s *Script) update(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("id"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request2.UpdateScript{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.svc.UpdateScript(uid, id, script)
	})
}

func (s *Script) updatecode(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("id"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request2.UpdateScriptCode{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.svc.UpdateScriptCode(uid, id, script)
	})
}

func (s *Script) sync(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("id"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		return s.svc.SyncScript(uid, id)
	})
}

func (s *Script) diff(c *gin.Context) {
	handle(c, func() interface{} {
		id := utils.StringToInt64(c.Param("id"))
		v1 := c.Param("v1")
		v2 := c.Param("v2")
		s1, err := s.svc.GetScriptCodeByVersion(id, v1, true)
		if err != nil {
			return err
		}
		s2, err := s.svc.GetScriptCodeByVersion(id, v2, true)
		if err != nil {
			return err
		}
		return gin.H{
			"diff": diff.Diff(s1.Code, s2.Code),
		}
	})
}
