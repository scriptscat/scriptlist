package subscribe

import (
	"context"

	"github.com/codfrm/cago/pkg/logger"
	issue2 "github.com/scriptscat/scriptlist/internal/api/issue"
	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/issue_repo"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/issue_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type Issue struct{}

func (s *Issue) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeIssueCreate(ctx, s.issueCreate); err != nil {
		return err
	}
	if err := producer.SubscribeCommentCreate(ctx, s.commentCreate); err != nil {
		return err
	}
	return nil
}

func (s *Issue) issueCreate(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue) error {
	list, err := script_repo.ScriptWatch().FindAll(ctx, issue.ScriptID, script_entity.ScriptWatchLevelIssue)
	if err != nil {
		logger.Ctx(ctx).Error("获取关注列表错误", zap.Int64("issue_id", issue.ID), zap.Error(err))
		return nil
	}
	// 评论者关注
	if _, err := issue_svc.Issue().Watch(ctx, issue.UserID, &issue2.WatchRequest{
		ScriptID: issue.ScriptID,
		IssueID:  issue.ID,
		Watch:    true,
	}); err != nil {
		return err
	}
	uids := make([]int64, 0)
	for _, v := range list {
		if v.Level == script_entity.ScriptWatchLevelIssueComment {
			// 关注issue评论
			if _, err := issue_svc.Issue().Watch(ctx, v.UserID, &issue2.WatchRequest{
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
			ScriptID: script.ID,
			IssueID:  issue.ID,
			Name:     script.Name,
			Title:    issue.Title,
			Content:  issue.Content,
		}), notice_svc.WithFrom(issue.UserID))
}

func (s *Issue) commentCreate(ctx context.Context, script *script_entity.Script, issue *issue_entity.ScriptIssue, comment *issue_entity.ScriptIssueComment) error {
	// 通知反馈关注人
	list, err := issue_repo.Watch().FindAll(ctx, issue.ID)
	if err != nil {
		logger.Ctx(ctx).Error("获取反馈关注人错误", zap.Int64("issue", issue.ID), zap.Error(err))
	} else {
		uids := make([]int64, 0)
		for _, v := range list {
			uids = append(uids, v.UserID)
		}
		return notice_svc.Notice().MultipleSend(ctx, uids, notice_svc.CommentCreateTemplate,
			notice_svc.WithParams(&template.IssueComment{
				ScriptID:  script.ID,
				IssueID:   issue.ID,
				CommentID: comment.ID,
				Name:      script.Name,
				Title:     issue.Title,
				Content:   comment.Content,
				Type:      comment.Type,
			}), notice_svc.WithFrom(comment.UserID))
	}
	return nil
}
