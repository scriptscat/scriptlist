package script_ctr

import (
	"context"
	"strconv"

	"github.com/cago-frame/cago/database/redis"
	"github.com/cago-frame/cago/pkg/limit"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"

	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type Favorite struct {
	limit limit.Limit
}

func NewFavorite() *Favorite {
	return &Favorite{
		limit: limit.NewCombinationLimit(limit.NewPeriodLimit(
			300, 2, redis.Default(), "limit:create:favorite:minute",
		), limit.NewPeriodLimit(
			3600, 5, redis.Default(), "limit:create:favorite:hour",
		)),
	}
}

// CreateFolder 创建收藏夹
func (f *Favorite) CreateFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error) {
	resp, err := f.limit.FuncTake(ctx, strconv.FormatInt(auth_svc.Auth().Get(ctx).UID, 10), func() (interface{}, error) {
		return script_svc.Favorite().CreateFolder(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return resp.(*api.CreateFolderResponse), nil
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

// FavoriteScriptList 获取收藏夹脚本列表
func (f *Favorite) FavoriteScriptList(ctx context.Context, req *api.FavoriteScriptListRequest) (*api.FavoriteScriptListResponse, error) {
	return script_svc.Favorite().FavoriteScriptList(ctx, req)
}
