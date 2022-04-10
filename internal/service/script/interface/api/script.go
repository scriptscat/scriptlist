package api

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
	"github.com/scriptscat/scriptlist/internal/infrastructure/config"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/interfaces/api/dto/request"
	"github.com/scriptscat/scriptlist/internal/pkg/cnt"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptlist/internal/service"
	service4 "github.com/scriptscat/scriptlist/internal/service/notify/service"
	"github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/vo"
	"github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/pkg/httputils"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/scriptscat/scriptlist/pkg/utils/diff"
	"github.com/sirupsen/logrus"
)

type Script struct {
	scriptSvc service2.Script
	scriptApp application.Script
	statisSvc service2.Statistics
	userSvc   service.User
	notifySvc service4.Sender
	watchSvc  application.ScriptWatch
}

func NewScript(svc service2.Script, app application.Script, statisSvc service2.Statistics, userSvc service.User, notify service4.Sender, watchSvc application.ScriptWatch, c *cron.Cron) *Script {
	// crontab 定时检查更新
	c.AddFunc("0 */6 * * *", func() {
		// 数据量大时可能要加入翻页，未来可能集群，要记得分布式处理
		list, err := app.FindSyncScript(request.AllPage)
		if err != nil {
			logrus.Errorf("Timing synchronization find script list: %v", err)
			return
		}
		for _, v := range list {
			if v.SyncMode != application.SyncModeAuto {
				continue
			}
			if err := svc.SyncScript(v.UserId, v.ID); err != nil {
				logrus.Errorf("Timing synchronization %v: %v", v.ID, err)
			}
		}
	})
	return &Script{
		scriptSvc: svc,
		scriptApp: app,
		statisSvc: statisSvc,
		userSvc:   userSvc,
		notifySvc: notify,
		watchSvc:  watchSvc,
	}
}

func (s *Script) Registry(ctx context.Context, r *gin.Engine) {
	tokenAuth := token.TokenAuth(false)
	r.Use(func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Writer.Status() != http.StatusNotFound {
			return
		}
		if strings.HasSuffix(ctx.Request.URL.Path, ".user.js") || strings.HasSuffix(ctx.Request.URL.Path, ".user.sub.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.downloadScript(ctx)
			}
		} else if strings.HasSuffix(ctx.Request.URL.Path, ".meta.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.getScriptMeta(ctx)
			}
		}
	})
	r.GET("/scripts/code/:id/*name", func(ctx *gin.Context) {
		if strings.HasSuffix(ctx.Request.URL.Path, ".user.js") || strings.HasSuffix(ctx.Request.URL.Path, ".user.sub.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.downloadScript(ctx)
			}
		} else if strings.HasSuffix(ctx.Request.URL.Path, ".meta.js") {
			tokenAuth(ctx)
			if !ctx.IsAborted() {
				s.getScriptMeta(ctx)
			}
		}
	})
	rg := r.Group("/api/v1/scripts")
	rg.GET("", s.list)
	rg.GET("/hot", s.hot)
	rg.POST("", token.UserAuth(true), s.add)
	rgg := rg.Group("/:script", token.UserAuth(true))
	rgg.PUT("", s.update)
	rgg.DELETE("", s.delete)
	rgg.PUT("/archive", s.archive)
	rgg.DELETE("/archive", s.unarchive)
	rgg.PUT("/code", s.updatecode)
	rgg.POST("/sync", s.sync)

	rgg = rg.Group("/:script", token.UserAuth(false))
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

	rgg = rg.Group("/:script/watch", token.UserAuth(true))
	rgg.GET("", s.iswatch)
	rgg.POST("", s.watch)
	rgg.DELETE("", s.unwatch)

	rg = r.Group("/api/v1/category")
	rg.GET("", s.category)

	r.Any("/api/v1/webhook/:uid", s.webhook)
}

// @Summary      删除脚本
// @Description  删除脚本
// @ID           script-delete
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId} [DELETE]
func (s *Script) delete(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, _ := token.UserInfo(ctx)
		id := utils.StringToInt64(ctx.Param("script"))
		return s.scriptApp.Delete(user, id)
	})
}

func (s *Script) hot(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		result, err := s.scriptApp.HotKeyword()
		if err != nil {
			return err
		}
		return result
	})
}

func (s *Script) iswatch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := token.UserId(c)
		level, err := s.watchSvc.IsWatch(script, uid)
		if err != nil {
			return err
		}
		return gin.H{"level": level}
	})
}

func (s *Script) watch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := token.UserId(c)
		return s.watchSvc.Watch(script, uid, application.ScriptWatchLevel(utils.StringToInt(c.PostForm("level"))))
	})
}

func (s *Script) unwatch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		script := utils.StringToInt64(c.Param("script"))
		uid, _ := token.UserId(c)
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
	httputils.Handle(c, func() interface{} {
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
			list, err := s.scriptApp.FindSyncPrefix(uid, "https://raw.githubusercontent.com/"+data.Repository.FullName)
			if err != nil {
				logrus.Errorf("Github hook FindSyncPrefix err: %v", err)
				return gin.H{
					"success": nil,
					"error":   nil,
				}
			}
			listtmp, err := s.scriptApp.FindSyncPrefix(uid, "https://github.com/"+data.Repository.FullName)
			if err != nil {
				logrus.Errorf("Github hook FindSyncPrefix err: %v", err)
				return gin.H{
					"success": nil,
					"error":   nil,
				}
			}
			list = append(list, listtmp...)
			success := make([]gin.H, 0)
			error := make([]gin.H, 0)
			for _, v := range list {
				if v.SyncMode != application.SyncModeAuto {
					continue
				}
				if err := s.scriptSvc.SyncScript(uid, v.ID); err != nil {
					logrus.Errorf("Github hook SyncScript: %v", err)
					error = append(error, gin.H{"id": v.ID, "name": v.Name, "err": err.Error()})
				} else {
					success = append(success, gin.H{"id": v.ID, "name": v.Name})
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
	//TODO: 暂时先允许刷吧
	uid, _ := token.UserId(ctx)
	id := utils.StringToInt64(ctx.Param("id"))
	version := ctx.Query("version")
	if id == 0 {
		id, version = s.parseScriptInfo(ctx.Request.URL.Path)
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 {
		return
	}
	if ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *vo.ScriptCode
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
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, GetStatisticsToken(ctx), true)
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Code)
}

func (s *Script) getScriptMeta(ctx *gin.Context) {
	uid, _ := token.UserId(ctx)
	id := utils.StringToInt64(ctx.Param("id"))
	if id == 0 {
		id, _ = s.parseScriptInfo(ctx.Request.URL.Path)
	}
	ua := ctx.GetHeader("User-Agent")
	if id == 0 || ua == "" {
		ctx.String(http.StatusNotFound, "脚本未找到")
		return
	}
	var code *vo.ScriptCode
	code, err := s.scriptSvc.GetLatestScriptCode(id, false)
	if err != nil {
		ctx.String(http.StatusBadGateway, err.Error())
		return
	}
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, GetStatisticsToken(ctx), false)
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Meta)
}

func (s *Script) list(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		page := &request.Pages{}
		if err := ctx.ShouldBind(page); err != nil {
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
		}, page)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) get(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		httputils.Handle(ctx, func() interface{} {
			uid, _ := token.UserId(ctx)
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
	httputils.Handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		page := &request.Pages{}
		if err := ctx.ShouldBind(page); err != nil {
			return err
		}
		list, err := s.scriptSvc.GetScriptCodeList(id, page)
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) versionsGet(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		httputils.Handle(ctx, func() interface{} {
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
	httputils.Handle(ctx, func() interface{} {
		list, err := s.scriptSvc.GetCategory()
		if err != nil {
			return err
		}
		return list
	})
}

func (s *Script) putScore(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, ok := token.UserInfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		id := utils.StringToInt64(ctx.Param("script"))
		score := &request.Score{}
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
	httputils.Handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		page := &request.Pages{}
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
	httputils.Handle(ctx, func() interface{} {
		uid, ok := token.UserId(ctx)
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
	httputils.Handle(ctx, func() interface{} {
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request.CreateScript{}
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

// @Summary      更新脚本
// @Description  更新脚本
// @ID           script-update
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @Success      200       {object}  request.UpdateScript
// @Failure      403
// @Router       /scripts/{scriptId} [PUT]
func (s *Script) update(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request.UpdateScript{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.scriptSvc.UpdateScript(uid, id, script)
	})
}

// @Summary      更新脚本代码
// @Description  更新脚本代码
// @ID           script-update-code
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @Success      200       {object}  request.UpdateScriptCode
// @Failure      403
// @Router       /scripts/{scriptId}/code [PUT]
func (s *Script) updatecode(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		script := &request.UpdateScriptCode{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		return s.scriptSvc.UpdateScriptCode(uid, id, script)
	})
}

func (s *Script) sync(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		id := utils.StringToInt64(ctx.Param("script"))
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		return s.scriptSvc.SyncScript(uid, id)
	})
}

func (s *Script) diff(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
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

// @Summary      设置脚本归档
// @Description  归档后无法再发issue、更新脚本
// @ID           script-archive
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId}/archive [PUT]
func (s *Script) archive(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		user, _ := token.UserInfo(c)
		id := utils.StringToInt64(c.Param("script"))
		return s.scriptApp.Archive(user, id, 1)
	})
}

// @Summary      取消脚本归档
// @Description  归档后无法再发issue、更新脚本
// @ID           script-un-archive
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId}/archive [DELETE]
func (s *Script) unarchive(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		user, _ := token.UserInfo(c)
		id := utils.StringToInt64(c.Param("script"))
		return s.scriptApp.Archive(user, id, 0)
	})
}