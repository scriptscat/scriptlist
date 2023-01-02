package consumer

import (
	"context"
	"errors"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

// 同步到es
type esSync struct {
}

func (e *esSync) Subscribe(ctx context.Context, bk broker.Broker) error {
	_, err := bk.Subscribe(ctx,
		producer.ScriptCreateTopic, e.scriptCreateHandler,
		broker.Group("es"),
	)
	if err != nil {
		return err
	}
	_, err = bk.Subscribe(ctx, producer.ScriptCodeUpdateTopic, e.scriptCodeUpdateHandler, broker.Group("es"))
	return err
}

// 消费脚本创建消息推送到elasticsearch
func (e *esSync) scriptCreateHandler(ctx context.Context, event broker.Event) error {
	return e.syncScript(ctx, event, false)
}

func (e *esSync) syncScript(ctx context.Context, event broker.Event, update bool) error {
	msg, err := producer.ParseScriptCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).Error("ParseScriptCreateMsg", zap.Error(err), zap.Binary("body", event.Message().Body))
		return err
	}
	if msg.Script == nil {
		return errors.New("script is nil")
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID), zap.Bool("update", update))
	search, err := script_repo.Migrate().Convert(ctx, msg.Script)
	if err != nil {
		logger.Error("迁移es获取数据失败", zap.Error(err))
		return err
	}
	if err := script_repo.Migrate().SaveToEs(ctx, search); err != nil {
		logger.Error("迁移es保存数据失败", zap.Error(err))
		return err
	}
	logger.Info("迁移es成功")
	return nil
}

// 消费脚本代码更新消息,更新es记录
func (e *esSync) scriptCodeUpdateHandler(ctx context.Context, event broker.Event) error {
	return e.syncScript(ctx, event, false)
}
