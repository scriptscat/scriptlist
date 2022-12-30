package user_ctr

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/user"
	"github.com/scriptscat/scriptlist/internal/service/user_svc"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// CurrentUser 获取当前登录的用户信息
func (u *User) CurrentUser(ctx context.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	resp, err := user_svc.User().UserInfo(ctx, user_svc.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	return &api.CurrentUserResponse{InfoResponse: resp}, nil
}

// Info 获取指定用户信息
func (u *User) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return user_svc.User().UserInfo(ctx, req.UID)
}
