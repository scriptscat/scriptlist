package errs

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

type RespondError struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func NewRespondError(status int, msg string) error {
	return &RespondError{
		Status: status,
		Msg:    msg,
	}
}

func (r *RespondError) Error() string {
	return r.Msg
}
