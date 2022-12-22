package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
)

type ScriptStatisticsMsg struct {
	ScriptID, ScriptCodeID, UserID    int64
	IP, UA, StatisticsToken, Download string
}

func PublishScriptStatistics(ctx context.Context, msg *ScriptStatisticsMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptStatisticsMsg(msg *broker2.Message) (*ScriptStatisticsMsg, error) {
	ret := &ScriptStatisticsMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
