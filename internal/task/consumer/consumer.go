package consumer

import (
	"context"

	"github.com/cago-frame/cago/configs"
	"github.com/scriptscat/scriptlist/internal/task/consumer/subscribe"
)

type Subscribe interface {
	Subscribe(ctx context.Context) error
}

// Consumer 消费者
func Consumer(ctx context.Context, cfg *configs.Config) error {
	subscribers := []Subscribe{
		&subscribe.EsSync{},
		&subscribe.Script{},
		&subscribe.Statistics{},
		&subscribe.Issue{},
		&subscribe.Access{},
	}
	for _, v := range subscribers {
		if err := v.Subscribe(ctx); err != nil {
			return err
		}
	}
	return nil
}
