package script_ctr

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Favorite struct {
}

func NewFavorite() *Favorite {
	return &Favorite{}
}

// CreateFolder 创建收藏夹
func (f *Favorite) CreateFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error) {
	return script_svc.Favorite().CreateFolder(ctx, req)
}

// DeleteFolder 删除收藏夹
func (f *Favorite) DeleteFolder(ctx context.Context, req *api.DeleteFolderRequest) (*api.DeleteFolderResponse, error) {
	return script_svc.Favorite().DeleteFolder(ctx, req)
}

// FavoriteFolderList 收藏夹列表
func (f *Favorite) FavoriteFolderList(ctx context.Context, req *api.FavoriteFolderListRequest) (*api.FavoriteFolderListResponse, error) {
	return script_svc.Favorite().FavoriteFolderList(ctx, req)
}

// FavoriteScript 收藏脚本
func (f *Favorite) FavoriteScript(ctx context.Context, req *api.FavoriteScriptRequest) (*api.FavoriteScriptResponse, error) {
	return script_svc.Favorite().FavoriteScript(ctx, req)
}

// UnfavoriteScript 取消收藏脚本
func (f *Favorite) UnfavoriteScript(ctx context.Context, req *api.UnfavoriteScriptRequest) (*api.UnfavoriteScriptResponse, error) {
	return script_svc.Favorite().UnfavoriteScript(ctx, req)
}
