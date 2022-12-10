package code

import "github.com/codfrm/cago/pkg/i18n"

func init() {
	i18n.Register(i18n.DefaultLang, zhCN)
}

var zhCN = map[int]string{
	UserNotFound: "用户不存在",
	UserIsBanned: "用户已被禁用",

	ScriptNameIsEmpty:    "脚本名不能为空",
	ScriptDescIsEmpty:    "脚本描述不能为空",
	ScriptVersionIsEmpty: "脚本版本不能为空",
	ScriptParseFailed:    "脚本解析失败",
	ScriptCreateFailed:   "脚本创建失败",
	ScriptUpdateFailed:   "脚本更新失败",
}
