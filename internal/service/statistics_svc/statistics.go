package statistics_svc

import (
	"context"
	"strconv"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	api "github.com/scriptscat/scriptlist/internal/api/statistics"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
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
	// Collect 统计数据收集
	Collect(ctx context.Context, req *api.CollectRequest) (*api.CollectResponse, error)
	// RealtimeChart 实时统计数据图表
	RealtimeChart(ctx context.Context, req *api.RealtimeChartRequest) (*api.RealtimeChartResponse, error)
	// Realtime 实时统计数据
	Realtime(ctx context.Context, req *api.RealtimeRequest) (*api.RealtimeResponse, error)
	// BasicInfo 基本统计信息
	BasicInfo(ctx context.Context, req *api.BasicInfoRequest) (*api.BasicInfoResponse, error)
	// UserOrigin 用户来源统计
	UserOrigin(ctx context.Context, req *api.UserOriginRequest) (*api.UserOriginResponse, error)
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
		ctx.Set("statistics_token", s.GetStatisticsToken(ctx))
		ctx.Next()
		id := ctx.GetInt64("id")
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

// Collect 统计数据收集
func (s *statisticsSvc) Collect(ctx context.Context, req *api.CollectRequest) (*api.CollectResponse, error) {
	// 判断本月是否超过限制
	ok, err := statistics_repo.StatisticsCollect().CheckLimit(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, i18n.NewError(ctx, code.StatisticsLimitExceeded)
	}
	if err := producer.PublishStatisticsCollect(ctx, &producer.StatisticsCollectMsg{
		SessionID:     req.SessionID,
		ScriptID:      req.ScriptID,
		VisitorID:     req.VisitorID,
		OperationPage: req.OperationPage,
		InstallPage:   req.InstallPage,
		Duration:      req.Duration,
		UA:            req.UA,
		IP:            req.IP,
		VisitTime:     req.VisitTime,
		ExitTime:      req.ExitTime,
		Version:       req.Version,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

// RealtimeChart 实时统计数据图表
func (s *statisticsSvc) RealtimeChart(ctx context.Context, req *api.RealtimeChartRequest) (*api.RealtimeChartResponse, error) {
	now := time.Now()
	list, err := statistics_repo.StatisticsCollect().RealtimeChart(ctx, req.ID, now)
	if err != nil {
		return nil, err
	}
	chart := &api.Chart{
		X: make([]string, 0),
		Y: make([]int64, 0),
	}
	listHash := make(map[int]int64)
	for _, v := range list {
		listHash[v.Time] = v.Num
	}
	for i := 0; i < 15; i++ {
		num, ok := listHash[now.Minute()]
		chart.X = append(chart.X, strconv.Itoa(i+1)+"分钟前")
		if ok {
			chart.Y = append(chart.Y, num)
		} else {
			chart.Y = append(chart.Y, 0)
		}
		now = now.Add(-time.Minute)
	}
	return &api.RealtimeChartResponse{
		Chart: chart,
	}, nil
}

// Realtime 实时统计数据
func (s *statisticsSvc) Realtime(ctx context.Context, req *api.RealtimeRequest) (*api.RealtimeResponse, error) {
	return nil, nil
}

// BasicInfo 基本统计信息
func (s *statisticsSvc) BasicInfo(ctx context.Context, req *api.BasicInfoRequest) (*api.BasicInfoResponse, error) {
	var quota int64 = 1000000
	usage, err := statistics_repo.StatisticsCollect().GetLimit(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if usage > quota {
		usage = quota
	}
	now := time.Now()
	uv := &api.Overview{
		Today:     s.IgnoreErrorCollectUv(ctx, req.ID, now.Add(-time.Hour*24), now),
		Yesterday: s.IgnoreErrorCollectUv(ctx, req.ID, now.Add(-time.Hour*48), now.Add(-time.Hour*24)),
		Week:      s.IgnoreErrorCollectUv(ctx, req.ID, now.Add(-time.Hour*24*7), now),
	}
	newUser := &api.Overview{
		Today:     s.IgnoreErrorFirstUserNumber(ctx, req.ID, now.Add(-time.Hour*24), now),
		Yesterday: s.IgnoreErrorFirstUserNumber(ctx, req.ID, now.Add(-time.Hour*48), now.Add(-time.Hour*24)),
		Week:      s.IgnoreErrorFirstUserNumber(ctx, req.ID, now.Add(-time.Hour*24*7), now),
	}
	return &api.BasicInfoResponse{
		Limit: &api.Limit{
			Quota: quota,
			Usage: usage,
		},
		PV: &api.Overview{
			Today:     s.IgnoreErrorCollectPv(ctx, req.ID, now.Add(-time.Hour*24), now),
			Yesterday: s.IgnoreErrorCollectPv(ctx, req.ID, now.Add(-time.Hour*48), now.Add(-time.Hour*24)),
			Week:      s.IgnoreErrorCollectPv(ctx, req.ID, now.Add(-time.Hour*24*7), now),
		},
		UV: uv,
		UseTime: &api.Overview{
			Today:     s.IgnoreErrorUseTimeAvg(ctx, req.ID, now.Add(-time.Hour*24), now),
			Yesterday: s.IgnoreErrorUseTimeAvg(ctx, req.ID, now.Add(-time.Hour*48), now.Add(-time.Hour*24)),
			Week:      s.IgnoreErrorUseTimeAvg(ctx, req.ID, now.Add(-time.Hour*24*7), now),
		},
		NewUser: newUser,
		OldUser: &api.Overview{
			Today:     uv.Today - newUser.Today,
			Yesterday: uv.Yesterday - newUser.Yesterday,
			Week:      uv.Week - newUser.Week,
		},
		Origin:          nil,
		Version:         nil,
		OperationDomain: nil,
		System:          nil,
		Browser:         nil,
	}, nil
}

func (s *statisticsSvc) IgnoreErrorFirstUserNumber(ctx context.Context, id int64, start, end time.Time) int64 {
	newNum, err := statistics_repo.StatisticsVisitor().FirstUserNumber(ctx, id, start, end)
	if err != nil {
		logger.Ctx(ctx).Error("statistics_repo.StatisticsVisitor().FirstUserNumber", zap.Error(err))
		return 0
	}
	return newNum
}

func (s *statisticsSvc) IgnoreErrorCollectPv(ctx context.Context, id int64, start time.Time, end time.Time) int64 {
	num, err := statistics_repo.StatisticsCollect().Pv(ctx, id, start, end)
	if err != nil {
		logger.Ctx(ctx).Error("statistics_repo.StatisticsCollect().Pv", zap.Error(err))
		return 0
	}
	return num
}

func (s *statisticsSvc) IgnoreErrorCollectUv(ctx context.Context, id int64, start time.Time, end time.Time) int64 {
	num, err := statistics_repo.StatisticsCollect().Uv(ctx, id, start, end)
	if err != nil {
		logger.Ctx(ctx).Error("statistics_repo.StatisticsCollect().Uv", zap.Error(err))
		return 0
	}
	return num
}

func (s *statisticsSvc) IgnoreErrorUseTimeAvg(ctx context.Context, id int64, start time.Time, end time.Time) int64 {
	num, err := statistics_repo.StatisticsCollect().UseTimeAvg(ctx, id, start, end)
	if err != nil {
		logger.Ctx(ctx).Error("statistics_repo.StatisticsCollect().UseTime", zap.Error(err))
		return 0
	}
	return num
}

// UserOrigin 用户来源统计
func (s *statisticsSvc) UserOrigin(ctx context.Context, req *api.UserOriginRequest) (*api.UserOriginResponse, error) {
	return nil, nil
}
