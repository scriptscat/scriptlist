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
}
