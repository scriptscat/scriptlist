package broker

import (
	"github.com/scriptscat/scriptlist/pkg/event"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const (
	EventScriptVersionUpdate = "event:script:version:update"
)

func PublishEventScriptVersionUpdate(script, code int64) error {
	return event.DefaultBroker.Publish(EventScriptVersionUpdate, &event.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(event.Ids{
			"code":   code,
			"script": script,
		}),
	})
}
