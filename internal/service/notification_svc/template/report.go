package template

import "fmt"

const (
	ReportCreateTitle   = "[{{.Value.Name}}] 收到新的举报"
	ReportCreateContent = `
  脚本 <b>{{.Value.Name}}</b> 收到了新的举报<br/>
  举报原因: {{.Value.Reason}}<br/>
  {{.Value.Content}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/report/{{.Value.ReportID}}">点击查看详情</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知
`
)

type ReportCreate struct {
	ScriptID int64  `json:"script_id"`
	ReportID int64  `json:"report_id"`
	Name     string `json:"name"`    // 脚本名
	Reason   string `json:"reason"`  // 举报原因
	Content  string `json:"content"` // 举报内容
}

func (r *ReportCreate) Link() string {
	return fmt.Sprintf("/script-show-page/%d/report/%d",
		r.ScriptID, r.ReportID)
}

const (
	ReportCommentTitle   = "[{{.Value.Name}}] 举报有新回复"
	ReportCommentContent = `
  {{.Value.Content}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/report/{{.Value.ReportID}}">点击查看详情</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知
`
)

type ReportComment struct {
	ScriptID  int64  `json:"script_id"`
	ReportID  int64  `json:"report_id"`
	CommentID int64  `json:"comment_id"`
	Name      string `json:"name"`    // 脚本名
	Content   string `json:"content"` // 评论内容
}

func (r *ReportComment) Link() string {
	return fmt.Sprintf("/script-show-page/%d/report/%d",
		r.ScriptID, r.ReportID)
}
