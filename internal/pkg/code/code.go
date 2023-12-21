package code

// user
const (
	UserIsBanned = iota + 100000
	UserNotFound
	UserNotPermission
	UserNotFollow
	UserNotFollowSelf
	UserExistFollow
	UserEmailNotVerified
	UserNotLogin
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
	ScriptScoreDeleted
	ScriptScoreNotFound
	ScriptChangePreReleaseNotLatest
	ScriptMustHaveVersion

	WebhookSecretError
	WebhookRepositoryNotFound
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

// resource
const (
	ResourceImageTooLarge = iota + 103000
	ResourceNotImage
	ResourceNotFound
)

// statistics
const (
	StatisticsLimitExceeded = iota + 104000
	StatisticsResultLimit
	StatisticsInfoUninitialized
	StatisticsWhitelistInvalid
	StatisticsWhitelistNotFound
)
