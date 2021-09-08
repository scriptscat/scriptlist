package errs

import "net/http"

var ErrUserIsBan = NewError(http.StatusForbidden, 2001, "用户封禁")
var ErrNotLogin = NewError(http.StatusForbidden, 2002, "请先登录")
