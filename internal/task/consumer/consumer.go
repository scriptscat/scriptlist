package consumer

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/task/producer"
)

// Consumer 消费者
func Consumer(ctx context.Context, broker broker.Broker) error {
	_, err := broker.Subscribe(ctx,
		producer.ScriptCreateTopic, scriptCreateHandler,
	)
	if err != nil {
		return err
	}
	return nil
}
