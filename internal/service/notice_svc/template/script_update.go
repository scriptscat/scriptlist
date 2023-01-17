package template

const (
	ScriptUpdateTitle    = "[{{.Value.Name}}]有新的版本:{{.Value.Version}}"
	ScriptUpdateTemplate = `
脚本{{.Value.Name}}更新
`
)

type ScriptUpdate struct {
	ID      int64
	Name    string // 脚本名
	Version string
}
