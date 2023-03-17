package code

import "github.com/codfrm/cago/pkg/i18n"

func init() {
	i18n.Register(i18n.DefaultLang, zhCN)
}

var zhCN = map[int]string{
	UserNotFound:         "用户不存在",
	UserIsBanned:         "用户已被禁用",
	UserNotPermission:    "用户不允许操作",
	UserNotFollow:        "用户未关注",
	UserNotFollowSelf:    "不能关注自己",
	UserExistFollow:      "已关注",
	UserEmailNotVerified: "邮箱未验证",

	ScriptNameIsEmpty:               "脚本名不能为空",
	ScriptDescIsEmpty:               "脚本描述不能为空",
	ScriptVersionIsEmpty:            "脚本版本不能为空",
	ScriptParseFailed:               "脚本解析失败",
	ScriptNotFound:                  "脚本不存在",
	ScriptIsDelete:                  "脚本被删除",
	ScriptVersionExist:              "版本已存在,更新脚本内容必须升级版本号",
	ScriptCreateFailed:              "脚本创建失败",
	ScriptUpdateFailed:              "脚本更新失败",
	ScriptNotAllowUrl:               "不允许的url域名,如果你需要添加,可以前往论坛反馈申请",
	ScriptIsArchive:                 "脚本已归档,无法进行此操作",
	ScriptScoreDeleted:              "评分已删除",
	ScriptScoreNotFound:             "评分不存在",
	ScriptChangePreReleaseNotLatest: "修改预发布版本失败,没有新的正式版本了",
	ScriptMustHaveVersion:           "脚本必须有有一个版本",

	WebhookSecretError:        "Webhook Secret 错误",
	WebhookRepositoryNotFound: "仓库不存在",

	IssueLabelNotExist:   "标签不存在",
	IssueNotFound:        "反馈不存在",
	IssueIsDelete:        "反馈被删除",
	IssueNoPermission:    "无权限操作",
	IssueCommentNotFound: "评论不存在",
	IssueLabelNotChange:  "标签未改变",

	ResourceImageTooLarge: "图片过大,不能超过1M",
	ResourceNotImage:      "不是图片",
	ResourceNotFound:      "资源不存在",

	StatisticsLimitExceeded:     "统计数据超过限制",
	StatisticsResultLimit:       "统计结果限制1000条数据",
	StatisticsInfoUninitialized: "统计信息未初始化",
	StatisticsWhitelistInvalid:  "统计白名单无效,不支持顶级域名: %s",
	StatisticsWhitelistNotFound: "不在白名单内",
}
