package statistics_repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/logger"
	"go.uber.org/zap"
)

type ScriptStatisticsType string

const (
	ViewScriptStatistics     ScriptStatisticsType = "view"
	DownloadScriptStatistics ScriptStatisticsType = "download"
	UpdateScriptStatistics   ScriptStatisticsType = "update"
)

// pf   statistics:script:@op:@id:day:uv:@date 30天过期

// hash statistics:script:@op:@id:day:pv @date @num 永不过期
// hash statistics:script:@op:@id:total:pv 永不过期

// set statistics:script:@op:@id:realtime:@time @num 一小时过期

// ScriptStatisticsRepo 统计平台数据库操作,与脚本统计不同,此处的纬度更丰富,且大多记录在redis中
type ScriptStatisticsRepo interface {
	// Save 数据落库
	//Save(ctx context.Context) error

	// Realtime 获取某操作的实时记录
	Realtime(ctx context.Context, scriptId int64, op ScriptStatisticsType) ([]int64, error)
	// DaysUvNum 获取某一段时间某一个操作的uv或者member数量
	// op: download, update, view
	DaysUvNum(ctx context.Context, scriptId int64, op ScriptStatisticsType, days int, t time.Time) (int64, error)
	// DaysPvNum 获取某一段时间某一个操作的pv数量
	// op: download, update, view
	DaysPvNum(ctx context.Context, scriptId int64, op ScriptStatisticsType, days int, t time.Time) (int64, error)
	// TotalPv 获取某操作的总pv数量
	TotalPv(ctx context.Context, scriptId int64, op ScriptStatisticsType) (int64, error)
	// IncrDownload 增加下载量,使用ip判断是否重复
	IncrDownload(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error)
	IncrUpdate(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error)
	IncrPageView(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error)
}

var defaultScriptStatistics ScriptStatisticsRepo

func ScriptStatistics() ScriptStatisticsRepo {
	return defaultScriptStatistics
}

func RegisterScriptStatistics(i ScriptStatisticsRepo) {
	defaultScriptStatistics = i
}

type scriptStatisticsRepo struct {
}

func NewScriptStatistics() ScriptStatisticsRepo {
	return &scriptStatisticsRepo{}
}

func (s *scriptStatisticsRepo) Realtime(ctx context.Context, scriptId int64, op ScriptStatisticsType) ([]int64, error) {
	var ret []int64
	t := time.Now().Unix() / 60
	for i := int64(0); i < 15; i++ {
		num, _ := redis.Ctx(ctx).Get(fmt.Sprintf("statistics:script:%s:%d:realtime:%d", op, scriptId, t-i)).Int64()
		ret = append(ret, num)
	}
	return ret, nil
}

func (s *scriptStatisticsRepo) DaysUvNum(ctx context.Context, scriptId int64, op ScriptStatisticsType, days int, t time.Time) (int64, error) {
	if days == 1 {
		ret, err := redis.Ctx(ctx).PFCount(fmt.Sprintf(
			"statistics:script:%s:%d:day:%s:%s", op, scriptId, "uv", t.Format("2006/01/02"))).Result()
		if err != nil {
			return 0, err
		}
		return ret, nil
	}
	key := fmt.Sprintf("statistics:script:cache:uv:%d:%s:%s", scriptId, op, t.Format("2006/01/02"))
	if redis.Ctx(ctx).Exists(key).Val() != 1 {
		var dayKey []string
		t := t
		for i := 1; i <= days; i++ {
			dayKey = append(dayKey,
				fmt.Sprintf("statistics:script:%s:%d:day:%s:%s", op, scriptId, "uv",
					t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")),
			)
		}
		redis.Ctx(ctx).PFMerge(key, dayKey...)
	}
	ret, err := redis.Ctx(ctx).PFCount(key).Result()
	if err != nil {
		return 0, err
	}
	redis.Ctx(ctx).Expire(key, time.Hour*24*15)
	return ret, nil
}

func (s *scriptStatisticsRepo) TotalPv(ctx context.Context, scriptId int64, op ScriptStatisticsType) (int64, error) {
	key := "statistics:script:" + string(op) + ":" + fmt.Sprintf("%d", scriptId) + ":total:pv"
	return redis.Ctx(ctx).Get(key).Int64()
}

func (s *scriptStatisticsRepo) DaysPvNum(ctx context.Context, scriptId int64, op ScriptStatisticsType, days int, t time.Time) (int64, error) {
	var num int64
	key := fmt.Sprintf("statistics:script:%s:%d:day:pv", op, scriptId)
	for i := 0; i < days; i++ {
		val, err := redis.Ctx(ctx).HGet(key, t.Add(-time.Hour*24*time.Duration(i)).Format("2006/01/02")).Int64()
		if err != nil {
			continue
		}
		num += val
	}
	return num, nil
}

func (s *scriptStatisticsRepo) IncrDownload(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error) {
	return s.save(ctx, scriptId, ip, statisticsToken, DownloadScriptStatistics)
}

func (s *scriptStatisticsRepo) IncrUpdate(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error) {
	return s.save(ctx, scriptId, ip, statisticsToken, UpdateScriptStatistics)
}

func (s *scriptStatisticsRepo) IncrPageView(ctx context.Context, scriptId int64, ip string, statisticsToken string) (bool, error) {
	return s.save(ctx, scriptId, ip, statisticsToken, ViewScriptStatistics)
}

func (s *scriptStatisticsRepo) save(ctx context.Context, scriptId int64, ip, statisticsToken string, op ScriptStatisticsType) (bool, error) {
	key := "statistics:script:" + string(op) + ":" + fmt.Sprintf("%d", scriptId)
	date := time.Now().Format("2006/01/02")
	ok := true
	// 有更新的不再统计下载
	switch op {
	case UpdateScriptStatistics:
		if err := redis.Ctx(ctx).Set(
			"statistics:script:update:"+fmt.Sprintf("%d", scriptId)+":"+ip, "1",
			time.Minute*60).Err(); err != nil {
			logger.Ctx(ctx).Error("更新统计保存失败", zap.Error(err))
		}
	case DownloadScriptStatistics:
		e, err := redis.Ctx(ctx).Exists("statistics:script:update:" +
			fmt.Sprintf("%d", scriptId) + ":" + ip).Result()
		if err != nil {
			logger.Ctx(ctx).Error("更新统计查询失败", zap.Error(err))
		} else if e == 1 {
			ok = false
		}
	}
	// 储存统计token计算uv
	if ok {
		redis.Ctx(ctx).PFAdd(key+fmt.Sprintf(":day:uv:%s", date), statisticsToken)
		redis.Ctx(ctx).Expire(key+fmt.Sprintf(":day:uv:%s", date), time.Hour*24*60)
	}
	// 日pv
	redis.Ctx(ctx).HIncrBy(key+":day:pv", date, 1)
	// 总pv
	redis.Ctx(ctx).Incr(key + ":total:pv")
	// 实时统计
	t := strconv.FormatInt(time.Now().Unix()/60, 10)
	redis.Ctx(ctx).Incr(key + ":realtime:" + t)
	redis.Ctx(ctx).Expire(key+":realtime:"+t, time.Hour)
	// 判断ip是否操作过了
	result, err := redis.Ctx(ctx).SetNX(key+":ip:exist:day:"+date+":"+ip, "1", time.Hour*16).Result()
	if err != nil {
		return false, err
	}
	return ok && result, nil
}
