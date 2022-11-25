package user

import (
	"context"

	api "github.com/scriptscat/scriptlist/internal/api/user"
)

type IUser interface {
	// Info 用户信息
	Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error)
}

type user struct {
}

var defaultUser = &user{}

func User() IUser {
	return defaultUser
}

// Info 用户信息
func (u *user) Info(ctx context.Context, req *api.InfoRequest) (*api.InfoResponse, error) {
	return nil, nil
}
