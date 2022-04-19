package errs

import "net/http"

var (
	ErrUserIsBan            = NewError(http.StatusForbidden, 2001, "用户封禁")
	ErrNotLogin             = NewError(http.StatusForbidden, 2002, "请先登录")
	ErrTokenNotFound        = NewError(http.StatusNotFound, 2003, "用户token未找到")
	ErrUserNotFound         = NewError(http.StatusNotFound, 2004, "用户不存在")
	ErrUserEmailNotVerified = NewBadRequestError(2005, "请先验证邮箱")
)
