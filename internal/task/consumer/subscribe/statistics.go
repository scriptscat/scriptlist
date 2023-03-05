package subscribe

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/mileusna/useragent"
	"github.com/scriptscat/scriptlist/internal/model/entity/statistics_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// Statistics 处理统计平台数据
type Statistics struct {
}

func (s *Statistics) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptStatistics(ctx, s.scriptStatistics); err != nil {
		return err
	}
	if err := producer.SubscribeStatisticsCollect(ctx, s.collect); err != nil {
		return err
	}
	return nil
}

func (s *Statistics) scriptStatistics(ctx context.Context, msg *producer.ScriptStatisticsMsg) error {
	switch msg.Download {
	case statistics_repo.DownloadScriptStatistics:
		if ok, err := statistics_repo.ScriptStatistics().IncrDownload(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计下载量失败", zap.Error(err))
			return err
		} else if ok {
			// 统计总量
			if err := script_repo.ScriptStatistics().IncrDownload(ctx, msg.ScriptID); err != nil {
				logger.Ctx(ctx).Error("统计总下载量失败", zap.Error(err))
				return err
			}
			// 统计当日
			if err := script_repo.ScriptDateStatistics().IncrDownload(ctx, msg.ScriptID, msg.Time); err != nil {
				logger.Ctx(ctx).Error("统计当日下载量失败", zap.Error(err))
				return err
			}
		}
	case statistics_repo.UpdateScriptStatistics:
		if ok, err := statistics_repo.ScriptStatistics().IncrUpdate(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计更新量失败", zap.Error(err))
			return err
		} else if ok {
			// 统计总量
			if err := script_repo.ScriptStatistics().IncrUpdate(ctx, msg.ScriptID); err != nil {
				logger.Ctx(ctx).Error("统计总更新量失败", zap.Error(err))
				return err
			}
			// 统计当日
			if err := script_repo.ScriptDateStatistics().IncrUpdate(ctx, msg.ScriptID, msg.Time); err != nil {
				logger.Ctx(ctx).Error("统计当日更新量失败", zap.Error(err))
				return err
			}
		}
	case statistics_repo.ViewScriptStatistics:
		if _, err := statistics_repo.ScriptStatistics().IncrPageView(ctx, msg.ScriptID, msg.IP, msg.StatisticsToken); err != nil {
			logger.Ctx(ctx).Error("统计浏览量失败", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Statistics) collect(ctx context.Context, msg *producer.StatisticsCollectMsg) error {
	// ip+用户提交的访客id生成后端存储的访客id
	vistitorId := fmt.Sprintf("%x", sha256.Sum256([]byte(msg.IP+msg.VisitorID)))
	operatorUrl, err := url.Parse(msg.OperationPage)
	if err != nil {
		logger.Ctx(ctx).Error("统计页url解析失败", zap.Error(err), zap.Any("msg", msg))
		return err
	}
	installUrl, err := url.Parse(msg.InstallPage)
	if err != nil {
		installUrl = &url.URL{Host: ""}
		logger.Ctx(ctx).Error("统计页url解析失败", zap.Error(err), zap.Any("msg", msg))
	}
	collect := &statistics_entity.StatisticsCollect{
		SessionID:     msg.SessionID,
		ScriptID:      msg.ScriptID,
		VisitorID:     vistitorId,
		OperationHost: operatorUrl.Host,
		OperationPage: msg.OperationPage,
		Duration:      msg.Duration,
		VisitTime:     msg.VisitTime,
		ExitTime:      msg.ExitTime,
	}
	if err := statistics_repo.StatisticsCollect().Create(ctx, collect); err != nil {
		logger.Ctx(ctx).Error("统计收集失败", zap.Error(err), zap.Any("msg", msg))
	}
	// 记录第一次访问时间
	key := fmt.Sprintf("statistics:visitor:%s", vistitorId)
	firstVisitTime, err := redis.Ctx(ctx).Get(key).Int64()
	if err != nil {
		firstVisitTime = msg.VisitTime
		if !redis.Nil(err) {
			logger.Ctx(ctx).Error("获取访客第一次访问时间失败", zap.Error(err), zap.Any("msg", msg))
		} else {
			if err := redis.Ctx(ctx).Set(key, firstVisitTime, 0).Err(); err != nil {
				logger.Ctx(ctx).Error("设置访客第一次访问时间失败", zap.Error(err), zap.Any("msg", msg))
			}
		}
	}
	ua := useragent.Parse(msg.UA)
	driverType := statistics_entity.DeviceTypeUnknown
	if ua.Mobile {
		driverType = statistics_entity.DeviceTypeMobile
	} else if ua.Desktop {
		driverType = statistics_entity.DeviceTypePC
	}
	if err := statistics_repo.StatisticsVisitor().Create(ctx, &statistics_entity.StatisticsVisitor{
		ScriptID:       msg.ScriptID,
		VisitorID:      vistitorId,
		UA:             msg.UA,
		IP:             msg.IP,
		Version:        msg.Version,
		InstallPage:    msg.InstallPage,
		FirstVisitTime: firstVisitTime,
		VisitTime:      msg.VisitTime,
		InstallHost:    installUrl.Host,
		DeviceType:     int64(driverType),
		BrowserType:    ua.Name,
	}); err != nil {
		logger.Ctx(ctx).Error("统计访客失败", zap.Error(err), zap.Any("msg", msg))
	}
	return nil
}
