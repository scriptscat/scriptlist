package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptlist/internal/infrastructure/middleware/token"
	"github.com/scriptscat/scriptlist/internal/infrastructure/persistence"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/csrf"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	service2 "github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/domain/entity"
	"github.com/scriptscat/scriptlist/internal/service/statistics/service"
	"github.com/scriptscat/scriptlist/pkg/httputils"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Statistics struct {
	db        *persistence.Repositories
	statisSvc service.Statistics
	scriptSvc service2.Script
}

func NewStatistics(db *persistence.Repositories, statisSvc service.Statistics, scriptSvc service2.Script, c *cron.Cron) *Statistics {
	ret := &Statistics{
		db:        db,
		statisSvc: statisSvc,
		scriptSvc: scriptSvc,
	}
	go statisSvc.Deal()
	c.AddFunc("0 */8 * * *", func() {
		if err := statisSvc.Deal(); err != nil {
			logrus.Errorf("statistics deal error: %v", err)
		}
	})
	return ret
}

// @Summary      脚本统计数据
// @Description  脚本统计数据,只允许管理员、脚本管理员访问
// @ID           script-statistics
// @Tags         script-statistics
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @Success      200
// @Failure      403
// @Router       /statistics/script/{scriptId} [GET]
func (s *Statistics) scriptStatistics(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, _ := token.UserInfo(ctx)
		id := utils.StringToInt64(ctx.Param("id"))
		script, err := s.scriptSvc.Info(id)
		if err != nil {
			return err
		}
		if script.UserID != user.UID || !user.IsAdmin.IsAdmin() {
			return errs.NewError(http.StatusForbidden, 1000, "没有权限访问")
		}
		now := time.Now().Add(-time.Hour * 24)
		lastweekly := time.Now().Add(-time.Hour * 24 * 7)
		return gin.H{
			"download": gin.H{
				"uv":            s.ignoreError(s.statisSvc.DownloadUv(script.ID, 7, now)),
				"uv-lastweekly": s.ignoreError(s.statisSvc.DownloadUv(script.ID, 7, lastweekly)),
				"pv":            s.ignoreError(s.statisSvc.DownloadPv(script.ID, 30, now)),
				"realtime":      s.ignoreError(s.statisSvc.RealtimeDownload(script.ID)),
			},
			"update": gin.H{
				"uv":            s.ignoreError(s.statisSvc.UpdateUv(script.ID, 7, now)),
				"uv-lastweekly": s.ignoreError(s.statisSvc.UpdateUv(script.ID, 7, lastweekly)),
				"pv":            s.ignoreError(s.statisSvc.UpdatePv(script.ID, 30, now)),
				"realtime":      s.ignoreError(s.statisSvc.RealtimeUpdate(script.ID)),
			},
			"member": gin.H{
				"num": s.ignoreError(s.statisSvc.WeeklyUv(script.ID)),
			},
		}
	})
}

// @Summary      脚本实时统计
// @Description  脚本实时统计,只允许管理员、脚本管理员访问
// @ID           script-realtime
// @Tags         script-statistics
// @Security     BearerAuth
// @param        scriptId  path  integer  true  "脚本id"
// @Success      200
// @Failure      403
// @Router       /statistics/script/{scriptId}/realtime [GET]
func (s *Statistics) scriptRealtime(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		user, _ := token.UserInfo(ctx)
		id := utils.StringToInt64(ctx.Param("id"))
		script, err := s.scriptSvc.Info(id)
		if err != nil {
			return err
		}
		if script.UserID != user.UID || !user.IsAdmin.IsAdmin() {
			return errs.NewError(http.StatusForbidden, 1000, "没有权限访问")
		}
		return gin.H{
			"download": s.ignoreError(s.statisSvc.RealtimeDownload(script.ID)),
			"update":   s.ignoreError(s.statisSvc.RealtimeUpdate(script.ID)),
		}
	})
}

func (s *Statistics) ignoreError(args interface{}, err error) interface{} {
	return args
}

func (s *Statistics) download(c *gin.Context) {
	uid, _ := token.UserId(c)
	id := utils.StringToInt64(c.PostForm("id"))
	version, ua, _csrf := c.PostForm("version"), c.GetHeader("User-Agent"), c.PostForm("_csrf")
	if id == 0 || ua == "" || _csrf == "" {
		return
	}
	//h := hmac.New(func() hash.Hash {
	//	return sha1.New()
	//}, []byte(csrf.Secret))
	//h.Write([]byte(csrf.Token(c)))
	//b, _ := base64.StdEncoding.DecodeString(_csrf)
	//if !hmac.Equal(h.Sum(nil), b) {
	//	return
	//}
	//TODO: 屏蔽刷下载
	return
	if ok, _ := s.db.Cache.Has("csrf:" + c.GetHeader("X-CSRF-Token")); ok {
		return
	}
	_ = s.db.Cache.Set("csrf:"+c.GetHeader("X-CSRF-Token"), "1", cache.WithTTL(time.Hour))

	var code *entity.ScriptCode
	var err error
	if version != "" {
		code, err = s.scriptSvc.GetScriptVersion(id, version)
	} else {
		code, err = s.scriptSvc.GetLatestVersion(id)
	}
	if err != nil {
		return
	}
	go func() {
		ok, err := s.statisSvc.Record(id, code.ID, uid, c.ClientIP(), ua, GetStatisticsToken(c), true)
		if err != nil {
			logrus.Errorf("statis download record %v: %v", id, err)
			return
		}
		if ok {
			if err := s.scriptSvc.Download(id); err != nil {
				logrus.Errorf("download record %v: %v", id, err)
			}
		}
	}()
}

func (s *Statistics) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/statistics")
	rgg := rg.Group("/script/:id", token.UserAuth(true))
	rgg.GET("", s.scriptStatistics)
	rgg.GET("/realtime", s.scriptRealtime)

	rgg = rg.Group("/script/:id")
	rgg.GET("/csrf", csrf.CsrfMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"csrf": csrf.Token(c)})
	})
	rgg.POST("/download", csrf.CsrfMiddleware(), s.download)
}
