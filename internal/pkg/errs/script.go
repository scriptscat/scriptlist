package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrScriptNotFound         = NewError(http.StatusNotFound, 3001, "脚本已删除或不存在")
	ErrScriptAudit            = NewError(http.StatusForbidden, 3002, "脚本审核中")
	ErrScriptCodeIsNil        = NewError(http.StatusNotFound, 3003, "没有任何脚本代码")
	ErrScriptCodeNotFound     = NewError(http.StatusNotFound, 3004, "没有找到脚本代码")
	ErrScriptForbidden        = NewError(http.StatusForbidden, 3005, "没有脚本访问权限")
	ErrScriptCodeExist        = NewError(http.StatusBadRequest, 3006, "脚本版本已经存在")
	ErrCodeDefinitionNotFound = NewError(http.StatusNotFound, 3007, "代码定义文件未找到")
	ErrScriptArchived         = NewError(http.StatusBadRequest, 3008, "脚本已归档")

	ErrScoreNotFound = NewError(http.StatusNotFound, 4001, "没有找到评分")
)

func NewErrScriptSyncNetwork(url string, err error) error {
	return NewError(http.StatusInternalServerError, 3007, fmt.Sprintf("脚本同步链接无法访问:%s %v", url, err))
}
