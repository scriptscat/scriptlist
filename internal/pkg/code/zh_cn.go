package code

import "github.com/codfrm/cago/pkg/i18n"

func init() {
	i18n.Register(i18n.DefaultLang, zhCN)
}

var zhCN = map[int]string{
	UserNotFound:      "用户不存在",
	UserIsBanned:      "用户已被禁用",
	UserNotPermission: "用户不允许操作",
	UserNotFollow:     "用户未关注",
	UserNotFollowSelf: "不能关注自己",
	UserExistFollow:   "已关注",

	ScriptNameIsEmpty:    "脚本名不能为空",
	ScriptDescIsEmpty:    "脚本描述不能为空",
	ScriptVersionIsEmpty: "脚本版本不能为空",
	ScriptParseFailed:    "脚本解析失败",
	ScriptNotFound:       "脚本不存在",
	ScriptNotActive:      "脚本被删除",
	ScriptVersionExist:   "版本已存在,更新脚本内容必须升级版本号",
	ScriptCreateFailed:   "脚本创建失败",
	ScriptUpdateFailed:   "脚本更新失败",
	ScriptNotAllowUrl:    "不允许的url域名,如果你需要添加,可以前往论坛反馈申请",
	ScriptIsArchive:      "脚本已归档,无法进行此操作",
}
