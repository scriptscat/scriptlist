package notice_svc

import (
	"github.com/scriptscat/scriptlist/internal/service/notice_svc/template"
)

// 模板id

const (
	ScriptUpdateTemplate = iota + 1
)

var templateMap = map[int]map[SenderType]string{
	ScriptUpdateTemplate: {
		MailSender: template.ScriptUpdateTemplate,
	},
}
