package user

import (
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model"
)

// CurrentUserRequest 获取当前登录的用户信息
type CurrentUserRequest struct {
	mux.Meta `path:"/users" method:"GET"`
}

type CurrentUserResponse struct {
	*InfoResponse `json:",inline"`
}

// InfoRequest 获取指定用户信息
type InfoRequest struct {
	mux.Meta `path:"/users/:uid/info" method:"GET"`
	UID      int64 `uri:"uid" binding:"required"`
}

type InfoResponse struct {
	UID         int64            `json:"uid"`
	Username    string           `json:"username"`
	Avatar      string           `json:"avatar"`
	IsAdmin     model.AdminLevel `json:"is_admin"`
	EmailStatus int64            `json:"email_status"`
}
