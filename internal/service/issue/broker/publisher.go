package broker

import (
	"github.com/scriptscat/scriptlist/internal/infrastructure/broker"
	"github.com/scriptscat/scriptlist/pkg/utils"
)

const (
	EventScriptIssueCreate        = "event:script:issue:create"
	EventScriptIssueCommentCreate = "event:script:issue:comment:create"
)

func PublishScriptIssueCreate(issue, script int64) error {
	return broker.DefaultBroker.Publish(EventScriptIssueCreate, &broker.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(broker.Ids{
			"issue":  issue,
			"script": script,
		}),
	})
}

func PublishScriptIssueCommentCreate(issue, comment int64) error {
	return broker.DefaultBroker.Publish(EventScriptIssueCommentCreate, &broker.Message{
		Header: nil,
		Body: utils.MarshalJsonByte(broker.Ids{
			"issue":   issue,
			"comment": comment,
		}),
	})
}
