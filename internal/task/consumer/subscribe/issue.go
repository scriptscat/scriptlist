package subscribe

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/codfrm/cago/pkg/logger"
	issue2 "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/issue_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type Issue struct {
}

func (s *Issue) Subscribe(ctx context.Context, broker broker.Broker) error {
	_, err := broker.Subscribe(ctx, producer.IssueCreateTopic, s.issueCreate)
	if err != nil {
		return err
	}
	_, err = broker.Subscribe(ctx, producer.CommentCreateTopic, s.commentCreate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Issue) issueCreate(ctx context.Context, event broker.Event) error {
	issue, err := producer.ParseIssueCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	list, err := script_repo.ScriptWatch().FindAll(ctx, issue.ScriptID, script_entity.ScriptWatchLevelIssue)
	if err != nil {
		logger.Ctx(ctx).Error("获取关注列表错误", zap.Int64("issue_id", issue.ID), zap.Error(err))
		return nil
	}
	uids := make([]int64, 0)
	for _, v := range list {
		if v.Level == script_entity.ScriptWatchLevelIssueComment {
			// 关注issue评论
			if _, err := issue_svc.Issue().Watch(auth_svc.Auth().SetCtxUid(ctx, v.UserID), &issue2.WatchRequest{
				ScriptID: v.ScriptID,
				IssueID:  issue.ID,
				Watch:    true,
			}); err != nil {
				logger.Ctx(ctx).Error("设置关注反馈评论错误",
					zap.Int64("user", v.UserID), zap.Int64("issue", issue.ID), zap.Error(err))
			}
		}
		uids = append(uids, v.UserID)
	}
	// 通知关注人
	return notice_svc.Notice().MultipleSend(ctx, uids, notice_svc.IssueCreateTemplate,
		notice_svc.WithParams(&template.IssueCreate{
			Name:    "",
			Title:   "",
			Content: "",
		}))
}

func (s *Issue) commentCreate(ctx context.Context, event broker.Event) error {
	comment, err := producer.ParseCommentCreateMsg(event.Message())
	if err != nil {
		logger.Ctx(ctx).
			Error("json.Unmarshal", zap.Error(err), zap.String("body", string(event.Message().Body)))
		return err
	}
	// 通知反馈关注人
	list, err := issue_repo.Watch().FindAll(ctx, comment.Issue.ID)
	if err != nil {
		logger.Ctx(ctx).Error("获取反馈关注人错误", zap.Int64("issue", comment.Issue.ID), zap.Error(err))
	} else {
		uids := make([]int64, 0)
		for _, v := range list {
			uids = append(uids, v.UserID)
		}
		return notice_svc.Notice().MultipleSend(ctx, uids, notice_svc.CommentCreateTemplate,
			notice_svc.WithParams(&template.IssueComment{
				Name: "",
			}))
	}
	return nil
}
