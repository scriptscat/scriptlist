package template

import "fmt"

const (
	IssueCreateTitle   = "[{{.Value.Name}}]{{.Value.Title}}"
	IssueCreateContent = `
  {{.Value.Content}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/issue/{{.Value.IssueID}}">点击查看原文</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知,或者取消对该脚本反馈的关注
`
)

type IssueCreate struct {
	ScriptID int64  `json:"script_id"`
	IssueID  int64  `json:"issue_id"`
	Name     string `json:"name"`    // 脚本名
	Title    string `json:"title"`   // 反馈标题
	Content  string `json:"content"` // 反馈内容
}

func (i *IssueCreate) Link() string {
	return fmt.Sprintf("/script-show-page/%d/issue/%d",
		i.ScriptID, i.IssueID)
}
