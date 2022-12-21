package consumer

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/codfrm/cago/database/elasticsearch"
	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/scriptscat/scriptlist/internal/model"
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

// 查询脚本下载量
func (e *esSync) queryDownload(ctx context.Context, id int64) (int64, int64, error) {
	return 0, 0, nil
}

// 查询脚本分数
func (e *esSync) queryScore(ctx context.Context, id int64) (float64, error) {
	return 0, nil
}

func (e *esSync) queryDomain(ctx context.Context, id int64) ([]string, error) {

	return nil, nil
}

func (e *esSync) syncScript(ctx context.Context, event broker.Event, update bool) error {
	msg, err := producer.ParseScriptCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).Error("ParseScriptCreateMsg", zap.Error(err), zap.Binary("body", event.Message().Body))
		return err
	}
	logger := logger.Ctx(ctx).With(zap.Int64("script_id", msg.Script.ID), zap.Bool("update", update))
	script := &model.ScriptSearch{
		ID:            msg.Script.ID,
		UserID:        msg.Script.UserID,
		Name:          msg.Script.Name,
		Description:   msg.Script.Description,
		Content:       msg.Script.Content,
		Changelog:     msg.Code.Changelog,
		TotalDownload: 0,
		TodayDownload: 0,
		Score:         0,
		Domain:        nil,
		Category:      nil,
		Public:        int(msg.Script.Public),
		Unwell:        int(msg.Script.Unwell),
		Createtime:    msg.Script.Createtime,
		Updatetime:    msg.Script.Updatetime,
	}
	r, err := script.Reader()
	if err != nil {
		return err
	}
	// 同步到es
	var resp *esapi.Response
	if update {
		resp, err = elasticsearch.Ctx(ctx).Update(
			script.CollectionName(), strconv.FormatInt(script.ID, 10), r,
		)
	} else {
		resp, err = elasticsearch.Ctx(ctx).Create(
			script.CollectionName(), strconv.FormatInt(script.ID, 10), r,
		)
	}
	if err != nil {
		logger.Error("insert error", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		logger.Error("insert error", zap.ByteString("body", b), zap.Int("status", resp.StatusCode))
		return err
	}
	logger.Info("insert success")
	return nil
}

// 消费脚本代码更新消息,更新es记录
func (e *esSync) scriptCodeUpdateHandler(ctx context.Context, event broker.Event) error {
	return e.syncScript(ctx, event, false)
}
