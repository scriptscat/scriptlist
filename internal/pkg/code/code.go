package code

// user
const (
	UserIsBanned = iota + 100000
	UserNotFound
	UserNotPermission
	UserNotFollow
	UserNotFollowSelf
	UserExistFollow
)

// script
const (
	ScriptNameIsEmpty = iota + 101000
	ScriptDescIsEmpty
	ScriptVersionIsEmpty
	ScriptParseFailed
	ScriptNotFound
	ScriptIsDelete
	ScriptVersionExist
	ScriptCreateFailed
	ScriptUpdateFailed
	ScriptNotAllowUrl
	ScriptIsArchive
)

// issue
const (
	IssueLabelNotExist = iota + 102000
	IssueNotFound
	IssueIsDelete
	IssueNoPermission
	IssueCommentNotFound
	IssueLabelNotChange
)
