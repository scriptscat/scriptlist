package subscriber

import (
	"context"

	application2 "github.com/scriptscat/scriptlist/internal/service/issue/application"
	service2 "github.com/scriptscat/scriptlist/internal/service/notify/service"
	"github.com/scriptscat/scriptlist/internal/service/script/application"
	"github.com/scriptscat/scriptlist/internal/service/script/broker"
	service3 "github.com/scriptscat/scriptlist/internal/service/user/service"
)

type UserSubscriber struct {
	notifySvc           service2.Sender
	scriptWatchSvc      application.ScriptWatch
	scriptIssueWatchSvc application2.ScriptIssueWatch
	scriptIssue         application2.Issue
	scriptSvc           application.Script
	userSvc             service3.User
}

func NewUserSubscriber(notifySvc service2.Sender, scriptWatchSvc application.ScriptWatch,
	scriptIssueWatchSvc application2.ScriptIssueWatch, scriptIssue application2.Issue, scriptSvc application.Script, userSvc service3.User) *UserSubscriber {
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
