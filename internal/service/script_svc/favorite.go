package script_svc

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
)

type FavoriteSvc interface {
	// CreateFolder 创建收藏夹
	CreateFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error)
	// DeleteFolder 删除收藏夹
	DeleteFolder(ctx context.Context, req *api.DeleteFolderRequest) (*api.DeleteFolderResponse, error)
	// FavoriteFolderList 收藏夹列表
	FavoriteFolderList(ctx context.Context, req *api.FavoriteFolderListRequest) (*api.FavoriteFolderListResponse, error)
	// FavoriteScript 收藏脚本
	FavoriteScript(ctx context.Context, req *api.FavoriteScriptRequest) (*api.FavoriteScriptResponse, error)
	// UnfavoriteScript 取消收藏脚本
	UnfavoriteScript(ctx context.Context, req *api.UnfavoriteScriptRequest) (*api.UnfavoriteScriptResponse, error)
}

type favoriteSvc struct {
}

var defaultFavorite = &favoriteSvc{}

func Favorite() FavoriteSvc {
	return defaultFavorite
}

// CreateFolder 创建收藏夹
func (f *favoriteSvc) CreateFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error) {
	return nil, nil
}

// DeleteFolder 删除收藏夹
func (f *favoriteSvc) DeleteFolder(ctx context.Context, req *api.DeleteFolderRequest) (*api.DeleteFolderResponse, error) {
	return nil, nil
}

// FavoriteFolderList 收藏夹列表
func (f *favoriteSvc) FavoriteFolderList(ctx context.Context, req *api.FavoriteFolderListRequest) (*api.FavoriteFolderListResponse, error) {
	return nil, nil
}

// FavoriteScript 收藏脚本
func (f *favoriteSvc) FavoriteScript(ctx context.Context, req *api.FavoriteScriptRequest) (*api.FavoriteScriptResponse, error) {
	return nil, nil
}

// UnfavoriteScript 取消收藏脚本
func (f *favoriteSvc) UnfavoriteScript(ctx context.Context, req *api.UnfavoriteScriptRequest) (*api.UnfavoriteScriptResponse, error) {
	return nil, nil
}
