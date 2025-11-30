package template

const (
	ScriptScoreTitle    = `收到评分:[{{.Value.Name}}]`
	ScriptScoreTemplate = `
  {{.Value.Name}} 被 {{.Value.Username}} 评分为 {{.Value.Score}} 分
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/comment">点击查看</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知
`
)

type ScriptScore struct {
	ScriptID int64
	Name     string // 脚本名
	Username string // 评分用户
	Score    int    // 分数
}

const (
	ScriptScoreReplyTitle    = `收到作者回复评分:[{{.Value.Name}}]`
	ScriptScoreReplyTemplate = `
 在 {{.Value.Name}} 收到 脚本作者 回复消息 : 
  <hr/>
     {{.Value.Content}}
  <hr/>
  <a href="{{.Config.Url}}/script-show-page/{{.Value.ScriptID}}/comment">点击查看</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知
`
)

type ScriptReplyScore struct {
	ScriptID int64
	Name     string // 脚本名
	Content  string // 分数
}
