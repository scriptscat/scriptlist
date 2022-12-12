package template

const (
	ScriptUpdateTemplate = `
脚本{{.Value.Name}}更新
`
)

type ScriptUpdate struct {
	Name string // 脚本名
}
