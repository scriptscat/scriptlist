package user

import "github.com/codfrm/cago/server/mux"

// CurrentUserRequest 获取当前登录的用户信息
type CurrentUserRequest struct {
	mux.Meta `path:"/user/info" method:"GET"`
}

type CurrentUserResponse struct {
	*InfoResponse `json:",inline"`
}

// InfoRequest 获取指定用户信息
type InfoRequest struct {
	mux.Meta `path:"/user/info/:uid" method:"GET"`
	UID      int64 `uri:"uid" binding:"required"`
}

type InfoResponse struct {
	UID      int64  `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
