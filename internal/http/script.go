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
	"github.com/robfig/cron/v3"
	service4 "github.com/scriptscat/scriptlist/internal/domain/notify/service"
	"github.com/scriptscat/scriptlist/internal/domain/script/repository"
	service3 "github.com/scriptscat/scriptlist/internal/domain/script/service"
	"github.com/scriptscat/scriptlist/internal/domain/user/service"
	request2 "github.com/scriptscat/scriptlist/internal/http/dto/request"
	"github.com/scriptscat/scriptlist/internal/http/dto/respond"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptlist/internal/service"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/scriptscat/scriptlist/pkg/utils/diff"
	"github.com/sirupsen/logrus"
)

type Script struct {
	scriptSvc service2.Script
	statisSvc service2.Statistics
	userSvc   service.User
	notifySvc service4.Sender
	watchSvc  service3.ScriptWatch
}

func NewScript(svc service2.Script, statisSvc service2.Statistics, userSvc service.User, notify service4.Sender, watchSvc service3.ScriptWatch, c *cron.Cron) *Script {
	// crontab 定时检查更新
	c.AddFunc("0 */6 * * *", func() {
		// 数据量大时可能要加入翻页，未来可能集群，要记得分布式处理
		list, err := svc.FindSyncScript(request2.AllPage)
		if err != nil {
			logrus.Errorf("Timing synchronization find script list: %v", err)
			return
		}
		for _, v := range list {
			if v.SyncMode != service3.SyncModeAuto {
				continue
			}
			if err := svc.SyncScript(v.UserId, v.ID); err != nil {
				logrus.Errorf("Timing synchronization %v: %v", v.ID, err)
			}
		}
	})
	return &Script{
		scriptSvc: svc,
		statisSvc: statisSvc,
		userSvc:   userSvc,
		notifySvc: notify,
		watchSvc:  watchSvc,
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
	r.GET("/scripts/code/:id/*name", func(ctx *gin.Context) {
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
	rgg := rg.Group("/:script", userAuth(true))
	rgg.PUT("", s.update)
	rgg.PUT("/code", s.updatecode)
	rgg.POST("/sync", s.sync)

	rgg = rg.Group("/:script", userAuth(false))
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

	rgg = rg.Group("/:script/watch", userAuth(true))
	rgg.GET("", s.iswatch)
	rgg.POST("", s.watch)
	rgg.DELETE("", s.unwatch)

	rg = r.Group("/api/v1/category")
	rg.GET("", s.category)

	r.Any("/api/v1/webhook/:uid", s.webhook)
}

func (s *Script) iswatch(c *gin.Context) {
	handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := userId(c)
		level, err := s.watchSvc.IsWatch(script, uid)
		if err != nil {
			return err
		}
		return gin.H{"level": level}
	})
}

func (s *Script) watch(c *gin.Context) {
	handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := userId(c)
		return s.watchSvc.Watch(script, uid, service3.ScriptWatchLevel(utils.StringToInt(c.PostForm("level"))))
	})
}

func (s *Script) unwatch(c *gin.Context) {
	handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := userId(c)
		return s.watchSvc.Unwatch(script, uid)
	})
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
			if data.Repository.FullName == "" {
				return errs.NewBadRequestError(1001, "仓库地址错误")
			}
			list, err := s.scriptSvc.FindSyncPrefix(uid, "https://raw.githubusercontent.com/"+data.Repository.FullName)
			if err != nil {
				logrus.Errorf("Github hook FindSyncPrefix err: %v", err)
				return gin.H{
					"success": nil,
					"error":   nil,
				}
			}
			listtmp, err := s.scriptSvc.FindSyncPrefix(uid, "https://github.com/"+data.Repository.FullName)
			if err != nil {
				logrus.Errorf("Github hook FindSyncPrefix err: %v", err)
				return gin.H{
					"success": nil,
					"error":   nil,
				}
			}
			list = append(list, listtmp...)
			var success []gin.H
			var error []gin.H
			for _, v := range list {
				if v.SyncMode != service3.SyncModeAuto {
					continue
				}
				if err := s.scriptSvc.SyncScript(uid, v.ID); err != nil {
					logrus.Errorf("Github hook SyncScript: %v", err)
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
	id := utils.StringToInt64(ctx.Param("id"))
	version := ctx.Query("version")
	if id == 0 {
		id, version = s.parseScriptInfo(ctx.Request.RequestURI)
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *respond.ScriptCode
	var err error
	if version != "" {
		code, err = s.scriptSvc.GetScriptCodeByVersion(id, version, true)
	} else {
		code, err = s.scriptSvc.GetLatestScriptCode(id, true)
	}
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Code)
}

func (s *Script) getScriptMeta(ctx *gin.Context) {
	uid, _ := userId(ctx)
	id := utils.StringToInt64(ctx.Param("id"))
	version := ctx.Query("version")
	if id == 0 {
		id, version = s.parseScriptInfo(ctx.Request.RequestURI)
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *respond.ScriptCode
	var err error
	if version != "" {
		code, err = s.scriptSvc.GetScriptCodeByVersion(id, version, false)
	} else {
		code, err = s.scriptSvc.GetLatestScriptCode(id, false)
	}
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, getStatisticsToken(ctx), false)
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Meta)
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
				id := utils.StringToInt64(v)
				if id > 0 {
					categorys = append(categorys, id)
				}
			}
		}
		list, err := s.scriptSvc.GetScriptList(&repository.SearchList{
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
			id := utils.StringToInt64(ctx.Param("script"))
			ret, err := s.scriptSvc.GetScript(id, "", withcode)
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
		id := utils.StringToInt64(ctx.Param("script"))
		list, err := s.scriptSvc.GetScriptCodeList(id)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) versionsGet(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handle(ctx, func() interface{} {
			id := utils.StringToInt64(ctx.Param("script"))
			version := ctx.Param("version")
			code, err := s.scriptSvc.GetScript(id, version, withcode)
			if err != nil {
				return err
			}
			return code
		})
	}
}

func (s *Script) category(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		list, err := s.scriptSvc.GetCategory()
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) putScore(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		user, ok := selfinfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		id := utils.StringToInt64(ctx.Param("script"))
		score := &request2.Score{}
		if err := ctx.ShouldBind(score); err != nil {
			return err
		}
		exist, err := s.scriptSvc.AddScore(user.UID, id, score)
		if err != nil {
			return err
		}
		if !exist {
			info, err := s.scriptSvc.GetScript(id, "", false)
			if err != nil {
				logrus.Errorf("GetScript: %v", err)
				return nil
			}
			cfg, err := s.userSvc.GetUserConfig(info.UserId)
			if err != nil {
				logrus.Errorf("GetUserConfig: %v", err)
				return nil
			}
			if n, ok := cfg.Notify[service.UserNotifyScore].(bool); ok && !n {
				return nil
			}
			sendUser, err := s.userSvc.SelfInfo(info.UserId)
			if err != nil {
				logrus.Errorf("SelfInfo: %v", err)
			} else {
				if err := s.notifySvc.SendEmail(sendUser.Email, "脚本有新的评分-"+info.Name,
					fmt.Sprintf("您的脚本【%s】有新的评分:<br/>%s:<br/>%s<br/><br/><a href='%s'>点我查看</a>或者复制链接:%s",
						info.Name, user.Username, score.Message,
						config.AppConfig.FrontendUrl+"script-show-page/"+strconv.FormatInt(info.ID, 10)+"/comment",
						config.AppConfig.FrontendUrl+"script-show-page/"+strconv.FormatInt(info.ID, 10)+"/comment",
					),
					"text/html"); err != nil {
					logrus.Errorf("sendemail: %v", err)
				}
			}
		}
		return nil
	})
}

func (s *Script) scoreList(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		page := &request2.Pages{}
		if err := ctx.ShouldBind(page); err != nil {
			return err
		}
		list, err := s.scriptSvc.ScoreList(id, page)
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
		id := utils.StringToInt64(ctx.Param("script"))
		ret, err := s.scriptSvc.UserScore(uid, id)
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
		ret, err := s.scriptSvc.CreateScript(uid, script)
		if err != nil {
			return err
		}
		return ret
	})
}

func (s *Script) update(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request2.UpdateScript{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.scriptSvc.UpdateScript(uid, id, script)
	})
}

func (s *Script) updatecode(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request2.UpdateScriptCode{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.scriptSvc.UpdateScriptCode(uid, id, script)
	})
}

func (s *Script) sync(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := userId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		return s.scriptSvc.SyncScript(uid, id)
	})
}

func (s *Script) diff(c *gin.Context) {
	handle(c, func() interface{} {
		id := utils.StringToInt64(c.Param("script"))
		v1 := c.Param("v1")
		v2 := c.Param("v2")
		s1, err := s.scriptSvc.GetScriptCodeByVersion(id, v1, true)
		if err != nil {
			return err
		}
		s2, err := s.scriptSvc.GetScriptCodeByVersion(id, v2, true)
		if err != nil {
			return err
		}
		return gin.H{
			"diff": diff.Diff(s1.Code, s2.Code),
		}
	})
}
