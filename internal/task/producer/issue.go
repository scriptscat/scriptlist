package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

func PublishIssueCreate(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue) error {
	body, err := json.Marshal(&IssueCreateMsg{
		script, issue,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, IssueCreateTopic, &broker2.Message{
		Body: body,
	})
}

type IssueCreateMsg struct {
	Script *script_entity.Script
	Issue  *issue_entity.ScriptIssue
}

func ParseIssueCreateMsg(msg *broker2.Message) (*IssueCreateMsg, error) {
	ret := &IssueCreateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeIssueCreate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue) error) error {
	_, err := broker.Default().Subscribe(ctx, IssueCreateTopic, func(ctx context.Context, event broker2.Event) error {
		msg, err := ParseIssueCreateMsg(event.Message())
		if err != nil {
			return err
		}
		return fn(ctx, msg.Script, msg.Issue)
	})
	return err
}

type CommentCreateMsg struct {
	Script  *script_entity.Script
	Issue   *issue_entity.ScriptIssue
	Comment *issue_entity.ScriptIssueComment
}

func PublishCommentCreate(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue, comment *issue_entity.ScriptIssueComment) error {
	body, err := json.Marshal(&CommentCreateMsg{
		Script:  script,
		Issue:   issue,
		Comment: comment,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, CommentCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseCommentCreateMsg(msg *broker2.Message) (*CommentCreateMsg, error) {
	ret := &CommentCreateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeCommentCreate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue, comment *issue_entity.ScriptIssueComment) error) error {
	_, err := broker.Default().Subscribe(ctx, CommentCreateTopic, func(ctx context.Context, event broker2.Event) error {
		msg, err := ParseCommentCreateMsg(event.Message())
		if err != nil {
			return err
		}
		return fn(ctx, msg.Script, msg.Issue, msg.Comment)
	})
	return err
}
