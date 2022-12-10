package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
)

func PublishScriptCreate(ctx context.Context, script *entity.Script) error {
	body, err := json.Marshal(script)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptCreateTopic, &broker2.Message{
		Body: body,
	})
}
