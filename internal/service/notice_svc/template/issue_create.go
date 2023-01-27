package template

const (
	IssueCreateTitle    = "[{{.Value.Name}}]有新的版本:{{.Value.Version}}"
	IssueCreateTemplate = `
脚本{{.Value.Name}}更新
`
)

type IssueCreate struct {
	Name    string // 脚本名
	Title   string
	Content string
}
