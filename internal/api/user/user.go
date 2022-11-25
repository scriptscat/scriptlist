package user

import "github.com/codfrm/cago/server/mux"

// InfoRequest 用户信息
type InfoRequest struct {
	mux.Route `path:"/user/info" method:"GET"`
}

type InfoResponse struct {
}
