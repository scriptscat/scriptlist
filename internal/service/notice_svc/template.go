package notice_svc

import (
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/sender"
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
)

// 模板id

const (
	ScriptUpdateTemplate = iota + 1
)

type Template struct {
	Title    string
	Template string
}

var templateMap = map[int]map[sender.Type]Template{
	// 脚本更新模板
	ScriptUpdateTemplate: {
		sender.MailSender: {
			Title:    template.ScriptUpdateTitle,
			Template: template.ScriptUpdateTemplate,
		},
	},
}
