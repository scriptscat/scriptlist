package consumer

import (
	"context"
	"fmt"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/pkg/producer"
)

// Consumer 消费者
func Consumer(ctx context.Context, broker broker.Broker) error {
	_, err := broker.Subscribe(ctx,
		producer.ScriptCreateTopic, scriptCreateHandler,
	)
	if err != nil {
		return err
	}
	err2 := producer.ScriptCreate(ctx, 123)
	fmt.Printf("%v", err2)
	return nil
}
