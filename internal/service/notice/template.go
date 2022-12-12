package notice

import "github.com/scriptscat/scriptlist/internal/service/notice/template"

// 模板id

const (
	ScriptUpdateTemplate = iota + 1
)

var templateMap = map[int]map[SenderType]string{
	ScriptUpdateTemplate: {
		MailSender: template.ScriptUpdateTemplate,
	},
}
