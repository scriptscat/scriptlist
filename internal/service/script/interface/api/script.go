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
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/repository"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/vo"
	service3 "github.com/scriptscat/scriptlist/internal/service/statistics/service"
	"github.com/scriptscat/scriptlist/internal/service/user/service"
	"github.com/scriptscat/scriptlist/pkg/httputils"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Script struct {
	scriptSvc service2.Script
	scriptApp application.Script
	scoreApp  application.Score
	statisSvc service2.Statistics
	userSvc   service.User
	notifySvc service4.Sender
	watchSvc  application.ScriptWatch
}

func NewScript(svc service2.Script, app application.Script, score application.Score, statisSvc service2.Statistics, userSvc service.User, notify service4.Sender, watchSvc application.ScriptWatch, c *cron.Cron) *Script {
	ret := &Script{
		scriptSvc: svc,
		scriptApp: app,
		scoreApp:  score,
		statisSvc: statisSvc,
		userSvc:   userSvc,
		notifySvc: notify,
		watchSvc:  watchSvc,
	}
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
			if err := svc.SyncScript(v.UserID, v.ID); err != nil {
				logrus.Errorf("Timing synchronization %v: %v", v.ID, err)
			}
		}
	})
	return ret
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
	// @require  https://scriptcat.org/lib/117/1.0.1/qwe1.js
	r.GET("/lib/:id/:version/:name", func(ctx *gin.Context) {
		tokenAuth(ctx)
		if !ctx.IsAborted() {
			s.downloadScript(ctx)
		}
	})
	rg := r.Group("/api/v1/scripts")
	r.GET("/api/v1/admin/refresh-search", token.UserAuth(true), s.refreshSearch)
	rg.GET("", s.list)
	rg.GET("/hot", s.hot)
	rg.POST("", token.UserAuth(true), s.add)
	rgg := rg.Group("/:script", token.UserAuth(true))
	rgg.PUT("", s.update)
	rgg.DELETE("", s.delete)
	rgg.GET("/setting", s.setting)
	rgg.PUT("/archive", s.archive)
	rgg.DELETE("/archive", s.unarchive)
	rgg.POST("/admin", s.admin)
	rgg.PUT("/code", s.updatecode)
	rgg.POST("/sync", s.sync)

	rgg = rg.Group("/:script", token.UserAuth(false))
	rgg.GET("", s.get(false))
	rgg.GET("/code", s.get(true))
	//rgg.GET("/diff/:v1/:v2", s.diff)
	rggg := rgg.Group("/versions")
	rggg.GET("", s.versions)
	rggg.GET("/:version", s.versionsGet(false))
	rggg.GET("/:version/code", s.versionsGet(true))

	rggg = rgg.Group("/score")
	rggg.GET("", s.scoreList)
	rggg.PUT("", s.putScore)
	rggg.GET("/self", s.selfScore)
	rggg.DELETE("/:scoreId", token.UserAuth(true), s.delScore)

	rgg = rg.Group("/:script/watch", token.UserAuth(true))
	rgg.GET("", s.iswatch)
	rgg.POST("", s.watch)
	rgg.DELETE("", s.unwatch)

	rg = r.Group("/api/v1/category")
	rg.GET("", s.category)

	r.Any("/api/v1/webhook/:uid", s.webhook)
}

// @Summary      删除脚本
// @Description  删除脚本,只有脚本管理员与超级版主以上的管理员可以操作
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
		level := struct {
			Level application.ScriptWatchLevel `json:"level"`
		}{}
		if err := c.BindJSON(&level); err != nil {
			return err
		}
		return s.watchSvc.Watch(script, uid, level.Level)
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
	if version == "" {
		version = ctx.Param("version")
	}
	if version == "latest" {
		version = ""
	}
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
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, GetStatisticsToken(ctx), service3.DOWNLOAD_STATISTICS)
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
	_ = s.statisSvc.Record(id, code.ID, uid, ctx.ClientIP(), ua, GetStatisticsToken(ctx), service3.UPDATE_STATISTICS)
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.WriteString(code.Meta)
}

// @Summary      脚本列表
// @Description  脚本列表
// @ID           script-list
// @Tags         script
// @Security     BearerAuth
// @param        page      query     integer  false  "页码"
// @param        count     query     integer  false  "页大小"
// @param        category  query     string   false  "分类id以','分割"
// @param        keyword   query     string   false  "搜索关键字"
// @param        sort      query     string   false  "排序方式"
// @Success      200       {object}  vo.Script
// @Failure      403
// @Router       /scripts [GET]
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

// @Summary      脚本信息
// @Description  脚本信息
// @ID           script-info
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @Success      200       {object}  vo.ScriptInfo
// @Failure      403
// @Router       /scripts/{scriptId} [GET]

// @Summary      脚本代码
// @Description  脚本代码
// @ID           script-code
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @Success      200       {object}  vo.ScriptInfo
// @Failure      403
// @Router       /scripts/{scriptId}/code [GET]
func (s *Script) get(withcode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		httputils.Handle(ctx, func() interface{} {
			uid, _ := token.UserId(ctx)
			id := utils.StringToInt64(ctx.Param("script"))
			ua := ctx.GetHeader("User-Agent")
			if id == 0 || ua == "" {
				return errs.ErrScriptNotFound
			}
			ret, err := s.scriptSvc.GetScript(id, "", withcode)
			if err != nil {
				return err
			}
			if ret.Unwell == 1 && uid == 0 {
				return errs.NewBadRequestError(1000, "该脚本含有不适内容，登录后设置才能访问")
			}
			if !withcode {
				_ = s.statisSvc.Record(id, ret.Script.ID, uid, ctx.ClientIP(), ua, GetStatisticsToken(ctx), service3.VIEW_STATISTICS)
			}
			return ret
		})
	}
}

// @Summary      版本列表
// @Description  版本列表
// @ID           script-version
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @Success      200       {object}  vo.ScriptCode
// @Failure      403
// @Router       /scripts/{scriptId}/versions [GET]
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

// @Summary      获取脚本指定版本
// @Description  获取脚本指定版本
// @ID           script-version-info
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @param        version   path      string   true  "版本号"
// @Success      200       {object}  vo.ScriptCode
// @Failure      403
// @Router       /scripts/{scriptId}/versions/{version} [GET]

// @Summary      获取脚本指定版本代码
// @Description  获取脚本指定版本代码
// @ID           script-version-code
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @param        version   path      string   true  "版本号"
// @Success      200       {object}  vo.ScriptCode
// @Failure      403
// @Router       /scripts/{scriptId}/versions/{version}/code [GET]
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

var _ = entity.ScriptCategoryList{}

// @Summary      脚本分类列表
// @Description  脚本分类列表
// @ID           script-category
// @Tags         script-category
// @Security     BearerAuth
// @Success      200  {object}  entity.ScriptCategoryList
// @Failure      403
// @Router       /category [GET]
func (s *Script) category(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		list, err := s.scriptSvc.GetCategory()
		if err != nil {
			return err
		}
		return list
	})
}

// @Summary      发表评分
// @Description  发表评分
// @ID           script-score-put
// @Tags         script-score
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @param        score     formData  integer  true   "分数"
// @param        message   formData  string   false  "评论"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId}/score [PUT]
func (s *Script) putScore(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, ok := token.UserInfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		if user.EmailStatus == 0 {
			return errs.ErrUserEmailNotVerified
		}
		id := utils.StringToInt64(ctx.Param("script"))
		score := &request.Score{}
		if err := ctx.ShouldBind(score); err != nil {
			return err
		}
		exist, err := s.scoreApp.AddScore(user.UID, id, score)
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

// @Summary      评分列表
// @Description  评分列表
// @ID           script-score-list
// @Tags         script-score
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @param        page      query     integer  false  "页码"
// @param        count     query     integer  false  "页大小"
// @Success      200       {object}  vo.ScriptScore
// @Failure      403
// @Router       /scripts/{scriptId}/score [GET]
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

// @Summary      自己的评分
// @Description  自己的评分
// @ID           script-score-self
// @Tags         script-score
// @Security     BearerAuth
// @param        scriptId  path      integer  true   "脚本id"
// @Success      200       {object}  vo.ScriptScore
// @Failure      403
// @Router       /scripts/{scriptId}/score/self [GET]
func (s *Script) selfScore(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		uid, ok := token.UserId(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		id := utils.StringToInt64(ctx.Param("script"))
		ret, err := s.scoreApp.UserScore(uid, id)
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary      删除评分
// @Description  删除评分, 只有管理员才能删除
// @ID           script-score-delete
// @Tags         script-score
// @Security     BearerAuth
// @param        scriptId  path      integer  true   "脚本id"
// @param        scoreId   path  integer  true  "评分id"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId}/score/{scoreId} [DELETE]
func (s *Script) delScore(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, ok := token.UserInfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		if !user.IsAdmin.IsAdmin() {
			return errs.NewError(http.StatusForbidden, 1000, "只有管理员才能删除评分")
		}
		scriptId := utils.StringToInt64(ctx.Param("script"))
		scoreId := utils.StringToInt64(ctx.Param("scoreId"))
		err := s.scoreApp.Delete(scriptId, scoreId)
		if err != nil {
			return err
		}
		return nil
	})
}

// @Summary      创建脚本
// @Description  创建脚本
// @ID           script-create
// @Tags         script
// @Security     BearerAuth
// @param        content      formData  string  true   "脚本详细描述"
// @param        code         formData  string  true   "脚本代码"
// @param        name         formData  string  false  "库的名字,当脚本类型为脚本调用库时必填"
// @param        description  formData  string  false  "库的描述,当脚本类型为脚本调用库时必填"
// @param        definition   formData  string  false  "库的定义文件,当脚本类型为脚本调用库时必填"
// @param        type         formData  int     true   "脚本类型 1 用户脚本 2 订阅脚本 3 脚本调用库"
// @param        public       formData  int     true   "公开类型 1 公开 2 半公开"
// @param        unwell       formData  int     true   "不适内容"
// @param        changelog    formData  string  true   "更新日志"
// @Success      200
// @Failure      403
// @Router       /scripts [POST]
func (s *Script) add(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, ok := token.UserInfo(ctx)
		if !ok {
			return errs.ErrNotLogin
		}
		if user.EmailStatus == 0 {
			return errs.ErrUserEmailNotVerified
		}
		script := &request.CreateScript{}
		if err := ctx.ShouldBind(script); err != nil {
			return err
		}
		ret, err := s.scriptSvc.CreateScript(user.UID, script)
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary      更新脚本配置
// @Description  更新脚本配置
// @ID           script-update
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path      integer  true  "脚本id"
// @param        name            formData  string   false  "库的名字,当脚本类型为脚本调用库时必填"
// @param        description     formData  string   false  "库的描述,当脚本类型为脚本调用库时必填"
// @param        sync_url        formData  string   false  "代码同步url"
// @param        content_url     formData  string   false  "详细描述同步url"
// @param        definition_url  formData  string   false  "定义文件同步url"
// @param        sync_mode       formData  int      false  "同步模式"
// @Success      200
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

//func (s *Script) diff(c *gin.Context) {
//	httputils.Handle(c, func() interface{} {
//		id := utils.StringToInt64(c.Param("script"))
//		v1 := c.Param("v1")
//		v2 := c.Param("v2")
//		s1, err := s.scriptSvc.GetScriptCodeByVersion(id, v1, true)
//		if err != nil {
//			return err
//		}
//		s2, err := s.scriptSvc.GetScriptCodeByVersion(id, v2, true)
//		if err != nil {
//			return err
//		}
//		return gin.H{
//			"diff": diff.Diff(s1.Code, s2.Code),
//		}
//	})
//}

// @Summary      设置脚本归档
// @Description  归档后无法再发issue、更新脚本,只有脚本管理员与超级版主以上的管理员可以操作
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
// @ID           script-unarchive
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

// @Summary      管理员管理
// @Description  管理员管理,允许管理员设置
// @ID           script-admin
// @Tags         script
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @param        action    formData  string   true  "操作:unwell 设置为不适脚本,unpublic 设置为非公开脚本,archive设置为归档脚本,delete 删除脚本"
// @Success      200
// @Failure      403
// @Router       /scripts/{scriptId}/admin [POST]
func (s *Script) admin(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		user, _ := token.UserInfo(c)
		if !user.IsAdmin.IsAdmin() {
			return errs.NewError(http.StatusForbidden, 1000, "没有权限")
		}
		id := utils.StringToInt64(c.Param("script"))
		switch c.PostForm("action") {
		case "unwell":
			return s.scriptApp.Unwell(user, id)
		case "unpublic":
			return s.scriptApp.Unpublic(user, id)
		case "archive":
			return s.scriptApp.Archive(user, id, 1)
		case "delete":
			return s.scriptApp.Delete(user, id)
		}
		return errs.NewError(http.StatusBadRequest, 1000, "操作错误")
	})
}

// @Summary      获取脚本设置
// @Description  获取脚本设置
// @ID           script-setting
// @Tags         script
// @Security     BearerAuth
// @param        scriptId        path      integer  true   "脚本id"
// @Success      200  {object}  vo.ScriptSetting
// @Failure      403
// @Router       /scripts/{scriptId}/setting [GET]
func (s *Script) setting(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := token.UserId(c)
		id := utils.StringToInt64(c.Param("script"))
		script, err := s.scriptApp.Info(id)
		if err != nil {
			return err
		}
		if script.UserID != uid {
			return errs.ErrScriptForbidden
		}
		return &vo.ScriptSetting{
			SyncUrl:       script.SyncUrl,
			ContentUrl:    script.ContentUrl,
			DefinitionUrl: script.DefinitionUrl,
			SyncMode:      script.SyncMode,
		}
	})
}

// @Summary      刷新搜索缓存
// @Description  刷新搜索缓存
// @ID           script-search-refresh
// @Tags         script
// @Security     BearerAuth
// @Success      200       {object}  vo.ScriptSetting
// @Failure      403
// @Router       /admin/refresh-search [GET]
func (s *Script) refreshSearch(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		user, _ := token.UserInfo(c)
		if user.IsAdmin != 1 {
			return nil
		}
		if err := s.scriptApp.RefreshGoFound(); err != nil {
			return err
		}
		return gin.H{"ok": "ok"}
	})
}
