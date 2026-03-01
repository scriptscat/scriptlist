package subscribe

import (
	"context"

	"github.com/cago-frame/cago/pkg/logger"
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/report_entity"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc/template"
	"github.com/scriptscat/scriptlist/internal/task/producer"
	"go.uber.org/zap"
)

type Report struct{}

func (s *Report) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeReportCreate(ctx, s.reportCreate); err != nil {
		return err
	}
	if err := producer.SubscribeReportCommentCreate(ctx, s.reportCommentCreate); err != nil {
		return err
	}
	return nil
}

func (s *Report) reportCreate(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport) error {
	uids := make([]int64, 0)

	// 通知脚本作者
	uids = append(uids, script.UserID)

	// 通知系统管理员
	admins, err := user_repo.User().FindAdmins(ctx)
	if err != nil {
		logger.Ctx(ctx).Error("获取管理员列表错误", zap.Error(err))
	} else {
		for _, admin := range admins {
			uids = append(uids, admin.ID)
		}
	}

	// 获取举报原因名称
	reasonName := report.Reason
	if reason, ok := report_entity.ReasonMap[report.Reason]; ok {
		reasonName = reason.Name
	}

	return notification_svc.Notification().MultipleSend(ctx, uids, notification_entity.ReportCreateTemplate,
		notification_svc.WithParams(&template.ReportCreate{
			ScriptID: script.ID,
			ReportID: report.ID,
			Name:     script.Name,
			Reason:   reasonName,
			Content:  report.Content,
		}), notification_svc.WithFrom(report.UserID))
}

func (s *Report) reportCommentCreate(ctx context.Context, script *script_entity.Script, report *report_entity.ScriptReport, comment *report_entity.ScriptReportComment) error {
	uids := make([]int64, 0)

	// 通知举报创建者
	uids = append(uids, report.UserID)

	// 通知脚本作者
	if script.UserID != report.UserID {
		uids = append(uids, script.UserID)
	}

	// 通知系统管理员
	admins, err := user_repo.User().FindAdmins(ctx)
	if err != nil {
		logger.Ctx(ctx).Error("获取管理员列表错误", zap.Error(err))
	} else {
		for _, admin := range admins {
			uids = append(uids, admin.ID)
		}
	}

	if err := notification_svc.Notification().MultipleSend(ctx, uids, notification_entity.ReportCommentTemplate,
		notification_svc.WithParams(&template.ReportComment{
			ScriptID:  script.ID,
			ReportID:  report.ID,
			CommentID: comment.ID,
			Name:      script.Name,
			Content:   comment.Content,
		}), notification_svc.WithFrom(comment.UserID)); err != nil {
		logger.Ctx(ctx).Error("发送举报评论通知错误", zap.Int64("report", report.ID), zap.Error(err))
	}

	return nil
}
