package script_svc

import (
	"context"
	"time"

	"github.com/scriptscat/scriptlist/internal/repository/user_repo"

	"github.com/cago-frame/cago/pkg/consts"
	"github.com/cago-frame/cago/pkg/i18n"
	"github.com/cago-frame/cago/pkg/utils/httputils"
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
	// FavoriteFolderDetail 收藏夹详情
	FavoriteFolderDetail(ctx context.Context, req *api.FavoriteFolderDetailRequest) (*api.FavoriteFolderDetailResponse, error)
	// EditFolder 编辑收藏夹
	EditFolder(ctx context.Context, req *api.EditFolderRequest) (*api.EditFolderResponse, error)
	// FavoriteFolderScripts 获取指定收藏夹的所有脚本
	FavoriteFolderScripts(ctx context.Context, folderId int64) (*api.FavoriteFolderDetailResponse, []*api.Script, error)
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
	if err := script_repo.ScriptFavoriteFolder().Delete(ctx, req.ID, auth_svc.Auth().Get(ctx).UID); err != nil {
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

	list, total, err := script_repo.ScriptFavoriteFolder().FindPage(ctx, req.UserID, self, httputils.PageRequest{Page: 1, Size: 100})
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
			Private:     2,
		}}
	}
	response := &api.FavoriteFolderListResponse{
		PageResponse: httputils.PageResponse[*api.FavoriteFolderItem]{
			List:  make([]*api.FavoriteFolderItem, 0),
			Total: total,
		},
	}
	for _, item := range list {
		user, err := user_repo.User().Find(ctx, item.UserID)
		if err != nil {
			return nil, err
		}
		response.List = append(response.List, &api.FavoriteFolderItem{
			ID:          item.ID,
			UserInfo:    user.UserInfo(),
			Name:        item.Name,
			Description: item.Description,
			Count:       item.Count,
			Private:     item.Private,
			Updatetime:  item.Updatetime,
		})
	}
	return response, nil
}

// FavoriteScript 收藏脚本
func (f *favoriteSvc) FavoriteScript(ctx context.Context, req *api.FavoriteScriptRequest) (*api.FavoriteScriptResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	// 先检查收藏夹存不存在
	folder, err := script_repo.ScriptFavoriteFolder().Find(ctx, req.FolderID)
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
	exist, err := script_repo.ScriptFavorite().FindByFavoriteAndScriptID(ctx, user.UID, req.FolderID, req.ScriptID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		if exist.Status == consts.ACTIVE {
			return nil, i18n.NewError(ctx, code.ScriptFavoriteExist)
		}
		exist.Status = consts.ACTIVE
		if err := script_repo.ScriptFavorite().Update(ctx, exist); err != nil {
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
	if err := script_repo.ScriptFavorite().Create(ctx, m); err != nil {
		return nil, err
	}
	return nil, nil
}

// UnfavoriteScript 取消收藏脚本
func (f *favoriteSvc) UnfavoriteScript(ctx context.Context, req *api.UnfavoriteScriptRequest) (*api.UnfavoriteScriptResponse, error) {
	if req.FolderID == 0 {
		list, err := script_repo.ScriptFavorite().FindByUserIDAndScriptID(ctx, auth_svc.Auth().Get(ctx).UID, req.ScriptID)
		if err != nil {
			return nil, err
		}
		for _, item := range list {
			item.Status = consts.DELETE
			if err := script_repo.ScriptFavorite().Update(ctx, item); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	record, err := script_repo.ScriptFavorite().FindByFavoriteAndScriptID(ctx, auth_svc.Auth().Get(ctx).UID, req.FolderID, req.ScriptID)
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
	if err := script_repo.ScriptFavorite().Update(ctx, record); err != nil {
		return nil, err
	}
	return nil, nil
}

// FavoriteScriptList 获取收藏夹脚本列表
func (f *favoriteSvc) FavoriteScriptList(ctx context.Context, req *api.FavoriteScriptListRequest) (*api.FavoriteScriptListResponse, error) {
	user := auth_svc.Auth().Get(ctx)
	userId := req.UserID
	self := false
	if req.UserID == 0 {
		if user == nil {
			return nil, i18n.NewError(ctx, code.ScriptFavoriteMustUserID)
		}
		userId = user.UID
	}
	if user != nil {
		self = user.UID == userId
	}
	var (
		list  []*script_entity.ScriptFavorite
		total int64
		err   error
	)
	if req.FolderID == 0 {
		// 获取所有收藏的脚本
		list, total, err = script_repo.ScriptFavorite().FindByUserUnique(ctx, userId, self, req.PageRequest)
	} else {
		var folder *script_entity.ScriptFavoriteFolder
		folder, err = script_repo.ScriptFavoriteFolder().Find(ctx, req.FolderID)
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
		// 获取指定收藏夹的脚本
		list, total, err = script_repo.ScriptFavorite().FindByFavoriteFolder(ctx, req.FolderID, req.PageRequest)
	}
	if err != nil {
		return nil, err
	}
	resp := &api.FavoriteScriptListResponse{
		PageResponse: httputils.PageResponse[*api.Script]{
			List:  make([]*api.Script, 0),
			Total: total,
		},
	}
	// 遍历脚本列表，填充脚本信息
	for _, item := range list {
		script, err := script_repo.Script().Find(ctx, item.ScriptID)
		if err != nil {
			return nil, err
		}
		if script == nil {
			continue // 如果脚本不存在，跳过
		}
		scriptItem, err := Script().ToScript(ctx, script, false, "")
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, scriptItem)
	}
	return resp, nil
}

// FavoriteFolderDetail 收藏夹详情
func (f *favoriteSvc) FavoriteFolderDetail(ctx context.Context, req *api.FavoriteFolderDetailRequest) (*api.FavoriteFolderDetailResponse, error) {
	folder, err := script_repo.ScriptFavoriteFolder().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderNotFound)
	}
	if folder.Private == 1 {
		user := auth_svc.Auth().Get(ctx)
		if user == nil || user.UID != folder.UserID {
			return nil, i18n.NewForbiddenError(ctx, code.ScriptFavoriteFolderNotFound)
		}
	}
	resp := &api.FavoriteFolderItem{
		ID:          folder.ID,
		Name:        folder.Name,
		Description: folder.Description,
		Count:       folder.Count,
		Private:     folder.Private,
		Updatetime:  folder.Updatetime,
	}
	user, err := user_repo.User().Find(ctx, folder.UserID)
	if err != nil {
		return nil, err
	}
	resp.UserInfo = user.UserInfo()
	return &api.FavoriteFolderDetailResponse{
		FavoriteFolderItem: resp,
	}, nil
}

// EditFolder 编辑收藏夹
func (f *favoriteSvc) EditFolder(ctx context.Context, req *api.EditFolderRequest) (*api.EditFolderResponse, error) {
	folder, err := script_repo.ScriptFavoriteFolder().Find(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if folder == nil {
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderNotFound)
	}

	// 如果是默认收藏夹，不能修改名字
	if folder.Name == "默认收藏夹" {
		if req.Name != "默认收藏夹" {
			return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderCannotEdit)
		}
	} else if req.Name == "默认收藏夹" {
		// 也不能修改名字为默认收藏夹
		return nil, i18n.NewError(ctx, code.ScriptFavoriteFolderCannotEdit)
	}

	if folder.UserID != auth_svc.Auth().Get(ctx).UID {
		return nil, i18n.NewForbiddenError(ctx, code.ScriptFavoriteFolderNotFound)
	}
	folder.Name = req.Name
	folder.Description = req.Description
	folder.Private = req.Private
	if err := script_repo.ScriptFavoriteFolder().Update(ctx, folder); err != nil {
		return nil, err
	}
	return &api.EditFolderResponse{}, nil
}

// FavoriteFolderScripts 获取指定收藏夹的所有脚本
func (f *favoriteSvc) FavoriteFolderScripts(ctx context.Context, folderId int64) (*api.FavoriteFolderDetailResponse, []*api.Script, error) {
	folder, err := f.FavoriteFolderDetail(ctx, &api.FavoriteFolderDetailRequest{ID: folderId})
	if err != nil {
		return nil, nil, err
	}
	if folder == nil {
		return nil, nil, i18n.NewError(ctx, code.ScriptFavoriteFolderNotFound)
	}
	if folder.Private == 1 {
		user := auth_svc.Auth().Get(ctx)
		if user == nil || user.UID != folder.UserID {
			return nil, nil, i18n.NewForbiddenError(ctx, code.ScriptFavoriteFolderNotFound)
		}
	}
	list, _, err := script_repo.ScriptFavorite().FindByFavoriteFolder(ctx, folderId, httputils.PageRequest{Size: -1})
	if err != nil {
		return nil, nil, err
	}
	resp := make([]*api.Script, 0)
	for _, item := range list {
		script, err := script_repo.Script().Find(ctx, item.ScriptID)
		if err != nil {
			return nil, nil, err
		}
		if script == nil {
			continue // 如果脚本不存在，跳过
		}
		scriptItem, err := Script().ToScript(ctx, script, false, "")
		if err != nil {
			return nil, nil, err
		}
		resp = append(resp, scriptItem)
	}
	return folder, resp, nil
}
