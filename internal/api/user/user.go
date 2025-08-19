package user

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
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

// GetUserDetailRequest 获取用户详细信息
type GetUserDetailRequest struct {
	mux.Meta `path:"/users/:uid/detail" method:"GET"`
	UID      int64 `uri:"uid" binding:"required"`
}

type BadgeItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetUserDetailResponse struct {
	*InfoResponse `json:",inline"`
	// 成就徽章
	Badge []*BadgeItem `json:"badge"`
	// 个人简介
	Description string `json:"description"`
	// 加入时间
	JoinTime int64 `json:"join_time"`
	// 最后活跃
	LastActive int64 `json:"last_active"`
	// 位置
	Location string `json:"location"`
	// 个人网站
	Website string `json:"website"`
	// 邮箱
	Email string `json:"email"`
	// 是否关注
	IsFollow bool `json:"is_follow"`
	// 粉丝
	Followers int64 `json:"followers"`
	// 关注
	Following int64 `json:"following"`
}

// UpdateUserDetailRequest 更新用户信息
type UpdateUserDetailRequest struct {
	mux.Meta `path:"/users/detail" method:"PUT"`
	// 个人简介
	Description string `json:"description"`
	// 位置
	Location string `json:"location"`
	// 个人网站
	Website string `json:"website"`
	// 邮箱
	Email string `json:"email" binding:"email"`
}

type UpdateUserDetailResponse struct {
}

// UpdateUserAvatarRequest 更新用户头像
type UpdateUserAvatarRequest struct {
	mux.Meta `path:"/users/avatar" method:"PUT"`
}

type UpdateUserAvatarResponse struct {
	Url string `json:"url"` // 新头像的URL
}

type LogoutRequest struct {
	mux.Meta `path:"/users/logout" method:"GET"`
}

type LogoutResponse struct{}

// ScriptRequest 用户脚本列表
type ScriptRequest struct {
	mux.Meta              `path:"/users/:uid/scripts" method:"GET"`
	httputils.PageRequest `form:",inline"`
	UID                   int64  `uri:"uid" binding:"required"`
	Keyword               string `form:"keyword"`
	ScriptType            int    `form:"script_type,default=0" binding:"oneof=0 1 2 3 4"` // 0:全部 1: 脚本 2: 库 3: 后台脚本 4: 定时脚本
	Sort                  string `form:"sort,default=today_download" binding:"oneof=today_download total_download score createtime updatetime"`
}

type ScriptResponse struct {
	httputils.PageResponse[*script.Script] `json:",inline"`
}

// GetFollowRequest 获取用户关注信息
type GetFollowRequest struct {
	mux.Meta `path:"/users/:uid/follow" method:"GET"`
	UID      int64 `uri:"uid" binding:"required"`
}

type GetFollowResponse struct {
	// 是否关注
	IsFollow bool `json:"is_follow"`
	// 粉丝
	Followers int64 `json:"followers"`
	// 关注
	Following int64 `json:"following"`
}

// FollowRequest 关注用户
type FollowRequest struct {
	mux.Meta `path:"/users/:uid/follow" method:"POST"`
	UID      int64 `uri:"uid" binding:"required"`
	Unfollow bool  `form:"unfollow"`
}

type FollowResponse struct {
}

// GetWebhookRequest 获取webhook配置
type GetWebhookRequest struct {
	mux.Meta `path:"/users/webhook" method:"GET"`
}

type GetWebhookResponse struct {
	Token string `json:"token"`
}

// RefreshWebhookRequest 刷新webhook配置
type RefreshWebhookRequest struct {
	mux.Meta `path:"/users/webhook" method:"PUT"`
}

type RefreshWebhookResponse struct {
	Token string `json:"token"`
}

// GetConfigRequest 获取用户配置
type GetConfigRequest struct {
	mux.Meta `path:"/users/config" method:"GET"`
}

type GetConfigResponse struct {
	// 邮件通知配置
	Notify *user_entity.Notify `json:"notify"`
}

// UpdateConfigRequest 更新用户配置
type UpdateConfigRequest struct {
	mux.Meta `path:"/users/config" method:"PUT"`
	Notify   *user_entity.Notify `json:"notify" binding:"required"`
}

type UpdateConfigResponse struct {
}

// SearchRequest 搜索用户
type SearchRequest struct {
	mux.Meta `path:"/users/search" method:"GET"`
	Query    string `form:"query" binding:"required" label:"关键词"`
}

type SearchResponse struct {
	Users []*InfoResponse `json:"users"`
}

// RefreshTokenRequest 刷新用户token
type RefreshTokenRequest struct {
	mux.Meta `path:"/users/refresh-token" method:"POST"`
}

type RefreshTokenResponse struct {
}
