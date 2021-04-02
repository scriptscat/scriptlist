package errs

import "net/http"

var ErrScriptNotFound = NewError(http.StatusNotFound, 1000, "脚本已删除或不存在")
var ErrScriptAudit = NewError(http.StatusForbidden, 1001, "脚本审核中")
var ErrScriptCodeIsNil = NewError(http.StatusNotFound, 1002, "没有任何脚本代码")
var ErrScriptCodeNotFound = NewError(http.StatusNotFound, 1003, "没有找到脚本代码")
