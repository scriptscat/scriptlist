package user_svc

import (
	"context"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/scriptscat/scriptlist/internal/api/script"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/model"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository/script_repo"
	"github.com/scriptscat/scriptlist/internal/repository/user_repo"
	"github.com/scriptscat/scriptlist/internal/service/auth_svc"
	"github.com/scriptscat/scriptlist/internal/service/script_svc"
)

type UserSvc interface {
	// UserInfo 获取用户信息
	UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error)
	// Script 用户脚本列表
	Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error)
}

type userSvc struct {
}

var defaultUser = &userSvc{}

func User() UserSvc {
	return defaultUser
}

// UserInfo 获取用户信息
func (u *userSvc) UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error) {
	user, err := user_repo.User().Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, i18n.NewError(ctx, code.UserNotFound)
	}
	return &api.InfoResponse{
		UserID:      user.UID,
		Username:    user.Username,
		Avatar:      user.Avatar(),
		IsAdmin:     model.AdminLevel(user.Adminid),
		EmailStatus: user.Emailstatus,
	}, nil
}

// Script 用户脚本列表
func (u *userSvc) Script(ctx context.Context, req *api.ScriptRequest) (*api.ScriptResponse, error) {
	self := auth_svc.Auth().Get(ctx).UID == req.UID
	resp, total, err := script_repo.Script().Search(ctx, &script_repo.SearchOptions{
		UserID:   req.UID,
		Keyword:  req.Keyword,
		Type:     req.ScriptType,
		Sort:     req.Sort,
		Category: make([]int64, 0),
		Self:     self,
	}, req.PageRequest)
	if err != nil {
		return nil, err
	}
	list := make([]*script.Script, 0)
	for _, item := range resp {
		data, err := script_svc.Script().ToScript(ctx, item, false, "")
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return &api.ScriptResponse{
		PageResponse: httputils.PageResponse[*script.Script]{
			List:  list,
			Total: total,
		},
	}, nil
}
