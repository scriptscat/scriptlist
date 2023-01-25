package user

import (
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/api/script"
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
	UserID      int64            `json:"user_id"`
	Username    string           `json:"username"`
	Avatar      string           `json:"avatar"`
	IsAdmin     model.AdminLevel `json:"is_admin"`
	EmailStatus int64            `json:"email_status"`
}

// ScriptRequest 用户脚本列表
type ScriptRequest struct {
	mux.Meta              `path:"/users/:uid/script" method:"GET"`
	httputils.PageRequest `form:",inline"`
	UID                   int64  `uri:"uid" binding:"required"`
	Keyword               string `form:"keyword"`
	ScriptType            int    `form:"script_type,default=0" binding:"oneof=0 1 2 3 4"` // 0:全部 1: 脚本 2: 库 3: 后台脚本 4: 定时脚本
	Sort                  string `form:"sort,default=today_download" binding:"oneof=today_download total_download score createtime updatetime"`
}

type ScriptResponse struct {
	httputils.PageResponse[*script.Script] `json:",inline"`
}
