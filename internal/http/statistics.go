package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/scriptscat/scriptweb/internal/domain/statistics/service"
)

type Statistics struct {
}

func NewStatistics(statisSvc service.Statistics, c *cron.Cron) *Statistics {
	ret := &Statistics{}
	c.AddFunc("0 0 1 * * *", func() {
		statisSvc.DealDay()
	})
	c.AddFunc("0 */20 * * * *", func() {
		statisSvc.DealRealtime()
	})
	return ret
}

func (s *Statistics) Registry(ctx context.Context, r *gin.Engine) {
	rg := r.Group("/api/v1/statistics")
	rg.GET("")
}
