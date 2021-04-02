package errs

import "net/http"

var ErrUserIsBan = NewError(http.StatusForbidden, 1000, "用户封禁")
var ErrNotLogin = NewError(http.StatusForbidden, 1000, "请先登录")
