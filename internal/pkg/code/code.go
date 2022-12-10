package code

const (
	UserIsBanned = iota + 100000
	UserNotFound

	// script
	ScriptNameIsEmpty = iota + 101000
	ScriptDescIsEmpty
	ScriptVersionIsEmpty
	ScriptParseFailed
	ScriptCreateFailed
	ScriptUpdateFailed
)
