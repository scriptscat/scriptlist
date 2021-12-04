package broker

import (
	"encoding/json"

	"github.com/scriptscat/scriptlist/pkg/event"
)

// SubscribeScriptIssueCreate 订阅issue创建事件,传递脚本id和issue id
func SubscribeScriptIssueCreate(h func(script, issue int64) error) (event.Subscriber, error) {
	return event.DefaultBroker.Subscribe(EventScriptIssueCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["script"], ids["issue"])
	})
}

// SubscribeScriptIssueCommentCreate 订阅issue评论事件,传递issue和评论
func SubscribeScriptIssueCommentCreate(h func(issue, comment int64) error) (event.Subscriber, error) {
	return event.DefaultBroker.Subscribe(EventScriptIssueCommentCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["issue"], ids["comment"])
	})
}
