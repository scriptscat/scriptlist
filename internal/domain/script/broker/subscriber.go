package broker

import (
	"encoding/json"

	"github.com/scriptscat/scriptlist/pkg/event"
)

func SubscribeEventScriptVersionUpdate(h func(script, code int64) error) (event.Subscriber, error) {
	return event.DefaultBroker.Subscribe(EventScriptVersionUpdate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["script"], ids["code"])
	})
}

func SubscribeEventScriptCreate(h func(script int64) error) (event.Subscriber, error) {
	return event.DefaultBroker.Subscribe(EventScriptCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["script"])
	})
}
