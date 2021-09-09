package errs

import "net/http"

var (
	ErrResourceNotImage = NewError(http.StatusBadRequest, 5001, "上传的资源不是图片")
	ErrResourceNotFound = NewError(http.StatusNotFound, 5002, "资源未找到")
)
