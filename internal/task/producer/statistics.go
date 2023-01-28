package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/repository/statistics_repo"
)

type ScriptStatisticsMsg struct {
	ScriptID, ScriptCodeID, UserID int64
	IP, UA, StatisticsToken        string
	Download                       statistics_repo.StatisticsType
	Time                           time.Time
}

func PublishScriptStatistics(ctx context.Context, msg *ScriptStatisticsMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptStatisticTopic, &broker2.Message{
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

func SubscribeScriptStatistics(ctx context.Context, fn func(ctx context.Context, msg *ScriptStatisticsMsg) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptStatisticTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptStatisticsMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m)
	}, opts...)
	return err
}
