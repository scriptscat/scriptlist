package producer

import (
	"context"
	"encoding/json"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
)

func PublishIssueCreate(ctx context.Context, issue *issue_entity.ScriptIssue) error {
	body, err := json.Marshal(issue)
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, IssueCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseIssueCreateMsg(msg *broker2.Message) (*issue_entity.ScriptIssue, error) {
	ret := &issue_entity.ScriptIssue{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type CommentCreateMsg struct {
	Issue   *issue_entity.ScriptIssue
	Comment *issue_entity.ScriptIssueComment
}

func PublishCommentCreate(ctx context.Context, issue *issue_entity.ScriptIssue, comment *issue_entity.ScriptIssueComment) error {
	body, err := json.Marshal(&CommentCreateMsg{
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
