package subscribe

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// EsSync 同步到es
type EsSync struct {
}

func (e *EsSync) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptCreate(ctx, e.scriptCreate, broker.Group("es")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptCodeUpdate(ctx, e.scriptCodeUpdate, broker.Group("es")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptStatistics(ctx, e.scriptStatistic, broker.Group("es")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptDelete(ctx, e.scriptDelete, broker.Group("es")); err != nil {
		return err
	}
	return nil
}

// 消费脚本创建消息推送到elasticsearch
func (e *EsSync) scriptCreate(ctx context.Context, script *script_entity.Script, code int64) error {
	return e.syncScript(ctx, script, code, false)
}

func (e *EsSync) syncScript(ctx context.Context, script *script_entity.Script, code int64, update bool) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", script.ID), zap.Bool("update", update))
	search, err := script_repo.Migrate().Convert(ctx, script)
	if err != nil {
		logger.Error("迁移es获取数据失败", zap.Error(err))
		return err
	}
	if err := script_repo.Migrate().Save(ctx, search); err != nil {
		logger.Error("迁移es保存数据失败", zap.Error(err))
		return err
	}
	logger.Info("迁移es成功")
	return nil
}

// 消费脚本代码更新消息,更新es记录
func (e *EsSync) scriptCodeUpdate(ctx context.Context, script *script_entity.Script, code int64) error {
	return e.syncScript(ctx, script, code, false)
}

func (e *EsSync) statisticSyncKey(scriptId int64) string {
	return fmt.Sprintf("script:es:sync:statistic:%d", scriptId)
}

func (e *EsSync) scriptStatistic(ctx context.Context, msg *producer.ScriptStatisticsMsg) error {
	num, err := redis.Ctx(ctx).HIncrBy(e.statisticSyncKey(msg.ScriptID), "num", 1).Result()
	if err != nil {
		return err
	}
	// 当囤了100条记录或者时间超过了5分钟,同步到es
	if num < 100 {
		t, err := redis.Ctx(ctx).HGet(e.statisticSyncKey(msg.ScriptID), "time").Int64()
		if err != nil {
			if !redis.Nil(err) {
				return err
			}
		}
		if time.Now().Unix()-t < 300 {
			return nil
		}
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.ScriptID), zap.String("download", string(msg.Download)))
	if err := redis.Ctx(ctx).HSet(e.statisticSyncKey(msg.ScriptID), "time", time.Now().Unix()).Err(); err != nil {
		logger.Error("数据设置失败", zap.Error(err))
	}
	if err := redis.Ctx(ctx).HDel(e.statisticSyncKey(msg.ScriptID), "num").Err(); err != nil {
		logger.Error("数据清理失败", zap.Error(err))
	}
	script, err := script_repo.Script().Find(ctx, msg.ScriptID)
	if err != nil {
		return err
	}
	if script == nil {
		return errors.New("script is nil")
	}
	esScript, err := script_repo.Migrate().Convert(ctx, script)
	if err != nil {
		return err
	}
	return script_repo.Migrate().Update(ctx, esScript)
}

// 消费脚本删除消息,删除es记录
func (e *EsSync) scriptDelete(ctx context.Context, msg *script_entity.Script) error {
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.ID))
	if err := script_repo.Migrate().Delete(ctx, msg.ID); err != nil {
		logger.Error("删除es数据失败", zap.Error(err))
		return err
	}
	logger.Info("删除es数据成功")
	return nil
}
