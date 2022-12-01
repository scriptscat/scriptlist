package consumer

import (
	"context"
	"fmt"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/scriptscat/scriptlist/pkg/producer"
)

// Consumer 消费者
func Consumer(ctx context.Context, cfg *configs.Config) error {
	_, err := broker.Default().Subscribe(ctx,
		producer.ScriptCreateTopic, scriptCreateHandler,
	)
	if err != nil {
		return err
	}
	err2 := producer.ScriptCreate(ctx, 123)
	fmt.Printf("%v", err2)
	return nil
}
