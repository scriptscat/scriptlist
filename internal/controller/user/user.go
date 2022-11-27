package user

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/user"
	service "github.com/scriptscat/scriptlist/internal/service/user"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// CurrentUser 获取当前登录的用户信息
func (u *User) CurrentUser(ctx context.Context, req *api.CurrentUserRequest) (*api.CurrentUserResponse, error) {
	resp, err := service.User().UserInfo(ctx, service.Auth().Get(ctx).UID)
	if err != nil {
		return nil, err
	}
	return &api.CurrentUserResponse{InfoResponse: resp}, nil
}

// Info 获取指定用户信息
func (u *User) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return service.User().UserInfo(ctx, req.UID)
}
