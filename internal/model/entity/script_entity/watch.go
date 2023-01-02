package script_entity

type ScriptWatchLevel int

const (
	ScriptWatchLevelNone ScriptWatchLevel = iota + 1
	ScriptWatchLevelVersion
	ScriptWatchLevelIssue
	ScriptWatchLevelIssueComment
)

type Watch struct {
	UserID int64 `json:"user_id"`
	// Watch级别,0:未监听 1:版本更新监听 2:新建issue监听 3:评论都监听
	Level ScriptWatchLevel `json:"level"`
}
