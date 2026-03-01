package producer

import (
	"context"
	"encoding/json"

	"github.com/cago-frame/cago/pkg/broker"
	broker2 "github.com/cago-frame/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
)

type ReportCreateMsg struct {
	Script *script_entity.Script
	Report *report_entity.ScriptReport
}

func PublishReportCreate(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport) error {
	body, err := json.Marshal(&ReportCreateMsg{
		Script: script,
		Report: report,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ReportCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseReportCreateMsg(msg *broker2.Message) (*ReportCreateMsg, error) {
	ret := &ReportCreateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeReportCreate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport) error) error {
	_, err := broker.Default().Subscribe(ctx, ReportCreateTopic, func(ctx context.Context, event broker2.Event) error {
		msg, err := ParseReportCreateMsg(event.Message())
		if err != nil {
			return err
		}
		return fn(ctx, msg.Script, msg.Report)
	})
	return err
}

type ReportCommentCreateMsg struct {
	Script  *script_entity.Script
	Report  *report_entity.ScriptReport
	Comment *report_entity.ScriptReportComment
}

func PublishReportCommentCreate(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport, comment *report_entity.ScriptReportComment) error {
	body, err := json.Marshal(&ReportCommentCreateMsg{
		Script:  script,
		Report:  report,
		Comment: comment,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ReportCommentCreateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseReportCommentCreateMsg(msg *broker2.Message) (*ReportCommentCreateMsg, error) {
	ret := &ReportCommentCreateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeReportCommentCreate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport, comment *report_entity.ScriptReportComment) error) error {
	_, err := broker.Default().Subscribe(ctx, ReportCommentCreateTopic, func(ctx context.Context, event broker2.Event) error {
		msg, err := ParseReportCommentCreateMsg(event.Message())
		if err != nil {
			return err
		}
		return fn(ctx, msg.Script, msg.Report, msg.Comment)
	})
	return err
}
