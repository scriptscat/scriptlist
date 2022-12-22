package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	entity "github.com/scriptscat/scriptlist/internal/model/entity/script"
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
