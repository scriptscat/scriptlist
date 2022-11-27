package user

import (
	"context"
	"strconv"

	"github.com/codfrm/cago/pkg/i18n"
	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/pkg/code"
	"github.com/scriptscat/scriptlist/internal/repository"
)

type IUser interface {
	// UserInfo 获取用户信息
	UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error)
}

type user struct {
}

var defaultUser = &user{}

func User() IUser {
	return defaultUser
}

// UserInfo 获取用户信息
func (u *user) UserInfo(ctx context.Context, uid int64) (*api.InfoResponse, error) {
	user, err := repository.User().Find(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, i18n.NewError(ctx, code.UserNotFound)
	}
	return &api.InfoResponse{
		UID:      user.UID,
		Username: user.Username,
		Avatar:   "/api/v2/user/avatar/" + strconv.FormatInt(user.UID, 10),
	}, nil
}
