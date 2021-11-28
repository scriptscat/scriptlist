package subscriber

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/scriptscat/scriptlist/internal/domain/issue/broker"
	service4 "github.com/scriptscat/scriptlist/internal/domain/issue/service"
	service2 "github.com/scriptscat/scriptlist/internal/domain/notify/service"
	"github.com/scriptscat/scriptlist/internal/domain/script/service"
	service3 "github.com/scriptscat/scriptlist/internal/domain/user/service"
	"github.com/scriptscat/scriptlist/internal/pkg/config"
	"github.com/scriptscat/scriptlist/pkg/event"
	"github.com/sirupsen/logrus"
)

type NotifySubscriber struct {
	notifySvc           service2.Sender
	scriptWatchSvc      service.ScriptWatch
	scriptIssueWatchSvc service4.ScriptIssueWatch
	scriptIssue         service4.Issue
	scriptSvc           service.Script
	userSvc             service3.User
}

func NewNotifySubscriber() *NotifySubscriber {
	return &NotifySubscriber{}
}

func (n *NotifySubscriber) Subscribe() error {

	if _, err := event.DefaultBroker.Subscribe(service.EventScriptVersionUpdate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return n.NotifyScriptUpdate(ids["script"], ids["code"])
	}); err != nil {
		return err
	}

	if _, err := event.DefaultBroker.Subscribe(broker.EventScriptIssueCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return n.NotifyScriptUpdate(ids["script"], ids["code"])
	}); err != nil {
		return err
	}

	if _, err := event.DefaultBroker.Subscribe(broker.EventScriptIssueCommentCreate, func(message *event.Message) error {
		ids := event.Ids{}
		if err := json.Unmarshal(message.Body, &ids); err != nil {
			return err
		}
		return n.NotifyScriptUpdate(ids["script"], ids["code"])
	}); err != nil {
		return err
	}

	return nil
}

func (n *NotifySubscriber) NotifyScriptUpdate(script, code int64) error {
	scriptInfo, err := n.scriptSvc.Info(script)
	if err != nil {
		return err
	}
	codeInfo, err := n.scriptSvc.GetCode(code)
	if err != nil {
		return err
	}
	list, err := n.scriptWatchSvc.WatchList(script)
	if err != nil {
		return err
	}
	user, err := n.userSvc.UserInfo(scriptInfo.UserId)
	if err != nil {
		return err
	}
	for uid, _ := range list {
		u, err := n.userSvc.SelfInfo(uid)
		if err != nil {
			continue
		}
		n.notifySvc.SendEmailFrom(user.Username, u.Email, "["+scriptInfo.Name+"]有新的版本: "+codeInfo.Version,
			fmt.Sprintf(scriptInfo.Name+"升级到了:%s<br/><a href=\"%s\">点击查看脚本页面</a>", codeInfo.Version,
				config.AppConfig.FrontendUrl+"script-show-page/"+strconv.FormatInt(scriptInfo.ID, 10),
			), "text/html")
	}

	return nil
}

func (n *NotifySubscriber) NotifyScriptIssueCreate(script, issue int64) error {
	scriptInfo, err := n.scriptSvc.Info(script)
	if err != nil {
		return err
	}
	issueInfo, err := n.scriptIssue.GetIssue(issue)
	if err != nil {
		return err
	}
	list, err := n.scriptWatchSvc.WatchList(script)
	if err != nil {
		return err
	}
	user, err := n.userSvc.UserInfo(issueInfo.UserID)
	if err != nil {
		return err
	}
	for uid, level := range list {
		if level < 1 {
			continue
		}
		if level == 2 {
			_ = n.scriptIssueWatchSvc.Watch(issue, uid)
		}
		_ = n.scriptIssueWatchSvc.Watch(issue, issueInfo.UserID)
		u, err := n.userSvc.SelfInfo(uid)
		if err != nil {
			continue
		}
		if err := n.notifySvc.SendEmailFrom(user.Username, u.Email, "["+scriptInfo.Name+"]"+issueInfo.Title, issueInfo.Content+
			fmt.Sprintf("<hr/><br/><a href=\"%s\">点击查看原文</a>",
				config.AppConfig.FrontendUrl+"script-show-page/"+strconv.FormatInt(issueInfo.ID, 10)+"/issue/"+strconv.FormatInt(issueInfo.ID, 10),
			), "text/html"); err != nil {
			logrus.Errorf("sendemail: %v", err)
		}
	}
	return nil
}
