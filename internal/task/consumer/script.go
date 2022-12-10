package consumer

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
)

func scriptCreateHandler(ctx context.Context, event broker.Event) error {
	logger.Ctx(ctx).Info("script create")
	return nil
}
