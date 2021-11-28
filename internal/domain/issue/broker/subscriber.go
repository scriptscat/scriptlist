package broker

import (
	"encoding/json"

	"github.com/scriptscat/scriptlist/pkg/event"
)

func SubscribeScriptIssueCreate(h func(issue, script int64) error) (event.Subscriber, error) {
	return event.DefaultBroker.Subscribe(EventScriptIssueCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["script"], ids["code"])
	})
}
