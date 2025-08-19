package script

import (
	"github.com/cago-frame/cago/pkg/utils/httputils"
	"github.com/cago-frame/cago/server/mux"
	"github.com/scriptscat/scriptlist/internal/model/entity/user_entity"
)

// CreateFolderRequest 创建收藏夹
type CreateFolderRequest struct {
	mux.Meta    `path:"/favorites/folders" method:"POST"`
	Name        string `json:"name" binding:"required,max=50" label:"收藏夹名称"`
	Description string `json:"description" binding:"max=200" label:"收藏夹描述"`
	Private     int32  `json:"private" binding:"omitempty,oneof=1 2" label:"私密收藏夹"` // 1私密 2公开
}

type CreateFolderResponse struct {
	ID int64 `json:"id"`
}

// DeleteFolderRequest 删除收藏夹
type DeleteFolderRequest struct {
	mux.Meta `path:"/favorites/folders/:id" method:"DELETE"`
	ID       int64 `uri:"id" binding:"required"`
}

type DeleteFolderResponse struct{}

// FavoriteFolderListRequest 收藏夹列表
type FavoriteFolderListRequest struct {
	mux.Meta `path:"/favorites/folders" method:"GET"`
	UserID   int64 `form:"user_id" label:"用户ID"` // 用户ID，0表示当前登录用户
}

type FavoriteFolderItem struct {
	ID                   int64 `json:"id"`
	user_entity.UserInfo `json:",inline"`
	Name                 string `json:"name"`
	Description          string `json:"description"` // 收藏夹描述
	Count                int64  `json:"count"`       // 收藏夹中脚本数量
	Private              int32  `json:"private"`     // 收藏夹类型 1私密 2公开
	Updatetime           int64  `json:"updatetime"`  // 收藏夹更新时间
}

type FavoriteFolderListResponse struct {
	httputils.PageResponse[*FavoriteFolderItem] `json:",inline"`
}

// FavoriteScriptRequest 收藏脚本
type FavoriteScriptRequest struct {
	mux.Meta `path:"/favorites/folders/:id/favorite" method:"POST"`
	ScriptID int64 `json:"script_id" binding:"required"`
	FolderID int64 `uri:"id" binding:"required" label:"收藏夹ID"` // 一次只能收藏到一个收藏夹
}

type FavoriteScriptResponse struct{}

// UnfavoriteScriptRequest 取消收藏脚本
type UnfavoriteScriptRequest struct {
	mux.Meta `path:"/favorites/folders/:id/favorite" method:"DELETE"`
	ScriptID int64 `form:"script_id" binding:"required"`
	FolderID int64 `uri:"id" label:"收藏夹ID"` // 一次只能从一个收藏夹移除，如果为0表示从所有收藏夹移除
}

type UnfavoriteScriptResponse struct{}

// FavoriteScriptListRequest 获取收藏夹脚本列表
type FavoriteScriptListRequest struct {
	mux.Meta              `path:"/favorites/scripts" method:"GET"`
	httputils.PageRequest `json:",inline"`
	FolderID              int64 `form:"folder_id" label:"收藏夹ID"` // 收藏夹ID，0表示所有的收藏
	UserID                int64 `form:"user_id" label:"用户ID"`    // 用户ID，0表示当前登录用户
}

type FavoriteScriptListResponse struct {
	httputils.PageResponse[*Script] `json:",inline"`
}

// FavoriteFolderDetailRequest 收藏夹详情
type FavoriteFolderDetailRequest struct {
	mux.Meta `path:"/favorites/folders/:id/detail" method:"GET"`
	ID       int64 `uri:"id" binding:"required"`
}

type FavoriteFolderDetailResponse struct {
	*FavoriteFolderItem `json:",inline"` // 收藏夹信息
}

// EditFolderRequest 编辑收藏夹
type EditFolderRequest struct {
	mux.Meta    `path:"/favorites/folders/:id" method:"PUT"`
	ID          int64  `uri:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=50" label:"收藏夹名称"`
	Description string `json:"description" binding:"max=200" label:"收藏夹描述"`
	Private     int32  `json:"private" binding:"omitempty,oneof=1 2" label:"私密收藏夹"` // 1私密 2公开
}

type EditFolderResponse struct {
}
