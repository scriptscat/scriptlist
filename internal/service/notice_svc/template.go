package notice_svc

import (
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/sender"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
)

// 模板id

type TemplateID int

const (
	ScriptUpdateTemplate TemplateID = iota + 1
	IssueCreateTemplate
	CommentCreateTemplate
	ScriptScoreTemplate
	AccessInviteTemplate
)

type Template struct {
	Title    string
	Template string
}

var templateMap = map[TemplateID]map[sender.Type]Template{
	// 脚本更新模板
	ScriptUpdateTemplate: {
		sender.MailSender: {
			Title:    template.ScriptUpdateTitle,
			Template: template.ScriptUpdateTemplate,
		},
	},
	// 问题创建模板
	IssueCreateTemplate: {
		sender.MailSender: {
			Title:    template.IssueCreateTitle,
			Template: template.IssueCreateTemplate,
		},
	},
	// 评论创建模板
	CommentCreateTemplate: {
		sender.MailSender: {
			Title:    template.IssueCommentTitle,
			Template: template.IssueCommentTemplate,
		},
	},
	ScriptScoreTemplate: {
		sender.MailSender: {
			Title:    template.ScriptScoreTitle,
			Template: template.ScriptScoreTemplate,
		},
	},
	AccessInviteTemplate: {
		sender.MailSender: {
			Title:    template.AccessInviteTitle,
			Template: template.AccessInviteTemplate,
		},
	},
}
