package broker

import (
	"github.com/scriptscat/scriptlist/pkg/event"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const (
	EventScriptIssueCreate        = "event:script:issue:create"
	EventScriptIssueCommentCreate = "event:script:issue:comment:create"
)

func PublishScriptIssueCreate(issue, script int64) error {
	return event.DefaultBroker.Publish(EventScriptIssueCreate, &event.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(event.Ids{
			"issue":  issue,
			"script": script,
		}),
	})
}

func PublishScriptIssueCommentCreate(issue, comment int64) error {
	return event.DefaultBroker.Publish(EventScriptIssueCommentCreate, &event.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(event.Ids{
			"issue":   issue,
			"comment": comment,
		}),
	})
}
