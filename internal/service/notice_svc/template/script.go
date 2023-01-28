package template

const (
	ScriptUpdateTitle    = "[{{.Value.Name}}]有新的版本:{{.Value.Version}}"
	ScriptUpdateTemplate = `
脚本{{.Value.Name}}更新到{{.Value.Version}}版本
<hr/>
<a href="{{.Config.Url}}/script-show-page/{{.Value.ID}}">点击查看脚本页面</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知,或者取消对该脚本的关注
`
)

type ScriptUpdate struct {
	ID      int64
	Name    string // 脚本名
	Version string
}
