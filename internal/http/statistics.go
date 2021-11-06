package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	service2 "github.com/scriptscat/scriptweb/internal/domain/script/service"
	"github.com/scriptscat/scriptweb/internal/domain/statistics/service"
	"github.com/scriptscat/scriptweb/internal/pkg/errs"
	"github.com/scriptscat/scriptweb/pkg/utils"
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
	c.AddFunc("0 0 */8 * * *", func() {
		statisSvc.Deal()
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

func (s *Statistics) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/statistics")
	rgg := rg.Group("/script", userAuth(true))
	rgg.GET("/:id", s.scriptStatistics)
	rgg.GET("/:id/realtime", s.scriptRealtime)

}
