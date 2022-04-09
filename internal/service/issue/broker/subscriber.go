package broker

import (
	"encoding/json"

	"github.com/scriptscat/scriptlist/internal/infrastructure/broker"
)

// SubscribeScriptIssueCreate 订阅issue创建事件,传递脚本id和issue id
func SubscribeScriptIssueCreate(h func(script, issue int64) error) (broker.Subscriber, error) {
	return broker.DefaultBroker.Subscribe(EventScriptIssueCreate, func(message *broker.Message) error {
		ids := broker.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["script"], ids["issue"])
	})
}

// SubscribeScriptIssueCommentCreate 订阅issue评论事件,传递issue和评论
func SubscribeScriptIssueCommentCreate(h func(issue, comment int64) error) (broker.Subscriber, error) {
	return broker.DefaultBroker.Subscribe(EventScriptIssueCommentCreate, func(message *broker.Message) error {
		ids := broker.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return h(ids["issue"], ids["comment"])
	})
}
