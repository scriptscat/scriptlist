package consumer

import (
	"context"
	"fmt"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

// 处理统计平台数据
type statistics struct {
}

func (e *statistics) Subscribe(ctx context.Context, bk broker.Broker) error {
	_, err := bk.Subscribe(ctx,
		producer.ScriptStatisticTopic, e.scriptStatistics,
	)
	return err
}

func (e *statistics) scriptStatistics(ctx context.Context, event broker.Event) error {
	msg, err := producer.ParseScriptStatisticsMsg(event.Message())
	if err != nil {
		return err
	}
	// TODO: 记录统计数据
	fmt.Sprintf("msg: %+v", msg)
	return nil
}
