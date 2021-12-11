package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptlist/internal/domain/script/entity"
	service2 "github.com/scriptscat/scriptlist/internal/domain/script/service"
	"github.com/scriptscat/scriptlist/internal/domain/statistics/service"
	"github.com/scriptscat/scriptlist/internal/pkg/cache"
	"github.com/scriptscat/scriptlist/internal/pkg/csrf"
	"github.com/scriptscat/scriptlist/internal/pkg/db"
	"github.com/scriptscat/scriptlist/internal/pkg/errs"
	"github.com/scriptscat/scriptlist/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Statistics struct {
	statisSvc service.Statistics
	scriptSvc service2.Script
}

func NewStatistics(statisSvc service.Statistics, scriptSvc service2.Script, c *cron.Cron) *Statistics {
	ret := &Statistics{
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

func (s *Statistics) scriptStatistics(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		user, _ := userId(ctx)
		id := utils.StringToInt64(ctx.Param("id"))
		script, err := s.scriptSvc.Info(id)
		if err != nil {
			return err
		}
		if script.UserId != user {
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

func (s *Statistics) scriptRealtime(ctx *gin.Context) {
	handle(ctx, func() interface{} {
		user, _ := userId(ctx)
		id := utils.StringToInt64(ctx.Param("id"))
		script, err := s.scriptSvc.Info(id)
		if err != nil {
			return err
		}
		if script.UserId != user {
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
	uid, _ := userId(c)
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
	if ok, _ := db.Cache.Has("csrf:" + c.GetHeader("X-CSRF-Token")); !ok {
		return
	}
	_ = db.Cache.Del("csrf:" + c.GetHeader("X-CSRF-Token"))

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
	_ = s.statisSvc.Record(id, code.ID, uid, c.ClientIP(), ua, getStatisticsToken(c), true)
}

func (s *Statistics) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/statistics")
	rgg := rg.Group("/script/:id", userAuth(true))
	rgg.GET("", s.scriptStatistics)
	rgg.GET("/realtime", s.scriptRealtime)

	rgg = rg.Group("/script/:id")
	rgg.GET("/csrf", csrf.CsrfMiddleware(), func(c *gin.Context) {
		db.Cache.Set("csrf:"+csrf.Token(c), "true", cache.WithTTL(time.Hour))
		c.JSON(http.StatusOK, gin.H{"csrf": csrf.Token(c)})
	})
	rgg.POST("/download", csrf.CsrfMiddleware(), s.download)
}
