package template

const (
	IssueCommentTitle    = "[{{.Value.Name}}]有新的版本:{{.Value.Version}}"
	IssueCommentTemplate = `
脚本{{.Value.Name}}更新
`
)

type IssueComment struct {
	Name    string // 脚本名
	Title   string
	Content string
}
