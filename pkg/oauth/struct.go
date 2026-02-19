package oauth

type ErrorRespond struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *ErrorRespond) Error() string {
	return e.Msg
}

type AccessTokenRespond struct {
	ErrorRespond
	AccessToken string `json:"access_token"` // #nosec G117 -- 这是响应字段定义，非硬编码token
}

type UserRespond struct {
	ErrorRespond
	User struct {
		UID      string `json:"uid"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Email    string `json:"email"`
	} `json:"user"`
}
