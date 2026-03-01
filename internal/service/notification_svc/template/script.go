package template

import "fmt"

const (
	ScriptUpdateTitle   = "[{{.Value.Name}}]有新的版本:{{.Value.Version}}"
	ScriptUpdateContent = `
脚本{{.Value.Name}}更新到{{.Value.Version}}版本
<hr/>
<a href="{{.Config.Url}}/script-show-page/{{.Value.ID}}">点击查看脚本页面</a><hr/>您可以在<a href="{{.Config.Url}}/users/notify">个人设置页面</a>中取消本邮件的通知,或者取消对该脚本的关注
`
)

type ScriptUpdate struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"` // 脚本名
	Version string `json:"version"`
}

func (s *ScriptUpdate) Link() string {
	return fmt.Sprintf("/script-show-page/%d", s.ID)
}

const (
	AccessInviteTitle   = `邀请您加入脚本:{{.Value.Name}}`
	AccessInviteContent = `
{{.Value.Username}}邀请您加入脚本:{{.Value.Name}}
<hr/>
<a href="{{.Config.Url}}/script/invite/?code={{.Value.Code}}">点击此链接加入</a>
`
)

type AccessInvite struct {
	Code     string `json:"code"`
	Name     string `json:"name"`     // 脚本名
	Username string `json:"username"` // 邀请人
}

func (s *AccessInvite) Link() string {
	return fmt.Sprintf("/script/invite/?code=%s", s.Code)
}

const (
	ScriptDeleteTitle   = "[{{.Value.Name}}] 您的脚本已被管理员删除"
	ScriptDeleteContent = `
您的脚本 <b>{{.Value.Name}}</b> 已被管理员删除<br/>
删除原因: {{.Value.Reason}}<br/>
<hr/>
如有疑问请联系管理员
`
)

type ScriptDelete struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func (s *ScriptDelete) Link() string {
	return ""
}
