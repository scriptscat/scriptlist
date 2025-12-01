package template

import (
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc/sender"
)

// Template 通知模板
type Template struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var TemplateMap = map[notification_entity.Type]map[sender.Type]Template{
	// 脚本更新模板
	notification_entity.ScriptUpdateTemplate: {
		sender.InAppSender: {
			Content: "script.update.content",
		},
		sender.MailSender: {
			Title:   ScriptUpdateTitle,
			Content: ScriptUpdateContent,
		},
	},
	// 问题创建模板
	notification_entity.IssueCreateTemplate: {
		sender.InAppSender: {
			Content: "issue.create.content",
		},
		sender.MailSender: {
			Title:   IssueCreateTitle,
			Content: IssueCreateContent,
		},
	},
	// 评论创建模板
	notification_entity.CommentCreateTemplate: {
		sender.MailSender: {
			Title:   IssueCommentTitle,
			Content: IssueCommentContent,
		},
	},
	notification_entity.ScriptScoreTemplate: {
		sender.MailSender: {
			Title:   ScriptScoreTitle,
			Content: ScriptScoreContent,
		},
	},
	notification_entity.AccessInviteTemplate: {
		sender.MailSender: {
			Title:   AccessInviteTitle,
			Content: AccessInviteContent,
		},
	},
	notification_entity.ScriptScoreReplyTemplate: {
		sender.MailSender: {
			Title:   ScriptScoreReplyTitle,
			Content: ScriptScoreReplyContent,
		},
	},
}
