package producer

import (
	"context"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
)

func ScriptCreate(ctx context.Context, test int) error {
	return broker.Default().Publish(ctx, ScriptCreateTopic, &broker2.Message{
		Header: nil,
		Body:   []byte("test"),
	})
}
