package script

import (
	"github.com/codfrm/cago/server/mux"
)

// CreateFolderRequest 创建收藏夹
type CreateFolderRequest struct {
	mux.Meta    `path:"/favorites/folders" method:"POST"`
	Name        string `json:"name" binding:"required,max=50" label:"收藏夹名称"`
	Description string `json:"description" binding:"max=200" label:"收藏夹描述"`
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
	ScriptID int64 `form:"script_id" label:"脚本ID"` // 可选，用于检查指定脚本是否在收藏夹中
}

type FavoriteFolderItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"` // 收藏夹描述
	Count       int    `json:"count"`       // 收藏夹中脚本数量
	Favorited   bool   `json:"favorited"`   // 如果传了script_id，表示该脚本是否在此收藏夹中
}

type FavoriteFolderListResponse struct {
	List []FavoriteFolderItem `json:"list"`
}

// FavoriteScriptRequest 收藏脚本
type FavoriteScriptRequest struct {
	mux.Meta `path:"/scripts/:id/favorite" method:"POST"`
	ScriptID int64 `uri:"id" binding:"required"`
	FolderID int64 `json:"folder_id" binding:"required" label:"收藏夹ID"` // 一次只能收藏到一个收藏夹
}

type FavoriteScriptResponse struct{}

// UnfavoriteScriptRequest 取消收藏脚本
type UnfavoriteScriptRequest struct {
	mux.Meta `path:"/scripts/:id/unfavorite" method:"POST"`
	ScriptID int64 `uri:"id" binding:"required"`
	FolderID int64 `json:"folder_id" binding:"required" label:"收藏夹ID"` // 一次只能从一个收藏夹移除
}

type UnfavoriteScriptResponse struct{}
