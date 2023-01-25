package auth

import "github.com/codfrm/cago/server/mux"

// OAuthCallbackRequest 第三方登录
type OAuthCallbackRequest struct {
	mux.Meta    `path:"/login/oauth" method:"GET"`
	Code        string `form:"code" binding:"required"`
	RedirectUri string `form:"redirect_uri"`
}

type OAuthCallbackResponse struct {
	RedirectUri string `json:"redirect_uri"`
	UID         int64  `json:"uid"`
}
