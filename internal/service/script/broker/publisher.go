package broker

import (
	"github.com/scriptscat/scriptlist/internal/infrastructure/broker"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const (
	EventScriptCreate        = "event:script:create"
	EventScriptVersionUpdate = "event:script:version:update"
)

func PublishEventScriptVersionUpdate(script, code int64) error {
	return broker.DefaultBroker.Publish(EventScriptVersionUpdate, &broker.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(broker.Ids{
			"code":   code,
			"script": script,
		}),
	})
}

func PublishEventScriptCreate(script int64) error {
	return broker.DefaultBroker.Publish(EventScriptCreate, &broker.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(broker.Ids{
			"script": script,
		}),
	})
}
