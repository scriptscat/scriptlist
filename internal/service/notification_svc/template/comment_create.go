package template

import (
	"fmt"

	"github.com/scriptscat/scriptlist/internal/model/entity/issue_entity"
)

const (
	IssueCommentTitle = `
  {{- if eq .Value.Type 1 -}}
    回复:[{{.Value.Name}}]{{.Value.Title}}
  {{- else if eq .Value.Type 4 -}}
    打开:[{{.Value.Name}}]{{.Value.Title}}
  {{- else if eq .Value.Type 5 -}}
    关闭:[{{.Value.Name}}]{{.Value.Title}}
  {{- end -}}
  `
	IssueCommentContent = `
  {{- if eq .Value.Type 1 -}}
    {{.Value.Content}}
  {{- else if eq .Value.Type 4 -}}
    打开了反馈
  {{- else if eq .Value.Type 5 -}}
    关闭了反馈
  {{- end -}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/issue/{{.Value.IssueID}}#comment-{{.Value.CommentID}}">点击查看原文</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知,或者取消对该脚本反馈评论的关注
`
)

type IssueComment struct {
	ScriptID  int64                    `json:"script_id"`
	IssueID   int64                    `json:"issue_id"`
	CommentID int64                    `json:"comment_id"`
	Name      string                   `json:"name"`    // 脚本名
	Title     string                   `json:"title"`   // 反馈标题
	Content   string                   `json:"content"` // 反馈内容
	Type      issue_entity.CommentType `json:"type"`    // 反馈类型
}

func (i *IssueComment) Link() string {
	return fmt.Sprintf("/script-show-page/%d/issue/%d#comment-%d",
		i.ScriptID, i.IssueID, i.CommentID)
}
