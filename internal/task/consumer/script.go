package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
	script2 "github.com/scriptscat/scriptlist/internal/repository/script"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type script struct {
	// 分类id
	bgCategory   *entity.ScriptCategoryList
	cronCategory *entity.ScriptCategoryList
}

func (s *script) Subscribe(ctx context.Context, broker broker.Broker) error {
	var err error
	s.bgCategory, err = script2.ScriptCategoryList().FindByName(ctx, "后台脚本")
	if err != nil {
		return err
	}
	if s.bgCategory == nil {
		s.bgCategory = &entity.ScriptCategoryList{
			Name:       "后台脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script2.ScriptCategoryList().Create(ctx, s.bgCategory); err != nil {
			return err
		}
	}
	s.cronCategory, err = script2.ScriptCategoryList().FindByName(ctx, "定时脚本")
	if err != nil {
		return err
	}
	if s.cronCategory == nil {
		s.cronCategory = &entity.ScriptCategoryList{
			Name:       "定时脚本",
			Createtime: time.Now().Unix(),
		}
		if err := script2.ScriptCategoryList().Create(ctx, s.cronCategory); err != nil {
			return err
		}
	}
	_, err = broker.Subscribe(ctx,
		producer.ScriptCreateTopic, s.scriptCreateHandler,
	)
	if err != nil {
		return err
	}
	_, err = broker.Subscribe(ctx, producer.ScriptCodeUpdateTopic, s.scriptCodeUpdate)
	return err
}

// 消费脚本创建消息,根据meta信息进行分类和发送邮件通知
func (s *script) scriptCreateHandler(ctx context.Context, event broker.Event) error {
	msg, err := producer.ParseScriptCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID))

	// 根据meta信息, 将脚本分类到后台脚本, 定时脚本, 用户脚本
	metaJson := make(map[string][]string)
	if err := json.Unmarshal([]byte(msg.Code.MetaJson), &metaJson); err != nil {
		logger.Error("json.Unmarshal", zap.Error(err), zap.String("meta", msg.Code.MetaJson))
		return err
	}

	if len(metaJson["background"]) > 0 || len(metaJson["crontab"]) > 0 {
		// 后台脚本
		if err := script2.ScriptCategory().LinkCategory(ctx, msg.Script.ID, s.bgCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}
	if len(metaJson["crontab"]) > 0 {
		// 定时脚本
		if err := script2.ScriptCategory().LinkCategory(ctx, msg.Script.ID, s.cronCategory.ID); err != nil {
			logger.Error("LinkCategory", zap.Error(err))
			return err
		}
	}

	// 发送邮件通知

	return nil
}

// 消费脚本代码更新消息,发送邮件通知给关注了的用户
func (s *script) scriptCodeUpdate(ctx context.Context, event broker.Event) error {
	msg, err := producer.ParseScriptCodeUpdateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID))

	logger.Info("update script code")
	// 发送邮件通知
	return nil
}
