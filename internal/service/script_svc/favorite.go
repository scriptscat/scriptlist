package script_svc

import (
	"context"
	"time"

	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/scriptscat/scriptlist/internal/api/script"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
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
	// FavoriteScriptList 获取收藏夹脚本列表
	FavoriteScriptList(ctx context.Context, req *api.FavoriteScriptListRequest) (*api.FavoriteScriptListResponse, error)
}

type favoriteSvc struct {
}

var defaultFavorite = &favoriteSvc{}

func Favorite() FavoriteSvc {
	return defaultFavorite
}

// CreateFolder 创建收藏夹
func (f *favoriteSvc) CreateFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error) {
	list, err := f.FavoriteFolderList(ctx, &api.FavoriteFolderListRequest{})
	if err != nil {
		return nil, err
	}
	if list.Total >= 10 {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderLimitExceeded)
	}
	return f.createFolder(ctx, req)
}

func (f *favoriteSvc) createFolder(ctx context.Context, req *api.CreateFolderRequest) (*api.CreateFolderResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	m := &script_entity.ScriptFavoriteFolder{
		ID:          0,
		Name:        req.Name,
		Description: req.Description,
		UserID:      user.UID,
		Private:     req.Private,
		Count:       0,
		Status:      consts.ACTIVE,
		Createtime:  time.Now().Unix(),
	}
	err := script_repo.ScriptFavoriteFolder().Create(ctx, m)
	if err != nil {
		return nil, err
	}
	return &api.CreateFolderResponse{ID: m.ID}, nil
}

// DeleteFolder 删除收藏夹
func (f *favoriteSvc) DeleteFolder(ctx context.Context, req *api.DeleteFolderRequest) (*api.DeleteFolderResponse, error) {
	if err := script_repo.NewScriptFavoriteFolder().Delete(ctx, req.ID, auth_svc.Auth().Get(ctx).UID); err != nil {
		return nil, err
	}
	return nil, nil
}

// FavoriteFolderList 收藏夹列表
func (f *favoriteSvc) FavoriteFolderList(ctx context.Context, req *api.FavoriteFolderListRequest) (*api.FavoriteFolderListResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	self := false
	if user != nil {
		if req.UserID == 0 {
			req.UserID = user.UID
		}
		self = user.UID == req.UserID
	}

	if req.UserID == 0 {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteMustUserID)
	}

	list, total, err := script_repo.NewScriptFavoriteFolder().FindPage(ctx, req.UserID, self, httputils.PageRequest{Page: 1, Size: 100})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 && self {
		// 如果没有收藏夹，创建一个默认收藏夹
		createReq := &api.CreateFolderRequest{
			Name:        "默认收藏夹",
			Description: "这是一个默认的收藏夹",
		}
		resp, err := f.createFolder(ctx, createReq)
		if err != nil {
			return nil, err
		}
		list = []*script_entity.ScriptFavoriteFolder{{
			ID:          resp.ID,
			Name:        "默认收藏夹",
			Description: createReq.Description,
			UserID:      req.UserID,
		}}
	}
	response := &api.FavoriteFolderListResponse{
		PageResponse: httputils.PageResponse[*api.FavoriteFolderItem]{
			List:  make([]*api.FavoriteFolderItem, 0),
			Total: total,
		},
	}
	for _, item := range list {
		response.List = append(response.List, &api.FavoriteFolderItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Count:       item.Count,
		})
	}
	return response, nil
}

// FavoriteScript 收藏脚本
func (f *favoriteSvc) FavoriteScript(ctx context.Context, req *api.FavoriteScriptRequest) (*api.FavoriteScriptResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	// 先检查收藏夹存不存在
	folder, err := script_repo.NewScriptFavoriteFolder().Find(ctx, req.FolderID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderNotFound)
	}
	// 限制每个收藏夹最多收藏100个脚本
	if folder.Count >= 100 {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteLimitExceeded)
	}
	// 检查有没有收藏过
	exist, err := script_repo.NewScriptFavorite().FindByFavoriteAndScriptID(ctx, user.UID, req.FolderID, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		if exist.Status == consts.ACTIVE {
			return nil, i18n.NewError(ctx, code.ScriptFavoriteExist)
		}
		exist.Status = consts.ACTIVE
		if err := script_repo.NewScriptFavorite().Update(ctx, exist); err != nil {
			return nil, err
		}
		return nil, nil
	}
	m := &script_entity.ScriptFavorite{
		UserID:           user.UID,
		ScriptID:         req.ScriptID,
		FavoriteFolderID: req.FolderID,
		Status:           consts.ACTIVE,
		Createtime:       time.Now().Unix(),
	}
	if err := script_repo.NewScriptFavorite().Create(ctx, m); err != nil {
		return nil, err
	}
	return nil, nil
}

// UnfavoriteScript 取消收藏脚本
func (f *favoriteSvc) UnfavoriteScript(ctx context.Context, req *api.UnfavoriteScriptRequest) (*api.UnfavoriteScriptResponse, error) {
	record, err := script_repo.NewScriptFavorite().FindByFavoriteAndScriptID(ctx, auth_svc.Auth().Get(ctx).UID, req.FolderID, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteNotFound)
	}
	if record.Status == consts.DELETE {
		return nil, nil
	}
	record.Status = consts.DELETE
	if err := script_repo.NewScriptFavorite().Update(ctx, record); err != nil {
		return nil, err
	}
	return nil, nil
}

// FavoriteScriptList 获取收藏夹脚本列表
func (f *favoriteSvc) FavoriteScriptList(ctx context.Context, req *api.FavoriteScriptListRequest) (*api.FavoriteScriptListResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	folder, err := script_repo.NewScriptFavoriteFolder().Find(ctx, req.FolderID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderNotFound)
	}
	if folder.Private == 1 {
		if user == nil {
			return nil, i18n.NewForbiddenError(ctx, code.ScriptFavoriteFolderNotFound)
		}
		if user.UID != folder.UserID {
			return nil, i18n.NewForbiddenError(ctx, code.ScriptFavoriteFolderNotFound)
		}
	}

	return nil, nil
}
