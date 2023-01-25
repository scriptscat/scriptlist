package code

// user
const (
	UserIsBanned = iota + 100000
	UserNotFound
	UserNotPermission
)

// script
const (
	ScriptNameIsEmpty = iota + 101000
	ScriptDescIsEmpty
	ScriptVersionIsEmpty
	ScriptParseFailed
	ScriptNotFound
	ScriptNotActive
	ScriptVersionExist
	ScriptCreateFailed
	ScriptUpdateFailed
	ScriptNotAllowUrl
	ScriptIsArchive
)
