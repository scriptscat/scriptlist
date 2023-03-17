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
	Download                       statistics_repo.ScriptStatisticsType
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

type StatisticsCollectMsg struct {
	SessionID     string
	ScriptID      int64
	VisitorID     string
	OperationHost string
	OperationPage string
	InstallPage   string
	Duration      int32
	UA            string
	IP            string
	VisitTime     int64
	ExitTime      int64
	Version       string
}

func PublishStatisticsCollect(ctx context.Context, msg *StatisticsCollectMsg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, StatisticCollectTopic, &broker2.Message{
		Body: body,
	})
}

func ParseStatisticsCollectMsg(msg *broker2.Message) (*StatisticsCollectMsg, error) {
	ret := &StatisticsCollectMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeStatisticsCollect(ctx context.Context, fn func(ctx context.Context, msg *StatisticsCollectMsg) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, StatisticCollectTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseStatisticsCollectMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m)
	}, opts...)
	return err
}
