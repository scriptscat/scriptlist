package template

import (
	"github.com/scriptscat/scriptlist/internal/model/entity/notification_entity"
	"github.com/scriptscat/scriptlist/internal/service/notification_svc/sender"
)

// Template 通知模板
type Template struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

var TplMap = map[notification_entity.Type]map[sender.Type]Template{
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
		sender.InAppSender: {
			Content: `
  {{- if eq .Value.Type 1 -}}
    issue.comment.reply.content
  {{- else if eq .Value.Type 4 -}}
    issue.comment.open.content
  {{- else if eq .Value.Type 5 -}}
    issue.comment.close.content
  {{- end -}}
`,
		},
		sender.MailSender: {
			Title:   IssueCommentTitle,
			Content: IssueCommentContent,
		},
	},
	notification_entity.ScriptScoreTemplate: {
		sender.InAppSender: {
			Content: "script.score.content",
		},
		sender.MailSender: {
			Title:   ScriptScoreTitle,
			Content: ScriptScoreContent,
		},
	},
	notification_entity.AccessInviteTemplate: {
		sender.InAppSender: {
			Content: "access.invite.content",
		},
		sender.MailSender: {
			Title:   AccessInviteTitle,
			Content: AccessInviteContent,
		},
	},
	notification_entity.ScriptScoreReplyTemplate: {
		sender.InAppSender: {
			Content: "script.score.reply.content",
		},
		sender.MailSender: {
			Title:   ScriptScoreReplyTitle,
			Content: ScriptScoreReplyContent,
		},
	},
}
