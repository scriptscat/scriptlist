package template

const (
	ScriptScoreTitle    = `收到评分:[{{.Value.Name}}]`
	ScriptScoreTemplate = `
  {{.Value.Name}} 被 {{.Value.Username}} 评分为 {{.Value.Score}} 分
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/score">点击查看</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知
`
)

type ScriptScore struct {
	ScriptID int64
	Name     string // 脚本名
	Username string // 评分用户
	Score    int    // 分数
}
