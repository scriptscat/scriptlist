package subscriber

import (
	"context"

	service4 "github.com/scriptscat/scriptlist/internal/domain/issue/service"
	service2 "github.com/scriptscat/scriptlist/internal/domain/notify/service"
	"github.com/scriptscat/scriptlist/internal/domain/script/broker"
	"github.com/scriptscat/scriptlist/internal/domain/script/service"
	service3 "github.com/scriptscat/scriptlist/internal/domain/user/service"
)

type UserSubscriber struct {
	notifySvc           service2.Sender
	scriptWatchSvc      service.ScriptWatch
	scriptIssueWatchSvc service4.ScriptIssueWatch
	scriptIssue         service4.Issue
	scriptSvc           service.Script
	userSvc             service3.User
}

func NewUserSubscriber(notifySvc service2.Sender, scriptWatchSvc service.ScriptWatch,
	scriptIssueWatchSvc service4.ScriptIssueWatch, scriptIssue service4.Issue, scriptSvc service.Script, userSvc service3.User) *UserSubscriber {
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
