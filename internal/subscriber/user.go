package subscriber

import (
	"context"

	service5 "github.com/scriptscat/scriptlist/internal/service/issue/service"
	service2 "github.com/scriptscat/scriptlist/internal/service/notify/service"
	"github.com/scriptscat/scriptlist/internal/service/script/broker"
	service6 "github.com/scriptscat/scriptlist/internal/service/script/service"
	service3 "github.com/scriptscat/scriptlist/internal/service/user/service"
)

type UserSubscriber struct {
	notifySvc           service2.Sender
	scriptWatchSvc      service6.ScriptWatch
	scriptIssueWatchSvc service5.ScriptIssueWatch
	scriptIssue         service5.Issue
	scriptSvc           service6.Script
	userSvc             service3.User
}

func NewUserSubscriber(notifySvc service2.Sender, scriptWatchSvc service6.ScriptWatch,
	scriptIssueWatchSvc service5.ScriptIssueWatch, scriptIssue service5.Issue, scriptSvc service6.Script, userSvc service3.User) *UserSubscriber {
	return &UserSubscriber{
		notifySvc:           notifySvc,
		scriptWatchSvc:      scriptWatchSvc,
		scriptIssueWatchSvc: scriptIssueWatchSvc,
		scriptIssue:         scriptIssue,
		scriptSvc:           scriptSvc,
		userSvc:             userSvc,
	}
}

func (n *UserSubscriber) Subscribe(ctx context.Context) error {

	if _, err := broker.SubscribeEventScriptCreate(n.ScriptCreate); err != nil {
		return err
	}

	return nil
}

func (n *UserSubscriber) ScriptCreate(script int64) error {

	return nil
}
