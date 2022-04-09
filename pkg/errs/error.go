package errs

import "net/http"

type JsonRespondError struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func NewError(status, code int, msg string) error {
	return &JsonRespondError{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
}

func (j *JsonRespondError) Error() string {
	return j.Msg
}

func NewBadRequestError(code int, err string) error {
	return &JsonRespondError{
		Status: http.StatusBadRequest,
		Code:   code,
		Msg:    err,
	}
}

func IsNotFound(err error) bool {
	e, ok := err.(*JsonRespondError)
	if !ok {
		return false
	}
	return e.Status == http.StatusNotFound
}
