package errs

import "net/http"

var ErrScriptNotFound = NewError(http.StatusNotFound, 3001, "脚本已删除或不存在")
var ErrScriptAudit = NewError(http.StatusForbidden, 3002, "脚本审核中")
var ErrScriptCodeIsNil = NewError(http.StatusNotFound, 3003, "没有任何脚本代码")
var ErrScriptCodeNotFound = NewError(http.StatusNotFound, 3004, "没有找到脚本代码")
