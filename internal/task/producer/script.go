package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

// 脚本相关消息生产者

type ScriptCreateMsg struct {
	Script *entity.Script
	Code   *entity.Code
}

func PublishScriptCreate(ctx context.Context, script *entity.Script, code *entity.Code) error {
	body, err := json.Marshal(&ScriptCreateMsg{
		Script: script,
		Code:   code,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptCreateMsg(msg *broker2.Message) (*ScriptCreateMsg, error) {
	ret := &ScriptCreateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeScriptCreate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, code *script_entity.Code) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptCreateTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptCreateMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m.Script, m.Code)
	}, opts...)
	return err
}

type ScriptCodeUpdateMsg struct {
	Script *entity.Script
	Code   *entity.Code
}

func PublishScriptCodeUpdate(ctx context.Context, script *entity.Script, code *entity.Code) error {
	body, err := json.Marshal(&ScriptCodeUpdateMsg{
		Script: script,
		Code:   code,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptCodeUpdateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptCodeUpdateMsg(msg *broker2.Message) (*ScriptCodeUpdateMsg, error) {
	ret := &ScriptCodeUpdateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeScriptCodeUpdate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, code *script_entity.Code) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptCodeUpdateTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptCodeUpdateMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m.Script, m.Code)
	}, opts...)
	return err
}

func PublishScriptDelete(ctx context.Context, script *entity.Script) error {
	body, err := json.Marshal(script)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptDeleteTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptDeleteMsg(msg *broker2.Message) (*entity.Script, error) {
	ret := &entity.Script{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeScriptDelete(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptDeleteTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptDeleteMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m)
	}, opts...)
	return err
}
