package consumer

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
)

type Subscribe interface {
	Subscribe(ctx context.Context, broker broker.Broker) error
}

// Consumer 消费者
func Consumer(ctx context.Context, broker broker.Broker) error {
	subscribers := []Subscribe{&esSync{}, &script{}, &statistics{}}
	for _, v := range subscribers {
		if err := v.Subscribe(ctx, broker); err != nil {
			return err
		}
	}
	return nil
}
