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
	UserWaitVerify
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

	ScriptDeleteReleaseNotLatest
	ScriptCategoryNotFound
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

// access
const (
	AccessAlreadyExist = iota + 10500
	AccessNotFound
)

// group
const (
	GroupNotFound = iota + 106000
	GroupMemberNotFound
	GroupMemberExist
)

// access invite
const (
	AccessInviteNotFound = iota + 107000
	AccessInviteIsAudit
	AccessInviteNotAudit
	AccessInviteNotPending
	AccessInviteExist
	AccessInviteExpired
	AccessInviteUsed
	AccessInviteInvalid
	AccessInviteUserError
)
