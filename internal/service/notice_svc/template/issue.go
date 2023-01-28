package template

const (
	IssueCreateTitle    = "[{{.Value.Name}}]{{.Value.Title}}"
	IssueCreateTemplate = `
  {{.Value.Content}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/issue/{{.Value.IssueID}}">点击查看原文</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知,或者取消对该脚本反馈的关注
`
)

type IssueCreate struct {
	ScriptID int64
	IssueID  int64
	Name     string // 脚本名
	Title    string // 反馈标题
	Content  string // 反馈内容
}
