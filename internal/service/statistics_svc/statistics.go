package statistics_svc

import (
	"context"
	"strconv"
	"time"

	"github.com/cago-frame/cago/pkg/logger"
	"github.com/cago-frame/cago/pkg/utils"
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// StatisticsSvc 统计平台
type StatisticsSvc interface {
	// ScriptRecord 脚本数据统计
	ScriptRecord(ctx context.Context, data *producer.ScriptStatisticsMsg) error
	// GetStatisticsToken 获取统计token
	GetStatisticsToken(ctx *gin.Context) string
	// Script 脚本统计数据
	Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error)
	// ScriptRealtime 脚本实时统计数据
	ScriptRealtime(ctx context.Context, req *api.ScriptRealtimeRequest) (*api.ScriptRealtimeResponse, error)
	// Middleware 中间件
	Middleware() gin.HandlerFunc
}

type statisticsSvc struct {
}

var defaultStatistics = &statisticsSvc{}

func Statistics() StatisticsSvc {
	return defaultStatistics
}

func (s *statisticsSvc) ScriptRecord(ctx context.Context, data *producer.ScriptStatisticsMsg) error {
	return producer.PublishScriptStatistics(ctx, data)
}

func (s *statisticsSvc) GetStatisticsToken(ctx *gin.Context) string {
	stk, _ := ctx.Cookie("_statistics")
	if stk == "" {
		stk = utils.RandString(32, 2)
		ctx.SetCookie("_statistics", stk, 3600*24*365, "/", "", false, true)
	}
	return stk
}

func (s *statisticsSvc) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sid := ctx.Param("id")
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		script, err := script_repo.Script().Find(ctx, id)
		if err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
		// 统计允许管理员查看
		if err := script.CheckPermission(ctx, model.Admin); err != nil {
			httputils.HandleResp(ctx, err)
			return
		}
	}
}

// Script 脚本统计数据
func (s *statisticsSvc) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	return &api.ScriptResponse{
		PagePv: &api.Overview{
			Today:     DaysPvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 1, time.Now()),
			Yesterday: DaysPvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 1, time.Now().AddDate(0, 0, -1)),
			Week:      DaysPvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 7, time.Now()),
		},
		PageUv: &api.Overview{
			Today:     DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 1, time.Now()),
			Yesterday: DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 1, time.Now().AddDate(0, 0, -1)),
			Week:      DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.ViewScriptStatistics, 7, time.Now()),
		},
		DownloadUv: &api.Overview{
			Today:     DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.DownloadScriptStatistics, 1, time.Now()),
			Yesterday: DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.DownloadScriptStatistics, 1, time.Now().AddDate(0, 0, -1)),
			Week:      DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.DownloadScriptStatistics, 7, time.Now()),
		},
		UpdateUv: &api.Overview{
			Today:     DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.UpdateScriptStatistics, 1, time.Now()),
			Yesterday: DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.UpdateScriptStatistics, 1, time.Now().AddDate(0, 0, -1)),
			Week:      DaysUvNumIgnoreError(ctx, req.ID, statistics_repo.UpdateScriptStatistics, 7, time.Now()),
		},
		UvChart: &api.DUChart{
			Download: s.daysData(ctx, req.ID, 30, time.Now(), statistics_repo.DownloadScriptStatistics, "uv"),
			Update:   s.daysData(ctx, req.ID, 30, time.Now(), statistics_repo.UpdateScriptStatistics, "uv"),
		},
		PvChart: &api.DUChart{
			Download: s.daysData(ctx, req.ID, 30, time.Now(), statistics_repo.DownloadScriptStatistics, "pv"),
			Update:   s.daysData(ctx, req.ID, 30, time.Now(), statistics_repo.UpdateScriptStatistics, "pv"),
		},
	}, nil
}

func DaysPvNumIgnoreError(ctx context.Context, scriptId int64, op statistics_repo.ScriptStatisticsType, days int, t time.Time) int64 {
	resp, err := statistics_repo.ScriptStatistics().DaysPvNum(ctx, scriptId, op, days, t)
	if err != nil {
		logger.Ctx(ctx).Error("DaysPvNumIgnoreError", zap.Error(err),
			zap.Int64("scriptId", scriptId), zap.Int("days", days), zap.Time("t", t))
	}
	return resp
}

func DaysUvNumIgnoreError(ctx context.Context, scriptId int64, op statistics_repo.ScriptStatisticsType, days int, t time.Time) int64 {
	resp, err := statistics_repo.ScriptStatistics().DaysUvNum(ctx, scriptId, op, days, t)
	if err != nil {
		logger.Ctx(ctx).Error("DaysUvNumIgnoreError", zap.Error(err),
			zap.Int64("scriptId", scriptId), zap.Int("days", days), zap.Time("t", t))
	}
	return resp
}

func (s *statisticsSvc) daysData(ctx context.Context, scriptId, days int64, date time.Time,
	op statistics_repo.ScriptStatisticsType, data string) *api.Chart {
	t := date.Add(-time.Hour * 24 * time.Duration(days))
	var x []string
	var y []int64
	for i := int64(0); i < days; i++ {
		t = t.Add(time.Hour * 24)
		day := t.Format("2006/01/02")
		var num int64
		switch data {
		case "uv":
			num, _ = statistics_repo.ScriptStatistics().DaysUvNum(ctx, scriptId, op, 1, t)
		case "pv":
			num, _ = statistics_repo.ScriptStatistics().DaysPvNum(ctx, scriptId, op, 1, t)
		}
		x = append(x, day)
		y = append(y, num)
	}
	return &api.Chart{
		X: x,
		Y: y,
	}
}

// ScriptRealtime 脚本实时统计数据
func (s *statisticsSvc) ScriptRealtime(ctx context.Context, req *api.ScriptRealtimeRequest) (*api.ScriptRealtimeResponse, error) {
	return &api.ScriptRealtimeResponse{
		Download: s.realtime(ctx, req.ID, statistics_repo.DownloadScriptStatistics),
		Update:   s.realtime(ctx, req.ID, statistics_repo.UpdateScriptStatistics),
	}, nil
}

func (s *statisticsSvc) realtime(ctx context.Context, scriptId int64, op statistics_repo.ScriptStatisticsType) *api.Chart {
	var nums []int64
	nums, _ = statistics_repo.ScriptStatistics().Realtime(ctx, scriptId, op)
	l := len(nums)
	var x = make([]string, l)
	var y = make([]int64, l)
	for n, v := range nums {
		x[l-n-1] = strconv.Itoa(n+1) + "分钟前"
		y[l-n-1] = v
	}
	return &api.Chart{
		X: x,
		Y: y,
	}
}
